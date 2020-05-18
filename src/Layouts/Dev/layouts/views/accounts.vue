<template>
  <div class="dev-accounts-wrapper">
    <div class="stats">
      <header class="secondary-font">Accounts in Numbers</header>
      <div class="cards-loader secondary-font" v-if="state.loading">
        <i class="fa fa-spinner fa-spin" />
        <p>Loading...</p>
      </div>
      <div class="cards" v-else>
        <div class="custom-card">
          <div class="card-content">
            <h3>{{Number(activeAccounts).toLocaleString()}}</h3>
            <h4>Active Accounts</h4>
          </div>
          <div class="icon">
            <i class="fa fa-user-check" />
          </div>
        </div>
        <div class="custom-card">
          <div class="card-content">
            <h3>{{Number(inactiveAccounts).toLocaleString()}}</h3>
            <h4>Inactive Accounts</h4>
          </div>
          <div class="icon">
            <i class="fa fa-user-times" />
          </div>
        </div>
      </div>
    </div>
    <div class="account-table">
      <header class="secondary-font custom-header">
        <h5>All Accounts</h5>
        <div class="add">
          <i class="fa fa-plus" />
        </div>
        <div class="refresh" @click="getData">
          <i class="fa fa-sync-alt" :class="{'fa-spin':state.loading}" />
        </div>
      </header>
      <div class="custom-table">
        <b-table
          hover
          responsive
          :items="items"
          :fields="fields"
          :busy="state.loading"
          head-variant="light"
          thead-class="table-header"
          tbody-tr-class="table-row"
          table-class="secondary-font"
          show-empty
        >
          <template v-slot:cell(updated_at)="data">
            <div class="d-flex align-items-center position-relative">
              <div class="edited-cell">{{data.value | dateFormatter}}</div>
              <i class="fa fa-ellipsis-v more-icon" />
            </div>
          </template>
          <template v-slot:cell(created_at)="data">
            <div class="edited-cell">{{data.value | dateFormatter}}</div>
          </template>
          <template v-slot:table-busy>
            <div class="table-loading">
              <i class="fa fa-spinner fa-spin" />
              <p>Loading...</p>
            </div>
          </template>
          <template v-slot:empty>
            <div class="table-empty">
              <p>No Accounts records available right now!</p>
            </div>
          </template>
        </b-table>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "developers-dashboard",
  data() {
    return {
      state: {
        loading: false
      },
      fields: [
        {
          key: "name",
          label: "account name",
          sortable: true,
          tdClass: "text-capitalize"
        },
        {
          key: "id",
          label: "Identifier",
          sortable: false
        },
        {
          key: "type",
          label: "account type",
          sortable: false,
          tdClass: "text-center"
        },
        {
          key: "created_at",
          label: "Creation Date",
          sortable: true
        },
        {
          key: "updated_at",
          label: "Last updated at",
          sortable: true
        }
      ],
      items: []
    };
  },
  beforeMount() {
    this.getData();
  },
  destroyed() {
    this.state.loading = false;
    delete this.items;
  },
  computed: {
    activeAccounts() {
      if (this.items.length < 1) return 0;
      return this.items.filter(item => item.active === true).length;
    },
    inactiveAccounts() {
      if (this.items.length < 1) return 0;
      return this.items.filter(item => item.active === false).length;
    }
  },
  filters: {
    dateFormatter: date => {
      if (!date) return "";
      return new Date(date).toLocaleDateString("en-EN", {
        year: "numeric",
        month: "long",
        day: "numeric"
      });
    }
  },
  methods: {
    async getData() {
      this.state.loading = true;
      const Total = await this.getTotals("/accounts?offset=0&limit=0");
      this.axios
        .get("/accounts?offset=0&limit=" + Total)
        .then(res => {
          this.items = res.data.Accounts;
          this.state.loading = false;
        })
        .catch(err => {
          console.log(err, err.response);
          this.state.loading = false;
        });
    },
    getTotals(endpoint) {
      return this.axios
        .get(endpoint)
        .then(res => res.data.Total)
        .catch(err => 0);
    }
  }
};
</script>

<style lang="scss">
.dev-accounts-wrapper {
  display: grid;
  grid-template-rows: auto auto;
  width: 100%;
  padding: 15px 0;
  header {
    font-size: 1.4rem;
    margin: 0.7rem 0;
    padding-right: 0.5rem;
    color: #3e4c52;
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
      max-width: 300px;
      min-width: 270px;
      height: fit-content;
      box-shadow: 0 1px 6px 0 rgba(32, 33, 36, 0.28);
      padding: 2rem;
      border-radius: 5px;
      display: flex;
      justify-content: center;
      align-items: center;
      flex-basis: calc(100% / 2 - 1rem);
      flex-wrap: nowrap;
      background: #0382b9;
      flex: 1;
      animation-name: fade;
      animation-duration: 500ms;
      animation-iteration-count: 1;

      &:not(:last-child) {
        margin-right: 2.5rem;
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
          white-space: nowrap;
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
  .cards-loader {
    padding: 4rem;
    display: flex;
    justify-content: center;
    align-items: center;
    user-select: none;
    animation-name: fade;
    animation-duration: 500ms;
    animation-iteration-count: 1;
    i {
      font-size: 2rem;
      margin-right: 0.5rem;
    }
    p {
      font-size: 1.3rem;
      margin: 0 !important;
    }
  }
  .account-table {
    margin-top: 2rem;

    .custom-table {
      box-shadow: 0 4px 5px 0 rgba(32, 33, 36, 0.09);
      border-radius: 5px;
      tbody tr {
        animation-name: fade;
        animation-duration: 500ms;
        animation-iteration-count: 1;
      }
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
      table.b-table[aria-busy="true"] {
        opacity: 0.8;
      }
      .table-loading {
        display: flex;
        justify-content: center;
        align-items: center;
        padding: 2.5rem;
        user-select: none;

        i {
          color: #010d19;
          font-size: 2rem;
          margin-right: 0.5rem;
        }
        p {
          font-size: 1.3rem;
          margin: 0 !important;
        }
      }
      .table-empty {
        padding: 3rem;
        text-align: center;
        p {
          font-size: 1.1rem;
          color: #33454c;
          text-transform: capitalize;
        }
      }
    }
  }
  @keyframes fade {
    0% {
      opacity: 0.1;
    }
    100% {
      opacity: 1;
    }
  }
}
</style>