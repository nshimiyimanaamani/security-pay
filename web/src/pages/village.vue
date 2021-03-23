<template>
  <div class="village-wrapper w-100" style="min-width: 700px">
    <vue-title title="Paypack | Village" />
    <form class="location-selector" @submit.prevent="loadData">
      <b-form-group class="control">
        <b-form-select
          class="br-2"
          v-model="select.cell"
          v-if="!isManager"
          :options="cellsOptions"
          required
        >
          <template v-slot:first>
            <option :value="null" disabled>select cell</option>
          </template>
        </b-form-select>
        <b-form-select class="br-2" v-model="select.village" :options="villageOptions" required>
          <template v-slot:first>
            <option :value="null" disabled>select village</option>
          </template>
        </b-form-select>
        <b-button type="submit" variant="info" :disabled="state.loading">Go</b-button>
      </b-form-group>
    </form>
    <div class="village-container mw-100 mt-3">
      <div class="container-title">
        <i class="fa fa-th-large width-3"></i>
        <h1 class="text-center">{{villageName}}</h1>
        <div class="control">
          <b-input
            type="search"
            size="sm"
            class="left br-2 primary-font"
            placeholder="search user..."
            v-model="search.name"
          />
          <b-button :disabled="state.loading">Refresh</b-button>
        </div>
      </div>
      <b-card body-bg-variant="white" v-show="houses.length" class="border-top-0 village-body">
        <b-card-group columns>
          <b-card body-bg-variant="white" no-body v-for="(house,index) in houses" :key="index">
            <user-card :house="house" :index="index" />
          </b-card>
        </b-card-group>
      </b-card>
      <section class="error" v-show="!houses.length">
        <div v-if="!state.loading" class="village-empty">
          <p class="primary-font">No House found in {{villageName}} village</p>
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
    "user-card": userCard,
  },
  data() {
    return {
      state: {
        loading: true,
      },
      search: {
        name: "",
      },
      houses: [],
      responseData: [],
      select: {
        cell: null,
        village: null,
      },
      villageName: "",
    };
  },
  computed: {
    activeVillage() {
      return this.$store.getters.getActiveVillage;
    },
    activeCell() {
      return this.$store.getters.getActiveCell;
    },
    location() {
      return this.$store.getters.location;
    },
    cellsOptions() {
      const { province, district, sector } = this.location;
      return this.$cells(province, district, sector);
    },
    villageOptions() {
      const { province, district, sector } = this.location;
      return this.$villages(province, district, sector, this.select.cell);
    },
    user() {
      return this.$store.getters.userDetails;
    },
    isManager() {
      return this.user.role.toLowerCase() === "basic";
    },
  },
  watch: {
    "search.name"() {
      handler: {
        const searchedName = this.search.name.toLowerCase();
        this.houses = this.responseData.filter((item) => {
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
    },
  },
  async mounted() {
    this.select.cell =
      (await this.location) && this.location.cell
        ? this.location.cell
        : this.activeCell;
    this.select.village =
      (await this.location) && this.location.village
        ? this.location.village
        : this.activeVillage;
    await this.$set(
      this,
      "villageName",
      this.select.village ? this.select.village : this.activeVillage
    );
    this.loadData();
  },
  methods: {
    async loadData() {
      await this.$set(
        this,
        "villageName",
        this.select.village ? this.select.village : this.activeVillage
      );
      this.state.loading = true;
      this.responseData = new Array();
      var total = await this.getTotal();
      this.axios
        .get(`/properties?village=${this.villageName}&offset=0&limit=${total}`)
        .then((res) => {
          this.responseData = res.data.Properties;
        })
        .catch((err) => {
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
        .get(`/properties?village=${this.villageName}&offset=0&limit=0`)
        .then((res) => {
          return res.data.Total;
        })
        .catch((err) => {
          return 0;
        });
    },
  },
};
</script>
<style lang="scss">
@import "../assets/css/village.scss";
</style>
