import Vue from 'vue'
import App from './App.vue'
import router from './router'

Vue.config.productionTip = false

/**
 * @todo Invite fred
 * @body This is just a test to check on the tdo
 */

new Vue({
 router,
 render: h => h(App)
}).$mount('#app')
