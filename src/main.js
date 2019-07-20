/* eslint-disable spaced-comment */
/* eslint-disable semi */
/* eslint-disable quotes */
import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import BootstrapVue from "bootstrap-vue";
import { store } from "./store";
import VueLoading from 'vuejs-loading-plugin'
// import Buefy from 'buefy'
import 'buefy/dist/buefy.css'
import ClipLoader from 'vue-spinner/src/ClipLoader.vue'
import pulseLoader from 'vue-spinner/src/PulseLoader.vue'

import "vue-select/dist/vue-select.css";
import vSelect from "vue-select";

Vue.component("v-select", vSelect);
Vue.component('clip-loader', ClipLoader);
Vue.component('pulse-loader', pulseLoader);

Vue.use(BootstrapVue);
Vue.use(VueLoading)
// Vue.use(Buefy);

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
