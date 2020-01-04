<template>
  <div id="agentView">
    <vue-title title="Paypack | Agents" />
    <nav class="py-2 mb-4">
      <b-container class="p-0">
        <b-row>
          <b-col lg="10" md="10" sm="11" xl="10" class="navs pl-4">
            <h2 class="m-0 text-white">P A Y P A C K</h2>
          </b-col>
          <b-col lg="1" md="1" sm="1" xl="1" class="navs mr-2" style="margin-left:auto">
            <b-button variant="info" @click="logout">Logout</b-button>
          </b-col>
        </b-row>
      </b-container>
    </nav>
    <user-table />
  </div>
</template>

<script>
import home from "../components/agents/home.vue";

export default {
  name: "agentsView",
  components: {
    "user-table": home
  },
  data() {
    return {};
  },
  methods: {
    logout() {
      this.$store.dispatch("logout");
    }
  }
};
</script>

<style lang="scss">
#agentView {
  font-family: "Raleway", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  width: 100vw;
  height: 100vh;
  overflow-x: hidden;
  nav {
    position: sticky;
    width: 100%;
    background: #3a82a1;
    .navs {
      width: fit-content;
      button {
        border: 1px solid white !important;
      }
    }
  }
}
</style>
