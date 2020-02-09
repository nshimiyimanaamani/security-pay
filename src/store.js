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
import state from "./store/state"
import mutations from "./store/mutations"
import actions from "./store/actions"
import getters from "./store/getters"

Vue.use(Vuex);

export const store = new Vuex.Store({
  state,
  mutations,
  actions,
  getters
});
