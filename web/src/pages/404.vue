<template>
  <b-container class="mw-100 error-page d-flex flex-column align-items-center">
    <h1>Oops! Something went wrong!</h1>
    <b-button class="back text-white p-3" size="sm" variant="info" @click="goBack">Return to Home</b-button>
  </b-container>
</template>

<script>
export default {
  name: "error-page",
  methods: {
    goBack() {
      this.$router.go(-2);
    }
  }
};
</script>

<style>
.error-page {
  width: 100%;
  min-height: 100vh;
  background-color: #f2eee8;
  font-family: inherit;
  justify-content: space-evenly;
  min-width: 250px;
}
.error-page h1 {
  font-weight: 800;
  color: #383a41;
  text-align: center;
  font-size: 2.5em;
  padding-top: 20px;
}
.error-page .back {
  width: 230px;
}
@media screen and (max-width: 400px) {
  .error-page h1 {
    padding-left: 20px;
    padding-right: 20px;
    font-size: 2em;
  }
  .error-page .back {
    margin: 0 auto;
    margin-top: 60px;
    margin-bottom: 50px;
    width: 200px;
  }
}
</style>