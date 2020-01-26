const actions = {
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
}
export default actions
