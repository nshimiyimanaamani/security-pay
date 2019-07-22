<template>
  <div class="dashboardWrapper">
    <div class="dashboardSidebar">
      <h1>PayPack</h1>
      <ul class="sidebarLinks">
        <router-link to="/dashboard">
          <li>Overview</li>
        </router-link>
        <router-link to="/cells">
          <li>
            Cells
            <b-dropdown>
              <template slot="button-content">
                <span class="fa fa-caret-down"></span>
              </template>
              <!-- <v-select></v-select> -->
              <b-dropdown-item
                @click="update({toUpdate: 'cell', changed: cell})"
                v-for="cell in sidebar.cells_array"
                :key="cell"
              >{{cell}}</b-dropdown-item>
            </b-dropdown>
          </li>
        </router-link>
        <router-link to="/village">
          <li>
            Village
            <b-dropdown>
              <template slot="button-content">
                <span class="fa fa-caret-down"></span>
              </template>
              <b-dropdown-item
                @click="update({toUpdate: 'village', changed: village})"
                v-for="village in sidebar.village_array"
                :key="village"
              >{{village}}</b-dropdown-item>
            </b-dropdown>
          </li>
        </router-link>
        <router-link to="/transactions">
          <li>Bank Accounts</li>
        </router-link>
        <router-link to="#">
          <li>Services</li>
        </router-link>
        <router-link to="#">
          <li>Penalties</li>
        </router-link>
      </ul>
      <button id="btn" @click="modalShow = !modalShow">
        <i class="fa fa-plus-circle" /> Add Properties
      </button>
    </div>
    <div class="top-nav">
      <div class="dropdown">
        <b-dropdown id="dropdown-left" class="m-md-6" text="YTD">
          <b-dropdown-item>First Action</b-dropdown-item>
          <b-dropdown-item>Second Action</b-dropdown-item>
          <b-dropdown-item>Third Action</b-dropdown-item>
        </b-dropdown>
        <i class="fa fa-caret-down"></i>
      </div>
      <div class="search">
        <input type="search" name="search" id="search" placeholder="Search..." />
      </div>
    </div>
    <div class="dashboardBody">
      <router-view />
    </div>
    <div class="AddPropertymodal" v-show="modalShow">
      <div class="form">
        <b-form class="form1" v-show="!toggleForms">
          <b-form-group label="Owner:">
            <div class="names">
              <b-form-input v-model="form.fname" required placeholder="First Name..."></b-form-input>
              <b-form-input v-model="form.lname" required placeholder="Last Name..."></b-form-input>
            </div>
          </b-form-group>
          <b-form-group label="Phone number:">
            <b-form-input type="number" v-model="form.phoneNo"></b-form-input>
          </b-form-group>

          <div class="buttons">
            <b-button type="submit" variant="primary" @click="search">Search</b-button>
            <b-button variant="danger" @click="cancel">cancel</b-button>
          </div>
        </b-form>

        <b-form @submit="onSubmit" v-show="toggleForms" class="form2">
          <b-form-group class="phone" label="First Name:">
            <b-form-input disabled v-model="form.fname" placeholder="First Name..."></b-form-input>
          </b-form-group>
          <b-form-group class="amount" label="Surname:">
            <b-form-input disabled v-model="form.lname" placeholder="surname..."></b-form-input>
          </b-form-group>
          <b-form-group class="phone" label="Phone number:">
            <b-form-input
              type="number"
              disabled
              v-model="form.phoneNo"
              placeholder="Phone Number..."
            ></b-form-input>
          </b-form-group>
          <b-form-group class="amount" label="Payment Due:">
            <b-form-input type="number" v-model="form.amount" placeholder="Amount..."></b-form-input>
          </b-form-group>
          <b-form-group label="cell:">
            <b-form-select v-model="form.cells" :options="this.getPropertyCell" required></b-form-select>
          </b-form-group>
          <b-form-group label="village:">
            <b-form-select v-model="form.village" :options="this.getPropertyVillage" required></b-form-select>
          </b-form-group>
          <div class="buttons">
            <b-button variant="danger" @click="cancel">cancel</b-button>
            <b-button type="submit" variant="primary">Submit</b-button>
          </div>
        </b-form>

        <div class="loader" v-show="loading">
          <pulse-loader :loading="loading" :color="color" :size="size"></pulse-loader>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { setInterval } from "timers";
import axios from "axios";
import { asyncLoading } from "vuejs-loading-plugin";
export default {
  data() {
    return {
      modalShow: false,
      toggleForms: false,
      userAvailable: false,
      loading: false,
      size: "10px",
      color: "#167df0",
      sidebar: {
        active_cell: this.getActiveCell,
        cells_array: [],
        active_village: "",
        village_array: []
      },
      form: {
        name: "",
        fname: "",
        lname: "",
        phoneNo: "",
        amount: "",
        cells: null,
        village: null,
        user_id: ""
      }
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    getActiveCell() {
      this.sidebar.active_cell = this.$store.getters.getActiveCell;
      return this.$store.getters.getActiveCell;
    },
    getActiveVillage() {
      this.sidebar.active_village = this.$store.getters.getActiveVillage;
      return this.$store.getters.getActiveVillage;
    },
    getCellsArray() {
      this.sidebar.cells_array = this.$store.getters.getCellsArray;
      return this.$store.getters.getCellsArray;
    },
    getVillageArray() {
      this.sidebar.village_array = this.$store.getters.getVillageArray;
      return this.$store.getters.getVillageArray;
    },
    getSector() {
      return this.$store.getters.getSector;
    },
    getPropertyCell() {
      return this.cells();
    },
    getPropertyVillage() {
      return this.village();
    }
  },
  mounted() {
    this.getActiveCell;
    this.getCellsArray;
    this.getActivevillage;
    this.getVillageArray;
  },
  methods: {
    update(res) {
      this.$store.dispatch("updatePlace", res).then(() => {
        this.getActiveCell;
        this.getCellsArray;
        this.getActiveVillage;
        this.getVillageArray;
      });
    },
    search(e) {
      e.preventDefault();
      this.loading = true;
      axios
        .get(
          `${this.endpoint}/properties/owners/search/?fname=${this.form.fname}&lname=${this.form.lname}&phone=${this.form.phoneNo}`
        )
        .then(res => {
          this.form.user_id = res.data.id;
          this.toggleForms = true;
          this.loading = false;
          this.$snotify.info(
            `Existing user. proceeding to property registration...`
          );
        })
        .catch(err => {
          this.$snotify.warning(`Oops! user not found. Creating user...`);
          axios
            .post(`${this.endpoint}/properties/owners/`, {
              fname: this.form.fname,
              lname: this.form.lname,
              phone: this.form.phoneNo
            })
            .then(res => {
              this.form.user_id = res.data.id;
              this.toggleForms = true;
              this.loading = false;
              this.$snotify.info(`User created. proceeding to registration...`);
            })
            .catch(err => {
              this.loading = false;
              this.$snotify.error(`Oops! user creation Failed`);
            });
        });
    },

    onSubmit(evt) {
      evt.preventDefault();
      this.loading = true;
      axios
        .post(`${this.endpoint}/properties/`, {
          cell: this.form.cells,
          owner: this.form.user_id,
          due: this.form.amount,
          sector: "remera",
          village: this.form.village
        })
        .then(res => {
          this.$snotify.success(`Property registered successfully`);
          this.loading = false;
          this.modalShow = !this.modalShow;
          (this.form.lname = ""),
            (this.form.fname = ""),
            (this.form.phoneNo = ""),
            (this.form.amount = ""),
            (this.form.cells = ""),
            (this.form.village = "");
          this.toggleForms = false;
        })
        .catch(err => {
          this.$snotify.error(`Oops! property registration Failed`);
          this.loading = false;
        });
    },

    cancel(e) {
      e.preventDefault();

      this.modalShow = !this.modalShow;
      (this.form.lname = ""),
        (this.form.fname = ""),
        (this.form.phoneNo = ""),
        (this.form.amount = ""),
        (this.form.cells = ""),
        (this.form.village = "");
      this.toggleForms = false;
    },

    cells() {
      let main_array = [{ text: "Select cell", value: null }];
      this.getCellsArray.forEach(element => {
        main_array = [...main_array, element];
      });
      return main_array;
    },

    village() {
      let main_array = [{ text: "Select village", value: null }];
      if (this.form.cells == null) {
        return main_array;
      } else if (
        this.form.cells != null &&
        this.getSector[this.form.cells] != undefined
      ) {
        this.getSector[this.form.cells].forEach(element => {
          main_array = [...main_array, element];
        });
        return main_array;
      }
    }
  }
};
</script>
<style scoped>
@import url("../assets/css/dashboardLayout.css");
</style>