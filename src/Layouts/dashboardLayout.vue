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
              <b-form-input v-model="form1.fname" required placeholder="First Name..."></b-form-input>
              <b-form-input v-model="form1.lname" required placeholder="Last Name..."></b-form-input>
            </div>
          </b-form-group>
          <b-form-group label="Phone number:">
            <b-form-input type="number" v-model="form1.phoneNo"></b-form-input>
          </b-form-group>

          <div class="buttons">
            <b-button type="submit" variant="primary" @click="search">Search</b-button>
            <b-button variant="danger" @click="cancel">cancel</b-button>
          </div>
        </b-form>

        <b-form @submit="onSubmit" v-show="toggleForms" class="form2">
         <b-form-group class="phone" label="First Name:">
            <b-form-input
              :disabled="userAvailable"
              v-model="form.fname"
              placeholder="Firse Name..."
            ></b-form-input>
          </b-form-group>
          <b-form-group class="amount" label="Surname:">
            <b-form-input :disabled="userAvailable" v-model="form.lname" placeholder="surname..."></b-form-input>
          </b-form-group>
          <b-form-group class="phone" label="Phone number:">
            <b-form-input
              type="number"
              :disabled="userAvailable"
              v-model="form.phone"
              placeholder="Phone Number..."
            ></b-form-input>
          </b-form-group>
          <b-form-group class="amount" label="Payment Due:">
            <b-form-input type="number" v-model="form.amount" placeholder="Amount..."></b-form-input>
          </b-form-group>
          <b-form-group label="cell:">
            <b-form-select v-model="form.cells" :options="cells()" required></b-form-select>
          </b-form-group>
          <b-form-group label="village:">
            <b-form-select v-model="form.village" :options="village()" required></b-form-select>
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
        phone: "",
        amount: "",
        cells: null,
        village: null,
        user_id: ""
      },
      form1: {
        fname: "Tucky",
        lname: "Bucky",
        phoneNo: "0784577882"
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
      this.userAvailable = false;
      this.loading = true;
      axios
        .get(
          `${this.endpoint}/properties/owners/search/?fname=${this.form1.fname}&lname=${this.form1.lname}&phone=${this.form1.phoneNo}`
        )
        .then(res => {
          this.userAvailable = false;
          this.form.name = `${res.data.fname} ${res.data.lname}`;
          this.form.phone = res.data.phone;
          this.form.user_id = res.data.id;
          this.toggleForms = true;
          console.log(res);
          // resole();
        })
        .catch(err => {
          console.log(err);
          this.toggleForms = true;
          this.loading = false;

          // reject();
        });
      // });
      // asyncLoading(search)
      //   .then(console.log("search finish"))
      //   .catch(console.log("search error"));
    },
    onSubmit(evt) {
      evt.preventDefault();
      axios
        .post(`${this.endpoint}/properties/`, {
          cell: this.form.cells,
          owner: this.form.user_id,
          due: this.form.amount,
          sector: "remera",
          village: this.form.village
        })
        .then(res => {
          console.log(res.data);
        })
        .catch(err => {
          console.log(err);
        });
    },
    cancel(e) {
      e.preventDefault();

      this.modalShow = !this.modalShow;
      (this.form.email = ""),
        (this.form.name = ""),
        (this.form.phone = ""),
        (this.form.amount = ""),
        (this.form.cells = ""),
        (this.form.village = "");
      this.toggleForms = false;
    },
    cells() {
      let main_array = [{ text: "Select cell", value: null }];
      for (const key in this.getSector) {
        if (this.getSector.hasOwnProperty(key)) {
          main_array = [...main_array, key];
        }
      }
      return main_array;
    },
    village() {
      let main_array = [{ text: "Select village", value: null }];
      if (
        this.getSector[this.form.cells] != undefined &&
        this.getSector[this.form.cells] != ""
      ) {
        this.getSector[this.form.cells].forEach(element => {
          main_array = [...main_array, element];
        });
      }
      return main_array;
    }
  }
};
</script>
<style scoped>
@import url("../assets/css/dashboardLayout.css");
</style>