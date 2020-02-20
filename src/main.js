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
import "chart.js/dist/Chart";
import axios from "axios";
import VueAxios from "vue-axios";
import titleComponent from "./components/title.vue";
import VueSimpleContextMenu from "./scripts/simplecontextMenu";
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

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
