<template>
  <div class="position-relative transaction-page">
    <vue-title title="Paypack | Transactions" />
    <div class="totals primary-font">
      <b-row>
        <b-col
          cols="6"
          class="text-white ml-auto py-2 text-overflow"
          style="font-size: 40px"
        >RWF &nbsp;{{total() | number}}</b-col>
      </b-row>
      <b-row class="text-white">
        <b-col cols="2" class="ml-auto">
          <p class="text-overflow">BK ACC.</p>
          <p>RWF {{bkTotal() | number}}</p>
        </b-col>
        <b-col cols="2" class="m-0">
          <p class="text-overflow">MTN MoMo</p>
          <p>RWF {{mtnTotal() | number}}</p>
        </b-col>
        <b-col cols="2" class="m-0">
          <p class="text-overflow">AIRTEL MONEY</p>
          <p>RWF {{airtelTotal() | number}}</p>
        </b-col>
      </b-row>
    </div>
    <div class="transaction-table max-width bg-light">
      <header class="primary-font mb-3 table-header">
        <h3>All Transactions</h3>
        <b-button variant="info" @click="requestItems" :disabled="loading" class="br-2">
          Refresh
          <i class="fa fa-sync-alt" />
        </b-button>
      </header>
      <b-table
        :items="table.items"
        :fields="table.fields"
        :busy="loading"
        head-variant="light"
        table-class="bg-white"
        thead-class="primary-font"
        tbody-class="secondary-font"
        show-empty
        responsive
        bordered
        hover
        striped
      >
        <template v-slot:cell(method)="data">
          <div :class="data.value=='momo-mtn-rw'? 'mtn' : data.value">
            <span>{{data.value=='momo-mtn-rw'? 'mtn' : data.value}}</span>
          </div>
        </template>
        <template
          v-slot:cell(owner_firstname)="data"
        >{{data.item.owner_firstname +" "+data.item.owner_lastname}}</template>
        <template v-slot:cell(amount)="data">
          <article>{{data.value | number}} Frw</article>
        </template>
        <template v-slot:cell(date_recorded)="data">
          <article>{{data.value | date}}</article>
        </template>
        <template v-slot:table-busy>
          <div class="table-loader">
            <i class="fa fa-spinner fa-spin" />
            <p>Loading...</p>
          </div>
        </template>
        <template v-slot:empty>
          <label class="table-empty">No records of transactions found!</label>
        </template>
      </b-table>
    </div>
  </div>
</template>

<script>
import loader from "../components/loader";
export default {
  components: {
    loader,
  },
  data() {
    return {
      loading: false,
      transactionData: [],
      table: {
        fields: [
          {
            key: "owner_firstname",
            label: "Names",
            sortable: true,
            tdClass: "table-name",
          },
          {
            key: "madefor",
            label: "Paid for",
            sortable: true,
          },
          {
            key: "method",
            label: "Paid With",
            sortable: true,
            thClass: "text-center",
            tdClass: "text-center",
          },
          {
            key: "amount",
            label: "Amount",
            sortable: true,
            thClass: "text-right",
            tdClass: "text-right",
          },
          {
            key: "date_recorded",
            label: "Date & Time",
            sortable: true,
            thClass: "text-right",
            tdClass: "text-right",
          },
        ],
        items: [],
      },
    };
  },
  computed: {
    activeCell() {
      return this.$store.getters.getActiveCell;
    },
    user() {
      return this.$store.getters.userDetails;
    },
  },
  filters: {
    number: (num) => {
      return Number(num).toLocaleString();
    },
    date: (date) => {
      return new Intl.DateTimeFormat("en-US", {
        hour12: false,
        year: "numeric",
        month: "long",
        day: "2-digit",
        hour: "numeric",
        minute: "numeric",
      }).format(new Date(date));
    },
  },
  mounted() {
    this.requestItems();
  },
  methods: {
    async requestItems() {
      this.loading = true;
      const total = await this.$getTotal("/transactions?offset=0&limit=0");
      this.axios
        .get("/transactions?offset=0&limit=" + total)
        .then((res) => {
          const data = res.data.Transactions.sort((a, b) => {
            return new Date(b.date_recorded) - new Date(a.date_recorded);
          });
          if (this.user.role === "basic") {
            this.table.items = data.filter(
              (item) => item.cell == this.activeCell
            );
          } else {
            this.table.items = data;
          }
        })
        .catch((err) => {
          this.table.items = [];
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
        })
        .finally(() => {
          this.loading = false;
        });
    },
    total() {
      if (this.table.items.length < 1) return 0;
      try {
        return this.table.items.reduce(
          (a, b) => Number(a.amount || a) + Number(b.amount)
        );
      } catch {
        return 0;
      }
    },
    mtnTotal() {
      if (this.table.items.length < 1) return 0;
      try {
        return this.table.items
          .filter((data) => data.method.includes("mtn"))
          .reduce((a, b) => Number(a.amount || a) + Number(b.amount));
      } catch {
        return 0;
      }
    },
    airtelTotal() {
      if (this.table.items.length < 1) return 0;
      try {
        return this.table.items
          .filter((data) => data.method.includes("airtel"))
          .reduce((a, b) => Number(a.amount || a) + Number(b.amount));
      } catch {
        return 0;
      }
    },
    bkTotal() {
      if (this.table.items.length < 1) return 0;
      try {
        return this.table.items
          .filter((data) => data.method.includes("bk"))
          .reduce((a, b) => Number(a.amount || a) + Number(b.amount));
      } catch {
        return 0;
      }
    },
  },
};
</script>
<style lang="scss" scoped>
@import "../assets/css/transactionTable.scss";
</style>
