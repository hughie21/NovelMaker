/* 
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-07-28
@Description: loading the languages to the i18n plugin.
*/

import en from '../language/en'
import zh_cn from '../language/zh_cn'
import {createI18n} from 'vue-i18n'

const message = {
    "en": en,
    "zh-CN": zh_cn
};

const languages = [
    {
        value: "en",
        label: "English"
    },
    {
        value: "zh-CN",
        label: "简体中文"
    }
];

const i18n = createI18n({
    locale: localStorage.getItem("lang"),
    legacy: false,
    messages: message
});

export default i18n;

export {languages}