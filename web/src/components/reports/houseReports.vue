<template>
  <div id="house-report">
    <header class="tabTitle">House Report</header>
    <div class="tabBody">
      <b-row class="controls" no-gutters>
        <b-input v-model="houseId" placeholder="Enter House ID..." class="text-uppercase mt-3 br-2" />
        <b-button variant="info" :disabled="houseId ? false : true" class="my-4 br-2" @click="generate">Generate House
          Report</b-button>
        <vue-load v-if="state.generating" label="Generating..." />
      </b-row>
      <b-row class="justify-content-center text-capitalize" no-gutters>
        <b-collapse id="housereport-collapse" v-model="state.showReport" class="w-100">
          <b-table-simple hover bordered small caption-top responsive v-if="showUserDetails">
            <caption>Details of {{ userDetails.names || 'user' }}:</caption>
            <b-tbody>
              <b-tr>
                <b-th>Names</b-th>
                <b-td>{{ userDetails.names }}</b-td>
              </b-tr>
              <b-tr>
                <b-th>Phone Number</b-th>
                <b-td>{{ userDetails.phone_number }}</b-td>
              </b-tr>
              <b-tr>
                <b-th>House ID</b-th>
                <b-td>{{ userDetails.house_id }}</b-td>
              </b-tr>
              <b-tr>
                <b-th>Location</b-th>
                <b-td>{{ userDetails.location }}</b-td>
              </b-tr>
              <b-tr>
                <b-th>Amount</b-th>
                <b-td>{{ userDetails.amount }}</b-td>
              </b-tr>
              <b-tr>
                <b-th>For Rent</b-th>
                <b-td>{{ userDetails.forRent }}</b-td>
              </b-tr>
              <b-tr>
                <b-th>Registered by</b-th>
                <b-td style="text-transform: none;">{{ userDetails.recorded_by }}</b-td>
              </b-tr>
              <b-tr>
                <b-th>Registered on</b-th>
                <b-td>{{ userDetails.created_at }}</b-td>
              </b-tr>
            </b-tbody>
          </b-table-simple>
          <b-table-simple hover bordered small caption-top responsive v-if="paymentDetails.length > 0">
            <caption>Payment History of {{ userDetails.names || 'user' }}:</caption>
            <b-tbody>
              <b-tr>
                <b-th>Month</b-th>
                <b-td>Status</b-td>
              </b-tr>
              <b-tr v-for="(item, index) in paymentDetails" :key="index">
                <b-th>{{ item.created_at }}</b-th>
                <b-td :title="item.status == 'Paid' ? 'Paid On ' + item.updated_at : 'Not Yet Paid'"
                  class="cursor-pointer">
                  {{ item.status }}
                  <small v-if="item.status == 'Paid'" style="text-size: 0.8em !important">
                    ({{ item.updated_at }})
                  </small>
                </b-td>
              </b-tr>
            </b-tbody>
          </b-table-simple>
          <div v-else class="w-100 d-flex justify-content-center align-items-center p-5 bg-light">
            <p>{{ state.error || 'Oops, something went wrong!' }}</p>
          </div>
        </b-collapse>
      </b-row>
      <b-row class="justify-content-end" v-if="showDownload" no-gutters>
        <b-button variant="info" class="br-2" @click="downloadReports">
          <i class="fa fa-download mr-1" />Download Report
        </b-button>
      </b-row>
    </div>
  </div>
</template>

<script>
import download from "../download scripts/downloadReports";
export default {
  name: "cellReports",
  data() {
    return {
      houseId: null,
      userDetails: [],
      paymentDetails: [],
      state: {
        showReport: false,
        generating: false,
        error: null
      }
    };
  },
  computed: {
    showDownload() {
      if (
        this.userDetails.length < 1 ||
        this.paymentDetails.length < 1 ||
        this.state.showReport === false ||
        this.state.generating === true ||
        this.state.error !== null
      )
        return false;
      return true;
    },
    showUserDetails() {
      return Object.keys(this.userDetails).length > 0;
    }
  },
  methods: {
    generate() {
      this.state.generating = true;
      this.userDetails = [];
      this.paymentDetails = [];
      this.state.showReport = false;
      this.axios
        .get("/properties/" + this.houseId.toUpperCase())
        .then(res => {
          this.userDetails = [res.data].map(item => {
            return {
              names: `${item.owner.fname} ${item.owner.lname}`,
              phone_number: item.owner.phone,
              house_id: item.id,
              location: `${item.address.sector}, ${item.address.cell}, ${item.address.village}`,
              amount: `${this.$options.filters.number(item.due)} Rwf`,
              forRent: item.occupied ? "Yes" : "No",
              recorded_by: item.recorded_by,
              created_at: this.$options.filters.date(item.created_at),
              updated_at: this.$options.filters.date(item.updated_at)
            };
          })[0];
          this.generatePayment();
        })
        .catch(err => {
          this.state.showReport = true;
          this.state.generating = false;
          try {
            this.state.error = err.response.data.error || err.response.data;
          } catch {
            this.state.error = "Failed to retrieve property details!";
          }
        });
    },
    generatePayment() {
      this.state.showReport = false;
      this.axios
        .get("/billing/invoices?property=" + this.houseId + "&months=12")
        .then(res => {
          this.paymentDetails = res.data.Invoices.map(item => {
            return {
              id: item.id,
              propertyID: item.property,
              amount: this.$options.filters.number(item.amount),
              status: item.status == "pending" ? "Not paid" : "Paid",
              created_at: this.$options.filters.date(item.created_at),
              updated_at: this.$options.filters.date(item.updated_at)
            };
          });
          this.state.showReport = true;
          this.state.generating = false;
        })
        .catch(err => {
          try {
            this.state.error = err.response.data.error || err.response.data;
          } catch {
            this.state.error = "Failed to retrieve user payment Details";
          }
          this.state.generating = false;
        });
    },
    downloadReports() {
      if (this.paymentDetails.length > 1 && this.userDetails) {
        const data = {
          config: {
            TITLE: String(
              `Report of ${this.userDetails.names}'s Property`
            ).toUpperCase(),
            name: `Report of ${this.userDetails.names}'s Property`,
            date: this.state.reportsDate
          },
          data: [
            {
              SHOWHEAD: "never",
              COLUMNS: [
                {
                  header: `key`,
                  dataKey: "key"
                },
                {
                  header: `value`,
                  dataKey: "value"
                }
              ],
              BODY: [this.userDetails].map(item => {
                return [
                  { key: "Names", value: item.names },
                  { key: "Phone Number", value: item.phone_number },
                  { key: "House Id", value: item.house_id },
                  { key: "Location", value: item.location },
                  { key: "Amount", value: item.amount },
                  { key: "For Rent", value: item.forRent },
                  { key: "Registered By", value: item.recorded_by },
                  { key: "Registered On", value: item.created_at }
                ];
              })[0]
            },
            {
              COLUMNS: [
                {
                  header: "Payment History",
                  dataKey: "key"
                },
                {
                  header: "Status",
                  dataKey: "value"
                }
              ],
              BODY: this.paymentDetails.map(obj => {
                return {
                  key: obj.created_at,
                  value: obj.status
                };
              })
            }
          ]
        };
        download(data);
      }
    },
    clear() {
      this.state.generating = false;
      this.state.error = null;
    }
  }
};
</script>

<style lang='scss'>
#house-report {
  &>header {
    text-align: center;
    font-size: 1.3rem;
    font-weight: bold;
    color: #384950;
  }

  .controls {
    display: flex;
    flex-direction: column;
    max-width: 500px;
    margin: auto;
  }
}
</style>