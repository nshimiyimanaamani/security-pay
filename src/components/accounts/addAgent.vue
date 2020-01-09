<template>
  <b-tabs content-class="mt-2" class="addAgent" fill lazy>
    <b-tab title="Create Agent" active>
      <b-form class="accountForm" @submit.prevent="create">
        <b-row>
          <b-col>
            <b-form-group id="input-group-1" label="First name:" label-for="input-1">
              <b-form-input id="input-1" required v-model="form.fname" placeholder="First name..."></b-form-input>
            </b-form-group>
          </b-col>
          <b-col>
            <b-form-group id="input-group-2" label="Last name:" label-for="input-2">
              <b-form-input id="input-2" required v-model="form.lname" placeholder="Last name..."></b-form-input>
            </b-form-group>
          </b-col>
        </b-row>
        <b-form-group id="input-group-3" label="Phone Number:" label-for="input-3">
          <b-form-input
            id="input-3"
            type="number"
            v-model="form.phone"
            required
            placeholder="Enter Phone number..."
          ></b-form-input>
        </b-form-group>

        <b-form-group id="input-group-4" label="Sector:" label-for="input-4">
          <b-form-select id="input-4" v-model="form.select.sector" style="font-size: 15px">
            <template v-slot:first>
              <option :value="null" disabled>-- Please select sector --</option>
            </template>
            <option value="Remera">Remera</option>
          </b-form-select>
        </b-form-group>

        <b-form-group id="input-group-5" label="Cell:" label-for="input-5">
          <b-form-select
            id="input-5"
            v-model="form.select.cell"
            :options="cellOptions"
            style="font-size: 15px"
          >
            <template v-slot:first>
              <option :value="null" disabled>-- Please select cell --</option>
            </template>
          </b-form-select>
        </b-form-group>

        <b-form-group id="input-group-5" label="Village:" label-for="input-5">
          <b-form-select
            id="input-5"
            v-model="form.select.village"
            :options="villageOptions"
            style="font-size: 15px"
          >
            <template v-slot:first>
              <option :value="null" disabled>-- Please select village --</option>
            </template>
          </b-form-select>
        </b-form-group>
        <b-form-group class="m-0">
          <b-button variant="info" class="float-right" type="submit">
            {{state.creating ? 'Creating' : "Create"}}
            <b-spinner v-show="state.creating" small type="grow"></b-spinner>
          </b-button>
        </b-form-group>
      </b-form>
    </b-tab>
    <b-tab title="List all Agents">
      <b-table
        id="data-table"
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
      :elementId="'agent-rightmenu'"
      :options="rightMenu.options"
      :ref="'agent_rightMenu'"
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
const { Village } = require("rwanda");
const jwt = require("jsonwebtoken");
export default {
  name: "add-agent",
  data() {
    return {
      form: {
        fname: null,
        lname: null,
        phone: null,
        newPwd: null,
        select: {
          sector: null,
          cell: null,
          village: null
        }
      },
      table: {
        fields: [
          { key: "telephone", label: "Phone Number" },
          { key: "sector", label: "sector" },
          { key: "cell", label: "cell" },
          { key: "village", label: "village" }
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
          { name: "Edit", slug: "edit" },
          { name: "Change Password", slug: "changePwd" },
          { name: "Update", slug: "update" },
          { name: "Delete", slug: "delete" }
        ]
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    cellOptions() {
      const sector = this.form.select.sector;
      if (sector) {
        return this.$store.getters.getCellsArray;
      }
    },
    villageOptions() {
      const sector = this.form.select.sector;
      const cell = this.form.select.cell;
      if (sector && cell) {
        return Village("Kigali", "Gasabo", sector, cell).sort();
      } else {
        return [];
      }
    }
  },
  methods: {
    create() {
      this.state.creating = true;
      const account = jwt.decode(sessionStorage.token).account;
      this.axios
        .post(this.endpoint + "/accounts/agents", {
          account: account,
          first_name: this.toCapital(this.form.fname),
          last_name: this.toCapital(this.form.lname),
          telephone: this.form.phone,
          cell: this.form.select.cell,
          village: this.form.select.village,
          sector: this.form.select.sector
        })
        .then(res => {
          this.state.creating = false;
          this.$snotify.info("Agent successfully created...");
          const message = `Password For ${res.data.first_name} ${last_name} is: ${res.data.password}`;
          this.$bvModal
            .msgBoxOk(message, {
              title: "Agent Details",
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
          const error = navigator.onLine
            ? err.response.data.error || err.response.data
            : "Please connect to the internet";
          this.$snotify.error(error);
          console.log(err.response);
          this.state.creating = false;
        });
    },
    loadData() {
      this.state.tableLoad = true;
      const promise = this.axios.get(
        this.endpoint + "/accounts/agents?offset=0&limit=10"
      );
      return promise
        .then(res => {
          this.state.tableLoad = false;
          return res.data.Agents;
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
      this.$refs.agent_rightMenu.showMenu(evt, house);
    },
    optionClicked(data) {
      console.log(data);
      if (data.option.slug == "delete") {
        this.state.tableLoad = true;
        this.axios
          .delete(this.endpoint + "/accounts/agents/" + data.item.telephone)
          .then(res => {
            this.loadData();
            this.state.tableLoad = false;
            this.$snotify.info("Agent deleted Succesfully");
            console.log(res.data);
          })
          .catch(err => {
            this.state.tableLoad = false;
            const error = navigator.onLine
              ? err.response.data.error || err.response.data
              : "Please connect to the internet";
            this.$snotify.error(error);
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
        .put(this.endpoint + "/accounts/agents/creds/" + data.telephone, {
          password: this.form.newPwd
        })
        .then(res => {
          this.$snotify.info(res.data.message);
          this.state.changing = false;
          this.change_pswd_modal.show = false;
        })
        .catch(err => {
          const error = navigator.onLine
            ? err.response.data.error || err.response.data
            : "Please connect to the internet";
          this.$snotify.error(error);
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
.addAgent .nav-link.active {
  background-color: white !important;
}
.addAgent .accountForm {
  width: auto;
  border: 1px solid #dee2e6;
  border-radius: 5px;
  padding: 1rem;
}
</style>