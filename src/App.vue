<template>
  <div id="app">
    <vue-snotify></vue-snotify>
    <b-alert :show="offline" variant="danger" class="text-center offline-indicator" dismissible>
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
    this.$store.dispatch("startup_function");
  },
  mounted() {
    setInterval(() => {
      navigator.onLine ? (this.offline = false) : (this.offline = true);
    }, 10000);
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
p {
  margin-bottom: 5px !important;
  font-size: 14px;
}
.alert-dismissible .close {
  padding: 6px 10px;
}
</style>
