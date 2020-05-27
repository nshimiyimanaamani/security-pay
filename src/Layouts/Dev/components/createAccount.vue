<template>
  <b-modal
    v-model="show"
    hide-footer
    title="Create Account"
    content-class="secondary-font"
    centered
    @hide="modalClosed"
  >
    <div class="dev-createAccount-modal">
      <b-form class="p-2" :class="{'blur':state.loading}" @submit.prevent="createAccount">
        <b-form-group label="Account Name">
          <b-input class="br-0" v-model="name" placeholder="Account name" required />
        </b-form-group>
        <b-form-group label="Account Type">
          <b-form-select v-model="type" class="br-0" required>
            <template v-slot:first>
              <b-form-select-option disabled value>Select account Type</b-form-select-option>
            </template>
            <b-form-select-option value="ben">Beneficiary</b-form-select-option>
            <b-form-select-option value="dev">Developers</b-form-select-option>
          </b-form-select>
        </b-form-group>
        <b-form-group label="Identifier" v-if="type == 'dev'">
          <b-input-group prepend="paypack.">
            <b-input class="br-0" v-model="id" :required="type=='dev'" />
          </b-input-group>
        </b-form-group>
        <section class="select-group" v-else>
          <b-form-group label="Province" class="flex-grow-1">
            <b-form-select v-model="province" class="br-0" :required="type=='ben'">
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
              :required="type=='ben'"
            >
              <template v-slot:first>
                <b-form-select-option value disabled>select district</b-form-select-option>
              </template>
              <b-form-select-option v-for="(item,i) in districts" :key="i" :value="item">{{item}}</b-form-select-option>
            </b-form-select>
          </b-form-group>
          <b-form-group label="Sector" class="flex-grow-1">
            <b-form-select
              v-model="sector"
              class="br-0"
              :disabled="!sectors.length>0"
              :required="type=='ben'"
            >
              <template v-slot:first>
                <b-form-select-option value disabled>select sector</b-form-select-option>
              </template>
              <b-form-select-option v-for="(item,i) in sectors" :key="i" :value="item">{{item}}</b-form-select-option>
            </b-form-select>
          </b-form-group>
        </section>

        <b-form-group label="Number of allowed users">
          <b-input
            class="br-0"
            type="number"
            v-model="num_of_seats"
            placeholder="Account number of seats"
            required
          />
        </b-form-group>
        <b-form-group class="mb-0">
          <b-button variant="info" class="br-0 float-right" type="submit" :disabled="disabled">
            Create
            <i class="fa fa-check" />
          </b-button>
        </b-form-group>
      </b-form>
      <div class="modal-loading" v-if="state.loading">
        <i class="fa fa-spinner fa-spin" />
        <p class="m-0">Creating...</p>
      </div>
    </div>
  </b-modal>
</template>

<script>
export default {
  name: "createAccountModal",
  data() {
    return {
      show: true,
      id: null,
      name: null,
      type: "",
      num_of_seats: null,
      province: "",
      district: "",
      sector: "",
      state: { loading: false }
    };
  },
  computed: {
    disabled() {
      if (!this.name) return true;
      if (!this.num_of_seats) return true;
      if (!this.type) return true;
      if (this.type == "ben")
        if (!this.province && !this.district && !this.sector) return true;
      if (this.type == "dev") if (!this.id) return true;
      return false;
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
    }
  },
  watch: {
    type() {
      handler: {
        this.$nextTick(() => {
          this.province = "";
          this.district = "";
          this.sector = "";
        });
      }
    }
  },
  methods: {
    createAccount() {
      if (this.disabled === false) {
        this.state.loading = true;
        const id = `${this.province}.${this.district}.${this.sector}`;
        if (this.type == "dev") id = "paypack." + this.id;
        const data = {
          id: id,
          name: this.name,
          number_of_seats: this.num_of_seats,
          type: this.type
        };
        this.axios
          .post("/accounts", data)
          .then(res => {
            console.log(res.data);
            this.state.loading = false;
            this.$emit("created");
          })
          .catch(err => {
            console.log(err, err.response);
            this.state.loading = false;
          });
      }
    },
    modalClosed() {
      this.state.loading = false;
      this.$emit("close");
    }
  }
};
</script>

<style lang="scss">
.dev-createAccount-modal {
  .br-0 {
    border-radius: 2px;
  }
  .blur {
    filter: blur(3px);
  }
  .modal-loading {
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    width: 100%;
    height: 100%;
    background: #00000080;
    color: white;
    display: flex;
    justify-content: center;
    align-items: center;

    i {
      font-size: 2rem;
      margin-right: 0.5rem;
    }
    p {
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