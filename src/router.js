/* eslint-disable semi */
/* eslint-disable quotes */
import Vue from "vue";
import Router from "vue-router";
import axios from "axios";
import { store } from "./store";
import login from "./pages/login.vue";
import register from "./pages/register.vue";
import startPage from "./Layouts/main.vue";
import dashboard from "./pages/dashboard.vue";
import transactions from "./pages/transactions.vue";
import dashboardLayout from "./Layouts/dashboardLayout.vue";
import accounts from "./pages/createAccount.vue";
import village from "./pages/village.vue";
import cells from "./pages/cells.vue";
import properties from "./pages/properties.vue";
import agentView from "./Layouts/agentView.vue";
import feedbacks from "./pages/feedbacks.vue";
import reports from "./pages/reports.vue";
import devLayout from "./Layouts/Dev/layouts/main.vue";
import DevAccounts from "./Layouts/Dev/layouts/views/account.vue";
import devStats from "./Layouts/Dev/layouts/views/stats.vue";
import notFound from "./pages/404.vue";

Vue.use(Router);
var jwt = require("jsonwebtoken");

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
  const decoded = jwt.decode(sessionStorage.token);

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
