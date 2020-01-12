<template>
  <div class="village-wrapper">
    <vue-title title="Paypack | Village" />
    <div class="container">
      <div class="container-title">
        <i class="fa fa-th-large"></i>
        <h1>{{activeVillage}}</h1>
        <span class="fa fa-cog"></span>
      </div>
      <b-card v-show="houses.length">
        <section v-for="(house,index) in houses" :key="index">
          <user-card :house="house" :index="index" />
        </section>
      </b-card>
      <section class="error" v-show="!houses.length">
        <article>
          <center>
            <div v-if="!state.loading">
              <i class="fa fa-exclamation-triangle"></i>
              <label for="error">No House found in {{activeVillage}}</label>
            </div>
            <div v-else-if="state.loading">
              <b-spinner variant="primary" small type="grow" label="Spinning"></b-spinner>
              <label>Loading...</label>
            </div>
          </center>
        </article>
      </section>
    </div>
  </div>
</template>

<script>
import userCard from "../components/usercard.vue";
export default {
  name: "village",
  components: {
    "user-card": userCard
  },
  data() {
    return {
      state: {
        loading: true,
        showCollapse: false,
        elId: null
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
        this.filterBy_village();
      }
    }
  },
  mounted() {
    this.loadData();
  },
  methods: {
    loadData() {
      this.state.loading = true;
      this.axios
        .get(this.endpoint + `/properties?sector=Remera&offset=1&limit=100`)
        .then(res => {
          this.state.loading = false;
          this.responseData = res.data.Properties;
          this.filterBy_village();
          console.log(res.data.Properties);
        })
        .catch(err => {
          if (navigator.onLine) {
            const error = isNullOrUndefined(err.response)
              ? "an error occured"
              : err.response.data.error || err.response.data;
            this.$snotify.error(error);
          } else {
            this.$snotify.error("Please connect to the internet");
          }
          this.state.loading = false;
          return [];
        });
    },
    filterBy_village() {
      this.$nextTick(() => {
        this.houses = this.responseData.filter(data => {
          return (
            data.address.village.toLowerCase() ==
            this.activeVillage.toLowerCase()
          );
        });
      });
    }
  }
};
</script>
<style>
@import url("../assets/css/village.css");
</style>
