<template>
  <b-row class="justify-content-center">
    <b-form class="px-3 w-100">
      <b-row>
        <b-col>
          <b-form-group id="input-group-1" label="First name:" label-for="input-1">
            <b-form-input
              id="input-1"
              v-model="house.owner.fname"
              required
              placeholder="First name..."
              size="sm"
            ></b-form-input>
          </b-form-group>
        </b-col>
        <b-col>
          <b-form-group id="input-group-2" label="Last name:" label-for="input-2">
            <b-form-input
              id="input-2"
              v-model="house.owner.lname"
              required
              placeholder="Last name..."
              size="sm"
            ></b-form-input>
          </b-form-group>
        </b-col>
      </b-row>
      <b-form-group id="input-group-3" label="Phone Number" label-for="input-3">
        <b-form-input
          size="sm "
          id="input-3"
          v-model="house.owner.phone"
          placeholder="Phone Number..."
          required
        />
      </b-form-group>
      <b-form-group label="Irakodeshwa ?">
        <b-form-radio-group v-model="rented" :options="query" name="radio-stacked" size="sm"></b-form-radio-group>
      </b-form-group>
      <b-form-group id="input-group-4" label-for="range-1" class="mb-4">
        <template v-slot:label>
          <b-row class="m-o align-items-center px-3">
            Due:
            <b-input
              v-model="house.due"
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

      <b-form-group
        id="input-group-5"
        :label="'Cell: '+house.address.cell.toUpperCase()"
        label-class="font-15"
        label-for="select-5"
      >
        <b-form-select
          id="select-5"
          class="font-15"
          :options="cellOptions"
          v-model="newAddress.cell"
          size="sm"
        >
          <template v-slot:first>
            <option :value="null" disabled>select new cell</option>
          </template>
        </b-form-select>
      </b-form-group>

      <b-form-group
        id="input-group-7"
        :label="'village: '+house.address.village.toUpperCase()"
        label-for="input-7"
        label-class="font-15"
      >
        <b-form-select
          id="input-7"
          v-model="newAddress.village"
          :options="villageOptions"
          size="sm"
          class="font-15"
        >
          <template v-slot:first>
            <option :value="null" disabled>select new village</option>
          </template>
        </b-form-select>
      </b-form-group>
      <b-form-group class="m-0">
        <b-button variant="info" class="float-right" @click.prevent="update">
          {{state.updating ? 'Updating' : "Update"}}
          <b-spinner v-show="state.updating" small type="grow"></b-spinner>
        </b-button>
      </b-form-group>
    </b-form>
  </b-row>
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
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    cellOptions() {
      return this.$store.getters.getCellsArray;
    },
    villageOptions() {
      const cell = this.newAddress.cell
        ? this.newAddress.cell
        : this.house.address.cell;
      if (cell) {
        return Village("Kigali", "Gasabo", this.activeSector, cell);
      } else {
        return [];
      }
    },
    activeSector() {
      return this.$store.getters.getActiveSector;
    }
  },
  methods: {
    due(house) {
      return Number(house.due).toLocaleString();
    },
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
        .put(this.endpoint + "/owners/" + this.house.owner.id, {
          fname: this.toCapital(this.house.owner.fname).trim(),
          lname: this.toCapital(this.house.owner.lname).trim(),
          phone: this.house.owner.phone.trim()
        })
        .then(res => {
          this.axios
            .put(this.endpoint + "/properties/" + this.house.id, {
              owner: {
                id: res.data.id
              },
              address: { cell: cell, village: village, sector: sector },
              recorded_by: this.house.recorded_by,
              due: String(this.house.due),
              occupied: this.rented
            })
            .then(response => {
              this.$snotify.info(response.data.message);
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
            })
            .finally(() => {
              this.state.updating = false;
              this.$emit("closeModal");
            });
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