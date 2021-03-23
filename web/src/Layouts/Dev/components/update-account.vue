<template>
  <div class="dev-updateAccount">
    <b-form class="py-2 px-4" :class="{'blur':state.loading}">
      <b-form-group label="Name">
        <b-form-input class="br-0" v-model="name" :placeholder="account.name" trim required />
      </b-form-group>
      <b-form-group label="Number of seats">
        <b-form-input
          type="number"
          v-model="num_of_seats"
          class="br-0"
          trim
          placeholder="--"
          required
        />
      </b-form-group>
      <b-form-group label="Type">
        <b-form-input class="br-0" v-model="type" :placeholder="account.type" trim required />
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
      name: null,
      num_of_seats: null,
      type: null,
      state: { loading: false }
    };
  },
  computed: {
    showButton() {
      if (this.name || this.type || this.num_of_seats) return false;
      return true;
    },
    accountId() {
      return this.account.id;
    }
  },
  methods: {
    updateAccount() {
      this.state.loading = true;
      const data = {
        name: this.name || this.account.name,
        number_of_seats: this.num_of_seats || 5,
        type: this.type || this.account.type
      };
      this.axios
        .put("/accounts/" + this.accountId, data)
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