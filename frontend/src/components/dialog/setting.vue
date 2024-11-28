<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-9
@LastEditors: Hughie
@LastEditTime: 2024-10-10
@Description: This is the dialog allow user to set up the config.
*/

import { visio, editTheme, editLang } from '../../assets/js/globals';
import { reactive, computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { ElMessage, ElNotification } from 'element-plus'
import { GetConfig, SetConfig } from '../../../wailsjs/go/core/App';
import { languages as langs } from '../../assets/js/i18n';

const { t, locale } = useI18n();
let defaultLang = localStorage.getItem('lang');
let defaultTheme =localStorage.getItem('theme');
const changeFlag = ref(false);
const generalSetting = reactive({
    language: defaultLang,
    theme: defaultTheme,
    windowSize: "normal",
    resPort: 7288
});
const windowSetting = reactive({
    GPU : true
})
const linuxSetting = reactive({
    GPUPolicy: "auto"
})
const activeBlock = ref([]);

const initConfig = async () => {
    function handleError(res) {
        if(res.Code == -1) {
            return false;
        }else if (res.Code == 1) {
            ElMessage({
                message: t('message.configLoadError') + ": " + res.Msg,
                type: "error"
            })
            return false;
        }else {
            return true;
        }
    }

    let size = await GetConfig("Appearance", "DefaultOpen").then(res => {
        if(handleError(res)){
            return res.Data;
        }
    });
    let port = await GetConfig("StaticResource", "Port").then(res => {
        if(handleError(res)){
            return res.Data;
        }
    });
    let windowGPU = await GetConfig("Window", "GPUAccelerate").then(res => {
        if(handleError(res)){
            return res.Data;
        }
    });
    let linuxGPU = await GetConfig("Linux", "GPUStrategy").then(res => {
        if(handleError(res)){
            return res.Data;
        }
    });
    generalSetting.windowSize = size;
    generalSetting.resPort = parseInt(port, 10);
    windowSetting.GPU = windowGPU == "true" ? true : false;
    linuxSetting.GPUPolicy = linuxGPU;
}

initConfig()

const themes = computed(() => [
    {label: t('dialog.setting.dark'), value: 'dark'},
    {label: t('dialog.setting.light'), value: 'light'}
]);

const GPUPolicy = computed(()=> [
    {label: t('dialog.setting.auto'), value: 'auto'},
    {label: t('dialog.setting.never'), value: 'never'},
    {label: t('dialog.setting.always'), value: 'always'}
])

const sizes = computed(() => [
    {label: t('dialog.setting.maximised'), value: 'maximised'},
    {label: t('dialog.setting.normal'), value: 'normal'},
    {label: t('dialog.setting.fullScreen'), value: 'fullscreen'},
]);

const saveConfig = async (sector, key, val) => {
    return SetConfig(sector, key, val).then(res => {
        if(res.Code == -1) {
            ElMessage({
                message: t('message.configSaveError') + ": " + res.Msg,
                type: "error"
            })
            return false;
        }
        return true;
    });
}

const handleSaveConfig = async () => {
    localStorage.setItem('lang', generalSetting.language);
    locale.value = generalSetting.language;
    localStorage.setItem('theme', generalSetting.theme);
    document.querySelector('html').setAttribute('class', generalSetting.theme);
    var theme = generalSetting.theme;
    var lang = generalSetting.language;
    editTheme.value = theme;
    editLang.value = lang;
    let flag = true;
    flag = await saveConfig("Appearance", "DefaultOpen", generalSetting.windowSize);
    flag = await saveConfig("StaticResource", "Port", generalSetting.resPort.toString());
    flag = await saveConfig("Window", "GPUAccelerate", windowSetting.GPU === true ? "true" : "false");
    flag = await saveConfig("Linux", "GPUStrategy", linuxSetting.GPUPolicy);
    if (flag && changeFlag.value) {
        ElNotification({
            title: t('message.saveSuccess'),
            message: t('message.configSaveInfo'),
            duration: 2000,
        })
    }
    changeFlag.value = false;
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
    <el-collapse v-model="activeBlock" accordion>
        <el-collapse-item :title="t('dialog.setting.general')" name="1">
            <el-form
            label-position="left"
            label-width="auto"
            :model="generalSetting"
            style="max-width: 400px">
                <el-form-item :label="t('dialog.setting.language')">
                    <el-select
                    v-model="generalSetting.language"
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
                    v-model="generalSetting.theme"
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
                <el-form-item :label="t('dialog.setting.windowSize')">
                    <el-select
                    @change="changeFlag = true"
                    v-model="generalSetting.windowSize"
                    style="width: 240px"
                    >
                    <el-option
                        v-for="item in sizes"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                    />
                    </el-select>
                </el-form-item>
                <el-form-item :label="t('dialog.setting.resPort')">
                    <el-input v-model="generalSetting.resPort" style="width: 240px"></el-input>
                </el-form-item>
            </el-form>
        </el-collapse-item>
        <el-collapse-item :title="t('dialog.setting.window')" name="2">
            <el-form
            label-position="left"
            label-width="auto"
            :model="windowSetting"
            style="max-width: 400px">
                <el-form-item :label="t('dialog.setting.windowGPU')">
                    <el-switch @change="changeFlag = true" v-model="windowSetting.GPU"></el-switch>
                </el-form-item>
            </el-form>
        </el-collapse-item>
        <el-collapse-item :title="t('dialog.setting.linux')" name="3">
            <el-form
            label-position="left"
            label-width="auto"
            :model="linuxSetting"
            style="max-width: 400px">
                <el-form-item :label="t('dialog.setting.linuxGPU')">
                    <el-select
                    v-model="linuxSetting.GPUPolicy"
                    size="large"
                    @change="changeFlag = true"
                    style="width: 240px"
                    >
                    <el-option
                        v-for="item in GPUPolicy"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                    />
                    </el-select>
                </el-form-item>
            </el-form>
        </el-collapse-item>
    </el-collapse>
        
    </el-dialog>
</template>

<style scoped>
.el-collapse-item__header {
    height: 3em !important;
    line-height: 3em !important;
    font-size: 1.2em !important;
    font-weight: bold !important;
}
</style>
