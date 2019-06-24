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
                Narrative3
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
                <div></div>
                <i v-if="trans.success == 'true'" class="fa fa-check-circle"></i>
                <i v-else-if="trans.success == 'false'" class="fa fa-times-circle"></i>
                <p>{{trans.date}}</p>
              </td>
              <td>
                <div v-if="trans.paymentMethod == 'mtn'" class="mtn">
                  <span>mtn</span>
                </div>
                <div v-else-if="trans.paymentMethod == 'airtel'" class="airtel">
                  <span>airtel</span>
                </div>
                <div v-else-if="trans.paymentMethod == 'bk'" class="bk">
                  <span>bk</span>
                </div>&nbsp;
                <p class="customerName">{{trans.payee}}</p>
              </td>
              <td>{{trans.narrator}}</td>
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
export default {
  data() {
    return {
      transactionData: [
        {
          success: "true",
          date: "12:35,22 March 2019",
          paymentMethod: "mtn",
          payee: "customer1",
          narrator: "Umutekano",
          amount: "4000 Rwf"
        },
        {
          success: "true",
          date: "12:35,22 March 2019",
          paymentMethod: "airtel",
          payee: "customer2",
          narrator: "Umutekano",
          amount: "3000 Rwf"
        },
        {
          success: "true",
          date: "12:35,22 March 2019",
          paymentMethod: "mtn",
          payee: "customer3",
          narrator: "Umutekano",
          amount: "500 Rwf"
        },
        {
          success: "true",
          date: "12:35,22 March 2019",
          paymentMethod: "mtn",
          payee: "customer4",
          narrator: "Umutekano",
          amount: "5000 Rwf"
        },
        {
          success: "true",
          date: "12:35,22 March 2019",
          paymentMethod: "bk",
          payee: "customer5",
          narrator: "Umutekano",
          amount: "2500 Rwf"
        },
        {
          success: "false",
          date: "12:35,22 March 2019",
          paymentMethod: "bk",
          payee: "customer6",
          narrator: "Umutekano",
          amount: "2000 Rwf"
        }
      ]
    };
  }
};
</script>

<style>
@import url("../assets/css/transactionTable.css");
</style>
