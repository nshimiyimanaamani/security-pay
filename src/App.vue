<template>
  <div id="app">
    <vue-snotify class="text-capitalize secondary-font fsize-sm"></vue-snotify>
    <b-alert :show="offline" variant="danger" class="offline-indicator secondary-font" dismissible>
      <b>OFFLINE!</b> Please check your internet connection...
    </b-alert>
    <div class="app-loading secondary-font" v-if="appLoading">
      <i class="fa fa-spinner fa-spin" />
      <h1>Initializing</h1>
    </div>
    <router-view v-else />
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
  computed: {
    appLoading() {
      return this.$store.getters.appLoading;
    }
  },
  async beforeMount() {
    this.showVersion = false;
    await this.$store.dispatch("startup_function");
    window.addEventListener("offline", e => (this.offline = true));
    window.addEventListener("online", e => (this.offline = false));
    this.axios.get("/version").then(res => {
      this.showVersion = res.data && res.data.version ? true : false;
      this.version = this.showVersion ? res.data.version || "" : "";
    });
  }
};
</script>

<style lang="scss">
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
.app-loading {
  position: absolute;
  top: 0;
  bottom: 0;
  right: 0;
  left: 0;
  z-index: 1000;
  background: white;
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  user-select: none;

  i {
    font-size: 3rem;
    margin-bottom: 2rem;
  }
  h1 {
    text-transform: uppercase;
    font-size: 2rem;
    letter-spacing: 0.5rem;
  }
}
</style>
