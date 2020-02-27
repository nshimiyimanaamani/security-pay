/* eslint-disable semi */
/* eslint-disable quotes */
import Vue from "vue";
import Router from "vue-router";
import axios from "axios";
import { store } from "./store";
const login = () =>
  import(/* webpackChunkName: "login-page" */ "./pages/login.vue");
const startPage = () =>
  import(/* webpackChunkName: "login-page" */ "./Layouts/main.vue");
const dashboard = () =>
  import(/* webpackChunkName: "dashboard-page" */ "./pages/dashboard.vue");
const transactions = () =>
  import(/* webpackChunkName: "dashboard-page" */ "./pages/transactions.vue");
const dashboardLayout = () =>
  import(
    /* webpackChunkName: "dashboard-page" */ "./Layouts/dashboardLayout.vue"
  );
const accounts = () =>
  import(/* webpackChunkName: "dashboard-page" */ "./pages/createAccount.vue");
const village = () =>
  import(/* webpackChunkName: "dashboard-page" */ "./pages/village.vue");
const cells = () =>
  import(/* webpackChunkName: "dashboard-page" */ "./pages/cells.vue");
const properties = () =>
  import(/* webpackChunkName: "dashboard-page" */ "./pages/properties.vue");
const feedbacks = () =>
  import(/* webpackChunkName: "dashboard-page" */ "./pages/feedbacks.vue");
const reports = () =>
  import(/* webpackChunkName: "dashboard-page" */ "./pages/reports.vue");
const agentView = () =>
  import(/* webpackChunkName: "agent-page" */ "./Layouts/agentView.vue");
const devLayout = () =>
  import(/* webpackChunkName: "dev-page" */ "./Layouts/Dev/layouts/main.vue");
const DevAccounts = () =>
  import(
    /* webpackChunkName: "dev-page" */ "./Layouts/Dev/layouts/views/account.vue"
  );
const devStats = () =>
  import(
    /* webpackChunkName: "dev-page" */ "./Layouts/Dev/layouts/views/stats.vue"
  );
const notFound = () =>
  import(/* webpackChunkName: "404-page" */ "./pages/404.vue");

Vue.use(Router);
import { decode } from "jsonwebtoken";
let router = new Router({
  mode: "history",
  routes: [
    {
      path: "/",
      component: startPage,
      children: [
        {
          path: "/",
          component: login,
          meta: {
            guest: true
          }
        }
      ]
    },
    {
      path: "/dashboard",
      component: dashboardLayout,
      children: [
        {
          path: "",
          name: "dashboard",
          component: dashboard,
          meta: {
            requireAuth: true,
            forAdmin: true,
            forDev: true
          }
        },
        {
          path: "/transactions",
          name: "transactions",
          component: transactions,
          meta: {
            requireAuth: true,
            forAdmin: true,
            forDev: true,
            forManager: true
          }
        },
        {
          path: "/village",
          name: "village",
          component: village,
          meta: {
            requireAuth: true,
            forAdmin: true,
            forDev: true,
            forManager: true
          }
        },
        {
          path: "/cells",
          name: "cells",
          component: cells,
          meta: {
            requireAuth: true,
            forAdmin: true,
            forDev: true,
            forManager: true
          }
        },
        {
          path: "/properties",
          name: "properties",
          component: properties,
          meta: {
            requireAuth: true,
            forAdmin: true,
            forDev: true,
            forManager: true
          }
        },
        {
          path: "/reports",
          name: "reports",
          component: reports,
          meta: {
            requireAuth: true,
            forAdmin: true,
            forDev: true,
            forManager: true
          }
        },
        {
          path: "/create",
          name: "createAccounts",
          component: accounts,
          meta: {
            requireAuth: true,
            forDev: true,
            forAdmin: true,
            forDev: true
          }
        },
        {
          path: "/feedbacks",
          name: "feedbacks",
          component: feedbacks,
          meta: {
            requireAuth: true,
            forAdmin: true,
            forDev: true,
            forManager: true
          }
        }
      ]
    },
    {
      path: "/agent",
      name: "agentView",
      component: agentView,
      meta: {
        requireAuth: true,
        agent: true
      }
    },
    {
      path: "/dev",
      component: devLayout,
      children: [
        {
          path: "",
          name: "dev-stats",
          component: devStats,
          meta: {
            requireAuth: true,
            forDev: true
          }
        },
        {
          path: "account",
          name: "accounts",
          component: DevAccounts,
          meta: {
            requireAuth: true,
            forDev: true
          }
        }
      ]
    },
    {
      path: "/error",
      name: "not-found",
      component: notFound,
      meta: {
        error: true
      }
    }
  ]
});

router.beforeEach((to, from, next) => {
  const decoded = decode(sessionStorage.token);

  if (to.matched.some(record => record.meta.requireAuth)) {
    if (decoded) {
      next();
      axios.defaults.headers.common["Authorization"] = sessionStorage.token;
      store.state.user = decoded;
      checkRoute(to, next, decoded.role);
    } else {
      next({ path: "/" });
      store.state.user = null;
      delete sessionStorage.token;
    }
  } else if (to.matched.some(record => record.meta.guest)) {
    if (decoded) {
      checkRoute(to, next, decoded.role);
    } else {
      next();
    }
  } else if (to.matched.some(record => record.meta.error)) {
    next();
  } else {
    next({
      path: "/error",
      params: {
        nextUrl: to.fullPath
      }
    });
  }
});
function checkRoute(to, next, role) {
  if (role == "min") {
    if (to.matched.some(record => record.meta.agent)) {
      next();
    } else {
      next({ path: "/agent" });
    }
  }
  if (role == "admin") {
    if (to.matched.some(record => record.meta.forAdmin)) {
      next();
    } else {
      next({ name: "dashboard" });
    }
  }
  if (role == "dev") {
    if (to.matched.some(record => record.meta.forDev)) {
      next();
    } else {
      next({ name: "dashboard" });
    }
  }
  if (role == "basic") {
    if (to.matched.some(record => record.meta.forManager)) {
      next();
    } else {
      next({ name: "cells" });
    }
  }
}
export default router;
