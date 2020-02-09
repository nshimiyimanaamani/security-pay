<template>
  <div style="position: relative">
    <vue-title title="Paypack | Transactions" />
    <div class="totals">
      <b-row>
        <b-col
          cols="6"
          class="text-white ml-auto py-2 text-overflow"
          style="font-size: 40px"
        >{{'RWF '+total()}}</b-col>
      </b-row>
      <b-row class="text-white">
        <b-col cols="2" class="ml-auto">
          <p class="text-overflow">BK Acc.</p>
          <p>RWF {{bkTotal()}}</p>
        </b-col>
        <b-col cols="2" class="m-0">
          <p class="text-overflow">MTN MoMo</p>
          <p>RWF {{mtnTotal()}}</p>
        </b-col>
        <b-col cols="2" class="m-0">
          <p class="text-overflow">AIRTEL MONEY</p>
          <p>RWF {{airtelTotal()}}</p>
        </b-col>
      </b-row>
    </div>
    <div class="container max-width">
      <b-table
        :items="table.items"
        :fields="table.fields"
        :busy="loading"
        show-empty
        responsive
        bordered
        small
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
          <article>{{Number(data.value).toLocaleString()}} Frw</article>
        </template>
        <template v-slot:cell(date_recorded)="data">
          <article>{{datify(data.value)}}</article>
        </template>
        <template v-slot:table-busy>
          <loader class="d-flex justify-content-center p-5" />
        </template>
        <template v-slot:empty>
          <label
            class="container w-100 font-14 text-center p-5 text-capitalize"
          >No records of transactions found!</label>
        </template>
      </b-table>
    </div>
  </div>
</template>

<script>
import loader from "../components/loader";
export default {
  components: {
    loader
  },
  data() {
    return {
      loading: false,
      color: "#3db3fa",
      size: "12px",
      transactionData: [],
      table: {
        fields: [
          {
            key: "owner_firstname",
            label: "Names",
            sortable: true,
            tdClass: "table-name"
          },
          {
            key: "madefor",
            label: "Payed for",
            sortable: true
          },
          {
            key: "method",
            label: "Payed With",
            sortable: true,
            thClass: "text-center",
            tdClass: "text-center"
          },
          {
            key: "amount",
            label: "Amount",
            sortable: false,
            thClass: "text-center",
            tdClass: "text-right"
          },
          {
            key: "date_recorded",
            label: "Date",
            sortable: true,
            thClass: "text-center",
            tdClass: "text-right"
          }
        ],
        items: []
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    activeCell() {
      return this.$store.getters.getActiveCell;
    },
    user() {
      return this.$store.getters.userDetails;
    }
  },
  mounted() {
    this.requestItems();
  },
  methods: {
    requestItems() {
      this.loading = true;
      this.axios
        .get(this.endpoint + "/transactions?offset=0&limit=1000")
        .then(res => {
          if (this.user.role.toLowerCase() == "basic") {
            this.table.items = res.data.Transactions.filter(
              item => item.cell == this.activeCell
            );
          } else {
            this.table.items = res.data.Transactions;
          }
        })
        .catch(err => {
          this.table.items = [];
          if (navigator.onLine) {
            const error = err.response
              ? err.response.data.error || err.response.data
              : "an error occured";
            this.$snotify.error(error);
          } else {
            this.$snotify.error("Please connect to the internet");
          }
        })
        .finally(() => {
          this.loading = false;
        });
    },
    total() {
      let total = 0;
      this.table.items.forEach(element => {
        total += element.amount;
      });
      return total;
    },
    mtnTotal() {
      let total = 0;
      const filtered = this.table.items.filter(data =>
        data.method.includes("mtn")
      );
      filtered.forEach(element => {
        total += element.amount;
      });
      return total;
    },
    airtelTotal() {
      let total = 0;
      const filtered = this.table.items.filter(data => data.method == "airtel");
      filtered.forEach(element => {
        total += element.amount;
      });
      return total;
    },
    bkTotal() {
      let total = 0;
      const filtered = this.table.items.filter(data => data.method == "bk");
      filtered.forEach(element => {
        total += element.amount;
      });
      return total;
    },
    datify(date) {
      const options = { year: "numeric", month: "long", day: "numeric" };
      return new Date(date).toLocaleDateString("en-EN", options);
    }
  }
};
</script>

<style scoped>
@import url("../assets/css/transactionTable.css");
</style>
