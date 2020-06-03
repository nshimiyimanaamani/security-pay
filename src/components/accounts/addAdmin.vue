<template>
  <b-tabs content-class="mt-2" class="addAdmin" fill lazy>
    <b-tab title="Create Admin" active>
      <b-form class="accountForm" @submit.prevent="create">
        <b-form-group id="input-group-1" label="Email:" label-for="input-1">
          <b-form-input
            id="input-1"
            type="email"
            v-model="form.email"
            class="br-2"
            required
            placeholder="Enter Email address..."
          />
        </b-form-group>

        <b-form-group>
          <b-button
            :disabled="form.email?false: true"
            variant="info"
            class="w-100 br-2"
            type="submit"
          >
            {{state.creating ? 'Creating' : "Create"}}
            <i
              v-show="state.creating"
              class="fa fa-spinner fa-spin"
            />
          </b-button>
        </b-form-group>
      </b-form>
    </b-tab>
    <b-tab title="List all Admins">
      <b-table
        id="Dev-table"
        bordered
        hover
        small
        responsive
        show-empty
        head-variant="light"
        thead-class="text-uppercase"
        :items="loadData"
        :fields="table.fields"
        :busy.sync="state.tableLoad"
        @row-contextmenu="menu"
      >
        <template v-slot:table-busy>
          <vue-load />
        </template>
      </b-table>
    </b-tab>
    <vue-menu
      :elementId="'admin-rightmenu'"
      :options="rightMenu.options"
      :ref="'admin_rightMenu'"
      @option-clicked="optionClicked"
    ></vue-menu>
    <b-modal v-model="change_pswd_modal.show" title="Change Password" hide-footer>
      <b-form @submit.prevent="changePassword">
        <b-form-group id="input-group-1" label="New Password:" label-for="input-1">
          <b-form-input
            id="input-1"
            required
            v-model="form.newPwd"
            class="br-2"
            placeholder="New Password..."
          ></b-form-input>
        </b-form-group>
        <b-form-group class="mb-0">
          <b-button
            variant="info"
            class="float-right br-2"
            :disabled="form.newPwd?false:true"
            type="submit"
          >
            {{state.changing ? 'changing' : "change"}}
            <i
              v-if="state.changing"
              class="fa fa-spinner fa-spin"
            />
          </b-button>
        </b-form-group>
      </b-form>
    </b-modal>
  </b-tabs>
</template>

<script>
export default {
  name: "add-agent",
  data() {
    return {
      form: {
        email: null,
        newPwd: null
      },
      table: {
        fields: [
          { key: "email", label: "Email", tdClass: "text-normal" },
          { key: "cell", label: "cell" },
          { key: "role", label: "Role" }
        ]
      },
      state: {
        creating: false,
        tableLoad: false,
        changing: false
      },
      change_pswd_modal: {
        show: false,
        data: null
      },
      rightMenu: {
        options: [{ name: "Change Password", slug: "changePwd" }]
      }
    };
  },
  computed: {
    user() {
      return this.$store.getters.userDetails;
    }
  },
  methods: {
    create() {
      this.state.creating = true;
      const account = this.user.account;
      this.axios
        .post("/accounts/admin", {
          account: account,
          email: this.form.email
        })
        .then(res => {
          this.state.creating = false;
          this.form = { email: null, password: null };
          this.$snotify.info("Admin successfully created...");
          const message = `Password For ${res.data.email} is: ${res.data.password}`;
          this.$bvModal
            .msgBoxOk(message, {
              title: "Admin Details",
              size: "md",
              buttonSize: "sm",
              okVariant: "info",
              headerClass: "p-2 border-bottom-0",
              footerClass: "p-2 border-top-0",
              centered: true
            })
            .then(res => {
              this.state.creating = false;
              this.form = { email: null, password: null };
            });
        })
        .catch(err => {
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
          this.state.creating = false;
          this.form = { email: null, password: null };
        });
    },
    async loadData() {
      this.state.tableLoad = true;
      const total = await this.$getTotal("/accounts/managers?offset=0&limit=0");
      const promise = this.axios.get(
        "/accounts/managers?offset=0&limit=" + total
      );
      return promise
        .then(res => {
          return res.data.Managers;
        })
        .catch(err => {
          return [];
        })
        .finally(() => {
          this.state.tableLoad = false;
        });
    },
    menu(house, index, evt) {
      evt.preventDefault();
      this.$refs.admin_rightMenu.showMenu(evt, house);
    },
    optionClicked(data) {
      if (data.option.slug == "delete") {
        this.state.tableLoad = true;
        this.axios
          .delete(`/accounts/${this.user.username}/users/${data.email}`)
          .then(res => {
            this.loadData();
            this.$snotify.info("Agent deleted Succesfully");
          })
          .finally(() => {
            this.state.tableLoad = false;
          });
      } else if (data.option.slug == "changePwd") {
        this.change_pswd_modal.show = true;
        this.change_pswd_modal.data = { ...data.item };
      }
    },
    changePassword() {
      this.state.changing = true;
      const data = this.change_pswd_modal.data;
      this.axios
        .put("/accounts/admin/creds/" + data.email, {
          password: this.form.newPwd
        })
        .then(res => {
          this.$snotify.info(res.data.message);
          this.loadData();
        })
        .finally(() => {
          this.state.changing = false;
          this.change_pswd_modal.show = false;
        });
    }
  }
};
</script>

<style lang="scss">
.addAdmin {
  .accountForm {
    width: auto;
    border: 1px solid #dee2e6;
    border-radius: 5px;
    padding: 1rem;

    & > .form-group {
      max-width: 300px;
      margin: auto;
      margin-bottom: 1rem;
    }
  }
}
</style>