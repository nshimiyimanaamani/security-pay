const actions = {
  //startup logic
  // -------------------------------------------
  startup_function({ commit }) {
    return new Promise(async (resolve, reject) => {
      await commit("on_startup");
      resolve();
    });
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
