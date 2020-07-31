<template>
  <b-form @submit.prevent="update">
    <b-row class="pb-1">
      <b-col>
        <b-form-group label="First name:">
          <b-form-input
            v-model="house.owner.fname"
            required
            placeholder="First name..."
            class="br-2"
          />
        </b-form-group>
      </b-col>
      <b-col>
        <b-form-group label="Last name:">
          <b-form-input
            v-model="house.owner.lname"
            required
            placeholder="Last name..."
            class="br-2"
          />
        </b-form-group>
      </b-col>
    </b-row>
    <b-form-group label="Phone Number" class="pb-1">
      <b-form-input
        v-model="house.owner.phone"
        placeholder="Phone Number..."
        required
        class="br-2"
      />
    </b-form-group>
    <b-form-group label="House is Rented?" class="pb-3">
      <b-form-radio-group v-model="rented" :options="query" />
    </b-form-group>
    <b-form-group class="pb-4">
      <template v-slot:label>
        <b-row class="align-items-center" no-gutters>
          Due:
          <b-input
            v-model="house.due"
            required
            step="100"
            min="500"
            size="sm"
            type="number"
            class="w-auto mx-1"
          />Rwf
        </b-row>
      </template>
      <div>
        <vue-slider
          v-model="house.due"
          :marks="slider.marks"
          :interval="500"
          :process="true"
          :tooltip="'none'"
          :min="500"
          :max="50000"
        >
          <template v-slot:label="{ active, value }">
            <div :class="['vue-slider-mark-label', 'custom-label', { active }]">{{ value/1000 }}K</div>
          </template>
        </vue-slider>
      </div>
    </b-form-group>

    <b-form-group class="pb-2">
      <template v-slot:label>
        <label class="m-0 d-flex justify-content-between">
          Cell
          <b-badge class="p-2">{{house.address.cell.toUpperCase()}}</b-badge>
        </label>
      </template>
      <b-form-select :options="cellOptions" v-model="newAddress.cell" class="br-2">
        <template v-slot:first>
          <option :value="null" disabled>select new cell</option>
        </template>
      </b-form-select>
    </b-form-group>

    <b-form-group class="pb-2">
      <template v-slot:label>
        <label class="m-0 d-flex justify-content-between">
          Village
          <b-badge class="p-2">{{house.address.village.toUpperCase()}}</b-badge>
        </label>
      </template>
      <b-form-select v-model="newAddress.village" :options="villageOptions" class="br-2">
        <template v-slot:first>
          <option :value="null" disabled>select new village</option>
        </template>
      </b-form-select>
    </b-form-group>
    <b-row class="d-flex justify-content-end" no-gutters>
      <b-button variant="info" type="submit" class="br-2">
        {{state.updating ? 'Updating' : "Update"}}
        <i
          v-show="state.updating"
          class="fa fa-spinner fa-spin"
        />
      </b-button>
    </b-row>
  </b-form>
</template>

<script>
const { Village } = require("rwanda");
export default {
  name: "updateHouse",
  props: {
    item: Object,
    option: Object
  },
  data() {
    return {
      house: null,
      rented: false,
      newAddress: {
        sector: null,
        cell: null,
        village: null
      },
      slider: {
        marks: val => val % 10000 === 0
      },
      query: [
        { text: "yego", value: "true" },
        { text: "oya", value: "false" }
      ],
      state: {
        updating: false
      }
    };
  },
  beforeMount() {
    const obj = JSON.parse(JSON.stringify(this.item));
    this.house = { ...obj };
  },
  computed: {
    cellOptions() {
      const { province, district, sector } = this.location;
      return this.$cells(province, district, sector);
    },
    villageOptions() {
      const cell = this.newAddress.cell
        ? this.newAddress.cell
        : this.house.address.cell;
      const { province, district, sector } = this.location;
      if (cell) return this.$villages(province, district, sector, cell);
      return [];
    },
    activeSector() {
      return this.$store.getters.getActiveSector;
    },
    location() {
      return this.$store.getters.location;
    },
    user() {
      return this.$store.getters.userDetails;
    }
  },
  methods: {
    update() {
      const sector = this.house.address.sector;
      const cell = this.newAddress.cell
        ? this.newAddress.cell
        : this.house.address.cell;
      const village = this.newAddress.village
        ? this.newAddress.village
        : this.house.address.village;
      this.state.updating = true;
      this.axios
        .put("/owners/" + this.house.owner.id, {
          fname: this.toCapital(this.house.owner.fname).trim(),
          lname: this.toCapital(this.house.owner.lname).trim(),
          phone: this.house.owner.phone.trim()
        })
        .then(res => {
          this.axios
            .put("/properties/" + this.house.id, {
              owner: {
                id: res.data.id
              },
              address: { cell: cell, village: village, sector: sector },
              recorded_by: this.house.recorded_by,
              due: String(this.house.due),
              occupied: this.rented,
              namespace: this.user.account
            })
            .then(response => {
              this.$snotify.info(response.data.message);
              this.state.updating = false;
              this.$emit("updated");
              this.$emit("closeModal");
            })
            .catch(err => {
              try {
                this.$snotify.error(
                  err.response.data.message || err.response.data
                );
              } catch {
                this.$snotify.error("Failed to update property");
              }
              this.state.updating = false;
              this.$emit("closeModal");
            });
        })
        .catch(err => {
          try {
            this.$snotify.error(err.response.data.message || err.response.data);
          } catch {
            this.$snotify.error("Failed to update owner");
          }
          this.state.updating = false;
          this.$emit("closeModal");
        });
    },
    toCapital(string) {
      string.toLowerCase();
      return string.charAt(0).toUpperCase() + string.slice(1);
    }
  }
};
</script>

<style>
</style>