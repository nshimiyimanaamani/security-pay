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

const {
  District,
  Sector,
  Cell,
  Village
} = require("rwanda");

Vue.use(Vuex);

export const store = new Vuex.Store({
  state: {
    token: sessionStorage.token ?
      sessionStorage.token : null,
    endpoint: process.env.VUE_APP_PAYPACK_API,
    active_sector: "remera",
    active_cell: "",
    cells_array: [],
    active_village: "",
    village_array: [],
    villageByCell: []
  },
  mutations: {
    on_startup(state) {
      state.cells_array = Cell("Kigali", "Gasabo", state.active_sector).sort();
      state.active_cell = state.cells_array[0];
      state.village_array = Village("Kigali", "Gasabo", state.active_sector, state.active_cell)
      state.active_village = state.village_array[0];
    },
    updatePlace(state, res) {
      if (res.toUpdate == "cell") {
        state.village_array = []
        state.village_array = Village("Kigali", "Gasabo", state.active_sector, res.changed)
        state.active_cell = res.changed
        state.active_village = state.village_array[0]
      } else if (res.toUpdate == "village") {
        state.active_village = res.changed
      }
    },
    villageByCell(state, cell) {
      state.villageByCell = state.sector[cell]
    },
    logout(state) {
      delete sessionStorage.token;
      delete sessionStorage.email;
      state.token = null;
    },
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
      commit("updatePlace", res);
    },
    villageByCell({
      commit
    }, data) {
      commit('villageByCell', data)
    },
    logout({
      commit
    }) {
      commit("logout");
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
    villageByCell: state => state.villageByCell
  }
});
