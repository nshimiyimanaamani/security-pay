<template>
  <div class="position-relative transaction-page bg-light">
    <vue-title title="Paypack | Transactions" />
    <div class="totals primary-font">
      <b-row>
        <b-col
          cols="6"
          class="text-white ml-auto py-2 text-overflow"
          style="font-size: 40px"
          >RWF &nbsp;{{ GrandTotal() | number }}</b-col
        >
        <b-col cols="6" md="1"></b-col>
      </b-row>
      <b-row class="text-white">
        <b-col cols="3" class="ml-auto">
          <p class="text-overflow">Total {{ selectedMonth }}</p>
          <p>RWF {{ MonthTotal() | number }}</p>
        </b-col>
        <b-col cols="3" class="m-0">
          <p class="text-overflow">MTN MoMo</p>
          <p>RWF {{ mtnTotal() | number }}</p>
        </b-col>
        <b-col cols="6" md="1"></b-col>
      </b-row>
    </div>
    <div class="transaction-table  bg-light">
      <header class="primary-font mb-3 table-header">
        <h3>All Transactions</h3>
        
        <fieldset class="control">
          <div v-show="!loading">
            <b-form-select
              class="br-2"
              v-model="select.year"
              :disabled="loading"
            >
              <template v-slot:first>
                <option :value="null" disabled>select year</option>
              </template>
              <b-select-option
                v-for="year in YearsOptions"
                :key="year"
                :value="year"
                >{{ year }}</b-select-option
              >
            </b-form-select>
            <b-form-select
              class="br-2"
              v-model="select.month"
              :disabled="loading"
            >
              <template v-slot:first>
                <option :value="null" disabled>select month</option>
              </template>
              <b-select-option
                v-for="(month, i) in MonthsOptions"
                :key="i"
                :value="i"
                >{{ month }}</b-select-option
              >
            </b-form-select>
          </div>
        </fieldset>
        <b-button
          variant="info"
          @click="requestItems"
          :disabled="loading"
          class="br-2"
        >
          Refresh
          <i class="fa fa-sync-alt" />
        </b-button>
      </header>
      <b-table
        :items="shownData"
        :fields="table.fields"
        :busy="loading"
        head-variant="light"
        table-class="bg-white"
        thead-class="primary-font"
        tbody-class="secondary-font"
        show-empty
        responsive
        bordered
        hover
        striped
      >
        <template v-slot:cell(method)="data">
          <div :class="data.value == 'momo-mtn-rw' ? 'mtn' : data.value">
            <span>{{ data.value == "momo-mtn-rw" ? "mtn" : data.value }}</span>
          </div>
        </template>
        <template v-slot:cell(owner_firstname)="data">
          {{
            data.item.owner_firstname + " " + data.item.owner_lastname
          }}</template
        >
        <template v-slot:cell(amount)="data">
          <article>{{ data.value | number }} Frw</article>
        </template>
        <template v-slot:cell(date_recorded)="data">
          <article>{{ data.value | date }}</article>
        </template>
        <template v-slot:table-busy>
          <div class="table-loader">
            <i class="fa fa-spinner fa-spin" />
            <p>Loading...</p>
          </div>
        </template>
        <template v-slot:empty>
          <label class="table-empty" v-if="!no_data"
            >No records of transactions for the month of
            {{ selectedMonth }}!</label
          >
          <label class="table-empty" v-else
            >No records of transactions found !</label
          >
        </template>
      </b-table>
      <div class="d-flex justify-content-end">
        <div style="">
          <p class="d-flex justify-content-between" style="font-size: 18px">
            <strong>Amount </strong>
            <span class="ml-5"> {{ MonthTotal() | number }} Rwf</span>
          </p>
          <br />
          <p class="d-flex justify-content-between" style="font-size: 18px">
            <strong>Fee </strong>
            <span> {{ FeeTotal() | number }} Rwf</span>
          </p>
          <hr />
          <p class="d-flex justify-content-between" style="font-size: 18px">
            <strong></strong>
            <span> {{ (totalMonth - totalFee) | number }} Rwf</span>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import loader from "../components/loader";
export default {
  components: {
    loader,
  },
  data() {
    return {
      loading: false,
      totalFee: 0,
      totalMonth: 0,
      transactionData: {},
      table: {
        fields: [
          {
            key: "owner_firstname",
            label: "Names",
            sortable: true,
            tdClass: "table-name",
          },
          {
            key: "madefor",
            label: "Paid for",
            sortable: true,
          },
          {
            key: "method",
            label: "Paid With",
            sortable: true,
            thClass: "text-center",
            tdClass: "text-center",
          },
          {
            key: "amount",
            label: "Amount",
            sortable: true,
            thClass: "text-right",
            tdClass: "text-right",
          },
          {
            key: "fee",
            label: "Fee",
            sortable: true,
            thClass: "text-right",
            tdClass: "text-right",
          },
          {
            key: "date_recorded",
            label: "Date & Time",
            sortable: true,
            thClass: "text-right",
            tdClass: "text-right",
          },
        ],
        items: [],
      },
      select: {
        year: new Date().getFullYear(),
        month: new Date().getMonth(),
      },
      a: 0,
    };
  },
  computed: {
    activeCell() {
      return this.$store.getters.getActiveCell;
    },
    user() {
      return this.$store.getters.userDetails;
    },
    isManager() {
      return this.user.role === "basic";
    },
    shownData() {
      const { year, month } = this.select;
      if (Object.keys(this.transactionData).length < 1) return [];
      else {
        const Data = this.transactionData[year]
          ? this.transactionData[year][month]
          : [];
        // this.getPercentage();
        this.totalMonth = 0;
        this.totalFee = 0;
        for (let i = 0; i < Data && Data.length; i++) {
          this.totalMonth += Data[i].amount;
          this.totalFee += Data[i].fee;
        }

        return Data;
      }
    },

    MonthsOptions() {
      return this.$store.getters.getMonths.sort(
        (a, b) => a.date_recorded - b.date_recorded
      );
    },
    YearsOptions() {
      if (Object.keys(this.transactionData).length < 1) return [];
      return Object.keys(this.transactionData).filter(
        (year) => Number(year) != NaN
      );
    },
    no_data() {
      return Object.keys(this.transactionData).length < 1;
    },
    selectedMonth() {
      return this.MonthsOptions[this.select.month];
    },
  },
  filters: {
    number: (num) => {
      return Number(num).toLocaleString();
    },
    date: (date) => {
      try {
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
        return date;
      }
    },
  },
  mounted() {
    this.requestItems();
  },
  methods: {
    async requestItems() {
      this.totalMonth = 0;
      this.totalFee = 0;
      this.loading = true;
      const total = await this.$getTotal("/transactions?offset=0&limit=0");
      this.axios
        .get("/transactions?offset=0&limit=" + total)
        .then((res) => {
          this.transactionData = {};
          let DataToClean = res.data.Transactions;
          if (this.isManager) {
            DataToClean = DataToClean.filter(
              (item) => item.cell == this.activeCell
            );
          }
          for (let { date_recorded: date } of DataToClean) {
            const year = new Date(date).getFullYear();
            const month = new Date(date).getMonth();
            if (!this.transactionData[year]) this.transactionData[year] = {};
            if (!this.transactionData[year][month])
              this.transactionData[year][month] = [];
            this.SortByMonth(this.transactionData, year, month, DataToClean, [
              "date_recorded",
            ]);
          }
          this.getPercentage();
        })
        .catch((err) => {
          this.table.items = [];
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
        })
        .finally(() => {
          this.loading = false;
        });
    },
    GrandTotal() {
      if (this.transactionData.length < 1) return 0;
      let total = 0;

      try {
        this.YearsOptions.forEach((year) => {
          total += Object.values(this.transactionData[year])
            .flat()
            .reduce((a, b) => a + Number(b.amount), 0);
        });
        return total;
      } catch {
        return 0;
      }
    },
    MonthTotal() {
      if (this.shownData && this.shownData.length < 1) return 0;
      try {
        return this.shownData.reduce((a, b) => Number(a) + Number(b.amount), 0);
      } catch {
        return 0;
      }
    },
    FeeTotal() {
      if (this.shownData && this.shownData.length < 1) return 0;
      try {
        return this.shownData.reduce((a, b) => Number(a) + Number(b.fee), 0);
      } catch {
        return 0;
      }
    },
    mtnTotal() {
      // console.log(this.shownData);
      if (this.shownData && this.shownData.length < 1) return 0;
      try {
        return this.shownData
          .filter((data) => data.method.includes("mtn"))
          .reduce((a, b) => Number(a) + Number(b.amount), 0);
      } catch {
        return 0;
      }
    },
    GetObjectValue(array, obj) {
      array.forEach((key) => {
        obj = obj[key];
      });
      return obj;
    },
    SortByMonth(o, year, month, array, key) {
      o[year][month] = array.filter(
        (item) =>
          new Date(this.GetObjectValue(key, item)).getFullYear() === year &&
          new Date(this.GetObjectValue(key, item)).getMonth() === month
      );
    },
    getPercentage() {
      // console.log("hello");
      // this.totalMonth = 0;
      // this.totalFee = 0;
      // for (let i = 0; i < this.shownData.length; i++) {
      //   this.totalMonth += this.shownData[i].amount;
      //   this.totalFee += this.shownData[i].fee;
      // }
    },
  },
};
</script>
<style lang="scss" scoped>
@import "../assets/css/transactionTable.scss";
</style>
