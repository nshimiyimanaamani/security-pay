<template>
  <div>
    <header class="d-flex justify-content-center font-20 text-uppercase">Create village Report</header>
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
      >Generate {{village ? village : 'Village'}} Report</b-button>
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
      village: null
    };
  },
  computed: {
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
    }
  }
};
</script>

<style>
</style>