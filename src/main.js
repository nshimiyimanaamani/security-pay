/* eslint-disable spaced-comment */
/* eslint-disable semi */
/* eslint-disable quotes */
import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import BootstrapVue from "bootstrap-vue";
import PortalVue from "portal-vue";
import { store } from "./store";
import Snotify from "vue-snotify";
import { SnotifyPosition } from "vue-snotify";
import "../node_modules/chart.js/dist/Chart.js";
import axios from "axios";
import VueAxios from "vue-axios";
import titleComponent from "./components/title.vue";
import VueSimpleContextMenu from "vue-simple-context-menu";
import VueSlider from "vue-slider-component";
import "vue-slider-component/theme/default.css";

Vue.component("VueSlider", VueSlider);
Vue.component("vue-title", titleComponent);
Vue.component("vue-simple-context-menu", VueSimpleContextMenu);

Vue.use(VueAxios, axios);
Vue.use(BootstrapVue);
Vue.use(PortalVue);
Vue.use(Snotify, {
  toast: {
    timeout: 3000,
    showProgressBar: false,
    closeOnClick: true,
    position: SnotifyPosition.rightTop
  }
});

axios.interceptors.request.use(
  config => {
    const Baseurl = process.env.VUE_APP_PAYPACK_API;
    config.url = `${Baseurl}${config.url}`;
    return config;
  },
  error => {
    if (navigator.onLine) {
      return Promise.reject(error);
    } else {
      Vue.prototype.$snotify.error("Please connect to the internet");
    }
  }
);
axios.interceptors.response.use(
  response => {
    return response;
  },
  error => {
    if (navigator.online) {
      if (error.response.status === 401) {
        store.dispatch("logout");
      } else {
        return Promise.reject(error);
      }
    } else {
      Vue.prototype.$snotify.error("Please connect to the internet");
    }
  }
);
Vue.config.productionTip = false;

/**
 * @todo Invite fred
 * @body This is just a test to check on the tdo
 */

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
