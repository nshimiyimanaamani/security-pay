<template>
  <div class="container">
    <div class="loading" v-if="this.loading">
      <div class="lds-roller">
        <div></div>
        <div></div>
        <div></div>
        <div></div>
        <div></div>
        <div></div>
        <div></div>
        <div></div>
      </div>
    </div>
    <div class="row" v-else-if="this.loading == false">
      <div class="col-md-12">
        <table class="table">
          <thead>
            <tr>
              <th>
                Today
                <i class="fa fa-caret-down"></i>
              </th>
              <th>
                Payment Method
                <i class="fa fa-caret-down"></i>
              </th>
              <th>
                Narrative
                <i class="fa fa-caret-down"></i>
              </th>
              <th>
                Amount
                <i class="fa fa-caret-down"></i>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(trans , index) in transactionData" :key="index">
              <td>
                <i class="fa fa-check-circle"></i>
                <p>12:35,22 March 2019</p>
              </td>
              <td>
                <div v-if="trans.method == 'mtn' || 'MTN'" class="mtn">
                  <span>mtn</span>
                </div>
                <div v-else-if="trans.method == 'airtel' || 'AIRTEL'" class="airtel">
                  <span>airtel</span>
                </div>
                <div v-else-if="trans.method == 'bk' || 'BK'" class="bk">
                  <span>bk</span>
                </div>&nbsp;
                <p class="customerName">Customer Name</p>
              </td>
              <td>Umutekano</td>
              <td>
                {{trans.amount}}
                <i class="fa fa-ellipsis-v" style="float:right;margin-right:10px"></i>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
export default {
  data() {
    return {
      transactionData: []
    };
  },
  computed: {
    loading() {
      return this.$store.state.loading;
    }
  },
  mounted() {
    this.$store.state.loading = true;
    axios
      .get(
        "https://paypack-backend-qahoqfdr3q-uc.a.run.app/api/transactions/?offset=0&limit=5"
      )
      .then(res => {
        this.transactionData = res.data.transactions;
        this.$store.state.loading = false;
      })
      .catch(err => {
        this.$store.state.loading = false;
        console.log(err);
      });
  }
};
</script>

<style>
@import url("../assets/css/transactionTable.css");
</style>
