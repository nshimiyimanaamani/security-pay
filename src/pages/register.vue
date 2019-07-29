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
      <div class="registerForm">
        <div class="registerUsername">
          <input type="email" name="username" id="username" placeholder="Email..." required />
        </div>
        <div class="registerPassword">
          <input type="password" name="password" id="password" placeholder="Password..." required />
        </div>
        <div class="registerBtn">
          <a>
            <button @click="register">
              Register
              <div class="lds-css ng-scope" v-if="this.loading">
                <div class="lds-rolling">
                  <div></div>
                </div>
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
  </div>
</template>

<script>
import axios from "axios";
export default {
  computed: {
    loading() {
      return this.$store.state.loading;
    },
    endpoint(){
      return this.$store.state.endPoint;
    }
  },
  methods: {
    register() {
      let email = document.querySelector("#username").value;
      let password = document.querySelector("#password").value;

      if (email != undefined && password != undefined) {
        this.$store.state.loading = true;
        axios
          .post(`${this.endpoint}/users/`, {
            email: `${email}`,
            password: `${password}`
          })
          .then(res => {
            console.log(res.data);
            let text = "User Successfully Registered";
            this.confirmation(text);
            this.$store.state.loading = false;
          })
          .catch(err => {
            console.log(err);
            let text = "User Not Registered";
            this.confirmation(text);
            this.$store.state.loading = false;
          });
      }
    },
    confirmation(text) {
      this.$bvModal.msgBoxOk(text, {
        title: "Confirmation",
        centered: true
      });
    }
  }
};
</script>
<style scoped>
@import url("../assets/css/register.css");
</style>
