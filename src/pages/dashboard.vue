<template>
  <b-container>
    <vue-title title="Paypack | Dashboard" />
    <b-row align-v="start" class="m-auto p-0" style="width: 100%;height: 50%">
      <b-col xl="6" lg="6" md="6" sm="12" class="column position-relative">
        <b-card-body>
          <b-card-header>
            <i class="fa fa-th-large"></i>
            <h1 class>{{activeSector}} COLLECTING ACCOUNT</h1>
            <i class="fa fa-cog"></i>
          </b-card-header>
          <div style="height: 85%">
            <bar-chart :chart-data="chart1Data" :options="optionsChart1" :style="style"></bar-chart>
          </div>
        </b-card-body>
        <!-- end of chart 1 -->
      </b-col>
      <b-col xl="6" lg="6" md="6" sm="12" class="column position-relative">
        <b-card-body class="chart-2">
          <b-card-header>
            <i class="fa fa-th-large"></i>
            <h1 class>{{activeSector}} TOTAL COLLECTED</h1>
            <i
              class="fa fa-refresh cursor-pointer"
              @click="loadData"
              :class="{'fa-spin':chart2.state.loading}"
            ></i>
          </b-card-header>
          <div class="chart position-relative" style="height: 85%">
            <div v-if="!chart2.state.loading" class="h-100">
              <doughnut-chart
                :chart-data="chart2.data"
                v-if="chart2.data"
                :options="optionsChart2"
                :style="style"
              ></doughnut-chart>
              <div
                class="center-text justify-content-center align-items-center w-100 h-100 d-flex position-absolute"
                v-if="!chart2.state.error"
              >{{chart2.percentage}}%</div>
            </div>

            <b-card
              no-body
              class="position-absolute bg-transparent border-0 chartLoader align-items-center"
              v-if="chart2.state.loading"
            >
              <loader />
            </b-card>
            <b-card
              no-body
              class="position-absolute bg-transparent border-0 chartLoader align-items-center text-uppercase"
              v-if="chart2.state.error"
            >{{chart2.state.errorMessage}}</b-card>
          </div>
        </b-card-body>
      </b-col>
    </b-row>
    <b-row align-v="end" class="m-auto p-0" style="width: 100%;height: 50%">
      <b-col xl="12" lg="12" md="12" sm="12" class="column">
        <b-card-body class="chart-3">
          <b-card-header>
            <i class="fa fa-th-large"></i>
            <h1 class>{{activeSector}} SECTOR</h1>
            <i class="fa fa-cog"></i>
          </b-card-header>
          <div style="height: 85%">
            <line-chart
              :chart-data="chart3Data"
              :options="optionsChart3"
              :style="style"
              :tooltipData="chart3AdditionalData"
            />
          </div>
        </b-card-body>
      </b-col>
    </b-row>
  </b-container>
</template>
<script>
import BarChart from "../components/BarChart.vue";
import DoughnutChart from "../components/DaughnutChart.vue";
import LineChart from "../components/MixedCharts.vue";
import loader from "../components/loader";
export default {
  name: "dashboard",
  components: {
    BarChart,
    DoughnutChart,
    LineChart,
    loader
  },
  data() {
    return {
      chart2: {
        data: null,
        percentage: null,
        state: {
          loading: true,
          error: false,
          errorMessage: null
        }
      },
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
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    activeSector() {
      this.fetchData();
      return this.$store.getters.getActiveSector;
    },
    cellArray() {
      return this.$store.getters.getCellsArray;
    },
    style() {
      return {
        height: "100%"
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
  beforeMount() {
    this.loadData();
    this.fetchData();
  },
  methods: {
    loadData() {
      this.chart2.data = null;
      this.chart2.state.loading = true;
      this.chart2.state.error = false;
      this.axios
        .get(this.endpoint + "/metrics/sectors/" + this.activeSector)
        .then(res => {
          const data = res.data.data;
          const percentage = (data.payed * 100) / (data.payed + data.pending);
          this.chart2.percentage = percentage.toFixed();
          this.chart2.data = {
            labels: Object.keys(data),
            datasets: [
              {
                label: res.data.label,
                backgroundColor: ["#008b8bb3", "#e4e4ec"],
                borderColor: "white",
                data: [data.payed, data.pending]
              }
            ]
          };
        })
        .catch(err => {
          this.chart2.state.error = true;
          this.chart2.state.errorMessage = err.response
            ? err.response.data.error || err.response.data
            : "Error";
        })
        .finally(() => (this.chart2.state.loading = false));
    },
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
            barPercentage: 0.95,
            categoryPercentage: 1,
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
.chartLoader {
  top: 25%;
  right: 0;
  left: 0;
}
</style>
