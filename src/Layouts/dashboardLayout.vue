<template>
  <div class="admin-wrapper d-flex">
    <div class="admin-sidebar" :class="{'active' : active}">
      <h1
        class="text-white m-0 w-100 d-flex justify-content-center align-items-center"
      >P A Y P A C K</h1>
      <hr class="m-0" />
      <ul class="sidebar-links font-14 p-0 mt-5">
        <router-link v-if="user.role.toLowerCase() !='basic'" to="/dashboard">
          <li class="hover-color">SECTOR</li>
        </router-link>
        <router-link v-if="user.role.toLowerCase() =='basic'" to="/cells">
          <li class="cursor-pointer hover-color">cells</li>
        </router-link>
        <li
          v-if="user.role.toLowerCase() != 'basic'"
          v-b-toggle.changecells
          class="cursor-pointer hover-color"
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
        <li v-b-toggle.changevillage class="cursor-pointer hover-color">Village</li>
        <b-collapse id="changevillage" accordion="changecells" role="tabpanel">
          <b-card no-body class="border-0 bg-transparent">
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
          <li class="hover-color">Bank Accounts</li>
        </router-link>
        <router-link to="/properties">
          <li class="hover-color">Properties</li>
        </router-link>
        <router-link to="/reports">
          <li class="hover-color">REPORTS</li>
        </router-link>
        <router-link to="/feedbacks">
          <li class="hover-color">Feedbacks</li>
        </router-link>
        <router-link
          to="/create"
          v-if="user.role.toLowerCase() =='dev'||user.role.toLowerCase() =='admin'"
        >
          <li class="hover-color">Accounts</li>
        </router-link>
        <router-link to="/message">
          <li class="hover-color">Messages</li>
        </router-link>
      </ul>
      <p class="text-center powered m-0 pb-1 app-color text-white" for="powered">
        Powered By
        <strong>Quarks Group.</strong>
      </p>
    </div>
    <div class="admin-content">
      <nav
        class="navbar navbar-expand-lg navbar-light bg-light border-bottom d-flex justify-content-end flex-nowrap"
      >
        <b-button class="ml-2 secondary-font br-2" variant="info" @click.prevent="logout">
          Logout
          <i class="fa fa-sign-out-alt" />
        </b-button>
      </nav>
      <div class="admin-body">
        <transition name="fade" :duration="250" mode="out-in">
          <router-view />
        </transition>
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
<style lang="scss" scoped>
@import "../assets/css/dashboardLayout.scss";
</style>
