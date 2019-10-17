/* eslint-disable operator-linebreak */
/* eslint-disable eqeqeq */
/* eslint-disable no-array-constructor */
/* eslint-disable camelcase */
/* eslint-disable space-before-function-paren */
/* eslint-disable spaced-comment */
/* eslint-disable standard/computed-property-even-spacing */
/* eslint-disable indent */
/* eslint-disable key-spacing */
/* eslint-disable semi */
/* eslint-disable quotes */
import Vue from "vue";
import Vuex from "vuex";
import {
  Promise
} from "core-js";

Vue.use(Vuex);

export const store = new Vuex.Store({
  state: {
    token: sessionStorage.getItem("token") ?
      sessionStorage.getItem("token") : null,
    endpoint: process.env.VUE_APP_PAYPACK_API,
    isOnline: navigator.onLine,
    active_sector: "remera",
    active_cell: "",
    cells_array: [],
    active_village: "",
    village_array: [],
    sector: {
      "rukiri I": [
        "amajyambere",
        "ubumwe",
        "agashyitsi",
        "gisimenti",
        "urumuri",
        "ukwezi",
        "izuba"
      ],
      "rukiri II": ["rebero", "Amahoro", "Ubumwe", "Ruturusu I", "ruturusu II"],
      Nyarutarama: [
        "Gishushu",
        "Kamahwa",
        "Juru",
        "Kibiraro I",
        "kibiraro II",
        "Kangondo I",
        "kangondo II"
      ],
      Nyabisindu: [
        "Kinunga",
        "Rugarama",
        "Kagara",
        "Gihogere",
        "Nyabisindu",
        "Amarembo I",
        "amarembo II"
      ]
    }
  },
  mutations: {
    on_startup(state) {
      Object.keys(state.sector).forEach(element => {
        state.cells_array = [...state.cells_array, element].sort();
      });
      state.active_cell = state.cells_array[0];
      state.sector[state.active_cell].forEach(element => {
        state.village_array = [...state.village_array, element].sort();
      });
      state.active_village = state.village_array[0];
    },
    updatePlace(state, res) {
      if (res.toUpdate == "cell") {
        //updating cell must update village too
        state.active_cell = state.cells_array[
            state.cells_array.indexOf(res.changed)
          ] ?
          (state.active_cell =
            state.cells_array[state.cells_array.indexOf(res.changed)]) :
          (state.active_cell = state.cells_array[0]);

        let village_array = new Array(); //start updating villages
        state.sector[state.active_cell].forEach(element => {
          village_array = [...village_array, element];
        });
        state.village_array = village_array;
        state.active_village = state.village_array[0];
      } else if (res.toUpdate == "village") {
        //updating village only
        let village_array = new Array();
        state.sector[state.active_cell].forEach(element => {
          village_array = [...village_array, element];
        });
        state.village_array = village_array;
        state.active_village =
          state.village_array[state.village_array.indexOf(res.changed)];
      }
    },
    logout(state) {
      delete sessionStorage.getItem("token");
      state.token = null;
    },
    checkConnection(state, data) {
      state.isOnline = data;
    }
  },
  actions: {
    startup_function({
      commit
    }) {
      commit("on_startup");
    },
    updatePlace({
      commit
    }, res) {
      return new Promise(resolve => {
        commit("updatePlace", res);
        resolve();
      });
    },
    logout({
      commit
    }) {
      commit("logout");
    },
    checkConnection({
      commit
    }) {
      return new Promise(resolve => {
        commit("checkConnection", navigator.onLine);
        resolve(navigator.onLine);
      });
    }
  },
  getters: {
    getEndpoint: state => state.endpoint,
    getSectorArray: state => state.sector,
    getCellsArray: state => state.cells_array,
    getActiveCell: state => state.active_cell,
    getActiveVillage: state => state.active_village,
    getVillageArray: state => state.village_array,
    getActiveSector: state => state.active_sector,
    token: state => state.token,
    isOnline: state => state.isOnline
  }
});
