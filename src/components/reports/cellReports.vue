<template>
  <div class="cell-reports">
    <header class="tabTitle">Cell Report</header>
    <div class="tabBody">
      <b-row class="controls m-0">
        <b-select
          id="input-1"
          v-model="cell"
          :options="cellOptions"
          class="w-100 mb-4"
          v-if="!isManager"
        >
          <template v-slot:first>
            <option :value="null" disabled>select cell</option>
          </template>
        </b-select>
        <selector
          :title="`Generate ${title || ''} Report`"
          :object="config"
          :disabled="cell?false:true"
          v-on:ok="generateAction"
          class="mb-3"
        />
      </b-row>
      <b-row no-gutters>
        <b-collapse id="sector-report-collapse" class="w-100" v-model="state.showReport">
          <div class="reports-card">
            <b-row no-gutters class="mb-2 justify-content-end">
              <b-badge
                variant="secondary"
                class="p-2 fsize-sm"
              >Report Date: &nbsp; {{state.reportsDate}}</b-badge>
            </b-row>
            <h5 class="text-uppercase bg-dark">{{cell}} cell</h5>
            <b-table
              id="cell-reports"
              :items="cellData"
              :fields="table.fields"
              :busy.sync="state.busy.table1"
              head-variant="secondary"
              small
              bordered
              responsive
              show-empty
            >
              <template v-slot:table-busy>
                <vue-load label="Generating..." class="p-3" />
              </template>
              <template v-slot:empty>{{state.error.table1 || 'No data available to display'}}</template>
            </b-table>
          </div>
          <div class="reports-card">
            <h5 class="text-uppercase bg-dark">villages</h5>
            <b-table
              id="cell-village-reports"
              :items="villageData"
              :fields="villageTable.fields"
              :busy.sync="state.busy.table2"
              head-variant="secondary"
              small
              bordered
              responsive
              show-empty
            >
              <template v-slot:table-busy>
                <vue-load label="Generating..." class="p-3" />
              </template>
              <template v-slot:empty>{{state.error.table2 || 'No data available to display'}}</template>
            </b-table>
          </div>
        </b-collapse>
      </b-row>
      <b-row v-if="canDownload" class="m-0 justify-content-end">
        <b-button @click="downloadReport" variant="info" class="br-2">
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
  name: "cellReports",
  components: {
    selector
  },
  data() {
    return {
      cell: null,
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
      cellData: [],
      villageData: [],
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
            label: "unPaid Amount",
            tdClass: "text-right",
            thClass: "text-center text-uppercase"
          }
        ]
      },
      villageTable: {
        fields: [
          {
            key: "name",
            label: "Village",
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
            label: "No of unPaid Houses",
            tdClass: "text-center",
            thClass: "text-center text-uppercase"
          },
          {
            key: "unpayedAmount",
            label: "unPaid Amount",
            tdClass: "text-right",
            thClass: "text-center text-uppercase"
          }
        ]
      }
    };
  },
  computed: {
    cellOptions() {
      const { province, district, sector } = this.location;
      return this.$cells(province, district, sector);
    },
    currentYear() {
      return new Date().getFullYear();
    },
    currentMonth() {
      return new Date().getMonth() + 1;
    },
    user() {
      return this.$store.getters.userDetails;
    },
    activeCell() {
      return this.$store.getters.getActiveCell;
    },
    title() {
      return this.cell ? this.cell : "Cell";
    },
    location() {
      return this.$store.getters.location;
    },
    isManager() {
      return this.user.role.toLowerCase() === "basic";
    },
    canDownload() {
      if (
        this.state.error.table1 ||
        this.state.error.table2 ||
        this.state.busy.table1 ||
        this.state.busy.table2 ||
        !this.state.showReport ||
        !this.villageData ||
        !this.cellData
      )
        return false;
      return true;
    },
    months() {
      return this.$store.getters.getMonths;
    }
  },
  watch: {
    cell() {
      handler: {
        this.clear();
      }
    }
  },
  mounted() {
    if (this.isManager) {
      this.cell = this.activeCell;
    }
  },
  methods: {
    generateAction() {
      this.clear();
      this.generateCell();
      this.generateVillage();
      this.state.showReport = true;
    },
    generateCell() {
      this.cellData = [];
      this.state.busy.table1 = true;
      const year = this.config.year;
      const month = this.config.month;
      this.state.reportsDate = `${this.months[month - 1]}, ${year}`;
      const first = this.axios.get(
        `/metrics/ratios/cells/${this.cell}?year=${year}&month=${month}`
      );
      const second = this.axios.get(
        `/metrics/balance/cells/${this.cell}?year=${year}&month=${month}`
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
          this.cellData = [items];
          this.state.busy.table1 = false;
          return [items];
        })
        .catch(err => {
          this.state.busy.table1 = false;
          this.state.error = true;
          try {
            this.state.error.table1 =
              err.response.data.error || err.response.response;
          } catch {
            this.state.error.table1 = "Failed to retrieve cell data!";
          }
          return [];
        });
    },
    generateVillage() {
      this.villageData = [];
      this.state.busy.table2 = true;
      const year = this.config.year;
      const month = this.config.month;
      this.state.reportsDate = `${this.months[month - 1]}, ${year}`;
      const first = this.axios.get(
        `/metrics/ratios/cells/all/${this.cell}?year=${year}&month=${month}`
      );
      const second = this.axios.get(
        `/metrics/balance/cells/all/${this.cell}?year=${year}&month=${month}`
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
          this.villageData = items;
          this.state.busy.table2 = false;
          return items;
        })
        .catch(err => {
          this.state.busy.table2 = false;
          try {
            this.state.errorMessage =
              err.response.data.error || err.response.response;
          } catch {
            this.state.error.table2 = "Failed to retrieve village data!";
          }
          return [];
        });
    },
    downloadReport() {
      if (this.cellData.length > 0 && this.villageData.length > 0) {
        const data = {
          config: {
            TITLE: String(`Monthly Report of ${this.cell}`).toUpperCase(),
            name: `${this.cell} Monthly Report of ${this.state.reportsDate}`,
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
              BODY: this.cellData
            },
            {
              COLUMNS: [
                {
                  header: `Village`,
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
              BODY: this.villageData
            }
          ]
        };
        download(data);
      }
    },
    clear() {
      this.state.showReport = false;
      this.state.reportsDate = null;
      this.cellData = [];
      this.villageData = [];
      this.state.error.table1 = null;
      this.state.error.table2 = null;
      this.state.busy.table1 = false;
      this.state.busy.table2 = false;
    }
  }
};
</script>

<style lang="scss">
.cell-reports {
  .controls {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-around;
  }
  .dropdown > button {
    padding: 0.7rem 1.5rem;
  }
}
</style>