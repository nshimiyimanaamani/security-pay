<template>
  <b-container class="sector-page p-3 h-100" fluid>
    <vue-title title="Paypack | Dashboard" />
    <div class="charts">
      <b-row class="has-columns">
        <!-- first chart -->
        <div class="chart-container">
          <header class="chart-header">
            <i class="fa fa-th-large" />
            <h1 class>{{activeSector || 'SECTOR'}} COLLECTING ACCOUNT</h1>
            <i class="fa fa-cog" />
          </header>
          <div class="chart-body">
            <div class="h-100">
              <bar-chart :chart-data="chart1Data" :options="options.chart1" :style="style" />
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
            <h1>{{activeSector||'SECTOR'}} TOTAL COLLECTED</h1>
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
            <vue-load v-if="chart2.state.loading" class="primary-font chart-loader" />
            <div
              class="chart-error primary-font"
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
            <h1>{{activeSector||""}} SECTOR</h1>
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
            <vue-load v-if="chart3.state.loading" class="primary-font chart-loader" />
            <div
              class="chart-error primary-font"
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
import yearSelectorVue from "../components/yearSelector.vue";
import options from "../components/scripts/chartOptions";
export default {
  name: "dashboard",
  components: {
    BarChart,
    DoughnutChart,
    LineChart,
    selector: yearSelectorVue,
  },
  data() {
    return {
      config: {
        configuring: false,
        year: new Date().getFullYear(),
        month: new Date().getMonth() + 1,
      },
      options: options,
      chart2: {
        data: null,
        percentage: null,
        state: {
          loading: true,
          error: false,
          errorMessage: null,
        },
      },
      chart3: {
        data: null,
        state: {
          loading: true,
          error: false,
          errorMessage: null,
        },
      },
      chart1Data: null,
      chart2Data: null,
      chart3Data: null,
      chart3AdditionalData: {
        abishyuye: this.getRandomInt(),
        abatarishyura: this.getRandomInt(),
      },
    };
  },
  computed: {
    activeSector() {
      this.fetchData();
      return this.$store.getters.getActiveSector;
    },
    cellArray() {
      const { province, district, sector } = this.location;
      return this.$cells(province, district, sector) || [];
    },
    location() {
      return this.$store.getters.location;
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
    },
    style() {
      return {
        position: "relative",
        height: "100%",
      };
    },
  },
  mounted() {
    this.loadData2();
    this.loadData3();
    this.fetchData();
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
      this.chart2.state.errorMessage = null;
      const year = this.config.year;
      const month = this.config.month;
      this.axios
        .get(
          `/metrics/ratios/sectors/${this.activeSector}?year=${year}&month=${month}`
        )
        .then((res) => {
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
                data: [data.payed, data.pending],
              },
            ],
          };
        })
        .catch((err) => {
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
          `/metrics/ratios/sectors/all/${this.activeSector}?year=${year}&month=${month}`
        )
        .then((res) => {
          const data = res.data.filter(
            (item) => this.cellArray.indexOf(item.label) != -1
          );
          if (data.length > 0) {
            let labels = data.map((item) => item.label);
            this.chart3.data = {
              labels: labels,
              datasets: [
                {
                  label: "Data",
                  barPercentage: 0.95,
                  categoryPercentage: 1,
                  backgroundColor: "#008b8bb3",
                  data: this.getDataByLabels(data),
                },
                {
                  label: "Data",
                  type: "line",
                  backgroundColor: "transparent",
                  borderColor: "#095252ad",
                  pointRadius: 5,
                  borderDash: [10],
                  data: this.getDataByLabels(data),
                },
              ],
            };
          } else {
            this.chart3.state.error = true;
            this.chart3.state.errorMessage = "NO DATA FOUND FOR THIS SECTOR";
          }
        })
        .catch((err) => {
          this.chart3.state.error = true;
          this.chart3.state.errorMessage = err.response
            ? err.response.data.error || err.response.data
            : "Error";
        })
        .finally(() => (this.chart3.state.loading = false));
    },
    fetchData() {
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
        type: "line",
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
            data: this.getRandomInt(num),
          },
        ],
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
            data: this.getRandomInt(num),
          },
        ],
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
    },
    getDataByLabels(data) {
      let array = [];
      data.forEach((item) => {
        const percentage =
          (item.data.payed * 100) / (item.data.payed + item.data.expired ? item.data.expired : item.data.pending);
        array.push(percentage.toFixed(2));
      });
      return array;
    },
  },
};
</script>
<style lang="scss">
@import "../assets/css/dashboard.scss";
</style>
