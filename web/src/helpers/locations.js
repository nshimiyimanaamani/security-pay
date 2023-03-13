import Vue from "vue";
import { Provinces, Districts, Sectors, Cells, Villages } from "rwanda";

Vue.prototype.$provinces = () => {
  return Provinces();
};
Vue.prototype.$districts = (province) => {
  if (!province) return Districts();
  return Districts(province);
};
Vue.prototype.$sectors = (province, district) => {
  return !province && !district ? Sectors(province, district) : Sectors();
  // if (!province && !district) return Sectors();
  // return Sectors(province, district);
};
Vue.prototype.$cells = (province, district, sector) => {
  return !province && !district && !sector
    ? Cells(province, district, sector)
    : Cells();
};
Vue.prototype.$villages = (province, district, sector, cell) => {
  return !province && !district && !sector && !cell
    ? Villages(province, district, sector, cell)
    : Villages();
  // if (!province && !district && !sector && !cell) return Villages();
  // return Villages(province, district, sector, cell);
};
