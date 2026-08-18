import { createApp } from 'vue'
import App from './App.vue'
import Router from './router/index';
import store from './store';
import ArcoVue from '@arco-design/web-vue';
import ArcoVueIcon from '@arco-design/web-vue/es/icon';
import '@arco-design/web-vue/dist/arco.css';
// import '@/assets/style/global.less';
import components from './components/icons';
import globalComponents from '@/components';
const app = createApp(App);
app.use(Router);
app.use(store);
app.use(ArcoVue, {});
app.use(ArcoVueIcon);
app.use(globalComponents);
// components
for (const i in components) {
 app.component(i, components[i])
}
app.mount('#app');
