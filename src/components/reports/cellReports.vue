<template>
  <div class="px-4 cell-reports">
    <header>Cell Report</header>
    <hr class="m-0 mt-1 mb-4" />
    <b-row class="controls m-0">
      <b-select
        id="input-1"
        v-model="cell"
        :options="cellOptions"
        class="w-100 mb-5"
        v-if="!isManager"
      >
        <template v-slot:first>
          <option :value="null" disabled>Please select cell</option>
        </template>
      </b-select>
      <selector
        :title="'Generate '+title+' Report'"
        :object="config"
        :disabled="cell?false:true"
        v-on:ok="generateAction"
      />
      <vue-load v-if="state.generating" class="w-100 bg-light mt-2" label="Generating..." />
    </b-row>
    <b-row class="m-0">
      <b-collapse id="sector-report-collapse" class="w-100 m-3" v-model="state.showReport">
        <b-card class="text-capitalize" v-if="!state.error">
          <b-card-title class="font-19 text-uppercase">{{cell}} cell</b-card-title>
          <hr />
          <b-table
            id="cell-reports"
            :items="generate"
            :fields="table.fields"
            :busy.sync="state.generating"
            :key="'cell-'+table.key"
            v-if="state.generate"
            small
            bordered
            responsive
            show-empty
          >
            <template v-slot:cell(unpayedAmount)="data">
              <b-card-text class="text-normal">{{Number(data.value).toLocaleString()}} Rwf</b-card-text>
            </template>
            <template v-slot:cell(payedAmount)="data">
              <b-card-text class="text-normal">{{Number(data.value).toLocaleString()}} Rwf</b-card-text>
            </template>
          </b-table>
          <b-card-title class="font-19 text-uppercase">villages</b-card-title>
          <hr />
          <b-table
            id="cell-village-reports"
            :items="generateVillage"
            :fields="villageTable.fields"
            :busy.sync="state.generating"
            :key="'village-'+villageTable.key"
            v-if="state.generate"
            small
            bordered
            responsive
            show-empty
          >
            <template v-slot:cell(unpayedAmount)="data">
              <b-card-text class="text-normal">{{Number(data.value).toLocaleString()}} Rwf</b-card-text>
            </template>
            <template v-slot:cell(payedAmount)="data">
              <b-card-text class="text-normal">{{Number(data.value).toLocaleString()}} Rwf</b-card-text>
            </template>
          </b-table>
        </b-card>
        <b-card v-if="state.error">
          <b-card-text>{{state.errorMessage}}</b-card-text>
        </b-card>
      </b-collapse>
    </b-row>
    <b-row v-if="!state.error && villageData && cellData" class="my-3 mr-1 justify-content-end">
      <b-button @click="downloadReport" size="sm" class="app-color">Download Report</b-button>
    </b-row>
  </div>
</template>

<script>
import download from "./downloadCellReport";
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
        generating: false,
        generate: false,
        showReport: false,
        error: false,
        errorMessage: null
      },
      config: {
        configuring: false,
        year: new Date().getFullYear(),
        month: new Date().getMonth() + 1
      },
      cellData: null,
      villageData: null,
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
        ],
        busy: false,
        key: 1
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
        ],
        busy: false,
        key: 1
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
      this.state.generate = true;
      this.table.key++;
      this.villageTable.key++;
    },
    generate() {
      this.cellData = null;
      this.state.generating = true;
      const year = this.config.year;
      const month = this.config.month;
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
          items.payedAmount = res[1].data.data.payed;
          items.unpayedAmount = res[1].data.data.pending;
          this.state.showReport = true;
          this.cellData = items;
          return [items];
        })
        .catch(err => {
          this.state.error = true;
          this.state.errorMessage = err.response.data.error
            ? err.response.data.error
            : err.response.response;
          this.cellData = null;
          if (this.villageData) {
            this.state.showReport = true;
          }
          return [];
        })
        .finally(() => {
          this.state.generating = false;
        });
    },
    generateVillage() {
      this.villageData = null;
      this.state.generating = true;
      const year = this.config.year;
      const month = this.config.month;
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
          this.state.showReport = true;
          this.villageData = items;
          return items;
        })
        .catch(err => {
          this.state.error = true;
          this.state.errorMessage = err.response.data.error
            ? err.response.data.error
            : err.response.response;
          this.villageData = null;
          if (this.cellData) {
            this.state.showReport = true;
          }
          return [];
        })
        .finally(() => {
          this.state.generating = false;
        });
    },
    downloadReport() {
      if (
        !this.state.generating &&
        this.cellData != null &&
        this.villageData != null
      ) {
        download(this.cellData, this.villageData, this.cell);
      }
    },
    clear() {
      this.state.showReport = false;
      this.state.generating = false;
      this.state.generate = false;
      this.state.error = false;
      this.state.errorMessage = null;
      this.cellData = null;
      this.villageData = null;
    }
  }
};
</script>

<style lang="scss">
.cell-reports {
  & > header {
    text-align: center;
    font-size: 1.3rem;
    font-weight: bold;
    color: #384950;
  }
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