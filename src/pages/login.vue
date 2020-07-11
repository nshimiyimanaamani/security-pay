<template>
  <div style="min-width: 300px; min-height:490px">
    <vue-title title="Paypack | Login" />
    <div class="loginPage">
      <b-row class="loginTitle">
        <div class="img">
          <img src="../../public/favicon.png" alt="paypack-logo" />
        </div>
        <div class="description">
          <b-row id="paypack" class="width-fit-content">PAYPACK</b-row>
          <b-row class="slogan">Easy way to collect and organise public fees</b-row>
        </div>
      </b-row>
      <b-row class="login-form">
        <b-form class="loginForm" @submit.prevent="login">
          <h3>LOGIN</h3>
          <b-form-group class="loginUsername mb-3">
            <label for="input">Username</label>
            <b-form-input
              id="username"
              v-model="form.email"
              required
              :disabled="loading"
              size="sm"
              trim
            ></b-form-input>
          </b-form-group>
          <b-form-group class="loginPassword mb-3">
            <label for="input">Password</label>
            <b-form-input
              type="password"
              id="password"
              v-model="form.password"
              :disabled="loading"
              size="sm"
              trim
              required
            ></b-form-input>
          </b-form-group>
          <a class="forgot" href="#">Forgot Password?</a>
          <div class="button-area">
            <b-button variant="info" class="loginBtn" type="submit" :disabled="loading">
              <span>{{loading ? 'Logging In ':'Login' }}</span>
              <i class="fa fa-spinner fa-spin" v-if="loading" />
            </b-button>
          </div>
        </b-form>
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
  destroyed() {
    this.loading = false;
    this.form.email = null;
    this.form.password = null;
  },
  methods: {
    login() {
      this.loading = true;
      const user = {
        username: this.form.email.trim(),
        password: this.form.password.trim()
      };
      this.axios
        .post("/accounts/login", user)
        .then(res => {
          sessionStorage.setItem("token", res.data.token);
          location.reload();
          localStorage.clear()
        })
        .catch(err => {
          this.loading = false;
          sessionStorage.removeItem("token");
          console.log(err, err.response, err.request);
          try {
            this.$snotify.error(err.response.data.error || err.response.data);
          } catch {
            this.$snotify.error(
              "Failed to login! Check your connection and try again"
            );
          }
        });
    }
  }
};
</script>
<style lang="scss">
.loginPage {
  color: #0000007a;
  min-width: 270px;
  min-height: 100vh;
  background: white;
  position: relative;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;

  .loginTitle {
    justify-content: center;
    align-items: center;
    flex: 2;
    width: auto;
    height: fit-content;
    user-select: none;
    .img {
      height: 6.5rem;
      width: 6.5rem;
      padding: 0;
      flex: 1;
      img {
        max-width: 100%;
        max-height: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
      }
    }
    .description {
      align-self: center;
      padding: 0;
      font-size: 14px;
      flex: 3;
      font-family: "Montserrat", sans-serif;
      #paypack {
        font-size: 50px;
        color: #017db3;
        line-height: 100%;
        letter-spacing: 5px;
        padding: 0;
        margin: 0;
        font-weight: 700;
      }
    }
    .slogan {
      letter-spacing: 0.4px;
      font-size: 16px;
      line-height: 1.4;
      color: #0c222b;
      padding: 0.2rem 0;
      margin: 0;
      font-weight: 300;
    }
  }
  .login-form {
    width: 100%;
    margin: auto;
    flex: 3;
    height: fit-content;
    display: flex;
    justify-content: center;
    align-items: center;

    .loginForm {
      background: white;
      padding: 2rem 1.2rem;
      border-radius: 6px;
      min-width: 250px;
      max-width: 380px;
      width: 100%;
      border: 1px solid #d0dae0;
      fieldset {
        background: #e8f0fe;
        padding: 0.2rem;
        border-bottom: 2px solid #017db366;
        margin-top: 1.5rem !important;
        margin-bottom: 0 !important;
        user-select: none;
      }
      h3 {
        text-align: center;
        margin-bottom: 2rem;
        font-size: 20px;
        font-weight: bold;
        font-family: "Montserrat", sans-serif;
        color: #017db3;
      }
    }

    label {
      font-size: 14px !important;
      text-transform: capitalize;
      letter-spacing: 0;
      color: #017db3;
      margin: 0;
      font-weight: 400;
      font-family: "Montserrat", sans-serif;
    }
    .forgot {
      font-size: 13px;
      margin-left: calc(100% - 126px);
      white-space: nowrap;
      color: #017db3;
      font-weight: 400;
      font-family: "Montserrat", sans-serif;
    }

    input {
      border-radius: 2px;
      padding: 0;
      font-weight: normal;
      font-size: 15px !important;
      height: fit-content;
      color: #000000b5 !important;
      -webkit-text-fill-color: #000000b5;
      letter-spacing: 0.7px;
      box-shadow: none !important;
      border: none !important;
      background: transparent !important;
      font-family: "Montserrat", sans-serif;
      &::placeholder {
        font-size: 13px;
        font-family: "Montserrat", sans-serif;
        color: #00000075;
        -webkit-text-fill-color: #00000087;
      }
    }

    .button-area {
      display: flex;
      flex-direction: column;
      align-items: center;
      .loginBtn {
        border: none;
        margin-top: 2rem;
        font-weight: 400;
        border-radius: 2px;
        width: fit-content;
        padding: 0.5rem 3rem;
        text-transform: uppercase;
      }
    }
  }
}
@media (max-width: 500px) {
  .loginTitle {
    flex-direction: column;
    flex: 1 !important;

    .img {
      flex: 0 !important;
    }
    .description {
      margin: 0 !important;
      margin-top: 15px !important;
      flex: 0 !important;
      display: flex;
      flex-direction: column;
      align-items: center;
      #paypack {
        font-size: 35px !important;
        letter-spacing: 3px !important;
        margin: auto;
      }
      .slogan {
        text-align: center;
        width: fit-content;
        margin: auto;
        font-size: 15px;
        padding-top: 0;
      }
    }
  }
  .loginForm {
    padding: 1.5rem 1rem !important;
    h3 {
      font-size: 18px !important;
    }
  }
}
</style>
