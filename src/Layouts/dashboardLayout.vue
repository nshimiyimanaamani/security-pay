<template>
  <div class="admin-wrapper d-flex">
    <div class="admin-sidebar" :class="{'active' : active}">
      <h1 class="text-white p-1 d-flex justify-content-center align-items-center">P A Y P A C K</h1>
      <hr class="m-0" />
      <ul class="sidebar-links font-14 p-0 mt-5">
        <router-link v-if="user.role.toLowerCase() !='basic'" to="/dashboard">
          <li>SECTOR</li>
        </router-link>
        <router-link v-if="user.role.toLowerCase() =='basic'" to="/cells">
          <li>cells</li>
        </router-link>
        <li
          v-if="user.role.toLowerCase() != 'basic'"
          v-b-toggle.changecells
          class="cursor-pointer"
        >Cells</li>
        <b-collapse
          id="changecells"
          accordion="changecells"
          role="tabpanel"
          v-if="user.role.toLowerCase() != 'basic'"
        >
          <b-card no-body class="border-0" style="background: transparent">
            <router-link to="/cells" v-for="cell in cellsOptions" :key="cell">
              <ul
                @click="update({toUpdate: 'cell', changed: cell})"
                class="text-white py-1 px-3 font-13 hover-color cursor-pointer"
              >{{cell}}</ul>
              <hr class="m-0" />
            </router-link>
          </b-card>
        </b-collapse>
        <li v-b-toggle.changevillage class="cursor-pointer">Village</li>
        <b-collapse id="changevillage" accordion="changecells" role="tabpanel">
          <b-card no-body class="border-0" style="background: transparent">
            <router-link to="/village" v-for="village in villageOptions" :key="village">
              <ul
                @click="update({toUpdate: 'village', changed: village})"
                class="text-white py-1 px-3 font-13 hover-color cursor-pointer"
              >{{village}}</ul>
              <hr class="m-0" />
            </router-link>
          </b-card>
        </b-collapse>
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
        <router-link
          to="/create"
          v-if="user.role.toLowerCase() =='dev'||user.role.toLowerCase() =='admin'"
        >
          <li>Accounts</li>
        </router-link>
      </ul>
      <p class="text-center powered m-0 pb-1 app-color text-white" for="powered">
        Powered By
        <strong>Quarks Group.</strong>
      </p>
    </div>
    <div class="admin-content">
      <nav
        class="navbar navbar-expand-lg navbar-light bg-light border-bottom d-flex justify-content-between"
      >
        <b-button size="sm" variant="info" @click="active=!active">
          <i class="fa fa-align-left"></i>
        </b-button>
        <b-button class="btn-info py-1 font-14" @click.prevent="logout">Logout</b-button>
      </nav>
      <div class="admin-body" :class="{'active':active}">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "dashboard-layout",
  data() {
    return {
      active: false
    };
  },
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
  mounted() {
    window.onresize = () => {
      if (window.innerWidth < 770) {
        this.active = true;
      } else if (window.innerWidth > 770) {
        this.active = false;
      }
    };
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
