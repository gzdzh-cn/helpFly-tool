import { App } from 'vue';
import Chart from './chart/index.vue';

export default {
  install(Vue: App) {
    Vue.component('Chart', Chart);
  },
};
