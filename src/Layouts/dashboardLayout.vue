<template>
  <div class="dashboardWrapper">
    <div class="dashboardSidebar">
      <h1>PayPack</h1>
      <ul class="sidebarLinks">
        <router-link to="/dashboard">
          <li>Overview</li>
        </router-link>
        <router-link to="/cells">
          <li>
            Cells
            <b-dropdown id="dropdown-right" class="m-md-2" v-if="this.$route.name == 'cells'">
              <template slot="button-content">
                <span class="fa fa-caret-down"></span>
              </template>
              <!-- <v-select></v-select> -->
              <b-dropdown-item>cell 1</b-dropdown-item>
              <b-dropdown-item>cell 2</b-dropdown-item>
              <b-dropdown-item>cell 4</b-dropdown-item>
              <b-dropdown-item>cell 5</b-dropdown-item>
              <b-dropdown-item>cell 6</b-dropdown-item>
              <b-dropdown-item>cell 7</b-dropdown-item>
              <b-dropdown-item>cell 8</b-dropdown-item>
              <b-dropdown-item>cell 9</b-dropdown-item>
              <b-dropdown-item>cell 10</b-dropdown-item>
            </b-dropdown>
          </li>
        </router-link>
        <router-link to="/village">
          <li>
            Village
            <b-dropdown id="dropdown-right" class="m-md-2" v-if="this.$route.name == 'village'">
              <template slot="button-content">
                <span class="fa fa-caret-down"></span>
              </template>
              <!-- <v-select></v-select> -->
              <b-dropdown-item>village 1</b-dropdown-item>
              <b-dropdown-item>village 2</b-dropdown-item>
              <b-dropdown-item>village 4</b-dropdown-item>
              <b-dropdown-item>village 5</b-dropdown-item>
              <b-dropdown-item>village 6</b-dropdown-item>
              <b-dropdown-item>village 7</b-dropdown-item>
              <b-dropdown-item>village 8</b-dropdown-item>
              <b-dropdown-item>village 9</b-dropdown-item>
              <b-dropdown-item>village 10</b-dropdown-item>
            </b-dropdown>
          </li>
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
        <b-dropdown id="dropdown-left" class="m-md-6" text="YTD">
          <b-dropdown-item>First Action</b-dropdown-item>
          <b-dropdown-item>Second Action</b-dropdown-item>
          <b-dropdown-item>Third Action</b-dropdown-item>
        </b-dropdown>
        <i class="fa fa-caret-down"></i>
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

export default {
  data() {
    return {};
  }
};
</script>
<style>
@import url("../assets/css/dashboardLayout.css");
</style>

