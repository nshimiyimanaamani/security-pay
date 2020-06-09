<template>
  <div id="payment-reports">
    <header>Payment Reports</header>
    <hr class="m-0 mt-1 mb-4" />
    <b-row class="m-0 buttons">
      <b-button variant="info">Generate All House Report</b-button>
      <b-button variant="info">Generate Paid House Report</b-button>
      <b-button variant="info">Generate Unpaid House Report</b-button>
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

<style lang="scss">
#payment-reports {
  & > header {
    text-align: center;
    font-size: 1.3rem;
    font-weight: bold;
    color: #384950;
  }
  .buttons {
    display: grid;
    grid-gap: 1rem;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));

    & > button {
      padding: 0.7rem 1rem;
    }
  }
}
</style>