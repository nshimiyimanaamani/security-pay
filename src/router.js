/* eslint-disable semi */
/* eslint-disable quotes */
import Vue from "vue";
import Router from "vue-router";
import login from "./pages/login.vue";
import register from "./pages/register.vue";
import startPage from "./Layouts/main.vue";
import dashboard from "./pages/dashboard.vue";
import transactions from "./pages/transactions.vue";
import dashboardLayout from "./Layouts/dashboardLayout.vue";
import village from "./pages/village.vue";

Vue.use(Router);

export default new Router({
  mode: "history",
  routes: [
    {
      path: "/",
      name: "main",
      component: startPage,
      children: [
        {
          path: "/",
          name: "login",
          component: login
        },
        {
          path: "/register",
          name: "register",
          component: register
        }
      ]
    },
    {
      path: "/dashboard",
      name: "dashboardLayout",
      component: dashboardLayout,
      children: [
        {
          path: "/dashboard",
          name: "dashboard",
          component: dashboard
        },
        {
          path: "/transactions",
          name: "transactions",
          component: transactions
        },
        {
          path: "/village",
          name: "village",
          component: village
        }
      ]
    }
  ]
});
