<template>
  <div id="app">
    <vue-snotify class="text-capitalize font-13"></vue-snotify>
    <b-alert
      :show="offline"
      variant="danger"
      class="text-center offline-indicator font-13 w-auto mx-3"
      style="z-index: 10001"
      dismissible
    >
      <b>OFFLINE!</b> Please check your internet connection...
    </b-alert>
    <router-view />
  </div>
</template>
<script>
export default {
  data() {
    return {
      offline: false
    };
  },
  beforeMount() {
    delete sessionStorage.token;
    this.$store.dispatch("startup_function");
  },
  mounted() {
    window.addEventListener("offline", e => (this.offline = true));
    window.addEventListener("online", e => (this.offline = false));
  }
};
</script>

<style>
@import url("./assets/css/main.css");
</style>
