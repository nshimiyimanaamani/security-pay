import Vue from "vue";
import { LayoutPlugin, CardPlugin, ModalPlugin } from "bootstrap-vue";
Vue.use(LayoutPlugin);
Vue.use(ModalPlugin);
Vue.use(CardPlugin);

import { SpinnerPlugin } from "bootstrap-vue";
Vue.use(SpinnerPlugin);

import {
  FormSelectPlugin,
  FormGroupPlugin,
  ButtonPlugin,
  FormPlugin
} from "bootstrap-vue";
Vue.use(ButtonPlugin);
Vue.use(FormSelectPlugin);
Vue.use(FormGroupPlugin);
Vue.use(FormPlugin);

import { CollapsePlugin, DropdownPlugin } from "bootstrap-vue";
Vue.use(DropdownPlugin);
Vue.use(CollapsePlugin);

import { TablePlugin } from "bootstrap-vue";
Vue.use(TablePlugin);

import { AlertPlugin } from "bootstrap-vue";
Vue.use(AlertPlugin);

import { BadgePlugin } from "bootstrap-vue";
Vue.use(BadgePlugin);
