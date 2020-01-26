<template>
  <div class="village-wrapper">
    <vue-title title="Paypack | Village" />
    <div class="container max-width">
      <div class="container-title">
        <i class="fa fa-th-large"></i>
        <h1>{{activeVillage}}</h1>
        <span class="fa fa-cog"></span>
      </div>
      <b-card body-bg-variant="white" v-show="houses.length" class="border-top-0 rounded-0">
        <b-card-group columns>
          <b-card body-bg-variant="white" no-body v-for="(house,index) in houses" :key="index">
            <user-card :house="house" :index="index" />
          </b-card>
        </b-card-group>
      </b-card>
      <section class="error" v-show="!houses.length">
        <article>
          <center>
            <div v-if="!state.loading">
              <i class="fa fa-exclamation-triangle"></i>
              <label for="error">No House found in {{activeVillage}}</label>
            </div>
            <loader :loading="state.loading" />
          </center>
        </article>
      </section>
    </div>
  </div>
</template>

<script>
import userCard from "../components/usercard.vue";
import loader from "../components/loader.vue";
export default {
  name: "village",
  components: {
    "user-card": userCard,
    loader
  },
  data() {
    return {
      state: {
        loading: true,
        showCollapse: false
      },
      houses: [],
      responseData: [],
      green: "#50a031",
      orange: "#f0a700",
      red: "#f3573c"
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    activeVillage() {
      return this.$store.getters.getActiveVillage;
    },
    activeCell() {
      return this.$store.getters.getActiveCell;
    }
  },
  watch: {
    activeVillage() {
      handler: {
        this.loadData();
      }
    }
  },
  mounted() {
    this.loadData();
  },
  methods: {
    loadData() {
      this.state.loading = true;
      this.houses = new Array();
      this.axios
        .get(
          this.endpoint +
            `/properties?village=${this.activeVillage}&offset=0&limit=1000`
        )
        .then(res => {
          this.houses = res.data.Properties;
        })
        .catch(err => {
          if (navigator.onLine) {
            const error = err.response
              ? err.response.data.error || err.response.data
              : "an error occured";
            this.$snotify.error(error);
          } else {
            this.$snotify.error("Please connect to the internet");
          }
          return [];
        })
        .finally(() => {
          this.state.loading = false;
        });
    }
  }
};
</script>
<style>
@import url("../assets/css/village.css");
</style>
