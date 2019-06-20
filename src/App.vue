<template>
  <div id="app">
    <router-view/>
  </div>
</template>
<script>
import "../node_modules/chart.js/dist/Chart.js";
export default {};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
@font-face {
  font-family: FRegular;
  src: url("./assets/webFonts/HelveticaNeueLt.ttf");
}
@font-face {
  font-family: FThin;
  src: url("./assets/webFonts/HelveticaNeue.ttf");
}
@font-face {
  font-family: Fbold;
  src: url("./assets/webFonts/HelveticaNeueMed.ttf");
}
@font-face {
  font-family: Fbolder;
  src: url("./assets/webFonts/HelveticaNeueBd.ttf");
}
@font-face {
  font-family: Fitalic;
  src: url("./assets/webFonts/HelveticaNeueIt.ttf");
}
@import url("../node_modules/@fortawesome/fontawesome-free/css/all.css");

html,
body {
  margin: 0;
  padding: 0;
}
</style>
