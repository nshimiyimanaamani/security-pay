/* eslint-disable semi */
/* eslint-disable quotes */
import Vue from "vue";
import Router from "vue-router";
import axios from 'axios';
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
  if (to.matched.some(record => record.meta.requireAuth)) {
    if (!sessionStorage.token) {
      next({
        path: '/',
        params: {
          nextUrl: to.fullPath
        }
      })
    } else {
      axios.defaults.headers.common['Authorization'] = sessionStorage.token;
      next()
    }
  } else if (to.matched.some(record => record.meta.guest)) {
    if (!sessionStorage.token) {
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
