<template>
  <b-container class="reports-page h-100" fluid>
    <vue-title title="Paypack | Reports" />
    <b-card no-body class="mh-100 h-100">
      <b-tabs
        pills
        card
        vertical
        lazy
        class="reports-tabs mh-100 h-100 secondary-font"
        content-class=" reports-content"
        active-tab-class="h-100 p-0"
        active-nav-item-class="app-color text-white"
        nav-class="report-navs"
      >
        <b-tab title="PAYMENT REPORTS">
          <payment-reports />
        </b-tab>
        <b-tab v-if="isAdmin" title="SECTOR REPORTS" active>
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
    },
    isAdmin() {
      return this.user.role === "admin" ? true : false;
    }
  }
};
</script>

<style lang="scss">
.reports-page {
  min-width: 500px;
  padding: 3rem;

  .reports-tabs {
    display: flex;
    flex-wrap: nowrap;
    overflow: auto;
    text-transform: uppercase;

    header.tabTitle {
      text-align: center;
      font-size: 1.3rem;
      font-weight: bold;
      color: #384950;
      padding: 1rem;
      z-index: 100;
      background: white;
      border-bottom: 1px solid #e5e5e5;
      position: sticky;
      top: 0;
    }
    .tabBody {
      margin: 0;
      width: 100%;
      padding: 1rem;

      .reports-card {
        & > h5 {
          padding: 0.5rem 0;
          font-size: 1rem;
          font-weight: bold;
          margin: 0;
          color: white;
          text-align: center;
        }
        table {
          min-width: max-content;
        }
      }
    }
  }
  .report-navs {
    a.nav-link {
      color: #017db3;
      margin: 1px 0;
      border-radius: 2px;
      transition-duration: 200ms;
      transition-timing-function: ease-in-out;
      &:hover {
        color: #fff;
        background-color: #0382b9;
      }
    }
  }

  .reports-content {
    overflow: auto;
    .active {
      transition: all 500ms;
    }
  }
}
</style>