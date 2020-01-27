<template>
  <div>
    <header class="d-flex justify-content-center font-20 text-uppercase">village Report</header>
    <hr class="m-0 mb-3" />
    <b-row class="px-3 align-items-center justify-content-between">
      <b-select
        size="sm"
        id="input-1"
        v-model="cell"
        :options="cellOptions"
        class="w-auto mr-2 flex-grow-1"
        v-if="user.role.toLowerCase() !='basic'"
      >
        <template v-slot:first>
          <option :value="null" disabled>Please select cell</option>
        </template>
      </b-select>
      <b-select
        size="sm"
        id="input-1"
        v-model="village"
        :options="villageOptions"
        class="w-auto mr-2 flex-grow-1"
      >
        <template v-slot:first>
          <option :value="null" disabled>Please select village</option>
        </template>
      </b-select>
      <b-button
        size="sm"
        variant="info"
        class="font-15 border-0 my-3"
        :disabled="village?false:true"
        @click="generateAction"
      >Generate {{village ? village : 'Village'}} Report</b-button>
      <div v-show="state.generating" class="w-100">
        <strong class="font-15">Generating&nbsp;</strong>
        <b-spinner small />
      </div>
    </b-row>
    <b-row>
      <b-collapse id="sector-report-collapse" class="w-100 m-3" v-model="state.showReport">
        <b-card class="text-capitalize" v-if="!state.error">
          <b-card-title class="font-20 text-uppercase">{{village}} village</b-card-title>
          <hr />
          <b-table
            id="village-reports"
            :items="generateVillage"
            :fields="table.fields"
            :busy.sync="state.generating"
            :key="'village-'+table.key"
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
    <b-row v-if="!state.error && downloadData" class="my-3 mr-1 justify-content-end">
      <b-button size="sm" class="app-color" @click="downloadReport">Download Report</b-button>
    </b-row>
  </div>
</template>

<script>
const { Village } = require("rwanda");
import download from "./downloadVillageReport";
export default {
  name: "VillageReports",
  data() {
    return {
      cell: null,
      village: null,
      state: {
        generating: false,
        showReport: false,
        generate: false,
        error: false,
        errorMessage: null
      },
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
      downloadData: null
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    activeSector() {
      return this.$store.getters.getActiveSector;
    },
    cellOptions() {
      return this.$store.getters.getCellsArray;
    },
    villageOptions() {
      if (this.cell) {
        return Village("Kigali", "Gasabo", this.activeSector, this.cell).sort();
      } else {
        return [];
      }
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
    }
  },
  watch: {
    village() {
      handler: {
        this.state.showReport = false;
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
    },
    generateVillage() {
      this.downloadData = null;
      this.state.generating = true;
      const first = this.axios.get(
        this.endpoint +
          `/metrics/ratios/villages/${this.village}?year=${this.currentYear}&month=${this.currentMonth}`
      );
      const second = this.axios.get(
        this.endpoint +
          `/metrics/balance/villages/${this.village}?year=${this.currentYear}&month=${this.currentMonth}`
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
            this.downloadData = items;
            return [items];
          })
        )
        .catch(err => {
          this.state.error = true;
          this.state.errorMessage = err.response.data.error
            ? err.response.data.error
            : err.response.response;
          this.downloadData = null;
          return [];
        })
        .finally(() => {
          this.state.generating = false;
        });
    },
    downloadReport() {
      if (!this.state.generating && this.downloadData != null) {
        download(this.downloadData, this.village);
      }
    },
    clear() {
      this.state.showReport = false;
      this.state.generating = false;
      this.state.generate = false;
      this.state.error = false;
      this.state.errorMessage = null;
      this.downloadData = null;
    }
  }
};
</script>

<style>
</style>