<template>
  <div class="px-4">
    <header class="d-flex justify-content-center font-19 text-uppercase">Cell Report</header>
    <hr class="m-0 mb-3" />
    <b-row class="px-3 align-items-center justify-content-between">
      <b-select
        size="sm"
        id="input-1"
        v-model="cell"
        :options="cellOptions"
        class="w-auto mr-3 flex-grow-1"
        v-if="user.role.toLowerCase() !='basic'"
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
      <div v-show="state.generating" class="w-100">
        <strong class="font-14">Generating&nbsp;</strong>
        <b-spinner small />
      </div>
    </b-row>
    <b-row>
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
            label: "No of Payed Houses",
            tdClass: "text-center",
            thClass: "text-center text-uppercase"
          },
          {
            key: "payedAmount",
            label: "Payed Amount",
            tdClass: "text-right",
            thClass: "text-center text-uppercase"
          },
          {
            key: "pending",
            label: "No of unpayed Houses",
            tdClass: "text-center",
            thClass: "text-center text-uppercase"
          },
          {
            key: "unpayedAmount",
            label: "unPayed Amount",
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
            label: "No of Payed Houses",
            tdClass: "text-center",
            thClass: "text-center text-uppercase"
          },
          {
            key: "payedAmount",
            label: "Payed Amount",
            tdClass: "text-right",
            thClass: "text-center text-uppercase"
          },
          {
            key: "pending",
            label: "No of unpayed Houses",
            tdClass: "text-center",
            thClass: "text-center text-uppercase"
          },
          {
            key: "unpayedAmount",
            label: "unPayed Amount",
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
      return this.$store.getters.getCellsArray;
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
    if (this.user.role.toLowerCase() == "basic") {
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
      const promise = this.axios.all([first, second]);
      return promise
        .then(
          this.axios.spread((...res) => {
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
        )
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
      const promise = this.axios.all([first, second]);
      return promise
        .then(
          this.axios.spread((...res) => {
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
        )
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

<style>
</style>