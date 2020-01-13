<template>
  <b-container>
    <vue-title title="Paypack | Cells" />
    <b-row align-v="start" class="m-auto p-0" style="width: 100%;height: 50%">
      <b-col xl="6" lg="6" md="6" sm="12" class="column">
        <b-card-body>
          <b-card-header>
            <i class="fa fa-th-large"></i>
            <h1 class>{{activeCell}} COLLECTING ACCOUNT</h1>
            <i class="fa fa-cog"></i>
          </b-card-header>
          <bar-chart
            v-if="chart1Data"
            :chart-data="chart1Data"
            :options="optionsChart1"
            :style="style"
          ></bar-chart>
        </b-card-body>
        <!-- end of chart 1 -->
      </b-col>
      <b-col xl="6" lg="6" md="6" sm="12" class="column">
        <b-card-body class="chart-2">
          <b-card-header>
            <i class="fa fa-th-large"></i>
            <h1 class>{{activeCell}} TOTAL COLLECTED</h1>
            <i class="fa fa-cog"></i>
          </b-card-header>
          <div class="chart" style="height:100%;position:relative">
            <doughnut-chart
              v-if="chart2Data"
              :chart-data="chart2Data"
              :options="optionsChart2"
              :style="style"
            ></doughnut-chart>
            <div class="center-text">{{percentage}}%</div>
          </div>
        </b-card-body>
      </b-col>
    </b-row>
    <b-row align-v="end" class="m-auto p-0" style="width: 100%;height: 50%">
      <b-col xl="12" lg="12" md="12" sm="12" class="column">
        <b-card-body class="chart-3">
          <b-card-header>
            <i class="fa fa-th-large"></i>
            <h1 class>{{activeCell}} CELL</h1>
            <i class="fa fa-cog"></i>
          </b-card-header>
          <line-chart
            v-if="chart3Data"
            :chart-data="chart3Data"
            :style="style"
            :tooltipData="chart3AdditionalData"
            :options="optionsChart3"
          />
        </b-card-body>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import BarChart from "../components/BarChart.vue";
import DoughnutChart from "../components/DaughnutChart.vue";
import LineChart from "../components/MixedCharts.vue";
export default {
  name: "cells",
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
    activeCell() {
      return this.$store.getters.getActiveCell;
    },
    cellArray() {
      return this.$store.getters.getCellsArray;
    },
    style() {
      return {
        height: "85%",
        "font-size": "15px"
      };
    },
    percentage() {
      const data = this.chart2Data.datasets[0].data;
      const percentage = (data[0] * 100) / (data[0] + data[1]);
      return percentage.toFixed();
    },
    optionsChart1() {
      return {
        tooltips: { enabled: false },
        hover: { mode: null },
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          yAxes: [
            {
              ticks: {
                min: 0,
                stepSize: 10,
                callback: (label, index, labels) => {
                  switch (label) {
                    case label:
                      return label + "M";
                  }
                }
              }
            }
          ],
          xAxes: [
            {
              gridLines: { display: false }
            }
          ]
        },
        legend: {
          display: false
        },
        layout: {
          padding: {
            left: 10,
            right: 20,
            top: 25,
            bottom: 0
          }
        }
      };
    },
    optionsChart2() {
      return {
        responsive: true,
        maintainAspectRatio: false,
        cutoutPercentage: 80,
        hover: { mode: null },
        tooltips: { enabled: false },
        legend: { display: false },
        layout: {
          padding: {
            left: 10,
            right: 10,
            top: 10,
            bottom: 10
          }
        }
      };
    },
    optionsChart3() {
      return {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          yAxes: [
            {
              ticks: {
                min: 0,
                stepSize: 10,
                callback: (label, index, labels) => {
                  switch (label) {
                    case label:
                      return label + "M";
                  }
                }
              }
            }
          ],
          xAxes: [
            {
              gridLines: { display: false }
            }
          ]
        },
        legend: {
          display: false
        },
        layout: {
          padding: {
            left: 10,
            right: 20,
            top: 25,
            bottom: 0
          }
        },
        tooltips: {
          displayColors: false,
          callbacks: {
            label: function(tooltipItem, data) {
              var label = [data.datasets[tooltipItem.datasetIndex].label] || "";

              if (label) {
                label = new Array();
                label.push(`Abishyuye: 30`);
                label.push("Abatarishyura: 40");
              }
              return label;
            }
          }
        }
      };
    }
  },
  watch: {
    activeCell() {
      handler: {
        this.fetchData();
      }
    }
  },
  beforeMount() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      window.Chart.defaults.global.defaultFontSize = 13.5;
      this.chart1Data = this.fillData(["BK", "MTN", "AIRTEL"]);
      this.chart2Data = this.fill2Data(["payed", "NotPayed"]);
      this.chart3Data = this.fillData(this.cellArray);
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
    fillData(labels) {
      let data = {
        labels: labels,
        datasets: [
          {
            label: "Data",
            barPercentage: 0.95,
            categoryPercentage: 1,
            backgroundColor: "#008b8bb3",
            data: this.getData(labels)
          }
        ]
      };
      return data;
    },
    fill2Data(labels) {
      let data = {
        labels: labels,
        datasets: [
          {
            label: "Data",
            backgroundColor: ["#008b8bb3", "#e4e4ec"],
            data: this.getData(labels)
          }
        ]
      };
      return data;
    },
    getRandomInt() {
      return Math.floor(Math.random() * (50 - 5 + 1)) + 5;
    },
    getData(labels) {
      let array = [];
      labels.forEach(element => {
        array.push(this.getRandomInt());
      });
      return array;
    }
  }
};
</script>
<style>
@import url("../assets/css/dashboard.css");
</style>
