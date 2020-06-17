<template>
  <div class="sector-reports">
    <header class="tabTitle">Sector Report</header>
    <div class="tabBody">
      <b-row class="justify-content-center flex-column" no-gutters>
        <selector
          :object="config"
          :title="'Generate '+activeSector+' Report'"
          v-on:ok="generateAction"
          class="w-auto m-auto date-selector pb-3"
        />
      </b-row>

      <b-row no-gutters>
        <b-collapse id="sectorreport-collapse" class="w-100" v-model="state.showReport">
          <div class="reports-card">
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
                :items="sectorData"
                :fields="table.fields"
                :busy="state.busy.table1"
                head-row-variant="secondary"
                small
                bordered
                hover
                responsive
                show-empty
              >
                <template v-slot:table-busy>
                  <vue-load label="Generating..." class="p-3" />
                </template>
                <template v-slot:empty>{{state.error.table1 || 'No data available to display'}}</template>
              </b-table>
            </div>
          </div>
          <div class="reports-card">
            <h5 class="bg-dark">cells Reports</h5>
            <div class="card--body">
              <b-table
                id="sector-cell-reports"
                :items="cellData"
                :fields="cellTable.fields"
                :busy="state.busy.table2"
                head-row-variant="secondary"
                small
                bordered
                hover
                responsive
                show-empty
              >
                <template v-slot:table-busy>
                  <vue-load label="Generating..." class="p-3" />
                </template>
                <template v-slot:empty>{{state.error.table1 || 'No data available to display'}}</template>
              </b-table>
            </div>
          </div>
        </b-collapse>
      </b-row>
      <b-row v-if="showDownload" class="py-3 justify-content-end" no-gutters>
        <b-button @click="downloadReport" variant="info" class="downloadBtn">
          <i class="fa fa-download mr-1" />Download Report
        </b-button>
      </b-row>
    </div>
  </div>
</template>

<script>
import download from "../download scripts/downloadReports";
import selector from "../reportsDateSelector";
export default {
  name: "sectorReports",
  components: {
    selector
  },
  data() {
    return {
      state: {
        showReport: false,
        reportsDate: null,
        busy: {
          table1: false,
          table2: false
        },
        error: {
          table1: null,
          table2: null
        }
      },
      config: {
        year: new Date().getFullYear(),
        month: new Date().getMonth() + 1
      },
      sectorData: [],
      cellData: [],
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
        ]
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
        ]
      }
    };
  },
  computed: {
    activeSector() {
      return this.$store.getters.getActiveSector;
    },
    months() {
      return this.$store.getters.getMonths;
    },
    showDownload() {
      if (
        this.state.error.table1 ||
        this.state.error.table2 ||
        this.state.busy.table1 ||
        this.state.busy.table2 ||
        !this.sectorData ||
        !this.cellData ||
        !this.state.showReport
      )
        return false;
      return true;
    }
  },
  methods: {
    generateAction() {
      this.clear();
      this.generateSector();
      this.generateCell();
    },
    generateSector() {
      this.sectorData = [];
      this.state.showReport = true;
      this.state.busy.table1 = true;
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
          items.payedAmount = `${Number(
            res[1].data.data.payed
          ).toLocaleString()} Rwf`;
          items.unpayedAmount = `${Number(
            res[1].data.data.pending
          ).toLocaleString()} Rwf`;
          this.state.showReport = true;
          this.sectorData = [items];
          this.state.busy.table1 = false;

          return [items];
        })
        .catch(err => {
          this.state.busy.table1 = false;
          try {
            this.state.error.table1 =
              err.response.data.error || err.response.response;
          } catch {
            this.state.error.table1 = "Failed to retrieve sector report data";
          }
          return [];
        });
    },
    generateCell() {
      this.cellData = [];
      this.state.showReport = true;
      this.state.busy.table2 = true;
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
                  unpayedAmount: `${Number(
                    element.data.pending
                  ).toLocaleString()} Rwf`,
                  payedAmount: `${Number(
                    element.data.payed
                  ).toLocaleString()} Rwf`
                });
              }
            });
          });
          this.state.busy.table2 = false;
          this.cellData = items;
          return items;
        })
        .catch(err => {
          this.state.busy.table2 = false;
          try {
            this.state.error.table2 =
              err.response.data.error || err.response.response;
          } catch {
            this.state.error.table2 = "Failed to retrieve Cells report data";
          }
          return [];
        });
    },
    downloadReport() {
      if (this.sectorData.length > 0 && this.cellData.length > 0) {
        const data = {
          config: {
            TITLE: String(
              `Monthly Report of ${this.activeSector}`
            ).toUpperCase(),
            name: `${this.activeSector} Monthly Report of ${this.state.reportsDate}`,
            date: this.state.reportsDate
          },
          data: [
            {
              COLUMNS: [
                {
                  header: `No of Properties`,
                  dataKey: "total"
                },
                {
                  header: `No of Paid Properties`,
                  dataKey: "payed"
                },
                { header: `Paid Amount`, dataKey: "payedAmount" },
                {
                  header: `No of Unpaid Properties`,
                  dataKey: "pending"
                },
                { header: `Unpaid Amount`, dataKey: "unpayedAmount" }
              ],
              BODY: this.sectorData
            },
            {
              COLUMNS: [
                {
                  header: `Cell`,
                  dataKey: "name"
                },
                {
                  header: `No of Properties`,
                  dataKey: "total"
                },
                {
                  header: `No of Paid Properties`,
                  dataKey: "payedAmount"
                },
                { header: `Paid Amount`, dataKey: "payedAmount" },
                {
                  header: `No of Unpaid Properties`,
                  dataKey: "pending"
                },
                { header: `Unpaid Amount`, dataKey: "unpayedAmount" }
              ],
              BODY: this.cellData
            }
          ]
        };
        download(data);
      }
    },
    clear() {
      this.state.showReport = false;
      this.state.reportsDate = null;
      this.state.error.table1 = null;
      this.state.error.table2 = null;
      this.state.busy.table1 = false;
      this.state.busy.table2 = false;
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
