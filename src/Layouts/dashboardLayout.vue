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
        <router-link to="#">
          <li>Bank Accounts</li>
        </router-link>
        <router-link to="#">
          <li>Services</li>
        </router-link>
        <router-link to="#">
          <li>Penalties</li>
        </router-link>
      </ul>
      <button>
        <i class="fa fa-plus-circle"/> Add Services
      </button>
    </div>
    <div class="dashboardBody">
      <router-view/>Lorem ipsum dolor sit amet consectetur adipisicing elit. Porro praesentium, accusamus est sapiente obcaecati totam nesciunt, alias, ratione numquam exercitationem officiis tenetur dolor fugiat? Iusto ipsam, ab adipisci eos, eveniet doloremque ad aspernatur dolorum veritatis sunt quisquam magnam illo sint! Fugiat voluptas harum quis repellendus facere distinctio rerum neque. Ad?
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

