import Vue from "vue";
var that = Vue.prototype;
const mutations = {

  reset_state(state) {
    state.user = null;
    state.province = null;
    state.district = null;
    state.active_sector = null;
    state.active_cell = null;
    state.cells_array = null;
    state.active_village = null;
    state.village_array = null;
    state.villageByCell = null;
    state.months = [
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
      "December",
    ];
  },

  async on_startup(state) {
    state.appLoading = true;
    if (sessionStorage.getItem("token")) {
      const user = await Vue.prototype.$decode(sessionStorage.getItem("token"));
      const roles = ["basic", "admin", "min"];
      if (user && user.role && roles.includes(user.role)) {
        const capitalize = (string) =>
          String(string).replace(/^./, String(string)[0].toUpperCase());

        const account = await user.account.toString().split(".");

        state.province = capitalize(account[0]);
        state.district = capitalize(account[1]);
        state.active_sector = capitalize(account[2]);
        state.cells_array = await Vue.prototype.$cells(
          state.province,
          state.district,
          state.active_sector
        );

        if (user.role === "basic") {
          const fullUser = await that.axios
            .get("accounts/managers/" + user.username, {
              headers: {
                Authorization: "Bearer " + sessionStorage.getItem("token"),
              },
            })
            .then((res) => res.data)
            .catch((err) => {
              console.log(err);
              state.appLoading = false;
              that.$snotify.error("Cant retrieve user details");
            });
          state.active_cell = await fullUser.cell;
        } else {
          state.active_cell = (await state.cells_array) && state.cells_array[0];
        }

        state.village_array = await Vue.prototype.$villages(
          state.province,
          state.district,
          state.active_sector,
          state.active_cell
        );

        state.active_village =
          (await state.village_array) && state.village_array[0];
      }
    }
    state.appLoading = false;
  },
  updatePlace(state, res) {
    if (res.toUpdate == "cell") {
      state.village_array = [];
      state.active_cell = res.changed;

      state.village_array = Vue.prototype.$villages(
        state.province,
        state.district,
        state.active_sector,
        state.active_cell
      );
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
    sessionStorage.clear();
    location.reload();
  },
};
export default mutations;
