<template>
  <div class="village-wrapper w-100" style="min-width: 700px">
    <vue-title title="Paypack | Village" />
    <div class="container mw-100">
      <div class="container-title d-flex justify-content-between align-items-center">
        <i class="fa fa-th-large width-3"></i>
        <h1 class="m-0 width-3 text-center">{{activeVillage}}</h1>
        <span class="d-flex flex-row-reverse align-items-center width-3 position-relative">
          <transition name="fade">
            <div v-if="state.showSearch">
              <b-input
                type="search"
                size="sm"
                style="left:-50px;top:-3px"
                class="position-absolute font-13 app-font left"
                placeholder="search user..."
              />
            </div>
          </transition>

          <i class="fa fa-cog cursor-pointer" @click="state.showSearch = !state.showSearch" />
        </span>
      </div>
      <b-card
        body-bg-variant="white"
        v-show="houses.length"
        class="border-top-0 rounded-0 village-body"
      >
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
              <label for="error" class="font-14">No House found in {{activeVillage}} village</label>
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
        showCollapse: false,
        showSearch: false
      },
      houses: [],
      responseData: [],
      green: "#50a031",
      orange: "#f0a700",
      red: "#f3573c"
    };
  },
  computed: {
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
        .get(`/properties?village=${this.activeVillage}&offset=0&limit=1000`)
        .then(res => {
          this.houses = res.data.Properties;
        })
        .catch(err => {
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
          return [];
        })
        .finally(() => {
          this.state.loading = false;
        });
    }
  }
};
</script>
<style lang="scss">
@import "../assets/css/village.scss";
</style>
