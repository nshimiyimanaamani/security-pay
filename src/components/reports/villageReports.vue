<template>
  <div id="village-reports">
    <header class="tabTitle">village Report</header>
    <div class="tabBody">
      <b-row class="controls mb-3" no-gutters>
        <b-select id="input-1" v-model="cell" :options="cellOptions" v-if="!isManager" class="br-2">
          <template v-slot:first>
            <option :value="null" disabled>Please select cell</option>
          </template>
        </b-select>
        <b-select id="input-1" v-model="village" :options="villageOptions" class="br-2">
          <template v-slot:first>
            <option :value="null" disabled>Please select village</option>
          </template>
        </b-select>
        <selector
          :disabled="village?false:true"
          :object="config"
          :title="`Generate ${title || ''} Report`"
          v-on:ok="generateAction"
          class="date-selector pb-3"
        />
        <vue-load v-if="state.generating" label="Generating..." />
      </b-row>
      <b-row no-gutters v-show="!state.generating">
        <b-collapse id="sector-report-collapse" class="w-100" v-model="state.showReport">
          <div class="reports-card" v-if="!state.error">
            <b-row no-gutters class="mb-2 justify-content-end">
              <b-badge
                variant="secondary"
                class="p-2 font-13"
              >Report Date: &nbsp; {{state.reportsDate}}</b-badge>
            </b-row>
            <h5 class="bg-dark text-uppercase">{{village || ''}} village</h5>
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
                <b-card-text class="text-normal">{{data.value | number}} Rwf</b-card-text>
              </template>
              <template v-slot:cell(payedAmount)="data">
                <b-card-text class="text-normal">{{data.value | number}} Rwf</b-card-text>
              </template>
            </b-table>
          </div>
          <b-card v-if="state.error">
            <b-card-text>{{state.errorMessage}}</b-card-text>
          </b-card>
        </b-collapse>

        <b-row v-if="canDownload" class="justify-content-end w-100 mt-3" no-gutters>
          <b-button variant="info" class="br-2" @click="downloadReport">
            <i class="fa fa-download mr-1" />Download Report
          </b-button>
        </b-row>
      </b-row>
    </div>
  </div>
</template>

<script>
import download from "./downloadVillageReport";
import selector from "../reportsDateSelector";
export default {
  name: "VillageReports",
  components: { selector },
  data() {
    return {
      cell: null,
      village: null,
      state: {
        generating: false,
        showReport: false,
        generate: false,
        error: false,
        errorMessage: null,
        reportsDate: null
      },
      config: {
        configuring: false,
        year: new Date().getFullYear(),
        month: new Date().getMonth() + 1
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
            label: "No of unPayed Houses",
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
      downloadData: null
    };
  },
  computed: {
    activeSector() {
      return this.$store.getters.getActiveSector;
    },
    cellOptions() {
      const { province, district, sector } = this.location;
      return this.$cells(province, district, this.activeSector);
    },
    villageOptions() {
      const { province, district, sector } = this.location;
      if (this.cell)
        return this.$villages(
          province,
          district,
          this.activeSector,
          this.cell
        ).sort();

      return [];
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
      return this.village ? this.village : "Village";
    },
    isManager() {
      return this.user.role.toLowerCase() === "basic";
    },
    location() {
      return this.$store.getters.location;
    },
    canDownload() {
      if (!this.state.error && this.downloadData) return true;
      return false;
    },
    months() {
      return this.$store.getters.getMonths;
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
    if (this.isManager) {
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
      const year = this.config.year;
      const month = this.config.month;
      this.state.reportsDate = `${this.months[month - 1]}, ${year}`;
      const first = this.axios.get(
        `/metrics/ratios/villages/${this.village}?year=${year}&month=${month}`
      );
      const second = this.axios.get(
        `/metrics/balance/villages/${this.village}?year=${year}&month=${month}`
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
          this.downloadData = items;
          return [items];
        })
        .catch(err => {
          this.state.error = true;
          this.state.showReport = true;
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
        download(this.downloadData, this.village, this.state.reportsDate);
      }
    },
    clear() {
      this.state.showReport = false;
      this.state.generating = false;
      this.state.generate = false;
      this.state.error = false;
      this.state.errorMessage = null;
      this.downloadData = null;
      this.state.reportsDate = null;
    }
  }
};
</script>

<style lang="scss">
#village-reports {
  .controls {
    display: flex;
    flex-direction: column;

    & > select {
      max-width: 500px;
      margin: 0 auto 1.5rem;
    }

    .date-selector {
      & > button {
        max-width: 500px;
        margin: auto;
        border-radius: 2px;
      }
      .dropdown-menu {
        max-width: 250px;
      }
    }
  }
}
</style>