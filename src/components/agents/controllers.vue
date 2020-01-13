<template>
  <b-container class="p-0 mr-1">
    <b-row class="py-2">
      <b-col class="px-1">
        <b-button v-b-modal.register-property class="py-1" variant="info">Register</b-button>
      </b-col>
      <b-col class="pl-1">
        <b-button class="py-1" variant="info" @click.prevent="refresh">Refresh</b-button>
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
                style="font-size: 15px"
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
                style="font-size: 15px"
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
            style="font-size: 15px"
          ></b-form-input>
        </b-form-group>
        <b-form-group
          id="input-group-4"
          :label="'Due: '+Number(form.due).toLocaleString()+' Rwf' "
          class="m-0"
          label-for="input-4"
        >
          <b-form-input
            id="input-4"
            v-model="form.due"
            type="range"
            min="500"
            max="10000"
            step="500"
          ></b-form-input>
        </b-form-group>
        <b-form-group id="input-group-8" class="float-right m-0 mt-3">
          <b-button type="submit" variant="primary" class="ml-2 px-3 py-1">
            {{state.loading ? 'Registering' : 'Register'}}
            <b-spinner v-show="state.loading" small type="grow"></b-spinner>
          </b-button>
          <b-button type="reset" variant="danger" class="px-3 py-1">Cancel</b-button>
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
      form: {
        fname: null,
        lname: null,
        phone: null,
        due: "500"
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
    }
  },
  mounted() {
    console.log();
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
            occupied: true,
            recorded_by: this.userDetails.username
          })
          .then(res => {
            this.state.loading = false;
            this.resetModal();
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
            this.state.loading = false;
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
          console.log(res.data);
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
                    console.log(res.data);
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
    .custom-range {
      border: none !important;
    }
  }
}
</style>