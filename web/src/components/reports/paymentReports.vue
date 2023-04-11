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
      <b-row v-if="state.showReport" no-gutters>
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
                <template v-slot:cell(index)="data">
                  <article class="text-center">{{ data.index + 1 }}</article>
                </template>
                <template v-slot:cell(amount)="data">
                  <article class="text-center">
                    {{ data.item.amount | number }}
                  </article>
                </template>
                <template v-slot:custom-foot>
                  <b-tr class="total">
                    <b-td></b-td>
                    <b-td></b-td>
                    <b-td></b-td>
                    <b-td></b-td>
                    <b-td class="text-center py-3">
                      <small
                        ><strong style=""
                          ><span style="color: #dc3545">Total </span>:
                          {{ totalAmount | number }} Rwf</strong
                        ></small
                      >
                    </b-td>
                  </b-tr>
                </template>
              </b-table>
              <b-pagination
                class="my-0"
                align="center"
                v-if="showPagination"
                :per-page="pagination.perPage"
                v-model="pagination.currentPage"
                :total-rows="pagination.totalRows"
                @input="pageChanged"
              ></b-pagination>
            </div>
          </div>
        </b-collapse>
      </b-row>

      <!-- dailyreporttable -->

      <b-row v-else no-gutters>
        <b-collapse
          id="sectorreport-collapse"
          class="w-100"
          v-model="state.showReport2"
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
                :items="dailyreports"
                :fields="table.dailyfields"
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
                <template v-slot:cell(index)="data">
                  <article class="text-center">{{ data.index + 1 }}</article>
                </template>
                <template v-slot:cell(created_at)="data">
                  <article class="text-center">
                    {{ data.item.created_at | date }}
                  </article>
                </template>
                <template v-slot:cell(amount)="data">
                  <article class="text-center">
                    {{ data.item.amount | number }}
                  </article>
                </template>
                <template v-slot:custom-foot>
                  <b-tr class="total">
                    <b-td></b-td>
                    <b-td></b-td>
                    <b-td></b-td>
                    <b-td class="text-center py-3">
                      <small
                        ><strong style=""
                          ><span style="color: #dc3545">Total </span>:
                          {{ houseTotal | number }} Houses</strong
                        ></small
                      >
                    </b-td>
                    <b-td class="text-center py-3">
                      <small
                        ><strong style=""
                          ><span style="color: #dc3545">Total </span>:
                          {{ dailyTotal | number }} Rwf</strong
                        ></small
                      >
                    </b-td>
                    <b-td></b-td>
                  </b-tr>
                </template>
              </b-table>
              <!-- <b-pagination class="my-0" align="center" v-if="showPagination" :per-page="pagination.perPage"
                v-model="pagination.currentPage" :total-rows="pagination.totalRows" @input="pageChanged"></b-pagination> -->
            </div>
          </div>
        </b-collapse>
      </b-row>

      <!-- dailyreporttable -->

      <b-row v-if="showDownload" class="py-3 justify-content-end" no-gutters>
        <b-button @click="downloadReport" variant="info" class="downloadBtn">
          <i class="fa fa-download mr-1" />Download Report
        </b-button>
      </b-row>
      <b-row v-if="showDownload2" class="py-3 justify-content-end" no-gutters>
        <b-button
          @click="downloadDailyReport"
          variant="info"
          class="downloadBtn"
        >
          <i class="fa fa-download mr-1" />Download Report
        </b-button>
      </b-row>
    </div>
  </div>
</template>

<script>
import download from "../download scripts/downloadReports";
export default {
  name: "paymentReports",
  data() {
    return {
      dropdownone: false,
      isLoading1: false,
      isLoadingdata: false,
      reportTitle: "",
      totalAmount: 0,
      dailyTotal: 0,
      houseTotal: 0,
      search: "",
      form: {
        select: {
          sector: null,
          cell: null,
          village: null,
        },
      },
      object: {
        frommonth: "",
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
        showReport2: false,
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
      dailyreports: [],
      table: {
        fields: [
          {
            key: "index",
            label: "No",
            tdClass: "",
            thClass: " text-uppercase",
          },
          {
            key: "fname",
            label: "Full Name",
            formatter: (value, key, item) => `${item.fname} ${item.lname}`,
            tdClass: "",
            thClass: " text-uppercase",
          },

          {
            key: "phone",
            label: "Phone",
            tdClass: "text-center",
            thClass: "text-center text-uppercase",
          },
          {
            key: "property_id",
            label: "Property",
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
        dailyfields: [
          {
            key: "index",
            label: "No",
            tdClass: "",
            thClass: " text-uppercase",
          },
          {
            key: "cell",
            label: "Cell",
            tdClass: "text-center",
            thClass: "text-center text-uppercase",
          },
          {
            key: "village",
            label: "Vilage",
            tdClass: "text-center",
            thClass: "text-center text-uppercase",
          },
          {
            key: "houses",
            label: "No Of House",
            tdClass: "text-right",
            thClass: "text-center text-uppercase",
          },
          {
            key: "amount",
            label: "Amount",
            tdClass: "text-center",
            thClass: "text-center text-uppercase",
          },
          {
            key: "created_at",
            label: "Date",
            tdClass: "text-right",
            thClass: "text-center text-uppercase",
          },
        ],
      },
      pagination: {
        perPage: 20,
        currentPage: 1,
        totalRows: 1,
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
    showDownload2() {
      if (this.dailyreports == null || !this.state.showReport2) return false;
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
      this.object.month = new Date().getMonth() + 1;
      return new Date().getMonth() + 1;
    },
    months() {
      return this.$store.getters.getMonths;
    },
    random() {
      return Math.floor(Math.random() * 101);
    },
    showPagination() {
      if (this.isLoadingdata) return false;
      // if (this.pagination.totalRows < this.pagination.perPage) return false;
      return true;
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
  filters: {
    number: (num) => {
      return Number(num).toLocaleString();
    },
    date: (date) => {
      try {
        return date.slice(0, 10);
        const o = new Date(date);
        return new Intl.DateTimeFormat("en-US", {
          hour12: false,
          year: "numeric",
          month: "long",
          day: "2-digit",
          hour: "numeric",
          minute: "numeric",
        }).format(o);
      } catch {
        return date.replace("T", " ").replace("Z", "");
      }
    },
  },
  created() {
    this.getCurrentMonth();
    this.getAllHouse();
  },
  methods: {
    getCurrentMonth() {
      this.object.month = new Date().getMonth() + 1;
    },
    async getAllHouse() {
      this.state.showReport = false;
      this.state.showReport2 = false;
      this.isLoading1 = true;
      this.isLoadingdata = true;
      this.reportTitle = "All House Report";
      // console.log("generate all house");
      this.loading = true;
      const yearString = this.object.year;
      var monthString = this.object.month;
      const yearDate = new Date(`${yearString}-${monthString}-01`);
      const today = new Date();
      const lastDayOfMonth = new Date(
        today.getFullYear(),
        today.getMonth() + 1,
        0
      );
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
            sector: this.form.select.sector || "",
            cell: this.form.select.cell || "",
            village: this.form.select.village || "",
            offset: (this.pagination.currentPage - 1) * this.pagination.perPage,
            limit: this.pagination.perPage,
            from: this.from,
            to: this.to,
          },
        });
        // this.accountant = data;

        this.reports = data.Payments;
        // const custrow = {};
        // custrow.fname = "Total Amount"
        // custrow.name = ""
        // custrow.name = ""
        // custrow.amount = data.amount
        // this.reports.push(custrow);
        this.totalAmount = data.amount;
        this.pagination.totalRows = data.Total;

        // console.log("report all houses", this.reports);
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
      this.state.showReport2 = false;
      this.isLoadingdata = true;
      this.reportTitle = "Paid House Report";
      // console.log("generate paid house");
      this.loading = true;
      const yearString = this.object.year;
      var monthString = this.object.month;
      const yearDate = new Date(`${yearString}-${monthString}-01`);
      const today = new Date();
      const lastDayOfMonth = new Date(
        today.getFullYear(),
        today.getMonth() + 1,
        0
      );
      this.from = yearDate;
      this.to = lastDayOfMonth;
      try {
        const { data } = await this.axios.get("payment/reports", {
          params: {
            status: "payed",
            sector: this.form.select.sector || "",
            cell: this.form.select.cell || "",
            village: this.form.select.village || "",
            offset: (this.pagination.currentPage - 1) * this.pagination.perPage,
            limit: this.pagination.perPage,
            from: this.from,
            to: this.to,
          },
        });
        // this.accountant = data;

        this.reports = data.Payments;
        // const custrow = {};
        // custrow.fname = "Total Amount"
        // custrow.name = ""
        // custrow.name = ""
        // custrow.amount = data.amount
        // this.reports.push(custrow);
        this.totalAmount = data.amount;
        this.pagination.totalRows = data.Total;
        // console.log("reports", this.reports);
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
      this.state.showReport2 = false;
      this.isLoadingdata = true;
      this.reportTitle = " Unpaid House Report";
      // console.log("generate unpaid house");
      this.loading = true;
      const yearString = this.object.year;
      var monthString = this.object.month;
      const yearDate = new Date(`${yearString}-${monthString}-01`);
      const today = new Date();
      const lastDayOfMonth = new Date(
        today.getFullYear(),
        today.getMonth() + 1,
        0
      );
      this.from = yearDate;
      this.to = lastDayOfMonth;
      try {
        const { data } = await this.axios.get("payment/reports", {
          params: {
            status: "pending",
            sector: this.form.select.sector || "",
            cell: this.form.select.cell || "",
            village: this.form.select.village || "",
            offset: (this.pagination.currentPage - 1) * this.pagination.perPage,
            limit: this.pagination.perPage,
            from: this.from,
            to: this.to,
          },
        });
        // this.accountant = data;
        this.reportTitle = "Unpaid House Report";
        this.reports = data.Payments;
        // const custrow = {};
        // custrow.fname = "Total Amount"
        // custrow.name = ""
        // custrow.name = ""
        // custrow.amount = data.amount
        // this.reports.push(custrow);
        this.totalAmount = data.amount;
        this.pagination.totalRows = data.Total;
        // console.log("reports", this.reports);
      } catch (error) {
        console.log(error);
      } finally {
        this.state.loading = false;
        this.dropdownone = false;
        this.state.showReport = true;
        this.isLoadingdata = false;
      }
    },
    async getDailyReport() {
      this.reportTitle = "Daily Report";
      this.state.showReport = false;
      this.state.showReport2 = false;
      this.dailyTotal = 0;
      this.isLoadingdata = true;
      // console.log("generate daily report");
      this.loading = true;

      this.from = this.object.frommonth;
      this.to = this.object.tomonth;
      try {
        const { data } = await this.axios.get("payment/summary/today", {
          params: {
            sector: this.form.select.sector || "",
            cell: this.form.select.cell || "",
            village: this.form.select.village || "",
            offset: (this.pagination.currentPage - 1) * this.pagination.perPage,
            limit: this.pagination.perPage,
            date: this.object.frommonth,
            // to: this.to,
          },
        });
        if (this.form.select.cell != null) {
          this.reportTitle = ` Daily Report Of  ${this.form.select.cell}`;
        } else {
          this.reportTitle = ` Daily Report Of  ${this.form.select.sector}`;
        }
        // console.log("data", data);
        this.dailyTotal = 0;
        this.houseTotal = 0;
        this.dailyreports = data.summaries;
        for (let i = 0; i < data.summaries.length; i++) {
          this.dailyTotal += data.summaries[i].amount;
          this.houseTotal += data.summaries[i].houses;
        }
        // this.pagination.totalRows = data.Total;
        // console.log("reports", this.reports);
      } catch (error) {
        console.log(error);
      } finally {
        this.state.loading = false;
        this.dropdownone = false;
        this.state.showReport = false;
        this.state.showReport2 = true;
        this.isLoadingdata = false;
      }
    },
    pageChanged() {
      if (this.reportTitle == "All House Report") {
        this.getAllHouse();
      } else if (this.reportTitle == "Paid House Report") {
        this.getPaidHouse();
      } else if (this.reportTitle == "Unpaid House Report") {
        this.getUnpaidHouse();
      } else {
        this.getDailyReport();
      }
    },

    downloadReport() {
      // console.log(this.reports);
      if (this.reports.length > 0) {
        const data = {
          config: {
            TITLE: String(` ${this.reportTitle}`).toUpperCase(),
            name: `${this.reportTitle} `,
            date: this.object.month,
          },
          data: [
            // {
            //   COLUMNS: [
            //     {
            //       header: `No of Properties`,
            //       dataKey: "total",
            //     },
            //     {
            //       header: `No of Paid Properties`,
            //       dataKey: "payed",
            //     },
            //     { header: `Paid Amount`, dataKey: "payedAmount" },
            //     {
            //       header: `No of Unpaid Properties`,
            //       dataKey: "pending",
            //     },
            //     { header: `Unpaid Amount`, dataKey: "unpayedAmount" },
            //   ],
            //   BODY: this.sectorData,
            // },
            {
              COLUMNS: [
                {
                  header: `Names`,
                  dataKey: `name`,
                },
                // {
                //   header: `First Name`,
                //   dataKey: "fname",
                // },
                // {
                //   header: `Last Name`,
                //   dataKey: "lname",
                // },
                {
                  header: `Phone`,
                  dataKey: "phone",
                },
                {
                  header: `Property`,
                  dataKey: "property_id",
                },
                { header: `Amount`, dataKey: `amount` },
                // {
                //   header: `No of Unpaid Properties`,
                //   dataKey: "pending",
                // },
                // { header: `Unpaid Amount`, dataKey: "unpayedAmount" },
              ],
              BODY: this.reports
                .map((report) => ({
                  name: `${report.fname} ${report.lname}`,
                  phone: report.phone,
                  property_id: report.property_id,
                  amount: report.amount,
                }))
                .concat({
                  name: "",
                  phone: "",
                  property_id: "",
                  amount: `Total: ${this.totalAmount} Rwf`,
                }),
            },
          ],
        };
        download(data);
      }
    },
    downloadDailyReport() {
      // console.log(this.dailyreports);
      if (this.dailyreports.length > 0) {
        const data = {
          config: {
            TITLE: String(` ${this.reportTitle}`).toUpperCase(),
            name: `${this.reportTitle} `,
            // date: this.object.month,
          },
          data: [
            {
              COLUMNS: [
                {
                  header: `Cell`,
                  dataKey: "cell",
                },
                {
                  header: `Village`,
                  dataKey: "village",
                },
                { header: `No Of Houses`, dataKey: "houses" },
                { header: `Amount`, dataKey: "amount" },
                { header: `Date`, dataKey: "created_at" },
                // {
                //   header: `No of Unpaid Properties`,
                //   dataKey: "pending",
                // },
                // { header: `Unpaid Amount`, dataKey: "unpayedAmount" },
              ],
              BODY: this.dailyreports,
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

.total {
  background: #b8daff;
  margin: 0 2px;
  padding: 10px;
}
</style>
