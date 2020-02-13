<template>
  <div class="dev-dashboard d-flex">
    <div class="dev-sidebar app-color d-flex flex-column" :class="{'active':active}">
      <h1
        class="text-uppercase text-white w-100 text-center p-1 d-flex border-bottom border-secondary"
      >
        paypack
        <code class="font-13 text-white">{dev}</code>
      </h1>

      <ul class="dev-sidebar-links p-0 w-100 m-0 mt-5">
        <router-link
          to="/dev"
          class="text-white text-decoration-none"
          exact-active-class="active-color"
        >
          <li class="text-uppercase p-3 hover-color">statistics</li>
        </router-link>
        <router-link
          to="/dev/account"
          class="text-white text-decoration-none"
          exact-active-class="active-color"
        >
          <li class="text-uppercase p-3 hover-color">accounts</li>
        </router-link>
      </ul>
    </div>
    <div class="dev-content" :class="{'active':active}">
      <nav
        class="dev-nav navbar navbar-expand-lg navbar-light bg-light border-bottom d-flex justify-content-between"
      >
        <b-button size="sm" variant="info" @click="active =!active">
          <i class="fa fa-align-left" />
        </b-button>
        <b-button size="sm" variant="info" @click="logout">Logout</b-button>
      </nav>
      <b-container class="dev-body">
        <router-view />
      </b-container>
    </div>
  </div>
</template>

<script>
export default {
  name: "dev-dashboard",
  data() {
    return {
      active: false
    };
  },
  methods: {
    logout() {
      this.$store.dispatch("logout");
      this.$destroy;
    }
  }
};
</script>

<style>
.dev-dashboard {
  width: 100vw;
}
.dev-dashboard .dev-sidebar {
  min-width: 200px;
  max-width: 200px;
  transition: all 0.5s;
}
.dev-dashboard .dev-sidebar.active {
  margin-left: -200px;
}
.dev-dashboard .dev-content {
  width: 100%;
  min-height: 100vh;
  transition: all 0.5s;
}
.dev-dashboard .dev-content nav {
  height: 40px;
  max-width: calc(100vw - 200px);
  transition: all 0.5s;
}
.dev-dashboard .dev-content .dev-body {
  height: calc(100vh - 40px);
  max-width: calc(100vw - 200px);
  transition: all 0.5s;
}
.dev-dashboard .dev-content.active .dev-body {
  max-width: 100vw;
}
.dev-dashboard .dev-content.active nav {
  max-width: 100vw;
}
.dev-dashboard .dev-sidebar h1 {
  font-size: 25px;
  height: 40px;
  letter-spacing: 5px;
  white-space: nowrap;
  user-select: none;
}
.dev-sidebar .dev-sidebar-links {
  list-style: none;
}
</style>