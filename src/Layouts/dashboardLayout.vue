<template>
  <div class="dashboardWrapper">
    <div class="dashboardSidebar">
      <h1>P A Y P A C K</h1>
      <ul class="sidebarLinks">
        <router-link to="/dashboard">
          <li>Overview</li>
        </router-link>
        <router-link to="/cells">
          <li>
            Cells
            <b-dropdown>
              <b-dropdown-item
                @click="update({toUpdate: 'cell', changed: cell})"
                v-for="cell in sidebar.cells_array"
                :key="cell"
              >{{cell}}</b-dropdown-item>
            </b-dropdown>
          </li>
        </router-link>
        <router-link to="/village">
          <li>
            Village
            <b-dropdown>
              <b-dropdown-item
                @click="update({toUpdate: 'village', changed: village})"
                v-for="village in sidebar.village_array"
                :key="village"
              >{{village}}</b-dropdown-item>
            </b-dropdown>
          </li>
        </router-link>
        <router-link to="/transactions">
          <li>Bank Accounts</li>
        </router-link>
        <router-link to="/reports">
          <li>Properties</li>
        </router-link>
        <router-link to="#">
          <li>Penalties</li>
        </router-link>
      </ul>
      <p class="text-center powered" for="powered">
          Powered By
          <strong>Quarks Group.</strong>
        </p>
    </div>
    <div class="rightSide">
      <div class="top-nav">
        <div class="logout">
          <b-button class="btn-info" @click.prevent="logout">Logout</b-button>
        </div>
      </div>
      <div class="dashboardBody">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      userAvailable: false,
      sidebar: {
        active_cell: this.getActiveCell,
        cells_array: [],
        active_village: "",
        village_array: []
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    getActiveCell() {
      this.sidebar.active_cell = this.$store.getters.getActiveCell;
      return this.$store.getters.getActiveCell;
    },
    getActiveVillage() {
      this.sidebar.active_village = this.$store.getters.getActiveVillage;
      return this.$store.getters.getActiveVillage;
    },
    getCellsArray() {
      this.sidebar.cells_array = this.$store.getters.getCellsArray;
      return this.$store.getters.getCellsArray;
    },
    getVillageArray() {
      this.sidebar.village_array = this.$store.getters.getVillageArray;
      return this.$store.getters.getVillageArray;
    },
    getSector() {
      return this.$store.getters.getSectorArray;
    },
    getPropertyCell() {
      return this.cells();
    },
    getPropertyVillage() {
      return this.village();
    }
  },
  mounted() {
    this.getActiveCell;
    this.getCellsArray;
    this.getActivevillage;
    this.getVillageArray;
  },
  methods: {
    update(res) {
      this.$store.dispatch("updatePlace", res).then(() => {
        this.getActiveCell;
        this.getCellsArray;
        this.getActiveVillage;
        this.getVillageArray;
      });
    },
    logout() {
      this.$store.dispatch("logout");
    }
  }
};
</script>
<style scoped>
@import url("../assets/css/dashboardLayout.css");
</style>