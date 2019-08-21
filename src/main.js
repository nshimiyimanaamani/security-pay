/* eslint-disable spaced-comment */
/* eslint-disable semi */
/* eslint-disable quotes */
import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import BootstrapVue from "bootstrap-vue";
import {
  store
} from "./store";
import VueLoading from 'vuejs-loading-plugin'
import ClipLoader from 'vue-spinner/src/ClipLoader.vue'
import pulseLoader from 'vue-spinner/src/PulseLoader.vue'
import Snotify from 'vue-snotify';
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import {
  SnotifyPosition
} from 'vue-snotify'
import "../node_modules/chart.js/dist/Chart.js";
import axios from 'axios'
import VueAxios from 'vue-axios'

Vue.component('clip-loader', ClipLoader);
Vue.component('pulse-loader', pulseLoader);

Vue.use(VueAxios, axios)
Vue.use(BootstrapVue);
Vue.use(VueLoading)
Vue.use(Snotify, {
  toast: {
    timeout: 3000,
    showProgressBar: false,
    closeOnClick: true,
    position: SnotifyPosition.rightTop
  }
})

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
