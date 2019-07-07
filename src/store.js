/* eslint-disable semi */
/* eslint-disable quotes */
import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

export const store = new Vuex.Store({
  state: {
    user1: "sector",
    user2: "cell",
    loading: false
  }
});
