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
  }
};
</script>

<style>
@import url("./assets/css/main.css");
</style>
