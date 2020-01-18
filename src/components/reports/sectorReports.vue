<template>
  <div class="px-4">
    <header class="d-flex justify-content-center font-20 text-uppercase">Create Sector Report</header>
    <hr class="m-0 mb-3" />
    <b-row>
      <b-button
        variant="info"
        class="font-15 border-0 my-2 py-2 mr-3"
      >Generate {{activeSector}} Report</b-button>
    </b-row>
  </div>
</template>

<script>
export default {
  name: "sectorReports",
  data() {
    return {};
  },
  computed: {
    activeSector() {
      return this.$store.getters.getActiveSector;
    }
  }
};
</script>

<style>
</style>