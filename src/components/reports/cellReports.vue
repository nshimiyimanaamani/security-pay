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
        class="font-15 border-0 my-3"
        :disabled="cell?false:true"
      >Generate {{cell ? cell : 'Cell'}} Report</b-button>
    </b-row>
  </div>
</template>

<script>
export default {
  name: "cellReports",
  data() {
    return {
      cell: null
    };
  },
  computed: {
    cellOptions() {
      return this.$store.getters.getCellsArray;
    }
  }
};
</script>

<style>
</style>