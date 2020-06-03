<template>
  <b-modal
    v-model="show"
    hide-footer
    title="Create Account for Developer"
    content-class="secondary-font"
    body-class="createDevAccount-modal"
    centered
    @hide="$emit('close')"
  >
    <b-form @submit.prevent="createAccount" :class="{'blur':loading}">
      <b-form-group label="E-mail">
        <b-form-input
          type="email"
          v-model="email"
          placeholder="Developer's E-mail address"
          trim
          required
        />
      </b-form-group>
      <b-form-group label="Password">
        <b-input
          type="password"
          v-model="password"
          placeholder="Developer's password"
          trim
          required
        />
      </b-form-group>
      <b-form-group label="Account">
        <b-row class="m-0 align-items-center flex-nowrap">
          <b-select v-model="account" class="flex-grow-1" :disabled="loadingAccounts">
            <template v-slot:first>
              <b-form-select-option value disabled>select an account</b-form-select-option>
            </template>
            <b-form-select-option
              v-for="(account,i) in accounts"
              :key="i"
              :value="account"
            >{{account}}</b-form-select-option>
          </b-select>
          <i
            class="fa fa-sync-alt mr-2 ml-3 cursor-pointer"
            :class="{'fa-spin': loadingAccounts}"
            @click="getAccounts"
          />
        </b-row>
      </b-form-group>
      <b-form-group class="mb-0">
        <b-button variant="info" class="float-right mt-3" type="submit" :disabled="disableButton">
          Create
          <i class="fa fa-check ml-1" />
        </b-button>
      </b-form-group>
    </b-form>
    <div v-if="loading" class="modal-loading">
      <vue-load label="Creating..." />
    </div>
  </b-modal>
</template>

<script>
export default {
  name: "CreateDevAccount",
  data() {
    return {
      show: true,
      account: "",
      email: null,
      password: null,
      loadingAccounts: false,
      loading: false,
      accounts: []
    };
  },
  computed: {
    disableButton() {
      if (this.account && this.email && this.password) return false;
      return true;
    }
  },
  beforeMount() {
    this.getAccounts();
  },
  methods: {
    async createAccount() {
      if (!this.email || !this.password || !this.account) return;
      this.loading = true;
      const data = {
        account: this.account,
        email: this.email,
        password: this.password
      };
      this.axios
        .post("/accounts/developers", data)
        .then(res => {
          console.log(res.data);
          this.$snotify.success("Account created successfully!");
          this.loading = false;
          this.$emit("created");
        })
        .catch(err => {
          console.log(err, err.response);
          try {
            this.$snotify.error(err.response.data.error);
          } catch {
            this.$snotify.error(
              "Oops! can't create account at this moment, try again later!"
            );
          }
          this.loading = false;
        });
    },
    async getAccounts() {
      this.loadingAccounts = true;
      const Total = await this.getTotal("/accounts?offset=0&limit=0");
      this.axios
        .get("/accounts?offset=0&limit=" + Total)
        .then(async res => {
          console.log(res);
          this.accounts = res.data.Accounts.map(item => item.id);
          await this.accounts;
          this.loadingAccounts = false;
        })
        .catch(err => {
          console.log(err, err.response);
          this.loadingAccounts = false;
        });
    },
    getTotal(endpoint) {
      return this.axios
        .get(endpoint)
        .then(res => res.data.Total)
        .catch(err => 0);
    }
  }
};
</script>

<style lang="scss">
.createDevAccount-modal {
  position: relative;
  .blur {
    filter: blur(3px);
  }
  .modal-loading {
    position: absolute;
    top: 0;
    bottom: 0;
    right: 0;
    left: 0;
    height: 100%;
    width: 100%;
    background: #00000080;
    color: white;
    display: flex;
  }
}
</style>