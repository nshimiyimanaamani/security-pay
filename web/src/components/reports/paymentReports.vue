<template>
  <div id="payment-reports">
    <header class="tabTitle">Payment Reports</header>
    <div class="tabBody">
      <b-row class="m-0 buttons">
        <div>
          <b-dropdown
            v-model="dropdownone"
            text="Generate All House Report"
            ref="dropdown"
            class="m-2"
            variant="info"
            :busy="isLoading1"
          >
            <b-dropdown-form style="width: 230px">
              <b-form class="accountForm">
                <b-form-group label="Sector:">
                  <b-form-select
                    class="br-2"
                    v-model="form.select.sector"
                    :options="sectorOptions"
                    required
                  >
                    <template v-slot:first>
                      <option :value="null" disabled>select sector</option>
                    </template>
                  </b-form-select>
                </b-form-group>

                <b-form-group label="Cell:">
                  <b-form-select
                    class="br-2"
                    v-model="form.select.cell"
                    :options="cellOptions"
                    required
                  >
                    <template v-slot:first>
                      <option :value="null" disabled>select cell</option>
                    </template>
                  </b-form-select>
                </b-form-group>

                <b-form-group label="Village:">
                  <b-select
                    v-model="form.select.village"
                    :options="villageOptions"
                    class="br-2"
                    required
                  >
                    <template v-slot:first>
                      <option :value="null" disabled>select village</option>
                    </template>
                  </b-select>
                </b-form-group>
                <!-- <b-form-group label="From Month">
                  <div class="input-date">
                    <input type="date" v-model="object.frommonth" />
                  </div>
                </b-form-group> -->
                <!-- <b-form-group label="To Month">
                  <div class="input-date">
                    <input type="date" v-model="object.tomonth" />
                  </div>
                </b-form-group> -->
                <b-form-group
                  label="Year"
                  :label-for="'dropdown-year_' + random"
                >
                  <b-form-select
                    :id="'dropdown-year_' + random"
                    v-model="object.year"
                    class="bg-light"
                    size="sm"
                  >
                    <option
                      v-for="(year, i) in currentYear - 2019"
                      :value="currentYear - i"
                      :key="`year` + year"
                    >
                      {{ currentYear - i }}
                    </option>
                  </b-form-select>
                </b-form-group>

                <b-form-group
                  label="Month"
                  :label-for="'dropdown-month_' + random"
                >
                  <b-form-select
                    :id="'dropdown-month' + random"
                    v-model="object.month"
                    class="bg-light"
                    size="sm"
                  >
                    <option v-for="i in 12" :value="i" :key="`month` + i">
                      {{ months[i - 1] }}
                    </option>
                  </b-form-select>
                </b-form-group>
              </b-form>
            </b-dropdown-form>
            <b-dropdown-divider></b-dropdown-divider>
            <b-dropdown-item no-hover no-active>
              <b-button variant="info" block @click="getAllHouse"
                >Generate</b-button
              >
            </b-dropdown-item>
          </b-dropdown>
          <!-- <b-button variant="info" @click="getAllHouse">Generate All House Report</b-button> -->
        </div>
        <!-- <div>
          <b-button variant="info">Generate Paid House Report</b-button>
        </div> -->
        <div>
          <b-dropdown
            v-model="dropdownone"
            text="Generate Paid House Report"
            ref="dropdown"
            class="m-2"
            variant="info"
          >
            <b-dropdown-form style="width: 248px">
              <b-form class="accountForm">
                <b-form-group label="Sector:">
                  <b-form-select
                    class="br-2"
                    v-model="form.select.sector"
                    :options="sectorOptions"
                    required
                  >
                    <template v-slot:first>
                      <option :value="null" disabled>select sector</option>
                    </template>
                  </b-form-select>
                </b-form-group>

                <b-form-group label="Cell:">
                  <b-form-select
                    class="br-2"
                    v-model="form.select.cell"
                    :options="cellOptions"
                    required
                  >
                    <template v-slot:first>
                      <option :value="null" disabled>select cell</option>
                    </template>
                  </b-form-select>
                </b-form-group>

                <b-form-group label="Village:">
                  <b-select
                    v-model="form.select.village"
                    :options="villageOptions"
                    class="br-2"
                    required
                  >
                    <template v-slot:first>
                      <option :value="null" disabled>select village</option>
                    </template>
                  </b-select>
                </b-form-group>
                <b-form-group
                  label="Year"
                  :label-for="'dropdown-year_' + random"
                >
                  <b-form-select
                    :id="'dropdown-year_' + random"
                    v-model="object.year"
                    class="bg-light"
                    size="sm"
                  >
                    <option
                      v-for="(year, i) in currentYear - 2019"
                      :value="currentYear - i"
                      :key="`year` + year"
                    >
                      {{ currentYear - i }}
                    </option>
                  </b-form-select>
                </b-form-group>

                <b-form-group
                  label="Month"
                  :label-for="'dropdown-month_' + random"
                >
                  <b-form-select
                    :id="'dropdown-month' + random"
                    v-model="object.month"
                    class="bg-light"
                    size="sm"
                  >
                    <option v-for="i in 12" :value="i" :key="`month` + i">
                      {{ months[i - 1] }}
                    </option>
                  </b-form-select>
                </b-form-group>
              </b-form>
            </b-dropdown-form>
            <b-dropdown-divider></b-dropdown-divider>
            <b-dropdown-item no-hover no-active>
              <b-button variant="info" block @click="getPaidHouse"
                >Generate</b-button
              >
            </b-dropdown-item>
          </b-dropdown>
          <!-- <b-button variant="info" @click="getAllHouse">Generate All House Report</b-button> -->
        </div>
        <!-- <div>
          <b-button variant="info">Generate Unpaid House Report</b-button>
        </div> -->
        <div>
          <b-dropdown
            v-model="dropdownone"
            text="Generate Unpaid House Report"
            ref="dropdown"
            class="m-2"
            variant="info"
          >
            <b-dropdown-form style="width: 270px">
              <b-form class="accountForm">
                <b-form-group label="Sector:">
                  <b-form-select
                    class="br-2"
                    v-model="form.select.sector"
                    :options="sectorOptions"
                    required
                  >
                    <template v-slot:first>
                      <option :value="null" disabled>select sector</option>
                    </template>
                  </b-form-select>
                </b-form-group>

                <b-form-group label="Cell:">
                  <b-form-select
                    class="br-2"
                    v-model="form.select.cell"
                    :options="cellOptions"
                    required
                  >
                    <template v-slot:first>
                      <option :value="null" disabled>select cell</option>
                    </template>
                  </b-form-select>
                </b-form-group>

                <b-form-group label="Village:">
                  <b-select
                    v-model="form.select.village"
                    :options="villageOptions"
                    class="br-2"
                    required
                  >
                    <template v-slot:first>
                      <option :value="null" disabled>select village</option>
                    </template>
                  </b-select>
                </b-form-group>
                <b-form-group
                  label="Year"
                  :label-for="'dropdown-year_' + random"
                >
                  <b-form-select
                    :id="'dropdown-year_' + random"
                    v-model="object.year"
                    class="bg-light"
                    size="sm"
                  >
                    <option
                      v-for="(year, i) in currentYear - 2019"
                      :value="currentYear - i"
                      :key="`year` + year"
                    >
                      {{ currentYear - i }}
                    </option>
                  </b-form-select>
                </b-form-group>

                <b-form-group
                  label="Month"
                  :label-for="'dropdown-month_' + random"
                >
                  <b-form-select
                    :id="'dropdown-month' + random"
                    v-model="object.month"
                    class="bg-light"
                    size="sm"
                  >
                    <option v-for="i in 12" :value="i" :key="`month` + i">
                      {{ months[i - 1] }}
                    </option>
                  </b-form-select>
                </b-form-group>
              </b-form>
            </b-dropdown-form>
            <b-dropdown-divider></b-dropdown-divider>
            <b-dropdown-item no-hover no-active>
              <b-button variant="info" block @click="getUnpaidHouse"
                >Generate</b-button
              >
            </b-dropdown-item>
          </b-dropdown>
          <!-- <b-button variant="info" @click="getAllHouse">Generate All House Report</b-button> -->
        </div>
      </b-row>
      <b-row class="my-4"></b-row>

      <b-row justify="center" style="margin: auto !important">
        <b-col class="text-center">
          <b-spinner v-if="isLoadingdata" label="Loading..."
            >Loading Data</b-spinner
          >
        </b-col>
      </b-row>
      <b-row no-gutters>
        <b-collapse
          id="sectorreport-collapse"
          class="w-100"
          v-model="state.showReport"
        >
          <div class="reports-card">
            <b-row no-gutters class="mb-2 justify-content-end">
              <!-- <b-badge variant="secondary" class="p-2 fsize-sm"
                >Report Date: &nbsp; {{ state.reportsDate }}
                </b-badge> -->
              <b-form-group>
                <b-form-input
                  required
                  v-model="search"
                  placeholder="Search Here..."
                  class="br-2"
                />
              </b-form-group>
            </b-row>
            <h5 class="bg-dark">{{ reportTitle }}</h5>
            <div class="card--body">
              <b-table
                id="sector-reports"
                :items="reports"
                :fields="table.fields"
                :busy="state.busy.table1"
                head-row-variant="secondary"
                :filter="search"
                small
                bordered
                hover
                responsive
                show-empty
              >
                <template v-slot:table-busy>
                  <vue-load label="Generating..." class="p-3" />
                </template>
                <template v-slot:empty>{{
                  state.error.table1 || "No data available to display"
                }}</template>
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
export default {
  name: "paymentReports",
  data() {
    return {
      dropdownone: false,
      isLoading1: false,
      isLoadingdata: false,
      reportTitle: "",
      search: "",
      form: {
        select: {
          sector: null,
          cell: null,
          village: null,
        },
      },
      object: {
        frommonth: null,
        tomonth: null,
        year: null,
        month: null,
      },
      from: null,
      to: null,
      state: {
        loading: false,
        tableLoad: false,
        changing: false,
        showReport: false,
        reportsDate: null,
        busy: {
          table1: false,
          table2: false,
        },
        error: {
          table1: null,
          table2: null,
        },
      },
      reports: [],
      table: {
        fields: [
          {
            key: "fname",
            label: "First Name",
            tdClass: "text-center",
            thClass: "text-center text-uppercase",
          },
          {
            key: "lname",
            label: "Last Name",
            tdClass: "text-center",
            thClass: "text-center text-uppercase",
          },
          {
            key: "phone",
            label: "Phone",
            tdClass: "text-center",
            thClass: "text-center text-uppercase",
          },
          {
            key: "amount",
            label: "Amount",
            tdClass: "text-center",
            thClass: "text-center text-uppercase",
          },
          // {
          //   key: "unpayedAmount",
          //   label: "unpaid Amount",
          //   tdClass: "text-right",
          //   thClass: "text-center text-uppercase"
          // }
        ],
      },
    };
  },
  computed: {
    sectorOptions() {
      return [this.activeSector];
    },
    cellOptions() {
      const sector = this.form.select.sector;
      const { province, district } = this.location;
      if (sector) return this.$cells(province, district, sector);
      return [];
    },
    villageOptions() {
      const sector = this.form.select.sector;
      const cell = this.form.select.cell;
      const { province, district } = this.location;
      if (sector && cell)
        return this.$villages(province, district, sector, cell);
      return [];
    },
    currentPwd() {
      if (this.change_pswd_modal.data)
        return this.change_pswd_modal.data.password;
      return null;
    },
    user() {
      return this.$store.getters.userDetails;
    },
    activeSector() {
      return this.$store.getters.getActiveSector;
    },
    location() {
      return this.$store.getters.location;
    },
    random() {
      return Math.floor(Math.random() * 101);
    },
    showDownload() {
      if (this.reports == null || !this.state.showReport) return false;
      return true;
    },
    currentYear() {
      this.object.year = new Date().getFullYear();
      return new Date().getFullYear();
    },
    currentMonth() {
      // const currentDate = new Date();
      // this.object.month = currentDate.toLocaleString("default", {
      //   month: "long",
      // });
      return new Date().getMonth() + 1;
    },
    months() {
      return this.$store.getters.getMonths;
    },
    random() {
      return Math.floor(Math.random() * 101);
    },
  },
  watch: {
    "form.select.sector"() {
      handler: {
        this.form.select.cell = null;
        this.form.select.village = null;
      }
    },
    "form.select.cell"() {
      handler: {
        this.form.select.village = null;
      }
    },
    "object.month"() {
      handler: {
        if (this.currentYear == this.object.year) {
          if (this.currentMonth < this.object.month) {
            this.$nextTick(() => {
              this.$set(this.object, "month", this.currentMonth);
            });
          }
        }
      }
    },
  },
  methods: {
    async getAllHouse() {
      this.state.showReport = false;
      this.isLoading1 = true;
      this.isLoadingdata = true;
      this.reportTitle = "Generate All House Report";
      console.log("generate all house");
      this.loading = true;
      const yearString = this.object.year;
      var monthString = this.object.month;
      const yearDate = new Date(`${yearString}-${monthString}-01`);
      const nextMonth = new Date(yearString, monthString + 1, 1);
      const lastDayOfMonth = new Date(nextMonth - 1);
      this.from = yearDate;
      this.to = lastDayOfMonth;
      // if (this.object.frommonth == null && this.object.tomonth == null) {
      //   const currentDate = new Date();
      //   this.from = currentDate;
      //   this.to = currentDate;
      // } else {
      //   this.from = this.object.frommonth;
      //   this.to = this.object.tomonth;
      // }
      try {
        const { data } = await this.axios.get("payment/reports", {
          params: {
            status: "",
            sector: this.form.select.sector,
            cell: this.form.select.cell,
            village: this.form.select.village,
            offset: 0,
            limit: 0,
            from: this.from,
            to: this.to,
          },
        });
        // this.accountant = data;

        this.reports = data.Payments;
        console.log("report all houses", this.reports);
      } catch (error) {
        console.log(error);
      } finally {
        this.state.loading = false;
        this.dropdownone = false;
        this.state.showReport = true;
        this.isLoadingdata = false;
      }
    },
    async getPaidHouse() {
      this.state.showReport = false;
      this.isLoadingdata = true;
      this.reportTitle = "Generate Paid House Report";
      console.log("generate paid house");
      this.loading = true;
      const yearString = this.object.year;
      var monthString = this.object.month;
      const yearDate = new Date(`${yearString}-${monthString}-01`);
      const nextMonth = new Date(yearString, monthString + 1, 1);
      const lastDayOfMonth = new Date(nextMonth - 1);
      this.from = yearDate;
      this.to = lastDayOfMonth;
      try {
        const { data } = await this.axios.get("payment/reports", {
          params: {
            status: "paid",
            sector: this.form.select.sector,
            cell: this.form.select.cell,
            village: this.form.select.village,
            offset: 0,
            limit: 0,
            from: this.from,
            to: this.to,
          },
        });
        // this.accountant = data;

        this.reports = data.Payments;
        console.log("reports", this.reports);
      } catch (error) {
        console.log(error);
      } finally {
        this.state.loading = false;
        this.dropdownone = false;
        this.state.showReport = true;
        this.isLoadingdata = false;
      }
    },
    async getUnpaidHouse() {
      this.state.showReport = false;
      this.isLoadingdata = true;
      this.reportTitle = "Generate Unpaid House Report";
      console.log("generate unpaid house");
      this.loading = true;
      const yearString = this.object.year;
      var monthString = this.object.month;
      const yearDate = new Date(`${yearString}-${monthString}-01`);
      const nextMonth = new Date(yearString, monthString + 1, 1);
      const lastDayOfMonth = new Date(nextMonth - 1);
      this.from = yearDate;
      this.to = lastDayOfMonth;
      try {
        const { data } = await this.axios.get("payment/reports", {
          params: {
            status: "pending",
            sector: this.form.select.sector,
            cell: this.form.select.cell,
            village: this.form.select.village,
            offset: 0,
            limit: 0,
            from: this.from,
            to: this.to,
          },
        });
        // this.accountant = data;
        this.reportTitle = "Generate Unpaid House Report";
        this.reports = data.Payments;
        console.log("reports", this.reports);
      } catch (error) {
        console.log(error);
      } finally {
        this.state.loading = false;
        this.dropdownone = false;
        this.state.showReport = true;
        this.isLoadingdata = false;
      }
    },
    downloadReport() {
      console.log("download");
      return;
      if (this.sectorData.length > 0 && this.cellData.length > 0) {
        const data = {
          config: {
            TITLE: String(
              `Monthly Report of ${this.activeSector}`
            ).toUpperCase(),
            name: `${this.activeSector} Monthly Report of ${this.state.reportsDate}`,
            date: this.state.reportsDate,
          },
          data: [
            {
              COLUMNS: [
                {
                  header: `No of Properties`,
                  dataKey: "total",
                },
                {
                  header: `No of Paid Properties`,
                  dataKey: "payed",
                },
                { header: `Paid Amount`, dataKey: "payedAmount" },
                {
                  header: `No of Unpaid Properties`,
                  dataKey: "pending",
                },
                { header: `Unpaid Amount`, dataKey: "unpayedAmount" },
              ],
              BODY: this.sectorData,
            },
            {
              COLUMNS: [
                {
                  header: `Cell`,
                  dataKey: "name",
                },
                {
                  header: `No of Properties`,
                  dataKey: "total",
                },
                {
                  header: `No of Paid Properties`,
                  dataKey: "payed",
                },
                { header: `Paid Amount`, dataKey: "payedAmount" },
                {
                  header: `No of Unpaid Properties`,
                  dataKey: "pending",
                },
                { header: `Unpaid Amount`, dataKey: "unpayedAmount" },
              ],
              BODY: this.cellData,
            },
          ],
        };
        download(data);
      }
    },
  },
};
</script>

<style lang="scss">
#payment-reports {
  & > header {
    text-align: center;
    font-size: 1.3rem;
    font-weight: bold;
    color: #384950;
  }
  .buttons {
    display: grid;
    grid-gap: 1rem;
    grid-template-columns: repeat(auto-fit, minmax(270px, 1fr));

    & > button {
      padding: 0.7rem 1rem;
    }
  }
}
.input-date {
  border: 1px solid #7f898d;
  color: #212529;
  padding: 3px 5px;
}
</style>
