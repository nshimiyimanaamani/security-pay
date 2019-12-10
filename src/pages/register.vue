<template>
  <div>
    <vue-title title="Paypack | Register" />
    <div class="registerPage">
      <div class="registerTitle">
        <p>Welcome To</p>
        <br />
        <span id="paypack">PayPack</span>
        <br />
        <span id="easy">
          Easy way to collect and organise
          <br />public fees
        </span>
      </div>
      <b-form class="registerForm" @submit.prevent="register">
        <b-form-group class="registerUsername">
          <b-form-input
            type="email"
            v-model="form.email"
            id="username"
            placeholder="Email..."
            required
          ></b-form-input>
        </b-form-group>
        <b-form-group class="registerPassword">
          <b-form-input
            type="password"
            v-model="form.password"
            id="password"
            placeholder="Password..."
            required
          ></b-form-input>
        </b-form-group>
        <div class="registerBtn">
          <a>
            <button>
              Register
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
        email: null,
        password: null,
        cell: "admin"
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    }
  },
  methods: {
    register() {
      const email = this.form.email;
      const key = this.form.password;
      if (email && key) {
        this.loading = true;
        this.axios
          .post(this.endpoint + "/users/", {
            email: email,
            password: key,
            cell: this.form.cell
          })
          .then(res => {
            this.$router.push("/");
            this.$snotify.success(`user successfully registered!`);
            this.loading = false;
          })
          .catch(err => {
            console.log(err);
            this.loading = false;
            this.$snotify.error(
              `user registration Failed! please try again Later `
            );
          });
      }
    }
  }
};
</script>
<style scoped>
@import url("../assets/css/register.css");
</style>