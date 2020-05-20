<template>
  <div class="dev-updateAccount">
    <b-form class="p-2" :class="{'blur':state.loading}">
      <b-form-group label="Password">
        <b-form-input
          class="br-0"
          type="password"
          v-model="newPassword"
          placeholder="new Password"
          trim
          required
        />
      </b-form-group>
      <b-form-group class="mb-0">
        <b-button
          variant="info"
          class="br-0 float-right"
          :disabled="showButton"
          @click="updateAccount"
        >
          Update
          <i class="fa fa-check ml-1" />
        </b-button>
      </b-form-group>
    </b-form>
    <div class="modal-loading" v-if="state.loading">
      <i class="fa fa-spinner fa-spin" />
      <p>Updating...</p>
    </div>
  </div>
</template>

<script>
export default {
  name: "devAccountUpdateModal",
  props: {
    account: {
      type: Object,
      required: true,
      default: null
    }
  },
  data() {
    return {
      newPassword: null,
      state: { loading: false }
    };
  },
  computed: {
    showButton() {
      if (this.newPassword) return false;
      return true;
    },
    accountId() {
      return this.account.email;
    }
  },
  methods: {
    updateAccount() {
      this.state.loading = true;
      const data = {
        password: this.newPassword
      };
      this.axios
        .put("/accounts/developers/creds/" + this.accountId, data)
        .then(res => {
          this.state.loading = false;
          this.$snotify.success(res.data.message);
          this.$emit("updated");
        })
        .catch(err => {
          console.log(err, err.response);
          this.state.loading = false;
        });
    }
  }
};
</script>

<style lang="scss">
.dev-updateAccount {
  .br-0 {
    border-radius: 2px;
  }
  .blur {
    filter: blur(3px);
  }
  .modal-loading {
    color: white;
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    margin: auto;
    width: 100%;
    height: 100%;
    background: #00000080;
    display: flex;
    justify-content: center;
    align-items: center;
    i {
      font-size: 2rem;
      margin-right: 0.5rem;
    }
    p {
      margin-bottom: 0 !important;
      font-size: 1.2rem;
    }
  }
}
</style>