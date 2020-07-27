<template>
  <b-tabs content-class="agentTabs-body" nav-class="agentTabs-nav" class="addAgent h-100" fill lazy>
    <b-tab title="Create Agent">
      <b-form class="accountForm" @submit.prevent="create">
        <b-row>
          <b-col>
            <b-form-group label="First name:">
              <b-form-input required v-model="form.fname" placeholder="First name..." class="br-2" />
            </b-form-group>
          </b-col>
          <b-col>
            <b-form-group label="Last name:">
              <b-form-input required v-model="form.lname" placeholder="Last name..." class="br-2" />
            </b-form-group>
          </b-col>
        </b-row>
        <b-form-group label="Phone Number:">
          <b-form-input
            type="number"
            v-model="form.phone"
            class="br-2"
            required
            placeholder="Enter Phone number..."
          />
        </b-form-group>

        <b-form-group label="Sector:">
          <b-form-select
            class="br-2"
            v-model="form.select.sector"
            :options="sectorOptions"
            required
          >
            <template v-slot:first>
              <option :value="null" disabled>select sector</option>
            </template>
          </b-form-select>
        </b-form-group>

        <b-form-group label="Cell:">
          <b-form-select class="br-2" v-model="form.select.cell" :options="cellOptions" required>
            <template v-slot:first>
              <option :value="null" disabled>select cell</option>
            </template>
          </b-form-select>
        </b-form-group>

        <b-form-group label="Village:">
          <b-select v-model="form.select.village" :options="villageOptions" class="br-2" required>
            <template v-slot:first>
              <option :value="null" disabled>select village</option>
            </template>
          </b-select>
        </b-form-group>
        <b-row class="justify-content-end" no-gutters>
          <b-button variant="info" type="submit" class="br-2">
            {{state.creating ? 'Creating' : "Create"}}
            <i
              class="fa fa-spinner fa-spin"
              v-if="state.creating"
            />
          </b-button>
        </b-row>
      </b-form>
    </b-tab>
    <b-tab title="List all Agents" class="h-100" active>
      <b-table
        id="data-table"
        striped
        hover
        small
        sticky-header
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
      class="primary-font"
    />
    <b-modal
      v-model="change_pswd_modal.show"
      title="Change Password"
      hide-footer
      content-class="primary-font"
      header-class="align-items-center"
    >
      <b-form @submit.prevent="changePassword">
        <b-form-group label="Current Password:">
          <b-form-input disabled v-model="currentPwd" />
        </b-form-group>
        <b-form-group label="New Password:">
          <b-form-input required v-model="form.newPwd" placeholder="New Password..."></b-form-input>
        </b-form-group>
        <b-row no-gutters class="justify-content-end">
          <b-button :disabled="form.newPwd ? false:true " variant="info" class="br-2" type="submit">
            {{state.changing ? 'Changing' : "Change"}}
            <i
              v-show="state.changing"
              class="fa fa-spinner fa-spin"
            />
          </b-button>
        </b-row>
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
        options: [
          { name: "Change Password", slug: "changePwd" },
          { name: "Delete", slug: "delete" }
        ]
      }
    };
  },
  computed: {
    sectorOptions() {
      return [this.activeSector];
    },
    cellOptions() {
      const sector = this.form.select.sector;
      const { province, district } = this.location;
      if (sector) return this.$cells(province, district, sector);
      return [];
    },
    villageOptions() {
      const sector = this.form.select.sector;
      const cell = this.form.select.cell;
      const { province, district } = this.location;
      if (sector && cell)
        return this.$villages(province, district, sector, cell);
      return [];
    },
    currentPwd() {
      if (this.change_pswd_modal.data)
        return this.change_pswd_modal.data.password;
      return null;
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
          this.state.creating = false;
          if (res.status === 201) {
            const el = this.$createElement;
            const messageNode = el("div", [
              el("p", { class: "text-center" }, [
                `Password For ${res.data.first_name} ${res.data.last_name} is`
              ]),
              el("br"),
              el("p", { class: "text-center" }, [` ${res.data.password}`])
            ]);
            this.$bvModal
              .msgBoxConfirm(messageNode, {
                okTitle: "YES",
                centered: true,
                cancelTitle: "NO",
                okVariant: "danger",
                modalClass: "__custom-modal",
                footerClass: "d-none",
                headerClass: "d-none",
                bodyClass: "__custom-modal-body",
                contentClass: "primary-font"
              })
              .then(() => {
                this.$snotify.info("Agent successfully created...");
              });
          }
        })
        .catch(err => {
          console.log(err, err.request, err.response);
          this.state.creating = false;
          try {
            this.$snotify.error(
              err.response.data.error || err.response.data || err
            );
          } catch {
            this.$snotify.error("Failed to create agent account!");
          }
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
          console.log(res);
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
            console.log(err, err.request, err.response);
            try {
              this.$snotify.error(
                err.response.data.error || err.response.data || err
              );
            } catch {
              this.$snotify.error("Failed to delete agent account!");
            }
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

<style lang="scss">
.addAgent {
  .nav-link.active {
    background-color: white !important;
  }

  .agentTabs-nav {
    height: 40px;
  }
  .agentTabs-body {
    height: calc(100% - 40px);
    overflow: auto;

    .b-table-sticky-header {
      max-height: 100%;
      margin: 0;
    }

    table {
      min-width: max-content;
      margin: 0;
    }
  }
  .accountForm {
    width: auto;
    border: 1px solid #fff;
    border-radius: 3px;
    padding: 1rem;
  }
}
</style>