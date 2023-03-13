import Vue from "vue";
import axios from "axios";
import VueAxios from "vue-axios";
import { store } from "../store";

const { VUE_APP_PAYPACK_API = "/api/" } = process.env;
const axiosInstance = axios.create({
  baseURL: VUE_APP_PAYPACK_API,
});
axiosInstance.interceptors.request.use(
  (config) => {
    if (sessionStorage.getItem("token"))
      axiosInstance.defaults.headers.common["Authorization"] =
        "Bearer " + sessionStorage.getItem("token");
    return config;
  },

  (error) => {
    if (navigator.onLine === false)
      Vue.prototype.$snotify.error("Please check internet connectivity!");
    return Promise.reject(error);
  }
);

axiosInstance.interceptors.response.use(
  (response) => {
    return Promise.resolve(response);
  },
  (error) => {
    if (error.response && error.response.status === 401)
      store.dispatch("logout");

    return Promise.reject(error);
  }
);
Vue.use(VueAxios, axiosInstance);

export default axiosInstance;
