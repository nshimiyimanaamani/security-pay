<template>
  <div class="px-4">
    <header class="d-flex justify-content-center font-19 text-uppercase">Sector Report</header>
    <hr class="m-0 mb-3" />
    <b-row class="justify-content-between">
      <selector
        :object="config"
        :title="'Generate '+activeSector+' Report'"
        v-on:ok="generateAction"
      />
      <div v-show="state.generating" class="w-auto px-3 align-self-center">
        <strong class="font-14">Generating&nbsp;</strong>
        <b-spinner small />
      </div>
    </b-row>
    <b-row>
      <b-collapse id="sectorreport-collapse" class="w-100" v-model="state.showReport">
        <b-card class="text-capitalize" v-if="!state.error">
          <b-card-title class="font-19 text-uppercase">Sector</b-card-title>
          <hr />
          <b-table
            id="sector-reports"
            :items="generate"
            :fields="table.fields"
            :busy="table.busy"
            :key="'sector-'+table.key"
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
          <b-card-title class="font-19 text-uppercase">cells</b-card-title>
          <hr />
          <b-table
            id="sector-cell-reports"
            :items="generateCell"
            :fields="cellTable.fields"
            :busy="cellTable.busy"
            :key="'cell-'+cellTable.key"
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
    <b-row v-if="!state.error && sectorData && cellData " class="my-3 justify-content-end">
      <b-button @click="downloadReport" size="sm" class="app-color">Download Report</b-button>
    </b-row>
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
        errorMessage: null
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
    
    activeSector() {
      return this.$store.getters.getActiveSector;
    },
    currentYear() {
      return new Date().getFullYear();
    },
    currentMonth() {
      return new Date().getMonth() + 1;
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
      const first = this.axios.get(
        
          `/metrics/ratios/sectors/${this.activeSector}?year=${year}&month=${month}`
      );
      const second = this.axios.get(
        
          `/metrics/balance/sectors/${this.activeSector}?year=${year}&month=${month}`
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
            this.sectorData = items;
            return [items];
          })
        )
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
      const first = this.axios.get(
        
          `/metrics/ratios/sectors/all/${this.activeSector}?year=${year}&month=${month}`
      );
      const second = this.axios.get(
        
          `/metrics/balance/sectors/all/${this.activeSector}?year=${year}&month=${month}`
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
            this.cellData = items;
            return items;
          })
        )
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
        download(this.sectorData, this.cellData, this.activeSector);
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
    }
  }
};
</script>

<style>
</style>