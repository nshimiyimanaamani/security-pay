<template>
  <div>
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
import axios from "axios";
export default {
  data() {
    return {
      loading: false,
      color: "#fff",
      size: "25px",
      form: {
        email: "",
        password: "",
        cell: "admin"
      }
    };
  },
  computed: {
    endpoint() {
      return process.env.VUE_APP_PAYPACK_API
    }
  },
  methods: {
    register() {
      if (this.form.email != "" && this.form.password != "") {
        this.loading = true;
        axios
          .post(`${this.endpoint}/users/`,{
            email: this.form.email,
            password: this.form.password,
            cell: this.form.cell
          })
          .then(res => {
            console.log(res.data);
            this.$router.push('/')
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