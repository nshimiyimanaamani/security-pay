<template>
  <b-container class="mw-100 sector-dashboard">
    <vue-title title="Paypack | Cells" />
    <b-row align-v="start" class="m-auto p-0 w-100 top">
      <b-col class="column position-relative">
        <b-card-body>
          <b-card-header class="align-items-center p-2 px-3" style="height: 40px">
            <i class="fa fa-th-large"></i>
            <h1 class="text-center">{{activeCell}} COLLECTING ACCOUNT</h1>
            <i class="fa fa-cog"></i>
          </b-card-header>
          <div style="height:calc(100% - 40px)">
            <div class="position-relative canvas">
              <bar-chart
                v-if="chart1Data"
                :chart-data="chart1Data"
                :options="options.chart1"
                :style="style"
              />
            </div>
          </div>
        </b-card-body>
        <!-- end of chart 1 -->
      </b-col>
      <b-col class="column position-relative">
        <b-card-body class="chart-2">
          <b-card-header class="align-items-center p-2 px-3" style="height: 40px">
            <i
              class="fa fa-refresh cursor-pointer"
              @click="loadData2"
              :class="{'fa-spin':chart2.state.loading}"
            />
            <h1 class="text-center">{{activeCell}} TOTAL COLLECTED</h1>
            <selector :object="config" v-on:ok="updated" />
          </b-card-header>
          <div class="chart position-relative" style="height:calc(100% - 40px)">
            <div v-if="!chart2.state.loading" class="position-relative canvas">
              <doughnut-chart
                :chart-data="chart2.data"
                v-if="chart2.data"
                :options="options.chart2"
                :style="style"
              />
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
            >
              <p>{{chart2.state.errorMessage}}</p>
            </b-card>
          </div>
        </b-card-body>
      </b-col>
    </b-row>
    <b-row align-v="end" class="m-auto p-0 w-100 h-50">
      <b-col class="column">
        <b-card-body class="chart-3">
          <b-card-header class="align-items-center p-2 px-3" style="height: 40px">
            <i
              class="fa fa-refresh cursor-pointer"
              @click="loadData3"
              :class="{'fa-spin':chart3.state.loading}"
            />
            <h1 class="text-center">{{activeCell}} CELL</h1>
            <selector :object="config" v-on:ok="updated" />
          </b-card-header>
          <div style="height:calc(100% - 40px)">
            <div v-if="!chart3.state.loading" class="h-100 position-relative canvas">
              <line-chart
                v-if="chart3.data"
                :chart-data="chart3.data"
                :style="style"
                :tooltipData="chart3AdditionalData"
                :options="options.chart3"
              />
            </div>

            <b-card
              no-body
              class="position-absolute bg-transparent border-0 chartLoader align-items-center"
              v-if="chart3.state.loading"
            >
              <loader />
            </b-card>
            <b-card
              no-body
              class="position-absolute bg-transparent border-0 chartLoader align-items-center text-uppercase"
              v-if="chart3.state.error"
            >
              <p>{{chart3.state.errorMessage}}</p>
            </b-card>
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
import yearSelectorVue from "../components/yearSelector.vue";
import options from "../components/scripts/chartOptions";
export default {
  name: "cells",
  components: {
    BarChart,
    DoughnutChart,
    LineChart,
    loader,
    selector: yearSelectorVue
  },
  data() {
    return {
      config: {
        configuring: false,
        year: new Date().getFullYear(),
        month: new Date().getMonth() + 1
      },
      options: options,
      chart2: {
        data: null,
        percentage: null,
        state: {
          loading: true,
          error: false,
          errorMessage: null
        }
      },
      chart3: {
        data: null,
        state: {
          loading: false,
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
    
    activeCell() {
      return this.$store.getters.getActiveCell;
    },
    cellArray() {
      return this.$store.getters.getCellsArray;
    },
    style() {
      return {
        height: "100%"
      };
    },
    months() {
      return this.$store.getters.getMonths;
    },
    percentage() {
      const data = this.chart2Data.datasets[0].data;
      const percentage = (data[0] * 100) / (data[0] + data[1]);
      return percentage.toFixed();
    },
    currentYear() {
      return new Date().getFullYear();
    },
    currentMonth() {
      return new Date().getMonth() + 1;
    }
  },
  watch: {
    activeCell() {
      handler: {
        this.fetchData();
        this.loadData2();
        this.loadData3();
      }
    }
  },
  beforeMount() {
    this.fetchData();
    this.loadData2();
    this.loadData3();
  },
  methods: {
    updated() {
      this.fetchData();
      this.loadData2();
      this.loadData3();
    },
    loadData2() {
      this.chart2.data = null;
      this.chart2.state.loading = true;
      this.chart2.state.error = false;
      const year = this.config.year;
      const month = this.config.month;
      this.axios
        .get(
          `/metrics/ratios/cells/${this.activeCell}?year=${year}&month=${month}`
        )
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
    loadData3() {
      this.chart3.data = null;
      this.chart3.state.loading = true;
      this.chart3.state.error = false;
      this.chart3.state.errorMessage = null;
      const year = this.config.year;
      const month = this.config.month;
      this.axios
        .get(
          `/metrics/ratios/cells/all/${this.activeCell}?year=${year}&month=${month}`
        )
        .then(res => {
          const data = res.data;
          if (data.length > 0) {
            let labels = data.map(item => item.label);
            this.chart3.data = {
              labels: labels,
              datasets: [
                {
                  label: "Data",
                  barPercentage: 0.95,
                  categoryPercentage: 1,
                  backgroundColor: "#008b8bb3",
                  data: this.getDataByLabels(data)
                },
                {
                  label: "Data",
                  type: "line",
                  backgroundColor: "transparent",
                  borderColor: "#095252ad",
                  pointRadius: 5,
                  borderDash: [10],
                  data: this.getDataByLabels(data)
                }
              ]
            };
          } else {
            this.chart3.state.error = true;
            this.chart3.state.errorMessage = "NO DATA FOUND FOR THIS SECTOR";
          }
        })
        .catch(err => {
          this.chart3.state.error = true;
          this.chart3.state.errorMessage = err.response
            ? err.response.data.error || err.response.data
            : "Error";
        })
        .finally(() => (this.chart3.state.loading = false));
    },
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
    },
    getDataByLabels(data) {
      let array = [];
      data.forEach(item => {
        const percentage =
          (item.data.payed * 100) / (item.data.payed + item.data.pending);
        array.push(percentage);
      });
      return array;
    }
  }
};
</script>
<style>
@import url("../assets/css/dashboard.css");
</style>
