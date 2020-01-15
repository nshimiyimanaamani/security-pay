<template>
  <div class="dashboardWrapper">
    <div class="dashboardSidebar">
      <h1>P A Y P A C K</h1>
      <hr />
      <ul v-if="user.role.toLowerCase() !='basic'" class="sidebarLinks">
        <router-link to="/dashboard">
          <li>Overview</li>
        </router-link>
        <router-link to="/cells">
          <li>
            Cells
            <b-dropdown>
              <b-dropdown-item
                @click="update({toUpdate: 'cell', changed: cell})"
                v-for="cell in cellsOptions"
                :key="cell"
              >{{cell}}</b-dropdown-item>
            </b-dropdown>
          </li>
        </router-link>
        <router-link to="/village">
          <li>
            Village
            <b-dropdown>
              <b-dropdown-item
                @click="update({toUpdate: 'village', changed: village})"
                v-for="village in villageOptions"
                :key="village"
              >{{village}}</b-dropdown-item>
            </b-dropdown>
          </li>
        </router-link>
        <router-link to="/transactions">
          <li>Bank Accounts</li>
        </router-link>
        <router-link to="/properties">
          <li>Properties</li>
        </router-link>
        <router-link to="/reports">
          <li>REPORTS</li>
        </router-link>
        <router-link to="/feedbacks">
          <li>Feedbacks</li>
        </router-link>
        <router-link to="/create" v-if="user.role == 'dev'">
          <li>Create Accounts</li>
        </router-link>
      </ul>
      <p class="text-center powered" for="powered">
        Powered By
        <strong>Quarks Group.</strong>
      </p>
    </div>
    <div class="rightSide">
      <div class="top-nav">
        <div class="logout">
          <b-button class="btn-info py-1" @click.prevent="logout">Logout</b-button>
        </div>
      </div>
      <div class="dashboardBody">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "dashboard-layout",
  computed: {
    cellsOptions() {
      return this.$store.getters.getCellsArray;
    },
    villageOptions() {
      return this.$store.getters.getVillageArray;
    },
    activeSector() {
      return this.$store.getters.getActiveSector;
    },
    user() {
      return this.$store.getters.userDetails;
    }
  },
  methods: {
    update(res) {
      this.$store.dispatch("updatePlace", res);
    },
    logout() {
      this.$store.dispatch("logout");
    }
  }
};
</script>
<style scoped>
@import url("../assets/css/dashboardLayout.css");
</style>