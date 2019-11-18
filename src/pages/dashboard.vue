<template>
  <div class="dashboard-container">
    <!-- start of chart 1 -->
    <div class="row">
      <b-card-body>
        <b-card-header>
          <i class="fa fa-th-large"></i>
          <h1 class>{{activeSector}} COLLECTING ACCOUNT</h1>
          <i class="fa fa-cog"></i>
        </b-card-header>
        <bar-chart :chart-data="chart1Data" :style="style"></bar-chart>
      </b-card-body>
      <!-- end of chart 1 -->

      <!-- start of chart 2 -->
      <b-card-body class="chart-2">
        <b-card-header>
          <i class="fa fa-th-large"></i>
          <h1 class>{{activeSector}} TOTAL COLLECTED</h1>
          <i class="fa fa-cog"></i>
        </b-card-header>
        <div class="chart" style="height:100%;position: relative">
          <doughnut-chart :chart-data="chart2Data" :style="style"></doughnut-chart>
          <div class="center-text">{{percentage}}%</div>
        </div>
      </b-card-body>
      <!-- end of chart 2 -->
    </div>
    <!-- start of chart 3 -->
    <div class="row">
      <b-card-body class="chart-3">
        <b-card-header>
          <i class="fa fa-th-large"></i>
          <h1 class>{{activeSector}} SECTOR</h1>
          <i class="fa fa-cog"></i>
        </b-card-header>
        <line-chart :chart-data="chart3Data" :style="style" :tooltipData="chart3AdditionalData" />
      </b-card-body>
    </div>
    <!-- end of chart 3 -->
  </div>
</template>

<script>
import BarChart from "../components/chart1.vue";
import DoughnutChart from "../components/chart2.vue";
import LineChart from "../components/chart3.vue";
export default {
  name: "dashboard",
  components: {
    BarChart,
    DoughnutChart,
    LineChart
  },
  data() {
    return {
      chart1Data: null,
      chart2Data: null,
      chart3Data: null,
      chart3AdditionalData: {
        abishyuye: this.getRandomInt(),
        abatarishyura: this.getRandomInt()
      }
    };
  },
  computed: {
    activeSector() {
      this.fetchData();
      return this.$store.getters.getActiveSector;
    },
    cellArray() {
      return this.$store.getters.getCellsArray;
    },
    style() {
      return {
        height: "85%",
        "font-size": "15px",
        content: "60%"
      };
    },
    percentage() {
      const data = this.chart2Data.datasets[0].data;
      const percentage = (data[0] * 100) / (data[0] + data[1]);
      return percentage.toFixed();
    }
  },
  beforeMount() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      window.Chart.defaults.global.defaultFontSize = 13.5;
      this.chart1Data = this.fillData(["BK Acc", "MTN", "AIRTEL"], 3);
      this.chart2Data = this.fill2Data(["abishyuye", "abasigaye"], 2);
      this.chart3Data = this.fillData(this.cellArray, this.cellArray.length);
      this.chart3Data.datasets.push({
        label: "Data",
        backgroundColor: "transparent",
        borderColor: "#095252ad",
        pointRadius: 5,
        borderDash: [10],
        data: this.chart3Data.datasets[0].data,
        type: "line"
      });
    },
    fillData(labels, num) {
      const data = {
        labels: labels,
        datasets: [
          {
            label: "Data",
            backgroundColor: "#008b8bb3",
            data: this.getRandomInt(num)
          }
        ]
      };
      return data;
    },
    fill2Data(labels, num) {
      const data = {
        labels: labels,
        datasets: [
          {
            label: "Data",
            backgroundColor: ["#008b8bb3", "#e4e4ec"],
            borderColor: "white",
            data: this.getRandomInt(num)
          }
        ]
      };
      return data;
    },
    rand() {
      return Math.floor(Math.random() * (50 - 5 + 1)) + 5;
    },
    getRandomInt(j) {
      let randArray = [];
      for (let i = 0; i < j; i++) {
        randArray.push(this.rand());
      }
      return randArray;
    }
  }
};
</script>
<style>
@import url("../assets/css/dashboard.css");
</style>
