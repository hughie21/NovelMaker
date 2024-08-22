/* 
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-08-16
@Description: Loading and initializing various plugins and data to the frontend.
*/

import {createApp} from 'vue'
import App from '../../App.vue'
import '../css/style.css';
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import i18n from './i18n';
import '@imengyu/vue3-context-menu/lib/vue3-context-menu.css'
import ContextMenu from '@imengyu/vue3-context-menu'
import {checkIfOpenFileDirectly} from './utils'

// initialize the theme and language settings
if(!localStorage.getItem('lock')){
    localStorage.setItem("lang", navigator.language);
    localStorage.setItem("lock", true);
    localStorage.setItem("theme", "dark");
    location.reload();
}

checkIfOpenFileDirectly();

const app = createApp(App);
app.use(ContextMenu).use(ElementPlus).use(i18n).mount('#app');
