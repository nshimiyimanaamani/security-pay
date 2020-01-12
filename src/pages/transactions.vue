<template>
  <div style="position: relative">
    <vue-title title="Paypack | Transactions" />
    <div class="totals">
      <b-row>
        <b-col cols="6" class="text-white ml-auto py-2" style="font-size: 40px">{{'RWF '+total()}}</b-col>
      </b-row>
      <b-row class="text-white">
        <b-col cols="2" class="ml-auto">
          <p>BK Acc.</p>
          <p>RWF {{bkTotal()}}</p>
        </b-col>
        <b-col cols="2" class="m-0">
          <p>MTN MoMo</p>
          <p>RWF {{mtnTotal()}}</p>
        </b-col>
        <b-col cols="2" class="m-0">
          <p>AIRTEL MONEY</p>
          <p>RWF {{airtelTotal()}}</p>
        </b-col>
      </b-row>
    </div>
    <div class="container">
      <b-table
        :items="table.items"
        :fields="table.fields"
        :busy="loading"
        show-empty
        bordered
        small
      >
        <template v-slot:cell(method)="data">
          <div :class="data.value">
            <span>{{data.value}}</span>
          </div>
        </template>
        <template v-slot:cell(recorded)="data">
          <article class="text-center">{{data.value.slice(0,10)}}</article>
        </template>
        <template v-slot:cell(amount)="data">
          <article class="text-center">{{Number(data.value).toLocaleString()}} Frw</article>
        </template>
        <template v-slot:cell(date_recorded)="data">
          <article class="text-center">{{datify(data.value)}}</article>
        </template>
        <template v-slot:table-busy>
          <div class="text-center my-2">
            <b-spinner class="align-middle"></b-spinner>
            <strong>Loading...</strong>
          </div>
        </template>
        <template v-slot:empty="scope">
          <label
            class="container"
            style="width: 100%;font-size: 17px;text-align: center;padding: 40px;"
          >No records of transactions found!</label>
        </template>
      </b-table>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      loading: false,
      color: "#3db3fa",
      size: "12px",
      transactionData: [],
      table: {
        fields: [
          {
            key: "madeby",
            label: "Payee",
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
            thClass: "text-center"
          },
          {
            key: "amount",
            label: "Amount",
            sortable: false,
            thClass: "text-center"
          },
          {
            key: "date_recorded",
            label: "Date",
            sortable: true,
            thClass: "text-center"
          }
        ],
        items: []
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    }
  },
  mounted() {
    this.requestItems();
  },
  methods: {
    requestItems() {
      this.loading = true;
      this.axios
        .get(this.endpoint + "/transactions?offset=0&limit=5")
        .then(res => {
          this.table.items = res.data.Transactions;
          console.log(this.table.items);
          this.loading = false;
        })
        .catch(err => {
          this.table.items = [];
          if (navigator.onLine) {
            const error = isNullOrUndefined(err.response)
              ? "an error occured"
              : err.response.data.error || err.response.data;
            this.$snotify.error(error);
          } else {
            this.$snotify.error("Please connect to the internet");
          }
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
      const filtered = this.table.items.filter(data => data.method == "mtn");
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
