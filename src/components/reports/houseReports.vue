<template>
  <div id="house-report">
    <header class="tabTitle">House Report</header>
    <div class="tabBody">
      <b-row class="controls" no-gutters>
        <b-input
          v-model="houseId"
          placeholder="Enter House ID..."
          class="text-uppercase mt-3 br-2"
        />
        <b-button
          variant="info"
          :disabled="houseId?false:true"
          class="my-4 br-2"
          @click="generate"
        >Generate House Report</b-button>
        <vue-load v-if="state.generating" label="Generating..." />
      </b-row>
      <b-row class="m-0 justify-content-center text-capitalize" no-gutters>
        <b-collapse id="housereport-collapse" v-model="state.showReport">
          <b-table-simple hover bordered small caption-top responsive v-if="userDetails">
            <caption>Details of {{userDetails.owner.fname+' '+userDetails.owner.lname}}:</caption>
            <b-tbody>
              <b-tr>
                <b-th>Names</b-th>
                <b-td>{{userDetails.owner.fname+' '+userDetails.owner.lname}}</b-td>
              </b-tr>
              <b-tr>
                <b-th>Phone Number</b-th>
                <b-td>{{userDetails.owner.phone}}</b-td>
              </b-tr>
              <b-tr>
                <b-th>House ID</b-th>
                <b-td>{{userDetails.id}}</b-td>
              </b-tr>
              <b-tr>
                <b-th>Location</b-th>
                <b-td>{{userDetails.address.sector+', '+userDetails.address.cell+', '+userDetails.address.village}}</b-td>
              </b-tr>
              <b-tr>
                <b-th>Amount</b-th>
                <b-td>{{userDetails.due | number}} Rwf</b-td>
              </b-tr>
              <b-tr>
                <b-th>For Rent</b-th>
                <b-td>{{userDetails.occupied?userDetails.occupied?"Yes":"No":'No'}}</b-td>
              </b-tr>
              <b-tr>
                <b-th>Registered by</b-th>
                <b-td style="text-transform: none;">{{userDetails.recorded_by}}</b-td>
              </b-tr>
              <b-tr>
                <b-th>Registered on</b-th>
                <b-td>{{userDetails.created_at | date}}</b-td>
              </b-tr>
            </b-tbody>
          </b-table-simple>
        </b-collapse>
      </b-row>
      <b-row class="justify-content-center text-capitalize">
        <vue-load v-if="state.generatingP" label="Generating..." />
        <b-collapse
          id="PaymentReport-collapse"
          class="flex-grow-1 mx-3"
          v-model="state.showPayment"
        >
          <b-table-simple hover bordered small caption-top responsive v-if="paymentDetails">
            <caption>Payment History of {{userDetails.owner.fname+' '+userDetails.owner.lname}}:</caption>
            <b-tbody>
              <b-tr>
                <b-th>Month</b-th>
                <b-td>Status</b-td>
              </b-tr>
              <b-tr v-for="(item,index) in paymentDetails" :key="index">
                <b-th>{{item.created_at | date}}</b-th>
                <b-td>{{item.status=="pending"?'Not Paid':'Paid'}}</b-td>
              </b-tr>
            </b-tbody>
          </b-table-simple>
        </b-collapse>
      </b-row>
      <b-row class="justify-content-end" v-if="paymentDetails" no-gutters>
        <b-button variant="info" @click="download">Download Report</b-button>
      </b-row>
    </div>
  </div>
</template>

<script>
import DownloadReport from "./downloadReport";
export default {
  name: "cellReports",
  data() {
    return {
      houseId: null,
      userDetails: null,
      paymentDetails: null,
      state: {
        showReport: false,
        showPayment: false,
        generatingP: false,
        generating: false
      }
    };
  },
  computed: {
    cellOptions() {
      return this.$store.getters.getCellsArray;
    }
  },
  mounted() {},
  methods: {
    generate() {
      this.state.generating = true;
      this.userDetails = null;
      this.paymentDetails = null;
      this.state.showReport = false;
      this.state.showPayment = false;
      this.axios
        .get("/properties/" + this.houseId.toUpperCase())
        .then(res => {
          this.state.showReport = true;
          this.userDetails = res.data;
          this.generatePayment();
        })
        .catch(err => {
          this.state.showReport = false;
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
        })
        .finally(() => {
          this.state.generating = false;
        });
    },
    generatePayment() {
      this.state.showPayment = false;
      this.state.generatingP = true;
      this.axios
        .get("/billing/invoices?property=" + this.houseId + "&months=12")
        .then(res => {
          this.paymentDetails = res.data.Invoices;
          this.state.showPayment = true;
        })
        .catch(err => {
          this.state.showPayment = true;
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
        })
        .finally(() => (this.state.generatingP = false));
    },
    download() {
      if (this.paymentDetails != null && this.userDetails != null) {
        DownloadReport(this.userDetails, this.paymentDetails);
      }
    }
  }
};
</script>

<style lang='scss'>
#house-report {
  & > header {
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