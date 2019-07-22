<template>
  <div class="container">
    <div class="row">
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
import { asyncLoading } from "vuejs-loading-plugin";
export default {
  data() {
    return {
      transactionData: [],
    };
  },
  computed: {
    endpoint() {
      return this.$store.state.endPoint;
    }
  },
  mounted() {
    const getTransactions = new Promise((resolve, reject) => {
      axios
        .get(`${this.endpoint}/transactions/?offset=0&limit=5`)
        .then(res => {
          if(res.status == 200){
          this.transactionData = res.data.transactions;
          resolve(res);
          }
        })
        .catch(err => {
          console.log(err);
          reject(err);
        });
    });
    asyncLoading(getTransactions)
    .then(res=>{console.log(res)})
    .catch(err=>{console.log(err)})
  }
};
</script>

<style>
@import url("../assets/css/transactionTable.css");
</style>
