<template>
  <b-container>
    <b-card no-body>
      <b-tabs pills card vertical lazy>
        <b-tab title="Payment Reports" active>
          <payment-reports />
        </b-tab>
        <b-tab title="Tab 2">
          <b-card-text>Tab contents 2</b-card-text>
        </b-tab>
        <b-tab title="Tab 3">
          <b-card-text>Tab contents 3</b-card-text>
        </b-tab>
      </b-tabs>
    </b-card>
  </b-container>
</template>

<script>
import paymentReports from "../components/reports/paymentReports.vue";
export default {
  name: "reports",
  components: {
    "payment-reports": paymentReports
  },
  data() {
    return {};
  }
};
</script>

<style>
.nav-pills .active {
  background: #3a82a1;
}
a {
  color: #3a82a1;
}
</style>