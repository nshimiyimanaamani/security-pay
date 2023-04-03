<template>
  <div class="admin-wrapper d-flex">
    <div class="admin-sidebar" v-if="showContent">
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
          v-if="isAdmin"
          to="/cells"
        >cells</router-link>

        <router-link tag="li" active-class="active-link" class="hover-color" to="/village">Villages</router-link>

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
        class="navbar navbar-expand-lg navbar-light bg-light border-bottom d-flex justify-content-between flex-nowrap"
      >
       <b-button class="ml-2 primary-font br-2" variant="info" @click="showContent = !showContent">
          <i class="fa fa-bars"></i>
        </b-button>
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
  data() {
    return {
      showContent: true,
    }
  },
  computed: {
    activeSector() {
      return this.$store.getters.getActiveSector;
    },
    user() {
      return this.$store.getters.userDetails;
    },
    isAdmin() {
      return this.user.role.toLowerCase() == "admin";
    },
    isLargeScreen() {
      if(window.innerWidth > 668) {
        this.showContent = true;
      } else {
        this.showContent = false;
      }
      // Example threshold for large screens
    },

  },
  methods: {
    logout() {
      this.$store.dispatch("logout");
    },
  },
};
</script>
<style lang="scss" scoped>
@import "../assets/css/dashboardLayout.scss";
</style>
