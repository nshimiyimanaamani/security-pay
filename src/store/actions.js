const actions = {
  //startup logic
  // -------------------------------------------
  startup_function({ commit }) {
    commit("on_startup");
  },
  updatePlace({ commit }, res) {
    commit("updatePlace", res);
  },
  //logout logic
  //------------------------------------------
  logout({ commit }) {
    commit("logout");
  }
};
export default actions;
