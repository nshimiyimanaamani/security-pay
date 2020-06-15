<template>
  <b-container class="accounts-page h-100" fluid>
    <vue-title title="Paypack | create accounts" />
    <header class="secondary-font">ACCOUNT REGISTRATION</header>
    <b-card no-body class="nav-controls">
      <b-tabs
        pills
        card
        vertical
        lazy
        class="text-uppercase secondary-font flex-nowrap h-100"
        content-class="text-capitalize"
        nav-wrapper-class="accounts-nav"
        active-nav-item-class="app-color text-white"
        active-tab-class="h-100 p-0"
      >
        <b-tab title="Agents Accounts" active>
          <add-agent />
        </b-tab>
        <b-tab title="Admins Accounts">
          <add-admin />
        </b-tab>
        <b-tab title="Managers Accounts">
          <add-manager />
        </b-tab>
      </b-tabs>
    </b-card>
  </b-container>
</template>

<script>
import addAgent from "../components/accounts/addAgent.vue";
import addManager from "../components/accounts/addManager.vue";
import addAdmin from "../components/accounts/addAdmin.vue";
export default {
  name: "createAccount",
  components: {
    "add-agent": addAgent,
    "add-manager": addManager,
    "add-admin": addAdmin
  },
  data() {
    return {
      form: {
        account: null,
        select: {
          sector: null,
          cell: null,
          village: null
        }
      }
    };
  },
  computed: {
    cellOptions() {
      const sector = this.form.select.sector;
      if (sector)
        return this.$cells(
          this.location.province,
          this.location.district,
          sector
        ).sort();
    },
    villageOptions() {
      const sector = this.form.select.sector;
      const cell = this.form.select.cell;
      if (cell) {
        return this.$villages(
          this.location.province,
          this.location.district,
          sector,
          cell
        ).sort();
      } else {
        return [];
      }
    },
    user() {
      return this.$store.getters.userDetails;
    },
    location() {
      return this.$store.getters.location;
    }
  }
};
</script>

<style lang="scss">
.accounts-page {
  padding-top: 1rem;
  padding-bottom: 2rem;
  header {
    font-size: 1.25rem;
    width: 100%;
    text-align: center;
    height: 50px;
    font-weight: 700;
  }
  .nav-controls {
   height: calc(100% - 50px);
  }
  .accounts-nav {
    .nav-item {
      margin: 5px 0;

      .nav-link {
        border-radius: 2px;
        color: #017db3;
        font-size: 14px;
        &:hover {
          color: white !important;
          background-color: #017db3 !important;
        }
      }
    }
  }
  .tab-content {
    .nav {
      background: #f8f9fa;
      .nav-item {
        margin: 0;

        .nav-link {
          border-radius: 2px;
          color: #017db3;
          font-size: 15px;
          &.active {
            background-color: #017db3 !important;
            color: white;
          }
          &:hover {
            color: white !important;
            background-color: #017db3 !important;
          }
        }
      }
    }
  }
}
</style>