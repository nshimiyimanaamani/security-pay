/* eslint-disable spaced-comment */
/* eslint-disable semi */
/* eslint-disable quotes */
import Vue from "vue";
import axios from "axios";
import App from "./App.vue";
import router from "./router";
import { store } from "./store";
import "./assets/css/main.scss";
import Scrollspy from "vue2-scrollspy";
import VueParticles from "vue-particles";
import title from "./helpers/title";
import { Provinces, Districts, Sectors, Cells, Villages } from "rwanda"
import VueAxios from "vue-axios";
import PortalVue from "portal-vue";
import { decode } from "jsonwebtoken";
import vueBoostrap from "bootstrap-vue";
import VueSlider from "vue-slider-component";
import titleComponent from "./components/title.vue";
import Snotify, { SnotifyPosition } from "vue-snotify";
import loadingComponent from "./components/loading.vue";
import VueSimpleContextMenu from "./scripts/simplecontextMenu";

Vue.component("VueSlider", VueSlider);
Vue.component("vue-title", titleComponent);
Vue.component("vue-load", loadingComponent);
Vue.component("vue-menu", VueSimpleContextMenu);
var VueScrollTo = require("vue-scrollto");

Vue.use(PortalVue);
Vue.use(vueBoostrap);
Vue.use(VueParticles);
Vue.use(VueScrollTo);
Vue.use(Scrollspy);

Vue.use(Snotify, {
  toast: {
    timeout: 5000,
    showProgressBar: false,
    closeOnClick: true,
    position: SnotifyPosition.rightTop
  }
});

Vue.config.productionTip = false;
Vue.filter("number", value => {
  if (!value) return 0;
  return Number(value).toLocaleString();
});
Vue.filter("date", date => {
  if (!date) return "";
  return new Date(date).toLocaleDateString("en-EN", {
    year: "numeric",
    month: "long",
    day: "numeric"
  });
});
Vue.prototype.$provinces = () => {
  return Provinces();
};
Vue.prototype.$districts = province => {
  if (!province) return Districts();
  return Districts(province);
};
Vue.prototype.$sectors = (province, district) => {
  if (!province && !district) return Sectors();
  return Sectors(province, district);
};
Vue.prototype.$cells = (province, district, sector) => {
  if (!province && !district && !sector) return Cells();
  return Cells(province, district, sector);
};
Vue.prototype.$villages = (province, district, sector, cell) => {
  if (!province && !district && !sector && !cell) return Villages();
  return Villages(province, district, sector, cell);
};
Vue.prototype.$isPhoneNumber = number => {
  const errors = {
    format: false
  };
  if (typeof number !== "string") {
    throw new Error("Input should be string");
  }

  const re = /^(\+?25)?(078|073|072)\d{7}$/;
  if (!re.test(number)) {
    return errors.format;
  }
  return true;
};
Vue.prototype.$capitalize = string => {
  if (!string) return string;
  return String(string).replace(/^./, String(string)[0].toUpperCase());
};
Vue.prototype.$decode = token => {
  return decode(token);
};
Vue.prototype.$getTotal = url => {
  return axiosInstance
    .get(url)
    .then(res => res.data.Total)
    .catch(err => 0);
};
Vue.prototype.$title = name => title(name);
// axios configs
// -----------------------------------------------------------------------------
const { VUE_APP_PAYPACK_API = "/api/" } = process.env;
const axiosInstance = axios.create({
  baseURL: VUE_APP_PAYPACK_API
});
axiosInstance.interceptors.request.use(
  config => {
    if (sessionStorage.getItem("token"))
      axiosInstance.defaults.headers.common["Authorization"] =
        "Bearer " + sessionStorage.getItem("token");
    return config;
  },

  error => {
    if (navigator.onLine === false)
      Vue.prototype.$snotify.error("Please check internet connectivity!");
    return Promise.reject(error);
  }
);

axiosInstance.interceptors.response.use(
  response => {
    return Promise.resolve(response);
  },
  error => {
    if (error.response && error.response.status === 401)
      store.dispatch("logout");

    return Promise.reject(error);
  }
);
Vue.use(VueAxios, axiosInstance);
new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");