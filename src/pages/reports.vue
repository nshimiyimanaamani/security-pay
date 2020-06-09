<template>
  <b-container class="reports-page h-100" fluid>
    <vue-title title="Paypack | Reports" />
    <b-card no-body class="mh-100">
      <b-tabs
        pills
        card
        vertical
        lazy
        class="reports-tabs text-uppercase flex-nowrap mh-100 overflow-auto"
        content-class="secondary-font overflow-auto reports-content"
        active-nav-item-class="app-color text-white"
        nav-class="report-navs"
      >
        <b-tab title="PAYMENT REPORTS">
          <payment-reports />
        </b-tab>
        <b-tab
          v-if="user.role.toLowerCase() =='admin' || user.role.toLowerCase() =='dev'"
          title="SECTOR REPORTS"
          active
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
        <!-- <b-tab title="DAILY REPORTS">
          <daily-reports />
        </b-tab>-->
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

<style lang="scss">
.reports-page {
  min-width: 500px;
  padding: 2rem 1rem 1rem;

  .report-navs {
    a.nav-link {
      color: #017db3;
      margin: 1px 0;
      transition-duration: 200ms;
      transition-timing-function: ease-in-out;
      &:hover {
        color: #fff;
        background-color: #0382b9;
      }
    }
  }

  .reports-content {
    .active {
      transition: all 500ms;
    }
  }
}
</style>