<template>
  <div class="dashboard-container">
    <div class="top-nav"></div>
    <div class="chart-1">
      <chart1/>
    </div>
    <div class="chart-2">
      <chart2/>
    </div>
    <div class="chart-3">
      <chart3/>
    </div>
  </div>
</template>

<script>
import HelloWorld from "../components/HelloWorld.vue";
import chart1 from "../components/chart1.vue";
import chart2 from "../components/chart2.vue";
import chart3 from "../components/chart3.vue";
export default {
  name: "dashboard",
  components: {
    HelloWorld,
    chart1,
    chart2,
    chart3
  }
};
</script>
<style>
@import url("../assets/css/dashboard.css");
</style>
