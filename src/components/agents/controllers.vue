<template>
  <div class="w-100 m-auto">
    <b-row class="m-0 justify-content-between align-items-center">
      <b-button v-b-modal.register-property class="font-14" variant="info">Register</b-button>
      <b-button class="font-14 ml-2" variant="info" @click.prevent="refresh">
        Refresh
        <i class="fa fa-sync-alt ml-1" />
      </b-button>
    </b-row>
    <b-modal id="register-property" ref="register-modal" hide-footer>
      <template v-slot:modal-title>Register Property</template>
      <b-form @reset="resetModal" @submit.prevent="addProperty">
        <b-row class="pl-3">
          <b-form-group
            id="input-group-1"
            label="First Name:"
            label-for="input-1"
            class="w-auto pr-3 flex-grow-1"
          >
            <b-form-input
              id="input-1"
              v-model="form.fname"
              required
              placeholder="Enter first name..."
              size="sm"
            ></b-form-input>
          </b-form-group>
          <b-form-group
            id="input-group-2"
            label="Last Names:"
            label-for="input-2"
            class="w-auto pr-3 flex-grow-1"
          >
            <b-form-input
              id="input-2"
              v-model="form.lname"
              required
              placeholder="Enter last name..."
              size="sm"
            ></b-form-input>
          </b-form-group>
        </b-row>

        <b-form-group id="input-group-3" label="Phone Number:" label-for="input-3">
          <b-form-input
            id="input-3"
            v-model="form.phone"
            type="number"
            required
            placeholder="Enter phone number..."
            size="sm"
          ></b-form-input>
        </b-form-group>
        <b-form-group label="House occupied?">
          <b-form-radio-group
            v-model="form.occupied"
            :options="occupiedOptions"
            name="radios-stacked"
            size="sm"
            :state="occupied ? null : false"
          />
        </b-form-group>
        <b-form-group id="input-group-4" class label-for="input-4">
          <template v-slot:label>
            <b-row class="m-o align-items-center px-3 flex-nowrap">
              Due:
              <b-input
                v-model="form.due"
                required
                step="500"
                min="500"
                size="sm"
                type="number"
                class="w-auto mx-1"
              />Rwf
            </b-row>
          </template>
          <b-form-input
            id="range-1"
            v-model="form.due"
            type="range"
            min="500"
            max="10000"
            step="500"
            size="sm"
            class="border-0"
          />
        </b-form-group>
        <hr />
        <b-row id="input-group-8" class="m-0 mt-3 justify-content-between">
          <b-button type="reset" variant="danger" size="sm" class="px-3 py-1 font-14">cancel</b-button>
          <b-button type="submit" size="sm" variant="info" class="px-3 py-1 font-14">
            {{state.loading ? 'Registering' : 'Register'}}
            <b-spinner v-show="state.loading" small type="grow" />
          </b-button>
        </b-row>
      </b-form>
    </b-modal>
  </div>
</template>

<script>
export default {
  name: "controllers",
  props: {
    user: Object
  },
  data() {
    return {
      occupiedOptions: [
        { text: "No", value: false },
        { text: "Yes", value: true }
      ],
      form: {
        fname: null,
        lname: null,
        phone: null,
        due: "500",
        occupied: null
      },
      slider: {
        marks: val => val % 10000 === 0
      },
      state: {
        loading: false
      },
      owner: null
    };
  },
  computed: {
    userDetails() {
      return this.$store.getters.userDetails;
    },
    occupied() {
      return Boolean(this.form.occupied !== null);
    }
  },
  methods: {
    async addProperty() {
      this.state.loading = true;
      const ownerId = await this.search();
      if (ownerId) {
        this.axios
          .post("/properties", {
            owner: {
              id: ownerId
            },
            address: {
              cell: this.user.cell,
              village: this.user.village,
              sector: this.user.sector
            },
            namespace: this.userDetails.account,
            due: this.form.due.toString(),
            occupied: this.form.occupied,
            recorded_by: this.userDetails.username
          })
          .then(res => {
            this.refresh();
            console.log(ownerId);
            this.$snotify.info(`Property Registered successfully!`);
          })
          .catch(err => {
            const error = err.response
              ? err.response.data.error || err.response.data
              : null;
            if (error) this.$snotify.error(error);
          })
          .finally(() => {
            this.state.loading = false;
            this.resetModal();
          });
      }
    },
    async search() {
      const fname = this.capitalize(this.form.fname.trim());
      const lname = this.capitalize(this.form.lname.trim());
      const phone = this.form.phone.trim();
      return this.axios
        .get(`/owners/search?fname=${fname}&lname=${lname}&phone=${phone}`)
        .then(res => res.data.id)
        .catch(err => {
          const message = `${fname} ${lname} is not a registered owner! Do you want to register this owner?`;
          return this.confirm(message).then(state => {
            if (state === true) {
              return this.axios
                .post("/owners", {
                  fname: fname,
                  lname: lname,
                  phone: phone
                })
                .then(async res => {
                  const id = await this.checkOwner(res.data.id);
                  return id;
                })
                .catch(err => {
                  this.$snotify.error("Failed to register owner!");
                  return null;
                });
            } else {
              this.resetModal();
              return null;
            }
          });
        });
    },
    checkOwner(id) {
      return new Promise(resolve => {
        var i = 0;
        var resolved = false;
        var interval = setInterval(() => {
          i++;
          this.axios.get(`/owners/${id}`).then(res => {
            if (res && res.status == 200) {
              resolve(res.data.id);
              clearInterval(interval);
              resolved = true;
            }
            if (i > 5) {
              clearInterval(interval);
              if (resolved === false) {
                resolve(null);
              }
            }
          });
        }, 1000);
      });
    },
    resetModal() {
      this.$refs["register-modal"].hide();
      this.form = {
        fname: null,
        lname: null,
        phone: null,
        due: "500"
      };
      this.owner = null;
    },
    refresh() {
      this.$emit("refresh");
    },
    capitalize(string) {
      string.toLowerCase();
      return string.charAt(0).toUpperCase() + string.slice(1);
    },
    confirm(message) {
      return this.$bvModal.msgBoxConfirm(message, {
        title: "Please Confirm",
        buttonSize: "sm",
        okVariant: "danger",
        okTitle: "YES",
        cancelTitle: "NO",
        footerClass: "p-3",
        hideHeaderClose: false,
        centered: true
      });
    }
  }
};
</script>
<style lang="scss" >
form {
  .form-group {
    margin-bottom: 0.7rem;
    label,
    button {
      font-size: 15px;
    }
  }
}
</style>