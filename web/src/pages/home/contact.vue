<script>
/**
 * Contact component
 */
export default {
  name: "Contact",
};
</script>

<template>
  <!-- CONTACT FORM START-->
  <section class="section" id="contact">
    <div class="container">
      <div class="row">
        <div class="col-lg-8 offset-lg-2">
          <h1 class="section-title text-center">Get In Touch</h1>
          <div class="section-title-border margin-t-20"></div>
          <p
            class="section-subtitle text-muted text-center font-secondary padding-t-30"
          >
            Got any query? visit us or send a message to our help desk.
          </p>
        </div>
      </div>
      <div class="row">
        <div class="col-lg-4">
          <div class="mt-4 pt-4">
            <p class="mt-4">
              <span class="h5">Office Address</span>
              <br />
              <span class="text-muted d-block mt-2">KG 1 Ave, Kigali, Rwanda</span>
            </p>
            <p class="mt-4">
              <span class="h5">Working Hours:</span>
              <br />
              <span class="text-muted d-block mt-2">8:00AM To 5:00PM</span>
            </p>
          </div>
        </div>
        <div class="col-lg-8">
          <div class="custom-form mt-4 pt-4">
            <div id="message"></div>
            <form
              method="post"
              action="php/contact.php"
              name="contact-form"
              id="contact-form"
            >
              <div class="row">
                <div class="col-lg-6">
                  <div class="form-group mt-2">
                    <input
                      name="name"
                      id="name"
                      type="text"
                      class="form-control"
                      placeholder="Your name*"
                    />
                  </div>
                </div>
                <div class="col-lg-6">
                  <div class="form-group mt-2">
                    <input
                      name="email"
                      id="email"
                      type="email"
                      class="form-control"
                      placeholder="Your email*"
                    />
                  </div>
                </div>
              </div>
              <div class="row">
                <div class="col-lg-12">
                  <div class="form-group mt-2">
                    <input
                      type="text"
                      class="form-control"
                      id="subject"
                      placeholder="Your Subject.."
                    />
                  </div>
                </div>
              </div>
              <div class="row">
                <div class="col-lg-12">
                  <div class="form-group mt-2">
                    <textarea
                      name="comments"
                      id="comments"
                      rows="4"
                      class="form-control"
                      placeholder="Your message..."
                    ></textarea>
                  </div>
                </div>
              </div>
              <div class="row">
                <div class="col-lg-12 text-right">
                  <b-button
                    variant="outline-primary"
                    id="submit"
                    name="send"
                  >Send message</b-button>
                  <div id="simple-msg"></div>
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </section>
  <!-- CONTACT FORM END-->
</template>
<style lang="scss" scoped>
@import "../../assets/scss/style.scss";
</style>