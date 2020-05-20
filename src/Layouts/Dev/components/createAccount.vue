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
      <b-form class="p-2" :class="{'blur':state.loading}">
        <b-form-group label="Identifier">
          <b-input-group prepend="paypack.">
            <b-input class="br-0" v-model="id" />
          </b-input-group>
        </b-form-group>
        <b-form-group label="Name">
          <b-input class="br-0" v-model="name" placeholder="Account name" />
        </b-form-group>
        <b-form-group label="Type">
          <b-input class="br-0" v-model="type" placeholder="Account type" />
        </b-form-group>
        <b-form-group label="Number of seats">
          <b-input class="br-0" v-model="num_of_seats" placeholder="Account number of seats" />
        </b-form-group>
        <b-form-group class="mb-0">
          <b-button
            variant="info"
            class="br-0 float-right"
            @click="createAccount"
            :disabled="disabled"
          >
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
      type: null,
      num_of_seats: null,
      state: { loading: false }
    };
  },
  computed: {
    disabled() {
      if (this.id && this.name && this.type && this.num_of_seats) return false;
      return true;
    }
  },
  methods: {
    createAccount() {
      if (this.id && this.name && this.type && this.num_of_seats) {
        this.state.loading = true;
        const data = {
          id: "paypack." + this.id,
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
}
</style>