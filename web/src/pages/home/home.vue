<template>
  <div>
    <Navbar />

    <div v-scroll-spy>
      <!--START HOME-->
      <section class="section bg-home height-100vh" id="home">
        <div class="bg-overlay"></div>
        <div class>
          <div class>
            <div class="container slidero">
              <div class="row">
                <div class="col-lg-8 offset-lg-2 text-white text-center">
                  <h1 class="home-title">PAYPACK</h1>
                  <p class="padding-t-15 home-desc mx-auto">
                    Easy way to collect and organise public fees
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
        <!-- Particles -->
        <vue-particles class="position-relative" color="#FFF"></vue-particles>
      </section>

      <!--END HOME-->

      <Features />
      <GetStarted />
      <Contact />
      <Footer />
    </div>
  </div>
</template>


<script>
import Navbar from "./navbar";
import Features from "./features";
import GetStarted from "./getstarted";
import Contact from "./contact";
import Footer from "./footer";

export default {
  components: {
    Navbar,
    Features,
    GetStarted,
    Contact,
    Footer,
  },
  beforeMount() {
    this.$title("Home");
  },
};
</script>

<style  lang="scss" scoped>
@import "../../assets/scss/style.scss";
::v-deep .modal-content {
  background-color: #000000;
  color: white;
  background: none;
  border: none;
}
::v-deep .modal-header {
  border: none;
}
::v-deep .modal-backdrop {
  opacity: 0.5;
}
</style>