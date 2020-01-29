import router from "../router.js";

const { Cell, Village } = require("rwanda");
const mutations = {
  on_startup(state) {
    state.cells_array = Cell("Kigali", "Gasabo", state.active_sector).sort();
    state.active_cell = state.cells_array[0];
    state.village_array = Village(
      "Kigali",
      "Gasabo",
      state.active_sector,
      state.active_cell
    );
    state.active_village = state.village_array[0];
  },
  updatePlace(state, res) {
    if (res.toUpdate == "cell") {
      state.village_array = [];
      state.village_array = Village(
        "Kigali",
        "Gasabo",
        state.active_sector,
        res.changed
      );
      state.active_cell = res.changed;
      state.active_village = state.village_array[0];
    } else if (res.toUpdate == "village") {
      state.active_village = res.changed;
    }
  },
  villageByCell(state, cell) {
    state.villageByCell = state.sector[cell];
  },
  logout(state) {
    delete sessionStorage.token;
    state.user = null;
    router.push("/");
  }
};
export default mutations;
