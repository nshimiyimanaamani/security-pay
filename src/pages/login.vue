<template>
  <div>
    <vue-title title="Paypack | Login" />
    <div class="loginPage p-3">
      <b-row
        class="loginTitle m-auto justify-content-sm-center"
        style="height: 50%;width: fit-content"
      >
        <b-col cols="4" class="align-self-center p-0">
          <img src="../../public/favicon.png" alt="paypack-logo" style="width:6.5rem;height:6.5rem" />
        </b-col>
        <b-col class="align-self-center p-0 ml-2" style="font-size: 15px;height:6.5rem">
          <b-row id="paypack" class="py-4">PAYPACK</b-row>
          <b-row class="py-2">Easy way to collect and organise public fees</b-row>
        </b-col>
      </b-row>
      <b-row class="m-2 justify-content-center">
        <b-col sm="8" md="7" lg="4" xl="3">
          <b-form class="loginForm" @submit.prevent="login">
            <b-form-group class="loginUsername mb-4">
              <b-form-input id="username" v-model="form.email" required placeholder="Username..."></b-form-input>
            </b-form-group>
            <b-form-group class="loginPassword mb-4">
              <b-form-input
                type="password"
                id="password"
                v-model="form.password"
                required
                placeholder="Password..."
              ></b-form-input>
            </b-form-group>
            <b-button class="loginBtn" type="submit" :disabled="loading">
              <span>{{loading ? 'Logging In ':'Login' }}</span>
              <b-spinner small type="grow" v-show="loading"></b-spinner>
            </b-button>
            <br />
            <span class="float-right mt-1">
              Forgot Password?
              <a href="#">Get HELP</a>
            </span>
          </b-form>
        </b-col>
      </b-row>
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      loading: false,
      form: {
        email: null,
        password: null
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    }
  },
  methods: {
    login() {
      const email = this.form.email;
      const key = this.form.password;
      if (email && key) {
        this.loading = true;
        this.axios
          .post(this.endpoint + "/accounts/login", {
            username: email,
            password: key
          })
          .then(res => {
            sessionStorage.setItem("token", res.data.token);
            this.$router.push("dashboard");
            this.loading = false;
          })
          .catch(err => {
            delete sessionStorage.token;
            this.loading = false;
            if (navigator.onLine) {
              const error = err.response
                ? err.response.data.error || err.response.data
                : "an error occured";
              this.$snotify.error(error);
            } else {
              this.$snotify.error("Please connect to the internet");
            }
          });
      }
    }
  }
};
</script>
<style scoped>
@import url("../assets/css/login.css");
</style>
