<template>
  <b-tabs content-class="mt-2" class="addAgent" fill lazy>
    <b-tab title="Create Agent" active>
      <b-form class="accountForm" @submit.prevent="create">
        <b-row>
          <b-col>
            <b-form-group id="input-group-1" label="First name:" label-for="input-1">
              <b-form-input
                id="input-1"
                required
                v-model="form.fname"
                placeholder="First name..."
                class="br-2"
              ></b-form-input>
            </b-form-group>
          </b-col>
          <b-col>
            <b-form-group id="input-group-2" label="Last name:" label-for="input-2">
              <b-form-input
                id="input-2"
                required
                v-model="form.lname"
                placeholder="Last name..."
                class="br-2"
              ></b-form-input>
            </b-form-group>
          </b-col>
        </b-row>
        <b-form-group id="input-group-3" label="Phone Number:" label-for="input-3">
          <b-form-input
            id="input-3"
            type="number"
            v-model="form.phone"
            class="br-2"
            required
            placeholder="Enter Phone number..."
          ></b-form-input>
        </b-form-group>

        <b-form-group id="input-group-4" label="Sector:" label-for="input-4">
          <b-form-select
            class="br-2"
            id="input-4"
            v-model="form.select.sector"
            :options="sectorOptions"
          >
            <template v-slot:first>
              <option :value="null" disabled>select sector</option>
            </template>
          </b-form-select>
        </b-form-group>

        <b-form-group id="input-group-5" label="Cell:" label-for="input-5">
          <b-form-select
            id="input-5"
            class="br-2"
            v-model="form.select.cell"
            :options="cellOptions"
          >
            <template v-slot:first>
              <option :value="null" disabled>select cell</option>
            </template>
          </b-form-select>
        </b-form-group>

        <b-form-group id="input-group-5" label="Village:" label-for="input-6">
          <b-select
            id="input-6"
            v-model="form.select.village"
            :options="villageOptions"
            class="br-2"
          >
            <template v-slot:first>
              <option :value="null" disabled>select village</option>
            </template>
          </b-select>
        </b-form-group>
        <b-form-group class="mb-0">
          <b-button variant="info" class="float-right" type="submit">
            {{state.creating ? 'Creating' : "Create"}}
            <i
              class="fa fa-spinner fa-spin"
              v-if="state.creating"
            />
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
        responsive
        head-variant="light"
        thead-class="text-uppercase"
        show-empty
        :items="loadData"
        :fields="table.fields"
        :busy.sync="state.tableLoad"
        @row-contextmenu="menu"
      >
        <template v-slot:cell(first_name)="data">{{data.item.first_name+" "+data.item.last_name}}</template>
        <template v-slot:table-busy>
          <vue-load />
        </template>
      </b-table>
    </b-tab>
    <vue-menu
      :elementId="'agent-rightmenu'"
      :options="rightMenu.options"
      :ref="'agent_rightMenu'"
      @option-clicked="optionClicked"
    ></vue-menu>
    <b-modal
      v-model="change_pswd_modal.show"
      title="Change Password"
      hide-footer
      header-class="align-items-center"
    >
      <b-form @submit.prevent="changePassword">
        <b-form-group id="input-group-1" label="Current Password:" label-for="input-1">
          <b-form-input id="input-1" disabled v-model="currentPwd"></b-form-input>
        </b-form-group>
        <b-form-group id="input-group-2" label="New Password:" label-for="input-2">
          <b-form-input id="input-2" required v-model="form.newPwd" placeholder="New Password..."></b-form-input>
        </b-form-group>
        <b-form-group class="m-0">
          <b-button
            :disabled="form.newPwd ? false:true "
            variant="info"
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
          { key: "first_name", label: "Names" },
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
        options: [{ name: "Change Password", slug: "changePwd" }]
      }
    };
  },
  computed: {
    sectorOptions() {
      return [this.activeSector];
    },
    cellOptions() {
      const sector = this.form.select.sector;
      if (sector) {
        return this.$cells(
          this.location.province,
          this.location.district,
          sector
        );
      }
    },
    villageOptions() {
      const sector = this.form.select.sector;
      const cell = this.form.select.cell;
      if (sector && cell) {
        return this.$villages(
          this.location.province,
          this.location.district,
          sector,
          cell
        ).sort();
      } else {
        return [];
      }
    },
    currentPwd() {
      if (this.change_pswd_modal.data) {
        return this.change_pswd_modal.data.password;
      } else {
        return null;
      }
    },
    user() {
      return this.$store.getters.userDetails;
    },
    activeSector() {
      return this.$store.getters.getActiveSector;
    },
    location() {
      return this.$store.getters.location;
    }
  },
  watch: {
    "form.select.sector"() {
      handler: {
        this.form.select.cell = null;
        this.form.select.village = null;
      }
    },
    "form.select.cell"() {
      handler: {
        this.form.select.village = null;
      }
    }
  },
  methods: {
    create() {
      this.state.creating = true;
      const account = this.user.account;
      this.axios
        .post("/accounts/agents", {
          account: account,
          first_name: this.toCapital(this.form.fname),
          last_name: this.toCapital(this.form.lname),
          telephone: this.form.phone,
          cell: this.form.select.cell,
          village: this.form.select.village,
          sector: this.form.select.sector
        })
        .then(res => {
          if (res) {
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
          }
        })
        .catch(err => {
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
        })
        .finally(() => {
          this.state.creating = false;
        });
    },
    async loadData() {
      this.state.tableLoad = true;
      const total = await this.$getTotal("/accounts/agents?offset=0&limit=0");
      const promise = this.axios.get(
        "/accounts/agents?offset=0&limit=" + total
      );
      return promise
        .then(res => {
          console.log(res)
          return res.data.Agents;
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

      this.$refs.agent_rightMenu.showMenu(evt, house);
    },
    optionClicked(data) {
      if (data.option.slug == "delete") {
        this.state.tableLoad = true;
        this.axios
          .delete("/accounts/agents/" + data.item.telephone)
          .then(res => {
            this.loadData();
            this.$snotify.info("Agent deleted Succesfully");
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
        .put("/accounts/agents/creds/" + data.telephone, {
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