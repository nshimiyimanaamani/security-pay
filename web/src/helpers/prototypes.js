import Vue from "vue";
import { decode } from "jsonwebtoken";
import axiosInstance from "./axios";

Vue.prototype.$isPhoneNumber = (number) => {
  const errors = {
    format: false,
  };
  if (typeof number !== "string") {
    throw new Error("Input should be string");
  }

  const re = /^(\+?25)?(078|073|072|079)\d{7}$/;
  if (!re.test(number)) {
    return errors.format;
  }
  return true;
};

Vue.prototype.$capitalize = (string) => {
  if (!string) return string;
  return String(string).replace(/^./, String(string)[0].toUpperCase());
};

Vue.prototype.$decode = (token) => {
  return decode(token);
};

Vue.prototype.$getTotal = (url) => {
  return axiosInstance
    .get(url)
    .then((res) => res.data.Total)
    .catch((err) => 0);
};

Vue.prototype.$title = (name) => {
  if(!name) return "Paypack";
  return `Paypack | ${name}`;
}
