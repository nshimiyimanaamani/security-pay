<template>
  <div id="app">
    <i class="fa fa-user-circle-o"></i>
    <router-view/>
  </div>
</template>
<script>
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

html,
body {
  margin: 0;
  padding: 0;
  font-family: FRegular;
}
</style>
