<template>
  <div id="app">
    <vue-snotify class="text-capitalize secondary-font font-13"></vue-snotify>
    <b-alert :show="offline" variant="danger" class="offline-indicator secondary-font" dismissible>
      <b>OFFLINE!</b> Please check your internet connection...
    </b-alert>
    <router-view />
    <!-- <div class="version text-truncate" v-if="showVersion">Version {{version}}</div> -->
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
      this.version = res.data.version || "";
      this.showVersion = res.data.version ? true : false;
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
.offline-indicator {
  z-index: 10001;
  position: absolute;
  width: fit-content;
  margin: auto;
  text-align: center;
  box-shadow: 0 1px 6px 0 rgba(32, 33, 36, 0.28);
  background-color: #f8d7da54;
}
</style>
