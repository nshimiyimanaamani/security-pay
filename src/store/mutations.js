import { Cell, Village } from "rwanda";
const mutations = {
  reset_state(state) {
    state = {
      user: null,
      active_sector: "Remera",
      active_cell: null,
      cells_array: null,
      active_village: null,
      village_array: null,
      villageByCell: null,
      months: [
        "January",
        "February",
        "March",
        "April",
        "May",
        "June",
        "July",
        "August",
        "September",
        "October",
        "November",
        "December"
      ]
    };
  },
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
  set_user(state, user) {
    state.user = user ? new Object(user) : null;
  },
  logout(state) {
    state.user = null;
  }
};
export default mutations;
