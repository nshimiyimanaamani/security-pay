<script>
export default {
  name: "Started-blog",
};
</script>

<template>
  <!--START GET STARTED-->
  <section
    class="section section-lg bg-get-start"
    :style="{ background: 'url(' + require('@/assets/images/img-2.jpg') + ')' }"
  >
    <div class="bg-overlay"></div>
    <div class="container">
      <div class="row">
        <div class="col-lg-8 offset-lg-2 text-center">
          <h1 class="get-started-title text-white">Let's Get Started</h1>
          <div class="section-title-border margin-t-20 bg-white"></div>
          <p
            class="section-subtitle font-secondary text-white text-center padding-t-30"
          >
           Interested in using our technology?
          </p>
          <b-button variant="white" to="/login" class="margin-t-20">
            Get Started
          </b-button>
        </div>
      </div>
    </div>
  </section>
  <!--END GET STARTED-->
</template>
<style lang="scss" scoped>
@import "../../assets/scss/style.scss";
</style>