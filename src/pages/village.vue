<template>
  <div class="village-wrapper">
    <div class="container">
      <div class="container-title">
        <i class="fa fa-th-large"></i>
        <h1>{{activeVillage}}</h1>
        <span class="fa fa-cog"></span>
      </div>
      <b-card>
        <b-card-body v-for="(house,index) in houses" :key="index">
          <b-card-body>
            <b-card-header>
              <p>{{house.owner}}</p>
              <p>{{house.id}}</p>
            </b-card-header>
            <b-card-footer>
              <b-card-group>
                <article>
                  <span v-show="house.percentage" class="completed">completed</span>
                  <span class="details">{{house.due}} /5 last months</span>
                  <b-progress :value="60" :max="100"></b-progress>
                </article>
                <i class="fa fa-ellipsis-v" v-b-toggle.collapse="''+index"></i>
              </b-card-group>
            </b-card-footer>
          </b-card-body>

          <!-- Elements to collapse -->
          <b-collapse :id="'' + index">
            <b-card-body>
              <b-card-text>
                Anim pariatur cliche reprehenderit, enim eiusmod high life accusamus terry
                richardson ad squid. 3 wolf moon officia aute, non cupidatat skateboard dolor
                brunch. Food truck quinoa nesciunt laborum eiusmod. Brunch 3 wolf moon
                tempor, sunt aliqua put a bird on it squid single-origin coffee nulla
                assumenda shoreditch et. Nihil anim keffiyeh helvetica, craft beer labore
                wes anderson cred nesciunt sapiente ea proident. Ad vegan excepteur butcher
                vice lomo. Leggings occaecat craft beer farm-to-table, raw denim aestheti
              </b-card-text>
            </b-card-body>
          </b-collapse>
        </b-card-body>
      </b-card>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
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
  mounted() {
    this.axios
      .get(this.endpoint + `/properties/sectors/remera?offset=1&limit=100`)
      .then(res => {
        console.log(res.data);
        this.houses = res.data.properties;
      })
      .catch(err => console.log(err));
  }
};
</script>
<style>
@import url("../assets/css/village.css");
</style>
