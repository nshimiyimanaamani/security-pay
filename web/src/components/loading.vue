<template>
  <div class="loading-component">
    <i class="fa fa-spinner fa-spin" />
    <p>{{label}}</p>
  </div>
</template>

<script>
export default {
  name: "LoadingComponent",
  props: {
    label: {
      type: String,
      default: "Loading..."
    }
  }
};
</script>

<style lang="scss">
.loading-component {
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 4rem;
  user-select: none;
  background-color: #f8f9fa;

  i {
    font-size: 2rem;
    margin-right: 0.5rem;
  }
  p {
    font-size: 1.2rem;
    text-transform: capitalize;
    margin: 0 !important;
  }
}
</style>