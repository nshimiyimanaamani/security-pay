import Vue from "vue";
import axios from "axios";

const req = axios.create({
  baseURL: process.env.VUE_APP_PAYPACK_API,
  headers: {
    "Content-Type": "application/json"
  }
});

axios.interceptors.request.use(
  config => {
    const baseURL = process.env.VUE_APP_PAYPACK_API;
    config.url = `${baseURL}${config.url}`;
    return config;
  },

  error => {
    if (!navigator.onLine)
      Vue.prototype.$snotify.error("Please connect to the internet");
    return Promise.reject(error);
  }
);
axios.interceptors.response.use(
  response => {
    return response;
  },
  error => {
    if (error.response && error.response.status === 401) {
      Vue.prototype.$store.dispatch("logout");
    }
    if (!navigator.onLine) {
      Vue.prototype.$snotify.error("Please connect to the internet");
    }
    if (navigator.onLine) {
      var message;
      if (error.response && error.response.data) {
        message =
          error.response.data.error ||
          error.response.data.message ||
          error.response.data;
      } else {
        message = null;
      }
      if (message) Vue.prototype.$snotify.error(message);
    }
  }
);
const requests = {
  get: {},
  post: {
    login(data) {
      return new Promise((resolve, reject) => {
        axios
          .post("/accounts/login", {
            username: data.username,
            password: data.key
          })
          .then(res => {
            sessionStorage.setItem("token", res.data.token);
            resolve(res.data.token);
          })
          .catch(err => {
            console.log(err);
            delete sessionStorage.token;
            const error = err.response
              ? err.response.data.error || err.response.data
              : null;
            if (error) Vue.prototype.$snotify.error(error);
            reject();
          });
      });
    }
  }
};
export default requests;
