<template>
  <b-container class="max-width" style="min-width:500px">
    <vue-title title="Paypack | Reports" />
    <b-card no-body>
      <b-tabs pills card vertical lazy class="font-13 text-uppercase">
        <b-tab title="PAYMENT REPORTS" active>
          <payment-reports />
        </b-tab>
        <b-tab
          v-if="user.role.toLowerCase() =='admin' || user.role.toLowerCase() =='dev'"
          title="SECTOR REPORTS"
        >
          <sector-reports />
        </b-tab>
        <b-tab title="CELL REPORTS">
          <cell-reports />
        </b-tab>
        <b-tab title="VILLAGE REPORTS">
          <village-reports />
        </b-tab>
        <b-tab title="HOUSE REPORTS">
          <house-reports />
        </b-tab>
        <b-tab title="DAILY REPORTS">
          <daily-reports />
        </b-tab>
      </b-tabs>
    </b-card>
  </b-container>
</template>

<script>
import paymentReports from "../components/reports/paymentReports";
import sectorReports from "../components/reports/sectorReports";
import cellReports from "../components/reports/cellReports";
import villageReports from "../components/reports/villageReports";
import houseReports from "../components/reports/houseReports";
import dailyReports from "../components/reports/dailyReports";
export default {
  name: "reports",
  components: {
    "payment-reports": paymentReports,
    "sector-reports": sectorReports,
    "cell-reports": cellReports,
    "village-reports": villageReports,
    "house-reports": houseReports,
    "daily-reports": dailyReports
  },
  data() {
    return {};
  },
  computed: {
    user() {
      return this.$store.getters.userDetails;
    }
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