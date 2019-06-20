<template>
  <div class="transactions-container">
    <div class="transaction-title">
      <div class="year-selector">
        YTD
        <i class="fa fa-caret-down"/>
      </div>
      <br>
      <div class="paragraph">RWF 9,974.55</div>
      <br>
      <div class="bottom-text">
        <div class="contributer1">
          <span class="separator">|</span>
          <span class="text">
            BK Acc
            <br>
            <span class="small">500 M</span>
          </span>
        </div>
        <div class="contributer2">
          <span class="separator">|</span>
          <span class="text">
            MTN MoMo
            <br>
            <span class="small">60 M</span>
          </span>
        </div>
        <div class="contributer3">
          <span class="separator">|</span>
          <span class="text">
            Airtel Money
            <br>
            <span class="small">79 M</span>
          </span>
        </div>
      </div>
      <!-- <div class="bottom-text">
        <span>
          BK Acc
          <br>
          <span class="small">12M</span>
        </span>
        <span class="separator">|</span>

        <span>
          MTN MoMo
          <br>
          <span class="small">12M</span>
        </span>
        <span class="separator">|</span>
        <span>
          Airtel Money
          <br>
          <span class="small">12M</span>
        </span>
      </div>-->
    </div>
    <div class="transaction-body">
      <transactionTable/>
    </div>
  </div>
</template>

<script>
import transactionTable from "../components/transactionTable.vue";
export default {
  components: { transactionTable }
};
</script>

<style>
@import url("../assets/css/transactions.css");
</style>
