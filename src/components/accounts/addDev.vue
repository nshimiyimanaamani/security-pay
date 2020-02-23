<template>
  <b-tabs content-class="mt-2" class="addDev" fill lazy>
    <b-tab title="Create Developer" active>
      <b-form class="accountForm" @submit.prevent="create">
        <b-form-group id="input-group-1" label="Email:" label-for="input-1">
          <b-form-input
            id="input-1"
            type="email"
            v-model="form.email"
            size="sm"
            required
            placeholder="Enter Email address..."
          />
        </b-form-group>
        <b-form-group id="input-group-2" label="Password:" label-for="input-2">
          <b-form-input
            id="input-2"
            type="password"
            size="sm"
            v-model="form.password"
            required
            placeholder="Enter Password..."
          />
        </b-form-group>

        <b-form-group class="m-0">
          <b-button
            variant="info"
            class="float-right"
            :disabled="(form.email && form.password)?false:true"
            size="sm"
            type="submit"
          >
            {{state.creating ? 'creating' : "create"}}
            <b-spinner v-show="state.creating" small type="grow"></b-spinner>
          </b-button>
        </b-form-group>
      </b-form>
    </b-tab>
    <b-tab title="List all Developers">
      <b-table
        id="Dev-table"
        bordered
        striped
        hover
        small
        show-empty
        responsive
        :items="loadData"
        :fields="table.fields"
        :busy.sync="state.tableLoad"
        @row-contextmenu="menu"
      >
        <template v-slot:table-busy>
          <div class="text-center my-2">
            <loader />
          </div>
        </template>
      </b-table>
    </b-tab>
    <vue-menu
      :elementId="'dev-rightmenu'"
      :options="rightMenu.options"
      :ref="'dev_rightMenu'"
      @option-clicked="optionClicked"
    ></vue-menu>
    <b-modal v-model="change_pswd_modal.show" title="Change Password" hide-footer>
      <b-form @submit.prevent="changePassword">
        <b-form-group id="input-group-1" label="New Password:" label-for="input-1">
          <b-form-input
            id="input-1"
            size="sm"
            required
            v-model="form.newPwd"
            placeholder="New Password..."
          ></b-form-input>
        </b-form-group>
        <b-form-group class="m-0">
          <b-button
            variant="info"
            :disabled="form.newPwd?false:true"
            size="sm"
            class="float-right"
            type="submit"
          >
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
import loader from "../loader";
export default {
  name: "add-dev",
  components: {
    loader
  },
  data() {
    return {
      form: {
        email: null,
        password: null,
        newPwd: null
      },
      table: {
        fields: [
          { key: "email", label: "Email", tdClass: "text-normal" },
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
  methods: {
    create() {
      this.state.creating = true;
      const account = jwt.decode(sessionStorage.token).account;
      this.axios
        .post("/accounts/developers", {
          account: account,
          email: this.form.email,
          password: this.form.password
        })
        .then(res => {
          this.$snotify.info("Developer successfully created...");
          const message = `Password For ${res.data.email} is: ${res.data.password}`;
          this.$bvModal
            .msgBoxOk(message, {
              title: "Developer Details",
              size: "sm",
              buttonSize: "md",
              okVariant: "info",
              headerClass: "p-2 border-bottom-0",
              footerClass: "p-2 border-top-0",
              centered: true
            })
            .then(res => console.log(""));
        })
        .catch(err => {
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
        })
        .finally(() => {
          this.state.creating = false;
          this.form = { email: null, password: null };
        });
    },
    loadData() {
      this.state.tableLoad = true;
      const promise = this.axios.get(
        "/accounts/developers?offset=0&limit=1000"
      );
      return promise
        .then(res => {
          return res.data.Developers;
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
      this.$refs.dev_rightMenu.showMenu(evt, house);
    },
    optionClicked(data) {
      if (data.option.slug == "delete") {
        this.state.tableLoad = true;
        this.axios
          .delete(`/accounts/${data.email}`)
          .then(res => {
            this.loadData();
            this.$snotify.info("Developer deleted Succesfully");
          })
          .catch(err => {
            const error = err.response
              ? err.response.data.error || err.response.data
              : null;
            if (error) this.$snotify.error(error);
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
        .put("/accounts/developers/creds/" + data.email, {
          password: this.form.newPwd
        })
        .then(res => {
          this.$snotify.info(res.data.message);
          this.loadData();
        })
        .catch(err => {
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
        })
        .finally(() => {
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
.addDev .nav-link.active {
  background-color: white !important;
}
.addDev .accountForm {
  width: auto;
  border: 1px solid #dee2e6;
  border-radius: 5px;
  padding: 1rem;
  margin: 1rem;
}
</style>