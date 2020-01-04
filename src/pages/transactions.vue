<template>
  <div style="position: relative">
    <vue-title title="Paypack | Transactions" />
    <div class="totals">
      <b-row>
        <b-col cols="6" class="text-white ml-auto py-2" style="font-size: 40px">RWF 9,986.55</b-col>
      </b-row>
      <b-row class="text-white">
        <b-col cols="2" class="ml-auto">
          <p>BK Acc.</p>
          <p>500 M</p>
        </b-col>
        <b-col cols="2" class="m-0">
          <p>MTN MoMo</p>
          <p>60 M</p>
        </b-col>
        <b-col cols="2" class="m-0">
          <p>AIRTEL MONEY</p>
          <p>79 M</p>
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
            key: "owner",
            label: "Name",
            sortable: true,
            thClass: "table-name"
          },
          {
            key: "method",
            label: "Payed With",
            sortable: true,
            thClass: "text-center table-icon"
          },
          {
            key: "amount",
            label: "Amount",
            sortable: false,
            thClass: "text-center table-amount"
          },
          {
            key: "recorded",
            label: "Date",
            sortable: true,
            thClass: "text-center table-date"
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
        .get(this.endpoint + "/transactions?offset=1&limit=100")
        .then(res => {
          if (res.status == 200) {
            this.table.items = res.data.transactions;
            this.loading = false;
          } else {
            this.$snotify.info(`Failed to load Transactions`);
          }
        })
        .catch(err => {
          this.loading = false;
          this.table.items = [];
          if (navigator.onLine && err.response.status == "404") {
            this.$snotify.info(`Failed to load Transactions`);
          } else if (navigator.onLine && err.response.status != "404") {
            this.$snotify.info(err.response.data.error);
          } else if (!navigator.onLine) {
            this.$snotify.info(`Please connect to the internet...`);
          }
        });
    }
  }
};
</script>

<style scoped>
@import url("../assets/css/transactionTable.css");
</style>
