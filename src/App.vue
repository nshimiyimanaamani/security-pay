<template>
  <div id="app">
    <vue-snotify></vue-snotify>
    <router-view />
  </div>
</template>
<script>
export default {
  beforeMount() {
    this.$store.dispatch("startup_function");
    this.authenticate();
    if (this.token) {
      this.axios.defaults.headers.common["Authorization"] = this.token;
    }
  },
  destroyed(){
    console.log('destroyed')
  },
  computed: {
    token() {
      return this.$store.getters.token;
    }
  },
  beforeUpdate() {
    this.authenticate();
  },
  methods: {
    authenticate() {
      if (this.token) {
        if (this.$route.name == "login") {
          this.$router.push("/dashboard");
          this.axios.defaults.headers.common["Authorization"] = this.token;
        }
      } else if (this.token == null) {
        this.$router.push("/");
        delete this.axios.defaults.headers.common["Authorization"];
      }
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
</style>
