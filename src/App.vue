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
  computed: {
    token() {
      return this.$store.getters.token;
    }
  },
  watch: {
    token() {
      if (!this.token) {
        this.$router.push("/");
      }
    }
  },
  mounted() {
    setInterval(() => {
      navigator.onLine ? (this.offline = false) : (this.offline = true);
    }, 10000);
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
  font-size: 15px;
}
.font-13 {
  font-size: 13px;
}
.cursor-pointer {
  cursor: pointer;
}
</style>
