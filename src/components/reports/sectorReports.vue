<template>
  <div class="px-4">
    <header class="d-flex justify-content-center font-20 text-uppercase">Sector Report</header>
    <hr class="m-0 mb-3" />
    <b-row class="justify-content-between">
      <b-button
        variant="info"
        class="font-15 border-0 my-2 py-2 mr-3 d-flex align-items-center"
        @click="generate"
      >Generate {{activeSector}} Report</b-button>
      <div v-show="state.generating" class="w-auto px-3">
        <strong class="font-15">Generating&nbsp;</strong>
        <b-spinner small />
      </div>
    </b-row>
    <b-row>
      <b-collapse id="sectorreport-collapse" class="w-100" v-model="state.showReport">
        <b-card class="text-capitalize" v-if="!state.error">
          <b-card-title class="font-weight-bold font-20">Sector Overall</b-card-title>
          <b-card-text>{{activeSector}} Sector is having {{population.total}} Houses with {{population.payed}} House{{population.payed>1?'s':''}} that finished paying and {{population.not_payed}} House{{population.not_payed>1?'s':''}} that haven't finished paying</b-card-text>
          <hr />
          <b-card-title class="font-weight-bold font-20">Respective Cells</b-card-title>
          <b-card-text v-for="(cell,i) in cellsInfo" :key="i">
            <b-card-text>
              <b>{{cell.label}}</b>
              cell is having {{cell.data.payed+cell.data.pending}} Houses with {{cell.data.payed}} House{{cell.data.payed>1?'s':''}} that finished paying and {{cell.data.pending}} House{{cell.pending>1?'s':''}} that haven't finished paying
              <hr />
            </b-card-text>
          </b-card-text>
        </b-card>
        <b-card v-if="state.error">
          <b-card-text>{{state.errorMessage}}</b-card-text>
        </b-card>
      </b-collapse>
    </b-row>
    <b-row v-if="!state.generating && !state.error && cellsInfo" class="my-3 justify-content-end">
      <b-button size="sm" class="app-color">Download Report</b-button>
    </b-row>
  </div>
</template>

<script>
export default {
  name: "sectorReports",
  data() {
    return {
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
      cellsInfo: null
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    activeSector() {
      return this.$store.getters.getActiveSector;
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
      this.state.generating = true;
      this.state.error = false;
      this.state.errorMessage = null;
      const first = this.axios.get(
        this.endpoint +
          `/metrics/ratios/sectors/${this.activeSector}?year=${this.currentYear}&month=${this.currentMonth}`
      );
      const second = this.axios.get(
        this.endpoint +
          `/metrics/ratios/sectors/all/${this.activeSector}?year=${this.currentYear}&month=${this.currentMonth}`
      );
      this.axios
        .all([first, second])
        .then(
          this.axios.spread((...res) => {
            this.population.total =
              res[0].data.data.payed + res[0].data.data.pending;
            this.population.payed = res[0].data.data.payed;
            this.population.not_payed = res[0].data.data.pending;
            this.cellsInfo = res[1].data;
            this.state.showReport = true;
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
          this.state.generating = false;
        });
    }
  }
};
</script>

<style>
</style>