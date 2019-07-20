<template>
  <div id="app">
    <router-view />
  </div>
</template>
<script>
import "../node_modules/chart.js/dist/Chart.js";
export default {
  computed: {
    user() {
      return this.$store.state.user1;
    }
  },
  beforeMount(){
    this.$store.dispatch('startup_function')
  }
};
</script>

<style>
@font-face {
  font-family: ubuntu;
  src: url("./assets/webFonts/Ubuntu-Light.ttf");
}
@font-face {
  font-family: ubuntuRegular;
  src: url("./assets/webFonts/Ubuntu-Regular.ttf");
}
@font-face {
  font-family: ubuntuBold;
  src: url("./assets/webFonts/Ubuntu-Bold.ttf");
}
@import url("../node_modules/@fortawesome/fontawesome-free/css/all.css");
@import url("../node_modules/bootstrap/dist/css/bootstrap.css");
@import url("../node_modules/bootstrap-vue/dist/bootstrap-vue.css");
html,
body {
  font-family: ubuntuRegular, Helvetica, Arial, sans-serif !important;
  color: #2c3e50;
  margin: 0;
  padding: 0;
  font-size: 14px;
}
* {
  outline: none;
  letter-spacing: 0.5px;
  -webkit-font-smoothing: antialiased !important;
  -moz-osx-font-smoothing: grayscale !important;
}
button,input{
  font-family: ubuntuRegular, Helvetica, Arial, sans-serif !important;
}
</style>
