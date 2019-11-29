<template>
  <div>
    <div class="loginPage">
      <div class="loginTitle">
        <p>Welcome To</p>
        <br />
        <span id="paypack">PayPack</span>
        <br />
        <span id="easy">
          Easy way to collect and organise
          <br />public fees
        </span>
      </div>
      <b-form class="loginForm" @submit.prevent="login()">
        <b-form-group class="loginUsername">
          <b-form-input type="email" id="username" v-model="form.email" required placeholder="Email..."></b-form-input>
        </b-form-group>
        <b-form-group class="loginPassword">
          <b-form-input type="password" id="password" v-model="form.password" required placeholder="password..."></b-form-input>
        </b-form-group>
        <div class="loginBtn">
          <a>
            <b-button variant="info" type="submit" :disabled="loading">
              <span>{{loading ? 'Logging In ':'Login' }}</span>
              <b-spinner small type="grow" v-show="loading"></b-spinner>
            </b-button>
          </a>
          <br />
          <span id="forgot">
            Forgot Password?
            <a href="#">Get HELP</a>
          </span>
        </div>
      </b-form>
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
          .post(this.endpoint + "/users/tokens", {
            email: email,
            password: key
          })
          .then(res => {
            sessionStorage.setItem("token", res.data.token);
            this.$router.push("dashboard");
            this.loading = false;
          })
          .catch(err => {
            console.warn(err);
            delete sessionStorage.getItem("token");
            this.loading = false;
          });
      }
    }
  }
};

</script>
<style scoped>
@import url("../assets/css/login.css");

</style>
