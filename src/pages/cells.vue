<template>
  <div class="dashboard-container">
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
import chart1 from "../components/cellsChart1.vue";
import chart2 from "../components/cellsChart2.vue";
import chart3 from "../components/cellsChart3.vue";
export default {
  name: "dashboard",
  components: {
    chart1,
    chart2,
    chart3
  }
};
</script>
<style>
@import url("../assets/css/dashboard.css");
</style>
