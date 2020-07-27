<template>
  <div class="village-wrapper w-100" style="min-width: 700px">
    <vue-title title="Paypack | Village" />
    <div class="container mw-100 mt-3">
      <div
        class="container-title d-flex justify-content-between align-items-center"
        style="height: 40px"
      >
        <i class="fa fa-th-large width-3"></i>
        <h1 class="m-0 width-3 text-center">{{activeVillage}}</h1>
        <span class="d-flex align-items-center width-3 position-relative justify-content-end">
          <transition name="fade">
            <div class="w-100" v-if="!state.showSearch">
              <b-input
                type="search"
                size="sm"
                class="left br-2 primary-font"
                placeholder="search user..."
                v-model="search.name"
              />
            </div>
          </transition>
          <i
            class="fa fa-search text-white ml-2 cursor-pointer"
            aria-hidden="true"
            @click="state.showSearch = !state.showSearch"
          />
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
        <div v-if="!state.loading" class="village-empty">
          <p class="primary-font">No House found in {{activeVillage}} village</p>
        </div>
        <vue-load v-if="state.loading" class="primary-font" />
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
        showSearch: false
      },
      search: {
        name: ""
      },
      houses: [],
      responseData: []
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
    },
    "search.name"() {
      handler: {
        const searchedName = this.search.name.toLowerCase();
        this.houses = this.responseData.filter(item => {
          const name = String(
            item.owner.fname + " " + item.owner.lname
          ).toLowerCase();
          return name.includes(searchedName);
        });
      }
    },
    responseData() {
      handler: {
        this.houses.length = 0;
        this.houses = this.responseData;
      }
    }
  },
  mounted() {
    this.loadData();
  },
  methods: {
    async loadData() {
      this.state.loading = true;
      this.responseData = new Array();
      var total = await this.getTotal();
      this.axios
        .get(
          `/properties?village=${this.activeVillage}&offset=0&limit=${total}`
        )
        .then(res => {
          this.responseData = res.data.Properties;
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
    },
    getTotal() {
      return this.axios
        .get(`/properties?village=${this.activeVillage}&offset=0&limit=0`)
        .then(res => {
          return res.data.Total;
        })
        .catch(err => {
          return 0;
        });
    }
  }
};
</script>
<style lang="scss">
@import "../assets/css/village.scss";
</style>
