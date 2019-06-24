/* eslint-disable spaced-comment */
/* eslint-disable semi */
/* eslint-disable quotes */
import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import BootstrapVue from "bootstrap-vue";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";
import "vue-select/dist/vue-select.css";
import vSelect from "vue-select";

Vue.component("v-select", vSelect);

Vue.use(BootstrapVue);

Vue.config.productionTip = false;

/**
 * @todo Invite fred
 * @body This is just a test to check on the tdo
 */

new Vue({
  router,
  render: h => h(App)
}).$mount("#app");
