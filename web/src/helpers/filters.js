import Vue from "vue";

Vue.filter("number", (value) => {
  if (!value) return 0;
  return Number(value).toLocaleString();
});

Vue.filter("date", (date) => {
  if (!date) return "";
  return new Date(date).toLocaleDateString("en-EN", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
});
