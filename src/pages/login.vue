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
          <b-form-input
            type="email"
            id="username"
            v-model="form.email"
            required
            placeholder="Email..."
          ></b-form-input>
        </b-form-group>
        <b-form-group class="loginPassword">
          <b-form-input
            type="password"
            id="password"
            v-model="form.password"
            required
            placeholder="password..."
          ></b-form-input>
        </b-form-group>
        <div class="loginBtn">
          <a>
            <button type="submit">
              Log In
              <div class="loading" v-show="loading">
                <clip-loader :loading="loading" :color="color" :size="size"></clip-loader>
              </div>
            </button>
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
      color: "#fff",
      size: "25px",
      form: {
        email: "",
        password: ""
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    token() {
      return this.$store.getters.token;
    }
  },
  methods: {
    login() {
      if (this.form.email != "" && this.form.password != "") {
        const data = {
          email: this.form.email,
          password: this.form.password
        };
        this.loading = true;
        this.$store
          .dispatch("login", data)
          .then(res => {
              this.loading = false
          })
          .catch(err => {
            console.log(err);
            this.loading = false
          });
      }
    }
  }
};
</script>
<style scoped>
@import url("../assets/css/login.css");
</style>
