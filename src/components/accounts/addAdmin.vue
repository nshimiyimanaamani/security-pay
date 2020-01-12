<template>
  <b-tabs content-class="mt-2" class="addAdmin" fill lazy>
    <b-tab title="Create Admin" active>
      <b-form class="accountForm" @submit.prevent="create">
        <b-form-group id="input-group-1" label="Email:" label-for="input-1">
          <b-form-input
            id="input-1"
            type="email"
            v-model="form.email"
            required
            placeholder="Enter Email address..."
          />
        </b-form-group>

        <b-form-group class="m-0">
          <b-button variant="info" class="float-right" type="submit">
            {{state.creating ? 'Creating' : "Create"}}
            <b-spinner v-show="state.creating" small type="grow"></b-spinner>
          </b-button>
        </b-form-group>
      </b-form>
    </b-tab>
    <b-tab title="List all Admins">
      <b-table
        id="Dev-table"
        bordered
        striped
        hover
        small
        show-empty
        :items="loadData"
        :fields="table.fields"
        :busy.sync="state.tableLoad"
        @row-contextmenu="menu"
      >
        <template v-slot:table-busy>
          <div class="text-center my-2">
            <b-spinner class="align-middle"></b-spinner>
            <strong>Loading...</strong>
          </div>
        </template>
      </b-table>
    </b-tab>
    <vue-simple-context-menu
      :elementId="'admin-rightmenu'"
      :options="rightMenu.options"
      :ref="'admin_rightMenu'"
      @option-clicked="optionClicked"
    ></vue-simple-context-menu>
    <b-modal v-model="change_pswd_modal.show" title="Change Password" hide-footer>
      <b-form @submit.prevent="changePassword">
        <b-form-group id="input-group-1" label="New Password:" label-for="input-1">
          <b-form-input id="input-1" required v-model="form.newPwd" placeholder="New Password..."></b-form-input>
        </b-form-group>
        <b-form-group class="m-0">
          <b-button variant="info" class="float-right" type="submit">
            {{state.changing ? 'Changing' : "Change"}}
            <b-spinner v-show="state.changing" small type="grow"></b-spinner>
          </b-button>
        </b-form-group>
      </b-form>
    </b-modal>
  </b-tabs>
</template>

<script>
const jwt = require("jsonwebtoken");
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
          { key: "email", label: "Email" },
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
        options: [
          { name: "Change Password", slug: "changePwd" },
          { name: "Delete", slug: "delete" }
        ]
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    userDetails() {
      return this.$store.getters.userDetails;
    }
  },
  methods: {
    create() {
      this.state.creating = true;
      const account = jwt.decode(sessionStorage.token).account;
      this.axios
        .post(this.endpoint + "/accounts/admin", {
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
            .then(res => console.log(""));
        })
        .catch(err => {
          if (navigator.onLine) {
            const error = isNullOrUndefined(err.response)
              ? "an error occured"
              : err.response.data.error || err.response.data;
            this.$snotify.error(error);
          } else {
            this.$snotify.error("Please connect to the internet");
          }
          this.state.creating = false;
        });
    },
    loadData() {
      this.state.tableLoad = true;
      const promise = this.axios.get(
        this.endpoint + "/accounts/managers?offset=0&limit=10"
      );
      return promise
        .then(res => {
          this.state.tableLoad = false;
          return res.data.Managers;
        })
        .catch(err => {
          this.state.tableLoad = false;
          return [];
        });
    },
    menu(house, index, evt) {
      evt.preventDefault();
      Object.defineProperty(event, "pageX", {
        value: event.pageX - 410,
        writable: true
      });
      Object.defineProperty(event, "pageY", {
        value: event.pageY - 110,
        writable: true
      });
      this.$refs.admin_rightMenu.showMenu(evt, house);
    },
    optionClicked(data) {
      if (data.option.slug == "delete") {
        this.state.tableLoad = true;
        this.axios
          .delete(
            this.endpoint +
              `/accounts/${this.userDetails.username}/users/${data.email}`
          )
          .then(res => {
            this.loadData();
            this.state.tableLoad = false;
            this.$snotify.info("Agent deleted Succesfully");
            console.log(res.data);
          })
          .catch(err => {
            if (navigator.onLine) {
              const error = isNullOrUndefined(err.response)
                ? "an error occured"
                : err.response.data.error || err.response.data;
              this.$snotify.error(error);
            } else {
              this.$snotify.error("Please connect to the internet");
            }
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
        .put(this.endpoint + "/accounts/admin/creds/" + data.email, {
          password: this.form.newPwd
        })
        .then(res => {
          this.$snotify.info(res.data.message);
          this.state.changing = false;
          this.change_pswd_modal.show = false;
          this.loadData();
        })
        .catch(err => {
          if (navigator.onLine) {
            const error = isNullOrUndefined(err.response)
              ? "an error occured"
              : err.response.data.error || err.response.data;
            this.$snotify.error(error);
          } else {
            this.$snotify.error("Please connect to the internet");
          }
          this.state.changing = false;
          this.change_pswd_modal.show = false;
        });
    },
    toCapital(string) {
      string.toLowerCase();
      return string.charAt(0).toUpperCase() + string.slice(1);
    }
  }
};
</script>

<style>
.addAdmin .nav-link.active {
  background-color: white !important;
}
.addAdmin .accountForm {
  width: auto;
  border: 1px solid #dee2e6;
  border-radius: 5px;
  padding: 1rem;
  margin: 1rem;
}
</style>