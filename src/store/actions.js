import request from "../scripts/requests";
import { decode } from "jsonwebtoken";
const actions = {
  //startup logic
  startup_function({ commit }) {
    commit("on_startup");
  },
  updatePlace({ commit }, res) {
    commit("updatePlace", res);
  },
  //end of startup logic
  //---------------------------------------------------------------------------
  //login logic
  login({ commit }, data) {
    return new Promise((resolve, reject) => {
      if (data.username && data.key) {
        request.post
          .login(data)
          .then(res => {
            const user = decode(res);
            commit("set_user", user);
            resolve();
          })
          .catch(() => {
            reject();
          });
      }
    });
  },
  //End of login logic
  //--------------------------------------------------------------------------
  //logout logic
  logout({ commit }) {
    delete sessionStorage.token;
    location.reload();
    commit("set_user", null);
    commit("reset_state");
  }

  //End of logout logic
};
export default actions;
