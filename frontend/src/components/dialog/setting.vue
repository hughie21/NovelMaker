<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-9
@LastEditors: Hughie
@LastEditTime: 2024-10-10
@Description: This is the dialog allow user to set up the config.
*/

import { visio, editTheme, editLang } from '../../assets/js/utils';
import { reactive } from 'vue';
import { useI18n } from 'vue-i18n';
import { languages as langs } from '../../assets/js/i18n';

const { t, locale } = useI18n();
const themes = [
    {label: t('dialog.setting.dark'), value: 'dark'},
    {label: t('dialog.setting.light'), value: 'light'}
];
let defaultLang = localStorage.getItem('lang');
let defaultTheme =localStorage.getItem('theme');
const setting = reactive({
    language: defaultLang,
    theme: defaultTheme
});


const handleSaveConfig = ()=> {
    localStorage.setItem('lang', setting.language);
    locale.value = setting.language;
    localStorage.setItem('theme', setting.theme);
    document.querySelector('html').setAttribute('class', setting.theme);
    var theme = setting.theme;
    var lang = setting.language;
    editTheme.value = theme;
    editLang.value = lang;
    visio.settingVisible = false;
}
</script>

<template>
    <el-dialog
    v-model="visio.settingVisible"
    :title="t('dialog.setting.title')"
    width="500"
    :before-close="handleSaveConfig"
    >
        <el-form
        label-position="left"
        label-width="auto"
        :model="setting"
        style="max-width: 400px">
            <el-form-item :label="t('dialog.setting.language')">
                <el-select
                v-model="setting.language"
                size="large"
                style="width: 240px"
                >
                <el-option
                    v-for="item in langs"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                />
                </el-select>
            </el-form-item>
            <el-form-item :label="t('dialog.setting.theme')">
                <el-select
                v-model="setting.theme"
                size="large"
                style="width: 240px"
                >
                <el-option
                    v-for="item in themes"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                />
                </el-select>
            </el-form-item>
        </el-form>
    </el-dialog>
</template>
