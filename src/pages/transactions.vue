<template>
  <div class="container">
    <b-table bordered :items="table.items" :fields="table.fields" small>
      <template slot="method" slot-scope="data">
        <div :class="data.value">
          <span>{{data.value}}</span>
        </div>
      </template>
    </b-table>
    <pulse-loader class="pulse" :loading="loading" :color="color" :size="size"></pulse-loader>
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
          { key: "property", label: "property ID", sortable: false },
          { key: "method", label: "Method of payment", sortable: true },
          { key: "amount", label: "Amount", sortable: false }
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
