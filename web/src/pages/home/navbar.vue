<template>
  <!-- STRAT NAVBAR -->
  <nav
    class="navbar navbar-expand-lg fixed-top navbar-custom sticky sticky-dark"
    id="navbar"
  >
    <b-container class="container">
      <!-- LOGO -->
      <b-navbar-brand class="logo text-uppercase" to="/">
        <i class="mdi mdi-alien"></i>PAYPACK
      </b-navbar-brand>

      <button
        class="navbar-toggler"
        type="button"
        data-toggle="collapse"
        data-target="#navbarCollapse"
        aria-controls="navbarCollapse"
        aria-expanded="false"
        aria-label="Toggle navigation"
        @click="toggleMenu()"
      >
        <i class="mdi mdi-menu"></i>
      </button>
      <div class="collapse navbar-collapse" id="navbarCollapse">
        <ul
          class="navbar-nav navbar-center"
          id="mySidenav"
          v-scroll-spy-active="{ class: 'active' }"
        >
          <li class="nav-item">
            <a
              v-scroll-to="{ el: '#home' }"
              class="nav-link"
              data-scroll-spy-id="home"
              href="javascript: void(0);"
              >Home</a
            >
          </li>
          <li class="nav-item">
            <a
              v-scroll-to="{ el: '#about', offset: -50 }"
              class="nav-link"
              data-scroll-spy-id="about"
              href="javascript: void(0);"
              >About</a
            >
          </li>
          <!-- <li class="nav-item">
            <a
              v-scroll-to="{ el: '#services', offset: -5 }"
              class="nav-link"
              data-scroll-spy-id="services"
              href="javascript: void(0);"
              >Services</a
            >
          </li> -->
          <li class="nav-item">
            <a
              v-scroll-to="{ el: '#contact' }"
              class="nav-link"
              data-scroll-spy-id="contact"
              href="javascript: void(0);"
              >Contact</a
            >
          </li>
        </ul>
        <div class="nav-button ml-auto">
          <ul class="nav navbar-nav navbar-right">
            <li>
              <b-button
                variant="primary"
                to="/login"
                pill tag="button"
                class="navbar-btn px-4 py-2"
              >
                Login
              </b-button>
            </li>
          </ul>
        </div>
      </div>
    </b-container>
  </nav>
  <!-- END NAVBAR -->
</template>

<script>
/**
 * Navbar component
 */
export default {
  mounted: () => {
    window.onscroll = function () {
      onwindowScroll();
    };
    var navbar = document.getElementById("navbar");
    function onwindowScroll() {
      if (
        document.body.scrollTop > 40 ||
        document.documentElement.scrollTop > 40
      ) {
        navbar.style.backgroundColor = "#272a33";
        navbar.style.padding = "10px";
      } else {
        navbar.style.backgroundColor = "";
        navbar.style.padding = "20px";
      }
    }
  },
  methods: {
    /**
     * Toggle menu
     */
    toggleMenu() {
      document.getElementById("navbarCollapse").classList.toggle("show");
    },
  },
};
</script>
<style lang="scss" scoped>
@import "../../assets/scss/style.scss";
</style>
