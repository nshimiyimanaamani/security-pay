<template>
  <div id="app">
    <vue-snotify class="text-capitalize"></vue-snotify>
    <b-alert
      :show="offline"
      variant="danger"
      class="text-center offline-indicator font-15"
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
    this.$store.dispatch("startup_function");
  },
  mounted() {
    window.addEventListener("offline", e => (this.offline = true));
    window.addEventListener("online", e => (this.offline = false));
  }
};
</script>

<style>
#app {
  height: 100vh;
  width: 100vw;
  overflow-x: auto;
}
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
.btn-info {
  background-color: #3a82a1 !important;
  border-color: #3a82a1 !important;
}
.active .page-link {
  background-color: #3a82a1 !important;
  border-color: #3a82a1 !important;
}
.font-15 {
  font-size: 15px !important;
}
.font-13 {
  font-size: 13px !important;
}
.font-20 {
  font-size: 20px !important;
}
.cursor-pointer {
  cursor: pointer;
}
.app-color {
  background: #3a82a1 !important;
  background-image: none;
}
.hover-color:hover {
  background-color: #49bdd6;
}
</style>
