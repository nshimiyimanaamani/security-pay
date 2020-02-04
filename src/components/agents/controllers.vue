<template>
  <b-container class="p-0 mr-1">
    <b-row class="py-2">
      <b-col class="px-1">
        <b-button v-b-modal.register-property class="py-1 font-15" size="sm" variant="info">Register</b-button>
      </b-col>
      <b-col class="pl-1">
        <b-button class="py-1 font-15" size="sm" variant="info" @click.prevent="refresh">Refresh</b-button>
      </b-col>
    </b-row>
    <b-modal id="register-property" ref="register-modal" scrollable hide-footer>
      <template v-slot:modal-title>Register Property</template>
      <b-form @reset="resetModal" @submit.prevent="addProperty">
        <b-row>
          <b-col lg="6" md="6" sm="auto">
            <b-form-group id="input-group-1" label="First Name:" label-for="input-1">
              <b-form-input
                id="input-1"
                v-model="form.fname"
                required
                placeholder="Enter first name..."
                size="sm"
              ></b-form-input>
            </b-form-group>
          </b-col>
          <b-col lg="6" md="6" sm="auto">
            <b-form-group id="input-group-2" label="Last Names:" label-for="input-2">
              <b-form-input
                id="input-2"
                v-model="form.lname"
                required
                placeholder="Enter last name..."
                size="sm"
              ></b-form-input>
            </b-form-group>
          </b-col>
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
        <b-form-group label="Inzu ituwemo?">
          <b-form-radio-group
            v-model="form.occupied"
            :options="occupiedOptions"
            name="radios-stacked"
            size="sm"
            :state="occupied ? null : false"
          />
        </b-form-group>
        <b-form-group id="input-group-4" class="m-2" label-for="input-4">
          <template v-slot:label>
            <b-row class="m-o align-items-center px-3">
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
          <div>
            <vue-slider
              v-model="form.due"
              :marks="slider.marks"
              :interval="500"
              :process="true"
              :tooltip="'none'"
              :min="500"
              :max="50000"
            >
              <template v-slot:label="{ active, value }">
                <div
                  :class="['vue-slider-mark-label', 'custom-label', { active }]"
                >{{ value/1000 }}K</div>
              </template>
            </vue-slider>
          </div>
        </b-form-group>
        <b-form-group id="input-group-8" class="float-right m-0 mt-3">
          <b-button
            type="submit"
            size="sm"
            variant="primary"
            class="ml-2 px-3 py-1 app-color font-15"
          >
            {{state.loading ? 'Registering' : 'Register'}}
            <b-spinner v-show="state.loading" small type="grow"></b-spinner>
          </b-button>
          <b-button type="reset" variant="danger" size="sm" class="px-3 py-1 font-15">cancel</b-button>
        </b-form-group>
      </b-form>
    </b-modal>
  </b-container>
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
        { text: "oya", value: false },
        { text: "yego", value: true }
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
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    userDetails() {
      return this.$store.getters.userDetails;
    },
    occupied() {
      return Boolean(this.form.occupied !== null);
    }
  },
  methods: {
    addProperty() {
      if (!this.owner) {
        this.search();
      } else {
        this.state.loading = true;
        this.axios
          .post(this.endpoint + "/properties", {
            owner: {
              id: this.owner.id
            },
            address: {
              cell: this.user.cell,
              village: this.user.village,
              sector: this.user.sector
            },
            due: this.form.due.toString(),
            occupied: this.form.occupied,
            recorded_by: this.userDetails.username
          })
          .then(res => {
            this.refresh();
            this.$snotify.info(`Property Registered successfully!`);
          })
          .catch(err => {
            if (navigator.onLine) {
              const error = err.response
                ? err.response.data.error || err.response.data
                : "an error occured";
              this.$snotify.error(error);
            } else {
              this.$snotify.error("Please connect to the internet");
            }
          })
          .finally(() => {
            this.state.loading = false;
            this.resetModal();
          });
      }
    },
    search() {
      const fname = this.capitalize(this.form.fname.trim());
      const lname = this.capitalize(this.form.lname.trim());
      const phone = this.form.phone.trim();
      this.state.loading = true;
      this.axios
        .get(
          this.endpoint +
            `/owners/search?fname=${fname}&lname=${lname}&phone=${phone}`
        )
        .then(res => {
          this.owner = { ...res.data };
          this.addProperty();
        })
        .catch(err => {
          if (!navigator.onLine) {
            this.state.loading = false;
            this.$snotify.error("Please connect to the internet");
          } else {
            this.state.loading = false;
            const message = `${fname} ${lname} is not a registered owner! Do you want to register this owner?"`;
            this.confirm(message).then(state => {
              if (state === true) {
                this.state.loading = true;
                this.axios
                  .post(this.endpoint + "/owners", {
                    fname: fname,
                    lname: lname,
                    phone: phone
                  })
                  .then(res => {
                    this.owner = { ...res.data };
                    this.addProperty();
                    this.resetModal();
                  })
                  .catch(err => {
                    if (navigator.onLine) {
                      this.resetModal();
                      const error = err.response
                        ? err.response.data.message || err.response.data
                        : "an error occured";
                      this.$snotify.error(error);
                    } else {
                      this.$snotify.error("Please connect to the internet");
                    }
                    this.state.loading = false;
                  });
              }
            });
          }
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