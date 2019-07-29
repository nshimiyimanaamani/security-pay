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
      <div class="loginForm">
        <div class="loginUsername">
          <input type="email" name="username" id="username" placeholder="Email..." />
        </div>
        <div class="loginPassword">
          <input type="password" name="password" id="password" placeholder="Password..." />
        </div>
        <div class="loginBtn">
          <a>
            <button @click="login">
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
      </div>
    </div>
    <b-modal id="error-message" hide-footer>
      <template slot="modal-title" class="text-center">Confirmation</template>
      <div class="d-block">
        <h3>user not found, try using right credentials</h3>
      </div>
      <b-button block @click="$bvModal.hide('error-message')">OK</b-button>
    </b-modal>
  </div>
</template>

<script>
import axios from "axios";
export default {
  data(){
    return {
      loading: false,
      color: "#fff",
      size: "25px"
    }
  },
  computed: {
    endpoint() {
      return this.$store.state.endPoint;
    }
  },
  methods: {
    login() {
      window.location.pathname = "/dashboard";
      let email = document.querySelector("#username").value;
      let password = document.querySelector("#password").value;
      // if (email != "" && password != "") {
      //   this.$store.state.loading = true;
      //   axios
      //     .post(`${this.endpoint}/users/tokens`, {
      //       email: `${email}`,
      //       password: `${password}`
      //     })
      //     .then(res => {
      //       console.log(res.data);
      //       window.location.pathname = "/dashboard";
      //       this.$store.state.loading = false;
      //     })
      //     .catch(err => {
      //       console.log(err);
      //       let text = "User Not Found";
      //       // this.confirmation(text);
      //       this.$bvModal.show("error-message");
      //       this.$store.state.loading = false;
      //     });
      // }
    },
    confirmation(text) {
      this.$bvModal.msgBoxOk(text, {
        title: "Confirmation"
      });
    }
  }
};
</script>
<style scoped>
@import url("../assets/css/login.css");
</style>
