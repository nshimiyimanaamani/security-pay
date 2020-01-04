<template>
  <b-container class="p-0 mr-1">
    <b-row class="py-2">
      <b-col class="px-1">
        <b-button v-b-modal.register-property class="py-1" variant="info">Register</b-button>
      </b-col>
      <b-col class="pl-1">
        <b-button class="py-1" variant="info">Refresh</b-button>
      </b-col>
    </b-row>
    <b-modal id="register-property" ref="register-modal" scrollable hide-footer>
      <template v-slot:modal-title>{{modal.title}}</template>
      <b-form @reset="resetModal">
        <b-row>
          <b-col lg="6" md="6" sm="auto">
            <b-form-group id="input-group-1" label="First Name:" label-for="input-1">
              <b-form-input
                id="input-1"
                v-model="modal.form.fname"
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
                v-model="modal.form.lname"
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
            v-model="modal.form.phone"
            type="number"
            required
            placeholder="Enter phone number..."
            style="font-size: 15px"
          ></b-form-input>
        </b-form-group>
        <b-form-group
          id="input-group-4"
          :label="'Due: '+Number(modal.form.due).toLocaleString()+' Rwf' "
          class="m-0"
          label-for="input-4"
        >
          <b-form-input
            id="input-4"
            v-model="modal.form.due"
            type="range"
            min="500"
            max="10000"
            step="500"
          ></b-form-input>
        </b-form-group>
        <b-form-group id="input-group-5" label="Sector:" label-for="input-5">
          <b-form-select id="input-5" v-model="modal.form.address.sector" style="font-size: 15px">
            <template v-slot:first>
              <option :value="null" disabled>-- Please select sector --</option>
            </template>
            <option value="Remera">Remera</option>
          </b-form-select>
        </b-form-group>

        <b-form-group id="input-group-6" label="Cell:" label-for="input-6">
          <b-form-select
            id="input-6"
            v-model="modal.form.address.cell"
            :options="cell_options"
            style="font-size: 15px"
          >
            <template v-slot:first>
              <option :value="null" disabled>-- Please select cell --</option>
            </template>
          </b-form-select>
        </b-form-group>

        <b-form-group id="input-group-7" label="Village:" label-for="input-7">
          <b-form-select
            id="input-7"
            v-model="modal.form.address.village"
            :options="village_options"
            style="font-size: 15px"
          >
            <template v-slot:first>
              <option :value="null" disabled>-- Please select village --</option>
            </template>
          </b-form-select>
        </b-form-group>
        <b-form-group id="input-group-8" class="float-right m-0 mt-3">
          <b-button type="reset" variant="danger" class="px-3 py-1">Cancel</b-button>
          <b-button type="submit" variant="primary" class="ml-2 px-3 py-1">{{modal.buttonTitle}}</b-button>
        </b-form-group>
      </b-form>
    </b-modal>
  </b-container>
</template>

<script>
const { District, Sector, Cell, Village } = require("rwanda");
export default {
  name: "controllers",
  data() {
    return {
      modal: {
        title: "Search Owner",
        loading: true,
        buttonTitle: "Search",
        form: {
          fname: null,
          lname: null,
          phone: null,
          ownerId: null,
          houseId: null,
          due: "500",
          address: {
            sector: null,
            cell: null,
            village: null
          }
        }
      }
    };
  },
  computed: {
    cell_options() {
      const sector = this.modal.form.address.sector;
      if (sector) {
        return Cell("Kigali", "Gasabo", sector);
      } else {
        return [];
      }
    },
    village_options() {
      const sector = this.modal.form.address.sector;
      const cell = this.modal.form.address.cell;
      if (sector && cell) {
        return Village("Kigali", "Gasabo", sector, cell);
      } else {
        return [];
      }
    }
  },
  mounted() {
    console.log();
  },
  methods: {
    resetModal() {
      this.$refs["register-modal"].hide();
      this.modal = {
        title: "Search Owner",
        loading: true,
        buttonTitle: "Search",
        form: {
          fname: null,
          lname: null,
          phone: null,
          ownerId: null,
          houseId: null,
          due: "500",
          address: {
            sector: null,
            cell: null,
            village: null
          }
        }
      };
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