<template>
  <div style="min-width: 300px; min-height:490px">
    <vue-title title="Paypack | Login" />
    <div class="loginPage p-3 position-relative h-100 d-flex flex-column align-items-center">
      <b-row class="loginTitle justify-content-center width-fit-content">
        <b-col cols="4" class="align-self-center p-0">
          <img src="../../public/favicon.png" alt="paypack-logo" style="width:6.5rem;height:6.5rem" />
        </b-col>
        <b-col class="align-self-center p-0 ml-2 font-14" style="height:6.5rem">
          <b-row id="paypack" class="py-4 width-fit-content">PAYPACK</b-row>
          <b-row class="py-2 slogan">Easy way to collect and organise public fees</b-row>
        </b-col>
      </b-row>
      <b-row class="login-form m-auto justify-content-center">
        <b-col>
          <b-form class="loginForm" @submit.prevent="login">
            <b-form-group class="loginUsername mb-3">
              <b-form-input
                id="username"
                v-model="form.email"
                required
                placeholder="Username..."
                size="sm"
              ></b-form-input>
            </b-form-group>
            <b-form-group class="loginPassword mb-4">
              <b-form-input
                type="password"
                id="password"
                v-model="form.password"
                size="sm"
                required
                placeholder="Password..."
              ></b-form-input>
            </b-form-group>
            <b-button
              size="sm"
              class="loginBtn app-color w-100 text-white align-baseline"
              type="submit"
              :disabled="loading"
            >
              <span>{{loading ? 'Logging In ':'Login' }}</span>
              <b-spinner small type="grow" v-show="loading"></b-spinner>
            </b-button>
            <br />
            <span class="float-right mt-1 font-13">
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

  methods: {
    login() {
      // this.loading = true;
      // const user = {
      //   username: this.form.email.trim(),
      //   key: this.form.password.trim()
      // };
      // this.$store
      //   .dispatch("login", user)
      //   .then(() => {
      //     location.reload();
      //   })
      //   .finally(() => {
      //     this.loading = false;
      //   });
      alert("App Temporary Down! Wait for further notice.");
    }
  }
};
</script>
<style>
.loginPage {
  color: #0000007a;
  min-width: 270px;
  min-height: 100vh;
}
.loginTitle {
  min-height: 250px;
}
.loginTitle #paypack {
  font-size: 50px;
  color: #3a82a1;
  line-height: 0%;
  letter-spacing: 5px;
}
.login-form {
  width: 100%;
  min-width: 300px;
  max-width: 400px;
}
@media (max-width: 389px) {
  .loginTitle .slogan {
    text-align: center;
    width: fit-content;
    margin: auto;
  }
  .loginTitle #paypack {
    margin: auto;
  }
}
</style>
