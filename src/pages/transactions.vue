<template>
  <div class="container">
    <b-table :items="table.items" :fields="table.fields" :busy="loading" show-empty bordered small>
      <template v-slot:cell(method)="data">
        <div :class="data.value">
          <span>{{data.value}}</span>
        </div>
      </template>
      <template v-slot:cell(recorded)="data">
        <article class="text-center">{{data.value.slice(0,10)}}</article>
      </template>
      <template v-slot:cell(amount)="data">
        <article class="text-center">{{data.value}} Frw</article>
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
      const promise = new Promise((resolve, reject) => {
        this.axios
          .get(this.endpoint + "/transactions/?offset=1&limit=100")
          .then(res => {
            if (res.status == 200) {
              resolve(res.data.transactions);
            }
          })
          .catch(err => {
            reject(err);
          });
      });
      promise
        .then(res => {
          this.table.items = res;
          console.log(res);
          this.loading = false;
        })
        .catch(err => {
          this.loading = false;
          this.table.items = [];
          console.log(err);
        });
    }
  }
};
</script>

<style scoped>
@import url("../assets/css/transactionTable.css");
</style>
