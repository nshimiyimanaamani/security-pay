<template>
  <div class="dashboardWrapper">
    <div class="dashboardSidebar">
      <h1>PayPack</h1>
      <ul class="sidebarLinks">
        <router-link to="#">
          <li>Overview</li>
        </router-link>
        <router-link to="#">
          <li>Cells</li>
        </router-link>
        <router-link to="#">
          <li>Village</li>
        </router-link>
        <router-link to="/transactions">
          <li>Bank Accounts</li>
        </router-link>
        <router-link to="#">
          <li>Services</li>
        </router-link>
        <router-link to="#">
          <li>Penalties</li>
        </router-link>
      </ul>
      <button id="btn">
        <i class="fa fa-plus-circle"/> Add Services
      </button>
    </div>
    <div class="top-nav">
      <div class="dropdown">
        <p>
          YTD
          <i class="fa fa-caret-down"></i>
        </p>
      </div>
      <div class="search">
        <input type="search" name="search" id="search" placeholder="Search...">
      </div>
    </div>
    <div class="dashboardBody">
      <router-view/>
    </div>
  </div>
</template>

<script>
// @ is an alias to /src

export default {};
</script>
<style>
@import url("../assets/css/dashboardLayout.css");
</style>

