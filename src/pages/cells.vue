<template>
  <b-container class="cells-page h-100 p-3" fluid>
    <vue-title title="Paypack | Cells" />
    <div class="charts">
      <b-row class="has-columns">
        <!-- first chart -->
        <div class="chart-container">
          <header class="chart-header">
            <i class="fa fa-th-large" />
            <h1 class>{{activeCell || 'CELL'}} COLLECTING ACCOUNT</h1>
            <i class="fa fa-cog" />
          </header>
          <div class="chart-body">
            <div class="h-100">
              <bar-chart
                v-if="chart1Data"
                :chart-data="chart1Data"
                :options="options.chart1"
                :style="style"
              />
            </div>
          </div>
        </div>
        <!-- second chart -->
        <div class="chart-container">
          <header class="chart-header">
            <i
              class="fa fa-refresh cursor-pointer"
              @click="loadData2"
              :class="{'fa-spin':chart2.state.loading}"
            />
            <h1>{{activeCell||'CELL'}} TOTAL COLLECTED</h1>
            <selector :object="config" v-on:ok="updated" />
          </header>
          <div class="chart-body px-5">
            <div v-if="!chart2.state.loading" class="h-100">
              <doughnut-chart
                :chart-data="chart2.data"
                v-if="chart2.data"
                :options="options.chart2"
                :style="style"
              />
              <div class="chart-center-text" v-if="!chart2.state.error">{{chart2.percentage}}%</div>
            </div>
            <vue-load v-if="chart2.state.loading" class="secondary-font chart-loader" />
            <div
              class="chart-error"
              v-if="chart2.state.error"
            >{{chart2.state.errorMessage||'No data available to display at this moment!'}}</div>
          </div>
        </div>
      </b-row>
      <b-row class="mb-3">
        <div class="chart-container">
          <header class="chart-header">
            <i
              class="fa fa-refresh cursor-pointer"
              @click="loadData3"
              :class="{'fa-spin':chart3.state.loading}"
            />
            <h1>{{activeCell||""}} CELL</h1>
            <selector :object="config" v-on:ok="updated" />
          </header>
          <div class="chart-body">
            <div v-if="!chart3.state.loading" class="h-100">
              <line-chart
                v-if="chart3.data"
                :chart-data="chart3.data"
                :tooltipData="chart3AdditionalData"
                :options="options.chart3"
                :style="style"
              />
            </div>
            <vue-load v-if="chart3.state.loading" class="secondary-font chart-loader" />
            <div
              class="chart-error"
              v-if="chart3.state.error"
            >{{chart3.state.errorMessage||'No data available to display at this moment!'}}</div>
          </div>
        </div>
      </b-row>
    </div>
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
            this.chart3.state.errorMessage = "NO DATA FOUND FOR THIS CELL";
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

<style lang="scss">
@import "../assets/css/dashboard.scss";
</style>
