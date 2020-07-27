<template>
  <div class="admin-wrapper d-flex">
    <div class="admin-sidebar">
      <h1
        class="text-white primary-font m-0 w-100 d-flex justify-content-center align-items-center"
      >P A Y P A C K</h1>
      <hr class="m-0" />
      <ul class="sidebar-links primary-font p-0 mt-5 scrollBar">
        <router-link
          tag="li"
          active-class="active-link"
          class="hover-color"
          v-if="isAdmin"
          to="/dashboard"
        >SECTOR</router-link>

        <router-link
          tag="li"
          active-class="active-link"
          class="hover-color"
          v-if="!isAdmin"
          to="/cells"
        >cells</router-link>

        <li
          v-if="isAdmin"
          v-b-toggle.changecells
          class="cursor-pointer hover-color"
          :class="{'active-link':highlightCellTag}"
        >Cells</li>
        <b-collapse
          id="changecells"
          class="sidebarCollapse"
          accordion="changecells"
          role="tabpanel"
          v-if="isAdmin"
        >
          <b-card tag="ul" no-body class="collapse-links">
            <router-link
              tag="li"
              class="hover-color"
              to="/cells"
              v-for="cell in cellsOptions"
              :key="cell"
              @click="update({toUpdate: 'cell', changed: cell})"
            >{{cell}}</router-link>
          </b-card>
        </b-collapse>

        <li
          v-b-toggle.changevillage
          class="cursor-pointer hover-color"
          :class="{'active-link':highlightVillageTag}"
        >Villages</li>
        <b-collapse
          id="changevillage"
          class="sidebarCollapse"
          accordion="changeVillage"
          role="tabpanel"
        >
          <b-card tag="ul" no-body class="collapse-links">
            <router-link
              tag="li"
              class="hover-color collapse-links"
              to="/village"
              v-for="village in villageOptions"
              :key="village"
              @click="update({toUpdate: 'village', changed: village})"
            >{{village}}</router-link>
          </b-card>
        </b-collapse>

        <router-link
          tag="li"
          active-class="active-link"
          class="hover-color"
          to="/transactions"
        >Bank Accounts</router-link>

        <router-link
          tag="li"
          active-class="active-link"
          class="hover-color"
          to="/properties"
        >Properties</router-link>

        <router-link tag="li" active-class="active-link" class="hover-color" to="/reports">REPORTS</router-link>

        <router-link
          tag="li"
          active-class="active-link"
          class="hover-color"
          to="/feedbacks"
        >Feedbacks</router-link>

        <router-link
          tag="li"
          active-class="active-link"
          class="hover-color"
          to="/create"
          v-if="isAdmin"
        >Accounts</router-link>

        <router-link tag="li" active-class="active-link" class="hover-color" to="/message">Messages</router-link>
      </ul>
      <p class="powered secondary-font">
        Powered By
        <strong>Quarks Group.</strong>
      </p>
    </div>
    <div class="admin-content">
      <nav
        class="navbar navbar-expand-lg navbar-light bg-light border-bottom d-flex justify-content-end flex-nowrap"
      >
        <b-button class="ml-2 primary-font br-2" variant="info" @click.prevent="logout">
          <i class="fa fa-sign-out-alt" />
          Logout
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
    },
    isAdmin() {
      if (this.user.role === "admin") return true;
      return false;
    },
    highlightCellTag() {
      return this.$route.path == "/cells";
    },
    highlightVillageTag() {
      return this.$route.path == "/village";
    },
  },
  methods: {
    update(res) {
      this.$store.dispatch("updatePlace", res);
    },
    logout() {
      this.$store.dispatch("logout");
    },
  },
  mounted() {
    console.log(this.$route);
  },
};
</script>
<style lang="scss" scoped>
@import "../assets/css/dashboardLayout.scss";
</style>
