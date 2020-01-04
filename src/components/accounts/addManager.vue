<template>
  <b-tabs content-class="mt-2" class="addManager" fill lazy>
    <b-tab title="Create Manager" active>
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
        <b-form-group id="input-group-2" label="Cell:" label-for="input-2">
          <b-form-select
            id="input-2"
            v-model="form.cell"
            :options="cellOptions"
            style="font-size: 15px"
          >
            <template v-slot:first>
              <option :value="null" disabled>Please select cell</option>
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
    <b-tab title="List all Managers">
      <b-table
        id="Dev-table"
        bordered
        striped
        hover
        small
        :items="loadData"
        :fields="table.fields"
        :busy.sync="state.tableLoad"
      >
        <template v-slot:table-busy>
          <div class="text-center my-2">
            <b-spinner class="align-middle"></b-spinner>
            <strong>Loading...</strong>
          </div>
        </template>
      </b-table>
    </b-tab>
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
        cell: null
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
        tableLoad: false
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    cellOptions() {
      return this.$store.getters.getCellsArray;
    }
  },
  methods: {
    create() {
      this.state.creating = true;
      const account = jwt.decode(sessionStorage.token).account;
      this.axios
        .post(this.endpoint + "/api/accounts/managers", {
          account: account,
          email: this.form.email,
          cell: this.form.cell
        })
        .then(res => {
          this.state.creating = false;
          this.form = { email: null, password: null };
          this.$snotify.info("Developer successfully created...");
        })
        .catch(err => {
          this.state.creating = false;
          const error = navigator.onLine
            ? err.response.data.error
            : "Please connect to the internet...";
          this.$snotify.error(error);
        });
    },
    loadData() {
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
    toCapital(string) {
      string.toLowerCase();
      return string.charAt(0).toUpperCase() + string.slice(1);
    }
  }
};
</script>

<style>
.addManager .nav-link.active {
  background-color: white !important;
}
.addManager .accountForm {
  width: auto;
  border: 1px solid #dee2e6;
  border-radius: 5px;
  padding: 1rem;
  margin: 1rem;
}
</style>