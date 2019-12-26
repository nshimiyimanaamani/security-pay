/* eslint-disable semi */
/* eslint-disable quotes */
import Vue from "vue";
import Router from "vue-router";
import axios from 'axios';
import {
  store
} from './store'
import login from "./pages/login.vue";
import register from "./pages/register.vue";
import startPage from "./Layouts/main.vue";
import dashboard from "./pages/dashboard.vue";
import transactions from "./pages/transactions.vue";
import dashboardLayout from "./Layouts/dashboardLayout.vue";
import village from "./pages/village.vue";
import cells from "./pages/cells.vue";
import reports from './pages/reports.vue'

Vue.use(Router);

var jwt = require("jsonwebtoken");

let router = new Router({
  mode: "history",
  routes: [{
      path: "/",
      name: "",
      component: startPage,
      children: [{
          path: "/",
          name: "login",
          component: login,
          meta: {
            guest: true
          }
        },
        {
          path: "/register",
          name: "register",
          component: register,
          meta: {
            guest: true
          }
        }
      ]
    },
    {
      path: "/dashboard",
      name: "dashboardLayout",
      component: dashboardLayout,
      children: [{
          path: "/dashboard",
          name: "dashboard",
          component: dashboard,
          meta: {
            requireAuth: true
          }
        },
        {
          path: "/transactions",
          name: "transactions",
          component: transactions,
          meta: {
            requireAuth: true
          }
        },
        {
          path: "/village",
          name: "village",
          component: village,
          meta: {
            requireAuth: true
          }
        },
        {
          path: "/cells",
          name: "cells",
          component: cells,
          meta: {
            requireAuth: true
          }
        },
        {
          path: '/reports',
          name: 'reports',
          component: reports,
          meta: {
            requireAuth: true
          }
        }
      ]
    }
  ]
});

router.beforeEach((to, from, next) => {
  const decoded = jwt.decode(sessionStorage.token);
  if (to.matched.some(record => record.meta.requireAuth)) {
    if (!decoded) {
      store.state.user = null
      next({
        path: '/',
        params: {
          nextUrl: to.fullPath
        }
      })
    } else {
      axios.defaults.headers.common['Authorization'] = sessionStorage.token;
      store.state.user = decoded
      next()
    }
  } else if (to.matched.some(record => record.meta.guest)) {
    if (!decoded) {
      delete sessionStorage.token
      next()
    } else {
      next({
        name: 'dashboard'
      })
    }
  } else {
    next()
  }
})

export default router
