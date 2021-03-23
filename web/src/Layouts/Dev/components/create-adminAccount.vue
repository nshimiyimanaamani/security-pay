<template>
  <b-modal
    hide-footer
    title="Create Administrator Account"
    content-class="primary-font"
    centered
    @hide="$emit('close')"
    v-model="show"
  >
    <div class="admin-createAccount">
      <b-form class="p-2" :class="{'blur':state.loading}" @submit.prevent="createAccount">
        <b-form-group label="E-mail">
          <b-form-input
            type="email"
            class="br-0"
            v-model="email"
            placeholder="Admin e-mail"
            trim
            required
          />
        </b-form-group>

        <b-form-group label="Account">
          <section class="d-flex align-items-center">
            <b-form-select v-model="account" class="br-0" required>
              <template v-slot:first>
                <b-form-select-option value disabled>select account</b-form-select-option>
              </template>
              <b-form-select-option
                v-for="(account,i) in accounts"
                :key="i"
                :value="account"
              >{{account}}</b-form-select-option>
            </b-form-select>
            <b-button variant="outline-info" class="ml-2" @click="getAccounts">
              <i class="fa fa-sync-alt" :class="{'fa-spin': state.loadingAccounts}" />
            </b-button>
          </section>
        </b-form-group>

        <b-form-group class="mb-0">
          <b-button variant="info" class="br-0 float-right" :disabled="showButton" type="submit">
            Create Account
            <i class="fa fa-check ml-1" />
          </b-button>
        </b-form-group>
      </b-form>

      <div class="modal-loading" v-if="state.loading">
        <i class="fa fa-spinner fa-spin" />
        <p>Creating account...</p>
      </div>
    </div>
  </b-modal>
</template>

<script>
export default {
  name: "createManagerAccount",
  data() {
    return {
      show: true,
      email: null,
      account: "",
      accounts: [],
      state: { loading: false, loadingAccounts: false }
    };
  },
  computed: {
    showButton() {
      if (this.account && this.email) return false;
      return true;
    }
  },
  mounted() {
    this.show = true;
    this.getAccounts();
  },
  methods: {
    createAccount() {
      if (!this.email && !this.account) return;
      this.state.loading = true;
      const data = { account: this.account, email: this.email };
      this.axios
        .post("/accounts/admin", data)
        .then(res => {
          this.state.loading = false;
          this.$bvModal
            .msgBoxOk(`Password: ${res.data.password}`, {
              title: "Administrator account created successfully!",
              headerClass:
                "pt-5 pb-1 border-bottom-0 justify-content-center primary-font",
              contentClass: "text-center p-0 pb-5 primary-font",
              footerClass: "d-none",
              centered: true
            })
            .then(value => {
              this.$emit("created");
              this.$emit("close");
            });
        })
        .catch(err => {
          console.log(err, err.response);
          try {
            this.$snotify.error(err.response.data.error);
          } catch {
            this.$snotify.error(
              "Administrator account coldn't be created! try again later"
            );
          }
          this.state.loading = false;
        });
    },
    async getAccounts() {
      this.state.loadingAccounts = true;
      const Total = await this.$getTotal("/accounts?offset=0&limit=0");
      this.axios
        .get("/accounts?offset=0&limit=" + Total)
        .then(async res => {
          const accounts = res.data.Accounts.map(item => item.id);
          this.accounts = await accounts.filter(
            (item, i) => accounts.indexOf(item) == i
          );
          this.state.loadingAccounts = false;
        })
        .catch(err => {
          console.log(err, err.response);
          try {
            this.$snotify.error(err.data.error);
          } catch {
            this.$snotify.error("Error! can't retrieve accounts");
          }
          this.state.loadingAccounts = false;
        });
    }
  }
};
</script>

<style lang="scss">
.admin-createAccount {
  .br-0 {
    border-radius: 1px;
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