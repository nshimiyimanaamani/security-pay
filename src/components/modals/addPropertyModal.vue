<template>
  <div class="add-property-modal" v-show="show">
    <!-- Modal content -->
    <b-card class="mb-2 modal-body">
      <h5 class="text-center mb-1">Add Property</h5>
      <hr />
      <b-form @submit.prevent="search_user" @reset="resetModal">
        <b-form-group
          id="input-group-1"
          class="mb-2"
          label="First Name:"
          label-for="input-1"
          description="Amazina ya nyiri inzu (*Ntabwo ari ukodesheje)"
        >
          <b-form-input
            id="input-1"
            v-model="form.fname"
            required
            placeholder="First name"
            :disabled="state.switch"
          ></b-form-input>
        </b-form-group>
        <b-form-group id="input-group-2" class="mb-2" label="Last Name:" label-for="input-2">
          <b-form-input
            id="input-2"
            v-model="form.lname"
            :disabled="state.switch"
            required
            placeholder="Last name"
          ></b-form-input>
        </b-form-group>
        <b-form-group id="input-group-3" class="mb-2" label="Phone Number:" label-for="input-3">
          <b-form-input
            id="input-3"
            v-model="form.phone"
            :state="checkNumber"
            :disabled="state.switch"
            required
            type="number"
            placeholder="Phone number"
          ></b-form-input>
          <b-form-invalid-feedback :state="checkNumber">Please use a valid Phone number!</b-form-invalid-feedback>
        </b-form-group>
        <b-form-group
          id="input-group-4"
          :label="'Due: '+ form.due +' Rwf'"
          label-for="range-1"
          v-show="state.switch"
          class="mb-2"
        >
          <b-form-input
            id="range-1"
            v-model="form.due"
            type="range"
            min="500"
            max="10000"
            step="500"
          ></b-form-input>
        </b-form-group>
        <b-form-group
          id="input-group-5"
          class="mb-2"
          label="Cell:"
          label-for="input-4"
          v-show="state.switch"
        >
          <b-form-select v-model="address.cell" :options="cellOptions" class="mb-0">
            <template v-slot:first>
              <option :value="null" disabled>select a cell</option>
            </template>
          </b-form-select>
        </b-form-group>
        <b-form-group
          id="input-group-6"
          label="Village:"
          label-for="input-5"
          v-show="state.switch"
          class="mb-3"
        >
          <b-form-select v-model="address.village" :options="villageOptions" class="mb-0">
            <template v-slot:first>
              <option :value="null" disabled>select a village</option>
            </template>
          </b-form-select>
        </b-form-group>
        <b-button type="submit" variant="primary">
          {{state.adding ? btnContent+'ing' : btnContent}}
          <b-spinner v-show="state.adding" small type="grow"></b-spinner>
        </b-button>
        <b-button type="reset" variant="danger">cancel</b-button>
      </b-form>
    </b-card>
  </div>
</template>

<script>
const { isPhoneNumber } = require("rwa-validator");
const { Village } = require("rwanda");
export default {
  name: "addPropertModal",
  props: {
    show: Boolean
  },
  data() {
    return {
      title: "Search House Owner",
      btnContent: "Search",
      form: {
        fname: null,
        lname: null,
        phone: null,
        id: null,
        due: "500"
      },
      address: {
        sector: null,
        cell: null,
        village: null
      },
      state: {
        adding: false,
        switch: false
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    sectorOptions() {
      return [this.activeSector];
    },
    cellOptions() {
      return this.$store.getters.getCellsArray;
    },
    villageOptions() {
      const cell = this.address.cell;
      if (cell) {
        return Village("Kigali", "Gasabo", "Remera", cell);
      } else {
        return [];
      }
    },
    checkNumber() {
      const n = this.form.phone;
      return n ? isPhoneNumber(n) : null;
    },
    userDetails() {
      return this.$store.getters.userDetails;
    }
  },
  methods: {
    search_user() {
      if (!this.state.switch) {
        const fname = this.capitalize(this.form.fname.trim());
        const lname = this.capitalize(this.form.lname.trim());
        const phone = this.form.phone.trim();
        this.state.adding = true;
        this.axios
          .get(
            this.endpoint +
              "/owners/search?fname=" +
              fname +
              "&lname=" +
              lname +
              "&phone=" +
              phone
          )
          .then(res => {
            this.state.adding = false;
            this.state.switch = true;
            this.btnContent = "Register";
            this.form.id = res.data.id;
          })
          .catch(err => {
            this.state.adding = false;
            if (navigator.onLine) {
              const message = `${fname} ${lname} is not registered! Do you want to register this user?`;
              this.confirm(message).then(state => {
                if (state === true) {
                  this.state.adding = true;
                  this.axios
                    .post(`${this.endpoint}/owners`, {
                      fname: fname,
                      lname: lname,
                      phone: phone
                    })
                    .then(res => {
                      this.state.adding = false;
                      this.state.switch = true;
                      this.btnContent = "Register";
                      this.form.id = res.data.id;
                      this.$snotify.info(
                        `User created. proceeding to registration...`
                      );
                    })
                    .catch(err => {
                      if (navigator.onLine) {
                        const error = err.response
                          ? err.response.data.message || err.response.data
                          : "an error occured";
                        this.$snotify.error(error);
                      } else {
                        this.$snotify.error("Please connect to the internet");
                      }
                      this.state.adding = false;
                    });
                }
              });
            } else if (!navigator.onLine) {
              this.$snotify.info(`Please connect to the internet...`);
            }
          });
      } else if (this.state.switch) {
        this.state.adding = true;
        this.axios
          .post(this.endpoint + "/properties", {
            owner: {
              id: this.form.id
            },
            address: {
              cell: this.address.cell,
              village: this.address.village,
              sector: "Remera"
            },
            due: this.form.due.toString(),
            occupied: true,
            recorded_by: this.userDetails.username
          })
          .then(res => {
            this.resetModal();
            this.$emit("refresh");
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
            this.state.adding = false;
          });
      }
    },
    resetModal() {
      this.$emit("closeModal");
      this.state.switch = false;
      this.state.adding = false;
      this.btnContent = "Search";
      this.form.fname = null;
      this.form.lname = null;
      this.form.phone = null;
      this.form.id = null;
      this.form.due = null;
      this.address.cell = null;
      this.address.village = null;
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

<style>
.add-property-modal {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  width: 100%;
  height: 100%;
  margin: auto;
  background: #000000cc;
  z-index: 100;
}

.add-property-modal .modal-body {
  position: sticky;
  -ms-flex: 1 1 auto;
  -webkit-box-flex: 1;
  flex: 1 1 auto;
  padding: 0;
  width: 40%;
  top: 5rem;
  margin: auto;
}

.modal-body form button {
  float: right;
  margin-left: 10px;
  padding: 3px 15px;
}
</style>