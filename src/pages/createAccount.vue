<template>
  <b-container class="max-width">
    <vue-title title="Paypack | create accounts" />
    <h3 class="d-flex justify-content-center mb-3">ACCOUNT REGISTRATION</h3>
    <b-card no-body class="nav-controls font-13">
      <b-tabs pills card vertical lazy class="text-uppercase" content-class="text-capitalize">
        <b-tab title="Add Agent" active>
          <add-agent />
        </b-tab>
        <b-tab v-if="user.role.toLowerCase() != 'admin'" title="Add Admin">
          <add-admin />
        </b-tab>
        <b-tab title="Add Managers">
          <add-manager />
        </b-tab>
        <b-tab v-if="user.role.toLowerCase() == 'dev'" title="Add Developers">
          <add-dev />
        </b-tab>
      </b-tabs>
    </b-card>
  </b-container>
</template>

<script>
import addAgent from "../components/accounts/addAgent.vue";
import addDev from "../components/accounts/addDev.vue";
import addManager from "../components/accounts/addManager.vue";
import addAdmin from "../components/accounts/addAdmin.vue";
const { Cell, Village } = require("rwanda");
export default {
  name: "createAccount",
  components: {
    "add-agent": addAgent,
    "add-dev": addDev,
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
      if (sector) {
        return this.$store.getters.getCellsArray;
      }
    },
    villageOptions() {
      const sector = this.form.select.sector;
      const cell = this.form.select.cell;
      if (sector && cell) {
        return Village("Kigali", "Gasabo", sector, cell).sort();
      } else {
        return [];
      }
    },
    user() {
      return this.$store.getters.userDetails;
    }
  }
};
</script>

<style>
.nav-controls .nav-link.active {
  background-color: #3a82a1;
}
.nav-controls a {
  color: #3a82a1;
}
</style>