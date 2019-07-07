<template>
  <div>
    <div class="loginPage">
      <div class="loginTitle">
        <p>Welcome To</p>
        <br>
        <span id="paypack">PayPack</span>
        <br>
        <span id="easy">
          Easy way to collect and organise
          <br>public fees
        </span>
      </div>
      <div class="loginForm">
        <div class="loginUsername">
          <input type="text" name="username" id="username" placeholder="Username...">
        </div>
        <div class="loginPassword">
          <input type="password" name="password" id="password" placeholder="Password...">
        </div>
        <div class="loginBtn">
          <a href="/dashboard">
            <button>Log In</button>
          </a>

          <br>
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
export default {};
</script>

<style>
</style>
<style>
@import url("../assets/css/login.css");
</style>

