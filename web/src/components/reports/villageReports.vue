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
            <h5 class="bg-dark text-uppercase">{{village || ''}} village</h5>
            <b-table
              id="village-reports"
              :items="villageData"
              :fields="table.fields"
              :busy.sync="state.busy"
              small
              bordered
              responsive
              show-empty
            >
              <template v-slot:table-busy>
                <vue-load label="Generating..." class="p-3" />
              </template>
              <template v-slot:empty>{{state.error || 'No data available to display'}}</template>
            </b-table>
          </div>
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
import download from "../download scripts/downloadReports";
import selector from "../reportsDateSelector";
export default {
  name: "VillageReports",
  components: { selector },
  data() {
    return {
      cell: null,
      village: null,
      state: {
        showReport: false,
        reportsDate: null,
        busy: false,
        error: null
      },
      config: {
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
        ]
      },
      villageData: []
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
      if (
        this.state.error ||
        this.state.busy ||
        this.villageData.length < 1 ||
        this.state.showReport === false
      )
        return false;
      return true;
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
      this.generateVillage();
    },
    generateVillage() {
      this.villageData = [];
      this.state.showReport = true;
      this.state.busy = true;
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
          items.payedAmount = `${Number(
            res[1].data.data.payed
          ).toLocaleString()} Rwf`;
          items.unpayedAmount = `${Number(
            res[1].data.data.pending
          ).toLocaleString()} Rwf`;
          this.villageData = [items];
          this.state.busy = false;
          return [items];
        })
        .catch(err => {
          try {
            this.state.errorMessage =
              err.response.data.error || err.response.response;
          } catch {
            this.state.error = "Failed to retrieve village data!";
          }
          this.state.busy = false;
          return [];
        });
    },
    downloadReport() {
      if (this.villageData.length > 0) {
        const data = {
          config: {
            TITLE: String(`Monthly Report of ${this.village}`).toUpperCase(),
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
              BODY: this.villageData
            }
          ]
        };
        download(data);
      }
    },
    clear() {
      this.state.showReport = false;
      this.state.error = null;
      this.state.reportsDate = null;
      this.state.busy = false;
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