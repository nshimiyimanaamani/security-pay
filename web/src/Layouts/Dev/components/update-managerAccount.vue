<template>
  <b-modal
    hide-footer
    title="Update Account"
    content-class="primary-font"
    centered
    @hide="$emit('close')"
    v-model="show"
  >
    <div class="manager-updateAccount">
      <b-form class="p-2" :class="{'blur':state.loading}" @submit.prevent="updateAccount">
        <b-form-group label="E-mail">
          <b-form-input
            type="email"
            class="br-0"
            v-model="newEmail"
            placeholder="new E-mail"
            trim
            required
          />
        </b-form-group>
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
          <b-button variant="info" class="br-0 float-right" :disabled="showButton" type="submit">
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
  </b-modal>
</template>

<script>
export default {
  name: "devAccountUpdateModal",
  props: {
    account: {
      type: Object,
      required: false,
      default: null
    }
  },
  data() {
    return {
      show: true,
      newPassword: null,
      newEmail: null,
      state: { loading: false }
    };
  },
  computed: {
    showButton() {
      if (this.newPassword && this.newEmail) return false;
      return true;
    },
    accountId() {
      return this.account.email;
    }
  },
  mounted() {
    this.show = true;
    this.newEmail = this.account.email;
  },
  methods: {
    updateAccount() {
      if (!this.newPassword && !this.newEmail) return;
      this.state.loading = true;
      const data = {
        email: this.newEmail,
        password: this.newPassword
      };
      this.axios
        .put("/accounts/managers/creds/" + this.accountId, data)
        .then(res => {
          this.state.loading = false;
          this.$snotify.success(res.data.message);
          this.$emit("updated");
          this.$emit("close");
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
.manager-updateAccount {
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