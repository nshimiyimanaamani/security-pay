<template>
  <b-container style="height: auto">
    <b-row>
      <controller :user="user" v-on:refresh="key++" />
    </b-row>
    <b-row>
      <b-col sm="12" md="12" lg="12" class="px-1 py-0 table-control">
        <user-table v-on:getInfo="getInfo" :user="user" :key="key" />
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
  },
  data() {
    return {
      user: null,
      key: 0
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    userDetails() {
      return this.$store.getters.userDetails;
    }
  },
  methods: {
    getInfo() {
      this.axios
        .get(this.endpoint + "/accounts/agents/" + this.userDetails.username)
        .then(res => {
          this.user = { ...res.data };
          this.key++;
        })
        .catch(err => {
          const error = navigator.onLine
            ? err.response.data.error || err.response.data
            : "Please connect to the internet";
          console.log(err);
          this.$snotify.error(error);
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