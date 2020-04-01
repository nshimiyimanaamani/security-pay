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
    <div class="version text-truncate" v-if="showVersion">Version {{version}}</div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      offline: false,
      showVersion: false,
      version: ""
    };
  },
  beforeMount() {
    this.$store.dispatch("startup_function");
  },
  mounted() {
    this.showVersion = false;
    window.addEventListener("offline", e => (this.offline = true));
    window.addEventListener("online", e => (this.offline = false));
    this.axios.get("/version").then(res => {
      this.version = res.data.version;
      this.showVersion = true;
    });
  }
};
</script>

<style>
.version {
  position: absolute;
  bottom: 0;
  right: 0;
  margin-right: 1rem;
  opacity: 0.7;
  font-size: 13px;
  letter-spacing: 1px;
  z-index: 1001;
  user-select: none;
  cursor: text;
}
</style>
