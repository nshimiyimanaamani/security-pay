<template>
  <div class="sector-reports">
    <header class="tabTitle">Sector Report</header>
    <div class="tabBody">
      <b-row class="justify-content-center flex-column m-0">
        <selector
          :object="config"
          :title="'Generate '+activeSector+' Report'"
          v-on:ok="generateAction"
          class="w-auto m-auto date-selector pb-3"
        />

        <vue-load v-if="state.generating" label="Generating..." />
      </b-row>

      <b-row class="m-0" v-show="!state.generating">
        <b-collapse id="sectorreport-collapse" class="w-100" v-model="state.showReport">
          <div class="reports-card" v-if="!state.error">
            <b-row no-gutters class="mb-2 justify-content-end">
              <b-badge
                variant="secondary"
                class="p-2 font-13"
              >Report Date: &nbsp; {{state.reportsDate}}</b-badge>
            </b-row>
            <h5 class="bg-dark">Sector Reports</h5>
            <div class="card--body">
              <b-table
                id="sector-reports"
                :items="generate"
                :fields="table.fields"
                :busy="table.busy"
                :key="'sector-'+table.key"
                v-if="state.generate"
                head-row-variant="secondary"
                small
                bordered
                hover
                responsive
                show-empty
              >
                <template v-slot:cell(unpayedAmount)="data">
                  <b-card-text class="text-normal">{{data.value | number}} Rwf</b-card-text>
                </template>
                <template v-slot:cell(payedAmount)="data">
                  <b-card-text class="text-normal">{{data.value | number}} Rwf</b-card-text>
                </template>
              </b-table>
            </div>
          </div>
          <div class="reports-card" v-if="!state.error">
            <h5 class="bg-dark">cells Reports</h5>
            <div class="card--body">
              <b-table
                id="sector-cell-reports"
                :items="generateCell"
                :fields="cellTable.fields"
                :busy="cellTable.busy"
                :key="'cell-'+cellTable.key"
                v-if="state.generate"
                head-row-variant="secondary"
                small
                bordered
                hover
                responsive
                show-empty
              >
                <template v-slot:cell(unpayedAmount)="data">
                  <b-card-text class="text-normal">{{data.value | number}} Rwf</b-card-text>
                </template>
                <template v-slot:cell(payedAmount)="data">
                  <b-card-text class="text-normal">{{data.value | number}} Rwf</b-card-text>
                </template>
              </b-table>
            </div>
          </div>
          <b-card v-if="state.error">
            <b-card-text>{{state.errorMessage}}</b-card-text>
          </b-card>
        </b-collapse>
      </b-row>
      <b-row v-if="!state.error && sectorData && cellData " class="m-0 py-3 justify-content-end">
        <b-button @click="downloadReport" variant="info" class="downloadBtn">Download Report</b-button>
      </b-row>
    </div>
  </div>
</template>

<script>
import download from "./downloadSectorReport";
import selector from "../reportsDateSelector";
export default {
  name: "sectorReports",
  components: {
    selector
  },
  data() {
    return {
      state: {
        generating: false,
        generate: false,
        showReport: false,
        error: false,
        errorMessage: null,
        reportsDate: null
      },
      config: {
        configuring: false,
        year: new Date().getFullYear(),
        month: new Date().getMonth() + 1
      },
      sectorData: null,
      cellData: null,
      table: {
        fields: [
          {
            key: "total",
            label: "No of Houses",
            tdClass: "text-center",
            thClass: "text-center text-uppercase"
          },
          {
            key: "payed",
            label: "No of Paid Houses",
            tdClass: "text-center",
            thClass: "text-center text-uppercase"
          },
          {
            key: "payedAmount",
            label: "Paid Amount",
            tdClass: "text-right",
            thClass: "text-center text-uppercase"
          },
          {
            key: "pending",
            label: "No of unpaid Houses",
            tdClass: "text-center",
            thClass: "text-center text-uppercase"
          },
          {
            key: "unpayedAmount",
            label: "unpaid Amount",
            tdClass: "text-right",
            thClass: "text-center text-uppercase"
          }
        ],
        busy: false,
        key: 1
      },
      cellTable: {
        fields: [
          {
            key: "name",
            label: "Cell",
            thClass: " text-uppercase",
            tdClass: "font-weight-bold text-uppercase"
          },
          {
            key: "total",
            label: "No of Houses",
            tdClass: "text-center",
            thClass: "text-center text-uppercase"
          },
          {
            key: "payed",
            label: "No of Paid Houses",
            tdClass: "text-center",
            thClass: "text-center text-uppercase"
          },
          {
            key: "payedAmount",
            label: "Paid Amount",
            tdClass: "text-right",
            thClass: "text-center text-uppercase"
          },
          {
            key: "pending",
            label: "No of unpaid Houses",
            tdClass: "text-center",
            thClass: "text-center text-uppercase"
          },
          {
            key: "unpayedAmount",
            label: "unpaid Amount",
            tdClass: "text-right",
            thClass: "text-center text-uppercase"
          }
        ],
        busy: false,
        key: 1
      }
    };
  },
  computed: {
    activeSector() {
      return this.$store.getters.getActiveSector;
    },
    months() {
      return this.$store.getters.getMonths;
    }
  },
  methods: {
    generateAction() {
      this.clear();
      this.state.generate = true;
      this.table.key++;
      this.cellTable.key++;
    },
    generate() {
      this.sectorData = null;
      this.state.generating = true;
      const year = this.config.year;
      const month = this.config.month;
      this.state.reportsDate = `${this.months[month - 1]}, ${year}`;
      const first = this.axios.get(
        `/metrics/ratios/sectors/${this.activeSector}?year=${year}&month=${month}`
      );
      const second = this.axios.get(
        `/metrics/balance/sectors/${this.activeSector}?year=${year}&month=${month}`
      );
      return Promise.all([first, second])
        .then(res => {
          const items = {};
          items.total = res[0].data.data.payed + res[0].data.data.pending;
          items.payed = res[0].data.data.payed;
          items.pending = res[0].data.data.pending;
          items.payedAmount = res[1].data.data.payed;
          items.unpayedAmount = res[1].data.data.pending;
          this.state.showReport = true;
          this.sectorData = items;
          return [items];
        })
        .catch(err => {
          this.state.error = true;
          this.state.errorMessage = err.response.data.error
            ? err.response.data.error
            : err.response.response;
          this.sectorData = null;
          if (this.cellData) {
            this.state.showReport = true;
          }
          return [];
        })
        .finally(() => {
          if (this.cellData) {
            this.state.generating = false;
          }
        });
    },
    generateCell() {
      this.cellData = null;
      this.state.generating = true;
      const year = this.config.year;
      const month = this.config.month;
      this.state.reportsDate = `${this.months[month - 1]}, ${year}`;
      const first = this.axios.get(
        `/metrics/ratios/sectors/all/${this.activeSector}?year=${year}&month=${month}`
      );
      const second = this.axios.get(
        `/metrics/balance/sectors/all/${this.activeSector}?year=${year}&month=${month}`
      );
      return Promise.all([first, second])
        .then(res => {
          var items = [];
          res[0].data.forEach(item => {
            res[1].data.forEach(element => {
              if (element.label == item.label) {
                items.push({
                  name: item.label || element.label,
                  total: item.data.payed + item.data.pending,
                  payed: item.data.payed,
                  pending: item.data.pending,
                  unpayedAmount: element.data.pending,
                  payedAmount: element.data.payed
                });
              }
            });
          });
          this.state.showReport = true;
          this.cellData = items;
          return items;
        })
        .catch(err => {
          this.state.error = true;
          this.state.errorMessage = err.response.data.error
            ? err.response.data.error
            : err.response.response;
          this.cellData = null;
          if (this.sectorData) {
            this.state.showReport = true;
          }
          return [];
        })
        .finally(() => {
          if (this.sectorData) {
            this.state.generating = false;
          }
        });
    },
    downloadReport() {
      if (
        !this.state.generating &&
        this.sectorData != null &&
        this.cellData != null
      ) {
        download(
          this.sectorData,
          this.cellData,
          this.activeSector,
          this.state.reportsDate
        );
      }
    },
    clear() {
      this.state.showReport = false;
      this.state.generate = false;
      this.state.generating = false;
      this.state.error = false;
      this.state.errorMessage = null;
      this.cellData = null;
      this.cellData = null;
      this.state.reportsDate = null;
    }
  }
};
</script>

<style lang="scss">
.sector-reports {
  .date-selector {
    & > button {
      border-radius: 3px;
      padding: 0.7rem 1rem;
    }
  }
  .reports-card {
    margin: 1rem 0;
    h4 {
      font-size: 1.2rem;
      margin: 0;
      text-align: center;
      padding: 0.8rem 1rem;
    }
    hr {
      margin-top: 0;
      margin-bottom: 1rem;
    }
    .table-responsive > table {
      min-width: max-content;
    }
  }
  .downloadBtn {
    border-radius: 2px;
  }
}
</style>
