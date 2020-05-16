<template>
  <div id="agentView">
    <vue-title title="Paypack | Agents" />
    <nav style="z-index:1000">
      <h2 class="m-0 text-white text-truncate">P A Y P A C K</h2>
      <b-dropdown variant="info" no-caret>
        <template v-slot:button-content>
          <i class="fa fa-bars" />
        </template>
        <b-dropdown-item v-b-modal.changePassword>Change Password</b-dropdown-item>
        <b-dropdown-item @click="logout">Logout</b-dropdown-item>
      </b-dropdown>
    </nav>
    <agent-home />
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
    "agent-home": home,
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
  font-family: "Montserrat", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  width: 100vw;
  height: 100vh;
  overflow-x: hidden;
  min-width: 230px;
  nav {
    position: sticky;
    top: 0;
    width: 100%;
    background: #017db3;
    height: 60px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 0.5rem;

    button {
      background: transparent !important;
      border: none !important;
      font-size: 25px;
    }
  }
}
</style>
