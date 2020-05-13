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
            <b-form-input id="username" v-model="form.email" required size="sm"></b-form-input>
          </b-form-group>
          <b-form-group class="loginPassword mb-3">
            <label for="input">Password</label>
            <b-form-input type="password" id="password" v-model="form.password" size="sm" required></b-form-input>
          </b-form-group>
          <a class="forgot" href="#">Forgot Password?</a>
          <div class="button-area">
            <b-button
              size="sm"
              class="loginBtn app-color w-100 text-white align-baseline"
              type="submit"
              :disabled="loading"
            >
              <span>{{loading ? 'Logging In ':'Login' }}</span>
              <b-spinner small type="grow" v-show="loading"></b-spinner>
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

  methods: {
    login() {
      this.loading = true;
      const user = {
        username: this.form.email.trim(),
        key: this.form.password.trim()
      };
      this.$store
        .dispatch("login", user)
        .then(() => {
          location.reload();
        })
        .finally(() => {
          this.loading = false;
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
      #paypack {
        font-size: 50px;
        color: #017db3;
        line-height: 0%;
        letter-spacing: 5px;
        padding: 1.5rem 0;
        margin: 0;
      }
    }
    .slogan {
      letter-spacing: 0.5px;
      font-size: 16px;
      line-height: 1.3;
      color: #11374a;
      padding: 0.5rem 0;
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
        background: #f5f8fa;
        padding: 0.2rem;
        border-bottom: 2px solid #017db366;
        margin-top: 1.5rem !important;
        margin-bottom: 0 !important;
        user-select: none;
      }
      h3 {
        text-align: center;
        margin-bottom: 1.2rem;
        font-size: 21px;
        font-weight: bold;
        color: #017db3;
      }
    }

    label {
      font-size: 15px !important;
      font-weight: normal;
      font-family: inherit;
      text-transform: capitalize;
      letter-spacing: 0;
      color: #017db3;
      margin: 0;
      font-weight: 300;
    }
    .forgot {
      font-size: 14px;
      color: #017db3;
      margin-left: calc(100% - 119px);
      white-space: nowrap;
      font-weight: 300;
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
      &::placeholder {
        font-size: 13px;
        color: #00000075;
        -webkit-text-fill-color: #00000087;
      }
    }

    .button-area {
      display: flex;
      flex-direction: column;
      align-items: center;
      .loginBtn {
        border-radius: 2px;
        padding: 0.5rem 3rem;
        width: fit-content !important;
        border: none;
        text-transform: uppercase;
        margin-top: 1.5rem;
        font-weight: 300;
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
    h3{
      font-size: 18px !important;
    }
    fieldset {
      background: #ebf1f5 !important;
    }
  }
}
</style>
