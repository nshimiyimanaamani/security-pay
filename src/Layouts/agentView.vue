<template>
  <div id="agentView">
    <vue-title title="Paypack | Agents" />
    <nav class="py-2" style="z-index:1000">
      <b-container class="p-0">
        <b-row>
          <b-col lg="10" md="10" sm="11" xl="10" class="navs pl-4">
            <h2 class="m-0 text-white">P A Y P A C K</h2>
          </b-col>
          <b-col lg="1" md="1" sm="1" xl="1" class="navs mr-2" style="margin-left:auto">
            <b-dropdown variant="info" no-caret>
              <template v-slot:button-content>
                <i class="fa fa-cog" />
              </template>
              <b-dropdown-item v-b-modal.changePassword>Change Password</b-dropdown-item>
              <b-dropdown-item @click="logout">Logout</b-dropdown-item>
            </b-dropdown>
          </b-col>
        </b-row>
      </b-container>
    </nav>
    <user-table />
    <b-modal
      id="changePassword"
      title="Change your Password"
      button-size="sm"
      @ok="change"
      @close="close"
    >
      <template v-slot:modal-footer="{ok}">
        <b-button
          size="sm"
          variant="info"
          :disabled="newPassword.length>5?false:true"
          class="px-3"
          @click="ok()"
        >
          OK
          <b-spinner v-show="state.changing" small type="grow" />
        </b-button>
      </template>
      <b-row>
        <b-form class="w-100 px-3">
          <b-form-group id="input-group-1" label="New Password" label-for="input-1">
            <b-input size="sm" id="input-1" v-model="newPassword" placeholder="Enter new password" />
          </b-form-group>
        </b-form>
      </b-row>
    </b-modal>
  </div>
</template>

<script>
import home from "../components/agents/home.vue";
import loader from "../components/loader";

export default {
  name: "agentsView",
  components: {
    "user-table": home,
    loader
  },
  data() {
    return {
      newPassword: "",
      state: {
        changing: false
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    user() {
      return this.$store.getters.userDetails;
    }
  },
  methods: {
    logout() {
      this.$store.dispatch("logout");
    },
    change(e) {
      this.state.changing = true;
      const password = this.newPassword;
      this.axios
        .put(`/accounts/agents/creds/${this.user.username}`, {
          password: password
        })
        .then(res => {
          this.$snotify.info("Password changed!");
          this.close();
        })
        .catch(err => {
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
          this.close();
        });
    },
    close() {
      this.newPassword = "";
      this.state.changing = false;
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
