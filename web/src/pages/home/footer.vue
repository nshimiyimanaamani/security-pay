<script>
/**
 * Footer component
 */
export default {
  name: "Footer",
};
</script>

<template>
  <div>
    <!--START FOOTER-->
    <footer class="footer">
      <div class="container">
        <div class="row">
          <div class="col-lg-4 margin-t-20">
            <h4>PAYPACK</h4>
            <div class="text-muted margin-t-20">
              <ul class="list-unstyled footer-list">
                <li>
                  <a href="#">Home</a>
                </li>
                 <li>
                  <a href="#">About us</a>
                </li>
                 <li>
                  <a href>Contact</a>
                </li>
              </ul>
            </div>
          </div>
          <!-- <div class="col-lg-4 margin-t-20">
            <h4>Information</h4>
            <div class="text-muted margin-t-20">
              <ul class="list-unstyled footer-list">
                <li>
                  <a href="#">About us</a>
                </li>
              </ul>
            </div>
          </div>
          <div class="col-lg-4 margin-t-20">
            <h4>Support</h4>
            <div class="text-muted margin-t-20">
              <ul class="list-unstyled footer-list">
                <li>
                  <a href>Contact</a>
                </li>
              </ul>
            </div>
          </div> -->
        </div>
      </div>
    </footer>
    <!--END FOOTER-->

    <!--START FOOTER ALTER-->
    <div class="footer-alt">
      <div class="container">
        <div class="row">
          <div class="col-lg-12">
            <div class="float-left pull-none">
              <p class="copy-rights text-muted mb-3 mb-sm-0">
                2019 - {{ new Date().getFullYear() }} © PAYPACK -
                <b-link href="https://quarksgroup.rw/" target="_blank"
                  >Quarks Group Ltd.</b-link
                >
              </p>
            </div>
            <div class="float-right pull-none">
              <ul class="list-inline social m-0">
                <li class="list-inline-item">
                  <a href class="social-icon">
                    <i class="mdi mdi-facebook"></i>
                  </a>
                </li>
                <li class="list-inline-item">
                  <a href class="social-icon">
                    <i class="mdi mdi-twitter"></i>
                  </a>
                </li>
                <li class="list-inline-item">
                  <a href class="social-icon">
                    <i class="mdi mdi-linkedin"></i>
                  </a>
                </li>
                <li class="list-inline-item">
                  <a href class="social-icon">
                    <i class="mdi mdi-google-plus"></i>
                  </a>
                </li>
                <li class="list-inline-item">
                  <a href class="social-icon">
                    <i class="mdi mdi-dribbble"></i>
                  </a>
                </li>
              </ul>
            </div>
            <div class="clearfix"></div>
          </div>
        </div>
      </div>
    </div>
    <!--START FOOTER ALTER-->
  </div>
</template>
<style lang="scss" scoped>
@import "../../assets/scss/style.scss";
</style>