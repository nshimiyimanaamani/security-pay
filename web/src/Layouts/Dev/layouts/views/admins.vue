<template>
  <div class="admins-wrapper">
    <div class="stats">
      <header class="primary-font">Administrators in Numbers</header>
      <vue-load class="bg-light primary-font" v-if="state.loading" />
      <div class="cards" v-else>
        <div class="custom-card">
          <div class="card-content">
            <h3>{{adminsTotal | number}}</h3>
            <h4>Administrators</h4>
          </div>
          <div class="icon">
            <i class="fa fa-laptop-code" />
          </div>
        </div>
        <div class="custom-card">
          <div class="card-content">
            <h3>{{accountsTotal | number}}</h3>
            <h4>Assigned Accounts</h4>
          </div>
          <div class="icon">
            <i class="fa fa-address-card" />
          </div>
        </div>
        <div class="custom-card">
          <div class="card-content">
            <h3>{{sectorsTotal | number}}</h3>
            <h4>Assigned Sectors</h4>
          </div>
          <div class="icon">
            <i class="fa fa-address-card" />
          </div>
        </div>
      </div>
    </div>
    <div class="account-table">
      <header class="primary-font custom-header">
        <h5>Administrators Accounts</h5>
        <div class="add" @click="state.show.createAccount_modal=true">
          <i class="fa fa-plus" />
        </div>
        <div class="refresh" @click="getData">
          <i class="fa fa-sync-alt" :class="{'fa-spin':state.loading}" />
        </div>
      </header>
      <div class="custom-table">
        <b-table
          hover
          show-empty
          responsive
          :items="items"
          :fields="fields"
          :busy.sync="state.loading"
          head-variant="light"
          thead-class="table-header primary-font"
          tbody-tr-class="table-row"
          tbody-class="secondary-class"
        >
          <template v-slot:cell(updated_at)="data">
            <div class="d-flex align-items-center position-relative">
              <div class="edited-cell">{{data.value | date}}</div>
              <i
                class="fa fa-ellipsis-v more-icon"
                @click.prevent.stop="showMenu($event,data.item)"
              />
            </div>
          </template>
          <template v-slot:cell(created_at)="data">{{data.value | date}}</template>
          <template v-slot:empty>
            <p class="custom-data">No accounts available to display at the moment!</p>
          </template>
          <template v-slot:table-busy>
            <vue-load />
          </template>
        </b-table>
        <vue-menu
          elementId="devManagers-left-menu"
          ref="devManagersLeftMenu"
          :options="menuOptions"
          @option-clicked="optionClicked"
        />
        <!-- <b-modal
          ref="dev-updateAccount-modal"
          hide-footer
          title="Update Account"
          content-class="primary-font"
          centered
          @hide="modalClosed"
        >
          <update-account :account="selectedAccount" v-if="selectedAccount" @updated="closeModal" />
        </b-modal>-->
        <create-account
          v-if="state.show.createAccount_modal"
          @close="closeModal"
          @created="getData"
        />
      </div>
    </div>
  </div>
</template>

<script>
import createAccount from "../../components/create-adminAccount";
export default {
  name: "devAdmins",
  components: {
    createAccount
  },
  data() {
    return {
      state: {
        loading: false,
        show: { createAccount_modal: false }
      },
      menuOptions: [{ slug: "update", name: "Update account" }],
      items: [],
      fields: [
        {
          key: "email",
          label: "Username",
          sortable: true
        },
        {
          key: "account",
          label: "assigned account",
          sortable: true
        },
        {
          key: "cell",
          label: "assigned cell",
          sortable: true
        },
        {
          key: "role",
          label: "account type",
          sortable: true,
          tdClass: "text-center text-capitalize"
        },
        {
          key: "created_at",
          label: "Creation Date",
          sortable: true
        },
        {
          key: "updated_at",
          label: "Last Updated at",
          sortable: true
        }
      ]
    };
  },
  computed: {
    adminsTotal() {
      return 0;
    },
    accountsTotal() {
      return 0;
    },
    sectorsTotal() {
      return 0;
    }
  },
  mounted() {
    this.getData();
  },
  methods: {
    async getData() {
      this.state.loading = true;
      const total = await this.$getTotal("/accounts/admins?offset=0&limit=0");

      this.axios
        .get("/accounts/admins?offset=0&limit=" + total)
        .then(res => {
          this.items = res.data.Administrators;
          this.state.loading = false;
          console.log(res);
        })
        .catch(err => {
          console.log(err, err.response, err.request);
          try {
            this.$snotify.error(err.response.data.error);
          } catch {
            this.$snotify.error("Error! can't retrieve accounts");
          }
          this.state.loading = false;
        });
    },
    showMenu(event, data) {
      this.$refs.devManagersLeftMenu.showMenu(event, data);
    },
    async optionClicked(data) {
      if (!data) return;
      if (data.option.slug == "update") {
        this.selectedAccount = data.item;
        await this.selectedAccount;
        this.$refs["dev-updateAccount-modal"].show();
        console.log(data.item);
      }
    },
    closeModal() {
      this.state.show.createAccount_modal = false;
    }
  }
};
</script>

<style lang="scss">
.admins-wrapper {
  display: grid;
  grid-template-rows: auto auto;
  width: 100%;
  padding: 15px 0;
  header {
    font-size: 1.4rem;
    margin: 0.7rem 0;
    color: #3e4c52;
    padding-right: 0.5rem;
    &.custom-header {
      display: flex;
      align-items: center;

      h5 {
        flex: 1;
        margin: 0;
        font-size: 1.4rem;
        color: #3e4c52;
      }
      .refresh,
      .add {
        background: #0382b9;
        height: 2rem;
        flex-basis: 2rem;
        border-radius: 50%;
        box-shadow: 0 2px 6px 0 rgba(32, 33, 36, 0.28);
        display: flex;
        justify-content: center;
        align-items: center;
        color: whitesmoke;
        margin-left: 1rem;
        cursor: pointer;

        i {
          font-size: 1.2rem;
        }
      }
    }
  }

  .cards {
    display: flex;
    flex-wrap: wrap;
    .custom-card {
      min-width: 270px;
      height: fit-content;
      max-width: 300px;
      background: white;
      box-shadow: 0 1px 6px 0 rgba(32, 33, 36, 0.28);
      padding: 2rem;
      border-radius: 5px;
      display: flex;
      justify-content: center;
      align-items: center;
      flex-basis: calc(100% / 3 - 1rem);
      flex-wrap: nowrap;
      background: #0382b9;

      &:not(:last-child) {
        margin-right: 1.5rem;
        margin-bottom: 1rem;
      }
      .card-content {
        flex: 1;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: flex-start;
        margin-right: 1rem;
        h3 {
          font-size: 2.5rem;
          color: white;
          margin: 0;
        }
        h4 {
          font-size: 1.2rem;
          margin: 0;
          color: white;
          font-weight: 100;
        }
      }
      .icon {
        border: 2px solid whitesmoke;
        border-radius: 50%;
        width: 5rem;
        height: 5rem;
        display: flex;
        justify-content: center;
        align-items: center;

        i {
          font-size: 2rem;
          color: whitesmoke;
        }
      }
    }
  }

  .account-table {
    margin-top: 2rem;

    .custom-table {
      box-shadow: 0 4px 5px 0 rgba(32, 33, 36, 0.09);
      border-radius: 5px;
      .table-header th {
        border-color: #e9ecef;
        font-weight: 600;
        padding-left: 1rem;
        color: #33454c;
        text-transform: capitalize;
      }
      .table-row {
        &:hover {
          background-color: ghostwhite;
          box-shadow: 0 3px 5px 0 rgba(32, 33, 36, 0.08);
        }
        td {
          padding-left: 1rem;
          .edited-cell {
            width: calc(100% - 30px);
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
          }
          .custom-data {
            text-align: center;
            padding: 4rem;
            font-size: 1.1rem;
            font-weight: bold;
            background: ghostwhite;
          }
        }
      }
      .more-icon {
        cursor: pointer;
        border-radius: 50%;
        width: 2.1rem;
        background: white;
        height: 2rem;
        position: absolute;
        right: 0;
        display: flex;
        justify-content: center;
        align-items: center;
        &:hover {
          box-shadow: 0 1px 3px 0 rgba(32, 33, 36, 0.22);
        }
      }
      table.table[aria-busy="true"] {
        opacity: 0.8;
      }
    }
  }
  @keyframes fade {
    5% {
      opacity: 0.1;
    }
    100% {
      opacity: 1;
    }
  }
  .vue-menu {
    font-family: "Montserrat", Arial, Helvetica, -apple-system,
      BlinkMacSystemFont, "Helvetica Neue", "Noto Sans", sans-serif !important;
    &--active {
      background-color: ghostwhite;
    }

    &__item {
      &:hover {
        background-color: #0382b9;
      }
    }
  }
}
</style>
</style>