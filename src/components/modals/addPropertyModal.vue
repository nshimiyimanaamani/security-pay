<template>
  <div class="add-property-modal secondary-font" v-show="show">
    <!-- Modal content -->
    <b-card class="mb-2 modal-content p-0" no-body>
      <header class="modal-header">
        <h5 class="text-center mb-1">ADD PROPERTY</h5>
        <i class="fa fa-times" @click="$emit('closeModal')" />
      </header>
      <b-form @submit.prevent="search_user" @reset="resetModal" class="modal-body p-4">
        <b-form-group
          id="input-group-1"
          label="First Name:"
          label-for="input-1"
          description="Amazina ya nyiri inzu (*Ntabwo ari ukodesheje)"
        >
          <b-form-input
            id="input-1"
            required
            v-model="form.fname"
            placeholder="First name"
            :disabled="state.switch"
            class="br-2"
            trim
          />
        </b-form-group>
        <b-form-group id="input-group-2" label="Last Name" label-for="input-2">
          <b-form-input
            id="input-2"
            v-model="form.lname"
            :disabled="state.switch"
            required
            trim
            class="br-2"
            placeholder="Last name"
          ></b-form-input>
        </b-form-group>
        <b-form-group id="input-group-3" label="Phone Number" label-for="input-3">
          <b-form-input
            id="input-3"
            v-model="form.phone"
            :state="checkNumber"
            :disabled="state.switch"
            required
            type="number"
            class="br-2"
            placeholder="Phone number"
          ></b-form-input>
          <b-form-invalid-feedback :state="checkNumber">Please use a valid Phone number!</b-form-invalid-feedback>
        </b-form-group>
        <!-- second Modal -->
        <b-form-group id="input-group-4" label-for="range-1" v-show="state.switch" class>
          <template v-slot:label>
            <b-row class="m-0 align-items-center">
              Due:
              <b-input
                v-model="form.due"
                required
                step="500"
                min="500"
                type="number"
                class="w-auto mx-1"
              />Rwf
            </b-row>
          </template>
          <vue-slider
            v-model="form.due"
            :marks="slider.marks"
            :interval="500"
            :process="true"
            :tooltip="'none'"
            :min="500"
            :max="50000"
            class="pt-2 pr-2 pb-4"
          >
            <template v-slot:label="{ active, value }">
              <div :class="['vue-slider-mark-label', 'custom-label', { active }]">{{ value/1000 }}K</div>
            </template>
          </vue-slider>
        </b-form-group>
        <b-form-group label="Inzu ituwemo?" v-if="state.switch">
          <b-form-radio-group
            v-model="form.occupied"
            :options="occupiedOptions"
            name="radios-stacked"
            :state="occupied?true:false"
          ></b-form-radio-group>
        </b-form-group>
        <b-form-group id="input-group-5" label="Cell:" label-for="input-4" v-show="state.switch">
          <b-form-select v-model="address.cell" :options="cellOptions" class="br-2">
            <template v-slot:first>
              <option :value="null" disabled>select a cell</option>
            </template>
          </b-form-select>
        </b-form-group>
        <b-form-group
          id="input-group-6"
          label="Village:"
          label-for="input-6"
          v-show="state.switch"
          class="mb-3"
        >
          <b-form-select v-model="address.village" :options="villageOptions" class="br-2">
            <template v-slot:first>
              <option :value="null" disabled>select a village</option>
            </template>
          </b-form-select>
        </b-form-group>
        <b-form-group class="mb-0 mt-4">
          <b-button :disabled="!clickable" type="submit" variant="info" class="br-2">
            {{state.adding ? btnContent+'ing' : btnContent}}
            <i
              v-if="state.adding"
              class="fa fa-spinner fa-spin"
            />
            <i v-if="!state.adding" class="fa fa-search" />
          </b-button>
          <b-button type="reset" class="br-2" variant="danger">cancel</b-button>
        </b-form-group>
      </b-form>
    </b-card>
  </div>
</template>

<script>
export default {
  name: "addPropertModal",
  props: {
    show: Boolean
  },
  data() {
    return {
      title: "Search House Owner",
      btnContent: "Search",
      occupiedOptions: [
        { text: "oya", value: false },
        { text: "yego", value: true }
      ],
      form: {
        fname: null,
        lname: null,
        phone: null,
        id: null,
        due: "500",
        occupied: null
      },
      slider: {
        marks: val => val % 10000 === 0
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
    sectorOptions() {
      return [this.activeSector];
    },
    activeSector() {
      return this.$store.getters.getActiveSector;
    },
    cellOptions() {
      return this.$store.getters.getCellsArray;
    },
    villageOptions() {
      const cell = this.address.cell;
      if (cell) return this.$villages("Kigali", "Gasabo", "Remera", cell);
      return [];
    },
    checkNumber() {
      return this.form.phone ? this.$isPhoneNumber(this.form.phone) : null;
    },
    userDetails() {
      return this.$store.getters.userDetails;
    },
    occupied() {
      return Boolean(this.form.occupied !== null);
    },
    clickable() {
      const { fname, lname, phone, id, due, occupied } = this.form;
      const { sector, cell, village } = this.address;
      if (this.state.switch) {
        if (
          fname &&
          lname &&
          phone &&
          due &&
          occupied != null &&
          cell &&
          village
        ) {
          return true;
        }
      } else {
        if (fname && lname && phone) {
          return true;
        }
      }
      return false;
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
            const message = `${fname} ${lname} is not registered! Do you want to register this user?`;
            this.confirm(message).then(state => {
              if (state === true) {
                this.state.adding = true;
                this.axios
                  .post(`/owners`, {
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
                    const error = err.response
                      ? err.response.data.message || err.response.data
                      : null;
                    if (error) this.$snotify.error(error);
                    this.state.adding = false;
                  });
              }
            });
          });
      } else if (this.state.switch) {
        this.state.adding = true;
        this.axios
          .post("/properties", {
            owner: {
              id: this.form.id
            },
            address: {
              cell: this.address.cell,
              village: this.address.village,
              sector: this.activeSector
            },
            due: this.form.due.toString(),
            occupied: this.form.occupied.toString(),
            recorded_by: this.userDetails.username
          })
          .then(res => {
            this.resetModal();
            this.$emit("refresh");
            this.$snotify.info(`Property Registered successfully!`);
          })
          .catch(err => {
            const error = err.response
              ? err.response.data.error || err.response.data
              : null;
            if (error) this.$snotify.error(error);
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
      this.form.due = "500";
      this.address.cell = null;
      this.address.village = null;
      this.form.occupied = null;
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

<style lang="scss">
.add-property-modal {
  position: fixed;
  top: 50px;
  width: calc(100% - 230px);
  height: calc(100% - 50px);
  left: 230px;
  right: 0;
  margin: auto;
  background: #000000cc;
  padding: 1rem;
  z-index: 1000;
  overflow-x: hidden;
  overflow-y: auto;

  .modal-content {
    position: absolute;
    width: auto;
    top: 2rem;
    left: 0;
    right: 0;
    margin: auto;
    min-width: 300px;
    max-width: 500px;

    .modal-header {
      width: 100%;
      height: 50px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0.5rem 1rem;
    }

    form button {
      float: right;
      margin-left: 10px;
    }
  }
}
</style>