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
  letter-spacing: 1px;
}
.font-13 {
  font-size: 13px !important;
  letter-spacing: 1px;
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
.max-width {
  max-width: 100vw;
}
.text-normal {
  text-transform: none;
}
.h-f-content {
  height: fit-content !important;
}
.w-f-content {
  width: fit-content !important;
}
.vue-simple-context-menu {
  position: fixed !important;
}
.vue-slider-rail {
  background-color: #ccc;
  border-radius: 15px;
}
.vue-slider-process {
  background: #3a82a1 !important;
  border-radius: 15px;
}
.vue-slider-dot-handle {
  cursor: pointer;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background-color: #fff;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
  -webkit-box-shadow: 0.5px 0.5px 2px 1px rgba(0, 0, 0, 0.32);
  box-shadow: 0.5px 0.5px 2px 1px rgba(0, 0, 0, 0.32);
}
.vue-slider-mark-step {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background-color: rgba(0, 0, 0, 0.16);
}
.vue-slider-mark-label {
  font-size: 12px;
  white-space: nowrap;
  margin-top: 6px !important;
}
</style>
