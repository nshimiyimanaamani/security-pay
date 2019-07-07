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
            <b-dropdown id="dropdown-right" class="m-md-2" v-if="this.$route.name == 'cells'">
              <template slot="button-content">
                <span class="fa fa-caret-down"></span>
              </template>
              <!-- <v-select></v-select> -->
              <b-dropdown-item>cell 1</b-dropdown-item>
              <b-dropdown-item>cell 2</b-dropdown-item>
              <b-dropdown-item>cell 4</b-dropdown-item>
              <b-dropdown-item>cell 5</b-dropdown-item>
              <b-dropdown-item>cell 6</b-dropdown-item>
              <b-dropdown-item>cell 7</b-dropdown-item>
              <b-dropdown-item>cell 8</b-dropdown-item>
              <b-dropdown-item>cell 9</b-dropdown-item>
              <b-dropdown-item>cell 10</b-dropdown-item>
            </b-dropdown>
          </li>
        </router-link>
        <router-link to="/village">
          <li>
            Village
            <b-dropdown id="dropdown-right" class="m-md-2" v-if="this.$route.name == 'village'">
              <template slot="button-content">
                <span class="fa fa-caret-down"></span>
              </template>
              <!-- <v-select></v-select> -->
              <b-dropdown-item>village 1</b-dropdown-item>
              <b-dropdown-item>village 2</b-dropdown-item>
              <b-dropdown-item>village 4</b-dropdown-item>
              <b-dropdown-item>village 5</b-dropdown-item>
              <b-dropdown-item>village 6</b-dropdown-item>
              <b-dropdown-item>village 7</b-dropdown-item>
              <b-dropdown-item>village 8</b-dropdown-item>
              <b-dropdown-item>village 9</b-dropdown-item>
              <b-dropdown-item>village 10</b-dropdown-item>
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
      <b-modal
        no-fade
        no-close-on-esc
        no-close-on-backdrop
        hide-header-close
        hide-footer
        id="modal"
        v-model="modalShow"
        title="Add a Property"
      >
        <b-form @submit="onSubmit">
          <b-form-group label="Owner:">
            <b-form-input v-model="form.name" required placeholder="Enter name . . ."></b-form-input>
          </b-form-group>
          <b-form-group label="Email address:">
            <b-form-input v-model="form.email" type="email" placeholder="Enter email"></b-form-input>
          </b-form-group>
          <b-form-group class="phone" label="Phone number:">
            <b-form-input type="number" v-model="form.phone"></b-form-input>
          </b-form-group>
          <b-form-group class="amount" label="Payment Due:">
            <b-form-input type="number" v-model="form.amount"></b-form-input>
          </b-form-group>
          <b-form-group label="cell:">
            <b-form-select v-model="form.cells" :options="cells" required></b-form-select>
          </b-form-group>
          <b-form-group label="village:">
            <b-form-select v-model="form.village" :options="village" required></b-form-select>
          </b-form-group>
          <div class="buttons">
            <b-button variant="danger" @click="cancel">cancel</b-button>
            <b-button type="submit" variant="primary">Submit</b-button>
          </div>
        </b-form>
      </b-modal>
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
  </div>
</template>

<script>
import { setInterval } from "timers";
import axios from "axios";
export default {
  data() {
    return {
      modalShow: false,
      form: {
        email: "",
        name: "",
        phone: "",
        amount: "",
        cells: null,
        village: null
      },
      village: [
        { text: "Select village", value: null },
        "village1",
        "village2",
        "village3",
        "village4"
      ],
      cells: [
        { text: "Select cells", value: null },
        "cell1",
        "cell2",
        "cell3",
        "cell4"
      ]
    };
  },
  methods: {
    onSubmit(evt) {
      evt.preventDefault();
      if (this.form.phone.length < 10) {
        var phone = document.querySelector(".phone input");
        phone.className = "phone-notify";
        setInterval(() => {
          phone.removeAttribute("class");
          phone.className = "form-control";
        }, 5000);
        return;
      }

      if (this.form.email == "") {
        this.form.email = "N/A";
      }
      console.log(this.form.phone.length);
      alert(JSON.stringify(this.form));
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
    }
  }
};
</script>
<style scoped>
@import url("../assets/css/dashboardLayout.css");
</style>