<template>
  <b-modal
    hide-footer
    title="Create Manager Account"
    content-class="primary-font"
    centered
    @hide="$emit('close')"
    v-model="show"
  >
    <div class="manager-createAccount">
      <b-form class="p-2" :class="{'blur':state.loading}" @submit.prevent="createAccount">
        <b-form-group label="E-mail">
          <b-form-input
            type="email"
            class="br-0"
            v-model="email"
            placeholder="Manager e-mail"
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
        <section class="select-group">
          <b-form-group label="Province" class="flex-grow-1">
            <b-form-select v-model="province" class="br-0" required>
              <template v-slot:first>
                <b-form-select-option value disabled>select province</b-form-select-option>
              </template>
              <b-form-select-option v-for="(item,i) in provinces" :key="i" :value="item">{{item}}</b-form-select-option>
            </b-form-select>
          </b-form-group>
          <b-form-group label="District" class="flex-grow-1">
            <b-form-select
              v-model="district"
              class="br-0"
              :disabled="!districts.length > 0"
              required
            >
              <template v-slot:first>
                <b-form-select-option value disabled>select district</b-form-select-option>
              </template>
              <b-form-select-option v-for="(item,i) in districts" :key="i" :value="item">{{item}}</b-form-select-option>
            </b-form-select>
          </b-form-group>
          <b-form-group label="Sector" class="flex-grow-1">
            <b-form-select v-model="sector" class="br-0" :disabled="!sectors.length>0" required>
              <template v-slot:first>
                <b-form-select-option value disabled>select sector</b-form-select-option>
              </template>
              <b-form-select-option v-for="(item,i) in sectors" :key="i" :value="item">{{item}}</b-form-select-option>
            </b-form-select>
          </b-form-group>
          <b-form-group label="Cell" class="flex-grow-1">
            <b-form-select v-model="cell" class="br-0" :disabled="!cells.length>0" required>
              <template v-slot:first>
                <b-form-select-option value disabled>select cell</b-form-select-option>
              </template>
              <b-form-select-option v-for="(item,i) in cells" :key="i" :value="item">{{item}}</b-form-select-option>
            </b-form-select>
          </b-form-group>
        </section>

        <b-form-group class="mb-0">
          <b-button variant="info" class="br-0 float-right" :disabled="showButton" type="submit">
            Create
            <i class="fa fa-check ml-1" />
          </b-button>
        </b-form-group>
      </b-form>
      <div class="modal-loading" v-if="state.loading">
        <i class="fa fa-spinner fa-spin" />
        <p>Creating...</p>
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
      province: "",
      district: "",
      sector: "",
      cell: "",
      accounts: [],
      state: { loading: false, loadingAccounts: false }
    };
  },
  computed: {
    showButton() {
      if (this.account && this.email && this.cell) return false;
      return true;
    },
    provinces() {
      this.province = "";
      this.district = "";
      this.sector = "";
      this.cell = "";
      return this.$provinces();
    },
    districts() {
      this.district = "";
      this.sector = "";
      this.cell = "";
      if (!this.province) return [];
      return this.$districts(this.province);
    },
    sectors() {
      this.sector = "";
      this.cell = "";
      if (!this.district) return [];
      return this.$sectors(this.province, this.district);
    },
    cells() {
      this.cell = "";
      if (!this.sector) return [];
      return this.$cells(this.province, this.district, this.sector);
    }
  },
  mounted() {
    this.show = true;
    this.getAccounts();
  },
  methods: {
    createAccount() {
      if (!this.email && !this.account && !this.cell) return;
      this.state.loading = true;
      const data = {
        account: this.account,
        email: this.email,
        cell: this.cell
      };
      this.axios
        .post("/accounts/managers", data)
        .then(res => {
          this.state.loading = false;
          this.$bvModal
            .msgBoxOk(`Password: ${res.data.password}`, {
              title: "Account created successfully!",
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
            this.$snotify(err.data.error);
          } catch {
            this.$snotify.error("Error! can't create account");
          }
          this.state.loadingAccounts = false;
        });
    }
  }
};
</script>

<style lang="scss">
.manager-createAccount {
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
  .select-group {
    width: 100%;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    grid-gap: 10px;
  }
}
</style>