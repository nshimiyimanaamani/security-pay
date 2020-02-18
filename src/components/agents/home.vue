<template>
  <b-container class="h-auto mw-100">
    <b-row class="px-4">
      <controller :user="user" v-on:refresh="key++" />
    </b-row>
    <b-row class="px-4">
      <user-table v-on:getInfo="getInfo" :user="user" :key="key" />
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
  },
  data() {
    return {
      user: null,
      key: 0
    };
  },
  computed: {
    userDetails() {
      return this.$store.getters.userDetails;
    }
  },
  created() {
    this.getInfo();
  },
  methods: {
    getInfo() {
      this.axios
        .get("/accounts/agents/" + this.userDetails.username)
        .then(res => {
          this.user = { ...res.data };
          this.key++;
        })
        .catch(err => {
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
        });
    }
  }
};
</script>

<style>
.table-control {
  overflow: auto;
  width: 100%;
}
</style>