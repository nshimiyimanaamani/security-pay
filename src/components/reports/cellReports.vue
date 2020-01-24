<template>
  <div>
    <header class="d-flex justify-content-center font-20 text-uppercase">Cell Report</header>
    <hr class="m-0 mb-3" />
    <b-row class="px-3 align-items-center justify-content-between">
      <b-select
        size="sm"
        id="input-1"
        v-model="cell"
        :options="cellOptions"
        class="w-auto mr-3 flex-grow-1"
      >
        <template v-slot:first>
          <option :value="null" disabled>Please select cell</option>
        </template>
      </b-select>
      <b-button
        size="sm"
        variant="info"
        class="font-15 border-0 my-2 d-flex align-items-center"
        :disabled="cell?false:true"
        @click="generate"
      >Generate {{cell ? cell : 'Cell'}} Report</b-button>
      <div v-show="state.generating" class="w-100">
        <strong class="font-15">Generating&nbsp;</strong>
        <b-spinner small />
      </div>
    </b-row>
    <b-row>
      <b-collapse id="sectorreport-collapse" class="w-100" v-model="state.showReport">
        <b-card class="text-capitalize mx-3" v-if="!state.error">
          <b-card-title class="font-weight-bold font-20">cell Overall</b-card-title>
          <b-card-text>{{cell}} cell is having {{population.total}} Houses with {{population.payed}} House{{population.payed>1?'s':''}} that finished paying and {{population.not_payed}} House{{population.not_payed>1?'s':''}} that haven't finished paying</b-card-text>
          <hr />
          <b-card-title class="font-weight-bold font-20">Respective villages</b-card-title>
          <b-card-text v-for="(village,i) in villagesInfo" :key="i">
            <b-card-text>
              <b>{{village.label}}</b>
              cell is having {{village.data.payed+village.data.pending}} Houses with {{village.data.payed}} House{{village.data.payed>1?'s':''}} that finished paying and {{village.data.pending}} House{{village.pending>1?'s':''}} that haven't finished paying
              <hr />
            </b-card-text>
          </b-card-text>
        </b-card>
        <b-card v-if="state.error" class="mx-3">
          <b-card-text>{{state.errorMessage}}</b-card-text>
        </b-card>
      </b-collapse>
    </b-row>
    <b-row
      v-if="!state.generating && !state.error && villagesInfo"
      class="my-3 justify-content-end"
    >
      <b-button size="sm" class="app-color mx-3">Download Report</b-button>
    </b-row>
  </div>
</template>

<script>
export default {
  name: "cellReports",
  data() {
    return {
      cell: null,
      state: {
        generating: false,
        showReport: false,
        error: false,
        errorMessage: null
      },
      population: {
        total: null,
        payed: null,
        not_payed: null
      },
      villagesInfo: null
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    cellOptions() {
      return this.$store.getters.getCellsArray;
    },
    currentYear() {
      return new Date().getFullYear();
    },
    currentMonth() {
      return new Date().getMonth() + 1;
    }
  },
  methods: {
    generate() {
      this.state.showReport = false;
      this.state.error = false;
      this.state.errorMessage = null;
      this.state.generating = true;
      this.villagesInfo = null;
      let axios = this.axios;
      const first = axios.get(
        this.endpoint +
          `/metrics/ratios/cells/${this.cell}?year=${this.currentYear}&month=${this.currentMonth}`
      );
      const second = axios.get(
        this.endpoint +
          `/metrics/ratios/cells/all/${this.cell}?year=${this.currentYear}&month=${this.currentMonth}`
      );
      axios
        .all([first, second])
        .then(
          this.axios.spread((...res) => {
            this.population.total =
              res[0].data.data.payed + res[0].data.data.pending;
            this.population.payed = res[0].data.data.payed;
            this.population.not_payed = res[0].data.data.pending;
            this.villagesInfo = res[1].data;
            console.log(res);
          })
        )
        .catch(err => {
          this.state.error = true;
          console.log(err);
          this.state.errorMessage = err.response.data.error
            ? err.response.data.error
            : err.response.response;
        })
        .finally(() => {
          this.state.showReport = true;
          this.state.generating = false;
        });
    }
  }
};
</script>

<style>
</style>