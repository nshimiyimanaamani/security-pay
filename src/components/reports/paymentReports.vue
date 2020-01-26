<template>
  <div class="px-4">
    <header class="d-flex justify-content-center font-20 text-uppercase">Payment Reports</header>
    <hr class="m-0 mb-3" />
    <b-row class="px-3">
      <b-button variant="info" class="font-15 border-0 my-2 py-2 mr-3">Generate All House Report</b-button>
    </b-row>
    <b-row class="px-3">
      <b-button variant="info" class="font-15 border-0 my-2 py-2 mr-3">Generate Payed House Report</b-button>
    </b-row>
    <b-row class="px-3">
      <b-button variant="info" class="font-15 border-0 my-2 mr-3 py-2">Generate Unpayed House Report</b-button>
    </b-row>
  </div>
</template>

<script>
export default {
  name: "paymentReports",
  data() {
    return {};
  }
};
</script>

<style>
</style>