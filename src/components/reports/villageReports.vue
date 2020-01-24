<template>
  <div>
    <header class="d-flex justify-content-center font-20 text-uppercase">village Report</header>
    <hr class="m-0 mb-3" />
    <b-row class="px-3 align-items-center justify-content-between">
      <b-select
        size="sm"
        id="input-1"
        v-model="cell"
        :options="cellOptions"
        class="w-auto mr-2 flex-grow-1"
      >
        <template v-slot:first>
          <option :value="null" disabled>Please select cell</option>
        </template>
      </b-select>
      <b-select
        size="sm"
        id="input-1"
        v-model="village"
        :options="villageOptions"
        class="w-auto mr-2 flex-grow-1"
      >
        <template v-slot:first>
          <option :value="null" disabled>Please select village</option>
        </template>
      </b-select>
      <b-button
        size="sm"
        variant="info"
        class="font-15 border-0 my-3"
        :disabled="village?false:true"
        @click="generate"
      >Generate {{village ? village : 'Village'}} Report</b-button>
      <div v-show="state.generating" class="w-100">
        <strong class="font-15">Generating&nbsp;</strong>
        <b-spinner small />
      </div>
    </b-row>
    <b-row>
      <b-collapse id="sectorreport-collapse" class="w-100" v-model="state.showReport">
        <b-card class="text-capitalize mx-3" v-if="!state.error">
          <b-card-title class="font-weight-bold font-20">village Overall</b-card-title>
          <b-card-text>
            <b>{{village}}</b>
            village is having {{population.total}} Houses with {{population.payed}} House{{population.payed>1?'s':''}} that finished paying and {{population.not_payed}} House{{population.not_payed>1?'s':''}} that haven't finished paying
          </b-card-text>
        </b-card>
        <b-card v-if="state.error" class="mx-3">
          <b-card-text>{{state.errorMessage}}</b-card-text>
        </b-card>
      </b-collapse>
    </b-row>
    <b-row
      v-if="!state.generating && !state.error && population.total"
      class="my-3 justify-content-end"
    >
      <b-button size="sm" class="app-color mx-3">Download Report</b-button>
    </b-row>
  </div>
</template>

<script>
const { Village } = require("rwanda");
export default {
  name: "cellReports",
  data() {
    return {
      cell: null,
      village: null,
      state: {
        showReport: false,
        generating: false,
        error: false,
        errorMessage: null
      },
      population: {
        total: null,
        payed: null,
        not_payed: null
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    activeSector() {
      return this.$store.getters.getActiveSector;
    },
    cellOptions() {
      return this.$store.getters.getCellsArray;
    },
    villageOptions() {
      if (this.cell) {
        return Village("Kigali", "Gasabo", this.activeSector, this.cell).sort();
      } else {
        return [];
      }
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
      this.clean();
      this.state.generating = true;
      let axios = this.axios;
      axios
        .get(
          this.endpoint +
            `/metrics/ratios/villages/${this.village}?year=${this.currentYear}&month=${this.currentMonth}`
        )
        .then(res => {
          this.population.total = res.data.data.payed + res.data.data.pending;
          this.population.payed = res.data.data.payed;
          this.population.not_payed = res.data.data.pending;
        })
        .catch(err => {
          this.clean();
          this.state.error = true;
          this.state.errorMessage = err.response.data.error
            ? err.response.data.error
            : err.response.response;
        })
        .finally(() => {
          this.state.generating = false;
          this.state.showReport = true;
        });
    },
    clean() {
      this.state = {
        showReport: false,
        generating: false,
        error: false,
        errorMessage: null
      };
      this.population = {
        total: null,
        payed: null,
        not_payed: null
      };
    }
  }
};
</script>

<style>
</style>