<template>
  <b-container style="height: auto">
    <b-row>
      <controller />
    </b-row>
    <b-row>
      <b-col sm="12" md="12" lg="12" class="px-1 py-0 table-control">
        <user-table />
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import userTable from "./table.vue";
import controllers from "./controllers.vue";
export default {
  name: "home",
  components: {
    "user-table": userTable,
    controller: controllers
  }
};
</script>

<style>
.table-control {
  overflow: auto;
  width: 100%;
}
</style>