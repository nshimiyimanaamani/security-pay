import Vue from "vue";
var that = Vue.prototype;
const mutations = {
  reset_state(state) {
    state = {
      user: null,
      province: null,
      district: null,
      active_sector: null,
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
  async on_startup(state) {
    state.appLoading = true;
    if (sessionStorage.getItem("token")) {
      const user = await Vue.prototype.$decode(sessionStorage.getItem("token"));
      if (
        user &&
        (user.role == "basic" || user.role === "admin" || user.role === "min")
      ) {
        const account = await user.account.toString().split(".");
        console.log(user, account);
        state.province = Vue.prototype.$capitalize(account[0]);
        state.district = Vue.prototype.$capitalize(account[1]);
        state.active_sector = Vue.prototype.$capitalize(account[2]);

        state.cells_array = await Vue.prototype.$cells(
          state.province,
          state.district,
          state.active_sector
        );

        if (user.role === "basic") {
          const fullUser = await that.axios
            .get("accounts/managers/" + user.username, {
              headers: {
                Authorization: "Bearer " + sessionStorage.getItem("token")
              }
            })
            .then(res => res.data)
            .catch(err => {
              state.appLoading = false;
              that.$snotify.error("Cant retrieve user details");
            });
          state.active_cell = await fullUser.cell;
        } else {
          state.active_cell = await state.cells_array[0];
        }

        state.village_array = await Vue.prototype.$villages(
          state.province,
          state.district,
          state.active_sector,
          state.active_cell
        );

        state.active_village = await state.village_array[0];
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
  }
};
export default mutations;
