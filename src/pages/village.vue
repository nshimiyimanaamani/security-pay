<template>
  <div class="village-wrapper">
    <div class="container">
      <div class="container-title">
        <i class="fa fa-th-large"></i>
        <h1>{{activeVillage}}</h1>
        <span class="fa fa-cog"></span>
      </div>
      <b-card v-show="houses.length">
        <section v-for="(house,index) in houses" :key="index">
          <b-card-header>
            <p>{{house.owner}}</p>
            <p>{{house.id}}</p>
          </b-card-header>
          <b-card-footer>
            <article>
              <span v-show="house.percentage" class="completed">completed</span>
              <span class="details">{{house.due}} /5 last months</span>
              <b-progress :value="60" :max="100"></b-progress>
            </article>
            <i class="fa fa-ellipsis-v" v-b-toggle.collapse="''+index"></i>
          </b-card-footer>
          <b-collapse :id="'' + index" class="more-data">
            <article>
              <label for="sector">Sector:</label>
              <p>{{house.sector}}</p>
            </article>
            <article>
              <label for="cell">Cell:</label>
              <p>{{house.cell}}</p>
            </article>
            <article>
              <label for="village">village:</label>
              <p>{{house.village}}</p>
            </article>
            <article>
              <label for="due">To Pay:</label>
              <p>{{house.due}} Rwf</p>
            </article>
          </b-collapse>
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
export default {
  data() {
    return {
      state: {
        loading: true
      },
      months: [
        "January",
        "February",
        "March",
        "April",
        "May",
        "June",
        "July",
        "August",
        "September",
        "October",
        "November",
        "December"
      ],
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
        .get(this.endpoint + `/properties/sectors/remera?offset=1&limit=100`)
        .then(res => {
          this.state.loading = false;
          this.responseData = res.data.properties;
          this.filterBy_village();
        })
        .catch(err => {
          this.state.loading = false;
          console.log(err);
          return [];
        });
    },
    filterBy_village() {
      this.$nextTick(() => {
        this.houses = this.responseData.filter(data => {
          return data.village === this.activeVillage;
        });
      });
    }
  }
};
</script>
<style>
@import url("../assets/css/village.css");
</style>
