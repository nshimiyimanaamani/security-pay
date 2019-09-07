<template>
  <div id="app">
    <vue-snotify></vue-snotify>
    <b-alert :show="!isOnline" variant="danger" class="text-center offline-indicator" dismissible>
      <b>OFFLINE!</b> Please check your internet connection...
    </b-alert>
    <router-view />
  </div>
</template>
<script>
export default {
  beforeMount() {
    this.$store.dispatch("startup_function");
    this.$store.dispatch("checkConnection");
    this.authenticate();
    if (this.token) {
      this.axios.defaults.headers.common["Authorization"] = this.token;
    }
  },
  computed: {
    token() {
      return this.$store.getters.token;
    },
    isOnline() {
      return this.$store.getters.isOnline;
    }
  },
  beforeUpdate() {
    this.authenticate();
  },
  methods: {
    authenticate() {
      if (this.token) {
        if (this.$route.name === "login" || this.$route.name === "register") {
          this.$router.push("/dashboard");
          this.axios.defaults.headers.common["Authorization"] = this.token;
        }
      } else if (!this.token) {
        if (this.$route.name === "register") {
          this.$route.push("/register");
        } else if (this.$route.name === "login") {
          this.$router.push("/");
        } else {
          this.$route.push("/");
        }
        delete this.axios.defaults.headers.common["Authorization"];
      }
    },
    checkConnection() {
      return navigator.onLine;
    }
  },
  watch: {
    token() {
      handler: {
        this.$forceUpdate();
      }
    }
  }
};
</script>

<style>
@import url("./assets/css/main.css");
.offline-indicator {
  position: absolute;
  z-index: 1000;
  left: 0;
  top: 5px;
  right: 0;
  width: 50%;
  padding: 6px 35px;
  margin: auto;
}
</style>
