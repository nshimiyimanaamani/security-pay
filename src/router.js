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
  import(/* webpackChunkName: "admin-page" */ "./pages/dashboard.vue");
const transactions = () =>
  import(/* webpackChunkName: "transactions" */ "./pages/transactions.vue");
const dashboardLayout = () =>
  import(/* webpackChunkName: "layouts" */ "./Layouts/dashboardLayout.vue");
const accounts = () =>
  import(/* webpackChunkName: "account" */ "./pages/createAccount.vue");
const village = () =>
  import(/* webpackChunkName: "village" */ "./pages/village.vue");
const cells = () => import(/* webpackChunkName: "cells" */ "./pages/cells.vue");
const properties = () =>
  import(/* webpackChunkName: "properties" */ "./pages/properties.vue");
const feedbacks = () =>
  import(/* webpackChunkName: "feedbacks" */ "./pages/feedbacks.vue");
const reports = () =>
  import(/* webpackChunkName: "reports" */ "./pages/reports.vue");
const agentView = () =>
  import(/* webpackChunkName: "agents" */ "./Layouts/agentView.vue");
const devLayout = () =>
  import(/* webpackChunkName: "developers" */ "./Layouts/Dev/layouts/main.vue");
const developers = () =>
  import(
    /* webpackChunkName: "developers" */ "./Layouts/Dev/layouts/views/developers.vue"
  );
const dev_accounts = () =>
  import(
    /* webpackChunkName: "developers" */ "./Layouts/Dev/layouts/views/accounts.vue"
  );
const dev_admins = () =>
  import(
    /* webpackChunkName: "developers" */ "./Layouts/Dev/layouts/views/admins.vue"
  );
const dev_managers = () =>
  import(
    /* webpackChunkName: "developers" */ "./Layouts/Dev/layouts/views/managers.vue"
  );
const notFound = () =>
  import(/* webpackChunkName: "404-page" */ "./pages/404.vue");
const message = () =>
  import(/* webpackChunkName: "message" */ "./pages/messages.vue");

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
            forAdmin: true
          }
        },
        {
          path: "/transactions",
          name: "transactions",
          component: transactions,
          meta: {
            requireAuth: true,
            forAdmin: true,
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
            forManager: true
          }
        },
        {
          path: "/create",
          name: "createAccounts",
          component: accounts,
          meta: {
            requireAuth: true,
            forAdmin: true,
            forManager: true
          }
        },
        {
          path: "/feedbacks",
          name: "feedbacks",
          component: feedbacks,
          meta: {
            requireAuth: true,
            forAdmin: true,
            forManager: true
          }
        },
        {
          path: "/message",
          name: "messages",
          component: message,
          meta: {
            requireAuth: true,
            forAdmin: true,
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
          name: "accounts",
          component: dev_accounts,
          meta: {
            requireAuth: true,
            forDev: true
          }
        },
        {
          path: "developers",
          name: "developers",
          component: developers,
          meta: {
            requireAuth: true,
            forDev: true
          }
        },
        {
          path: "admins",
          name: "devAdmins",
          component: dev_admins,
          meta: {
            requireAuth: true,
            forDev: true
          }
        },
        {
          path: "managers",
          name: "devManagers",
          component: dev_managers,
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
      store.state.user = decoded;
      next();
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
      next({ path: "/dev" });
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
