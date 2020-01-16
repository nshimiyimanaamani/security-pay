<template>
  <b-container>
    <b-card no-body>
      <b-tabs pills card vertical lazy>
        <b-tab  title="Payment Reports" active>
          <payment-reports />
        </b-tab>
        <b-tab title="Daily Reports">
          <b-card-text>Daily Reports under construction</b-card-text>
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
.nav-pills .nav-link.active {
  background-color: #3a82a1 !important;
}
a {
  color: #3a82a1;
}
</style>