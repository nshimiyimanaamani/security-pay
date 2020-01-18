<template>
  <div>
    <header class="d-flex justify-content-center font-20 text-uppercase">Create House Report</header>
    <hr class="m-0 mb-3" />
    <b-row class="px-3 align-items-center justify-content-between">
      <b-input
        size="sm"
        v-model="houseId"
        placeholder="Enter House ID..."
        class="w-auto mr-3 flex-grow-1 text-uppercase"
      ></b-input>
      <b-button
        size="sm"
        variant="info"
        :disabled="houseId?false:true"
        class="font-15 border-0 my-3"
        @click="generate"
      >Generate House Report</b-button>
    </b-row>
    <b-row class="justify-content-center text-capitalize">
      <div v-show="state.generating" class="w-100 px-3">
        <strong class="font-15">Generating&nbsp;</strong>
        <b-spinner small />
      </div>
      <b-collapse id="housereport-collapse" class="flex-grow-1 mx-3" v-model="state.showReport">
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
              <b-td>{{Number(userDetails.due).toLocaleString()}} Rwf</b-td>
            </b-tr>
            <b-tr>
              <b-th>Occupied</b-th>
              <b-td>{{userDetails.occupied}}</b-td>
            </b-tr>
            <b-tr>
              <b-th>For Rent</b-th>
              <b-td>{{userDetails.for_rent}}</b-td>
            </b-tr>
            <b-tr>
              <b-th>Recorded by</b-th>
              <b-td>{{userDetails.recorded_by}}</b-td>
            </b-tr>
            <b-tr>
              <b-th>Created on</b-th>
              <b-td>{{new Date(userDetails.created_at).toLocaleString('en-EN', { year: 'numeric', month: 'long', day: 'numeric' })}}</b-td>
            </b-tr>
          </b-tbody>
        </b-table-simple>
      </b-collapse>
    </b-row>
  </div>
</template>

<script>
export default {
  name: "cellReports",
  data() {
    return {
      houseId: null,
      userDetails: null,
      state: {
        showReport: false,
        generating: false
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    cellOptions() {
      return this.$store.getters.getCellsArray;
    }
  },
  methods: {
    generate() {
      this.state.generating = true;
      this.userDetails = null;
      this.axios
        .get(this.endpoint + "/properties/" + this.houseId.toUpperCase())
        .then(res => {
          this.state.showReport = true;
          this.userDetails = res.data;
        })
        .catch(err => {
          if (navigator.onLine) {
            this.state.showReport = false;
            const error = err.response
              ? err.response.data.error || err.response.data
              : "an error occured";
            this.$snotify.error(error);
          } else {
            this.$snotify.error("Please connect to the internet...");
          }
        })
        .finally(() => {
          this.state.generating = false;
        });
    }
  }
};
</script>

<style>
</style>