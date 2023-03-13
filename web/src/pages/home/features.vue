<script>
/**
 * Features component
 */
export default {
  name: "Features",
};
</script>

<template>
  <!--START FEATURES-->
  <section class="section" id="about">
    <div class="container">
       <div class="row">
        <div class="col-lg-8 offset-lg-2">
          <h1 class="section-title text-center">About Paypack</h1>
          <div class="section-title-border margin-t-20"></div>
        </div>
      </div>
      <div class="row align-items-center margin-t-30">
        <div class="col-lg-6 order-2 order-lg-1">
          <div class="features-box mt-5 mt-lg-0">
            <h3>
              Community based security fee collection is one of our local community's bedrock. 
            </h3>
            <p class="text-muted web-desc">
              So far the fees have been collected and accounted for through a manual and paper-based method that is archaic by today’s technological standards. And like many manual methods, this system is full of inefficiencies that negatively affect all involved stakeholders notably the beneficiaries of the services (Clients), those who provide the services (service providers), and the local government authorities and national level regulators (Authorities and Regulators).

            </p>
          </div>
        </div>
        <div class="col-lg-6 order-1 order-lg-2">
          <div class="features-img mx-auto mr-lg-0">
            <img
              src="@/assets/images/growth-analytics.svg"
              alt="macbook image"
              class="img-fluid"
            />
          </div>
        </div>
      </div>
    </div>
  </section>
  <!--END FEATURES-->
</template>
<style lang="scss" scoped>
@import "../../assets/scss/style.scss";
</style>