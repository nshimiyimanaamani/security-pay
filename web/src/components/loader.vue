<template>
  <div v-if="loading">
    <b-spinner variant="black" small />&nbsp;
    <strong>Loading...</strong>
  </div>
</template>

<script>
export default {
  name: "loader",
  props: {
    loading: {
      type: Boolean,
      default: true
    }
  }
};
</script>

<style>
</style>