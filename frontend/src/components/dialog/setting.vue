<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-9
@LastEditors: Hughie
@LastEditTime: 2024-10-10
@Description: This is the dialog allow user to set up the config.
*/

import { visio, editTheme, editLang, autoSave, generalSetting, linuxSetting, windowSetting, epubLayout } from '../../assets/js/globals';
import { TimerContext } from '../../assets/js/utils';
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { ElMessage, ElNotification } from 'element-plus';
import { GetConfig, SetConfig } from '../../../wailsjs/go/core/App';
import { languages as langs } from '../../assets/js/i18n';

const { t, locale } = useI18n();

const textDirect = computed(()=>[
    {
        value: 'ltr',
        label: t('dialog.setting.textDirec.ltr')
    },
    {
        value: 'rtl',
        label: t('dialog.setting.textDirec.rtl')
    },
    {
        value: 'auto',
        label: t('dialog.setting.textDirec.auto')
    }
])
const layout = computed(()=>[
    {
        value: "reflowable",
        label: t('dialog.setting.layout.reflow')
    },
    {
        value: "pre-paginated",
        label: t('dialog.setting.layout.prePaginated')
    }
])
const flow = computed(()=>[
    {
        value: "paginated",
        label: t('dialog.setting.flow.paginated')
    },
    {
        value: "scrolled-continuous",
        label: t('dialog.setting.flow.scrolledContinuous')
    },
    {
        value: "scrolled-doc",
        label: t('dialog.setting.flow.scrolledDoc')
    },
    {
        value: "auto",
        label: t('dialog.setting.flow.auto')
    }
])
const spread = computed(()=>[
    {
        value: "auto",
        label: t('dialog.setting.spread.auto')
    },
    {
        value: "landscape",
        label: t('dialog.setting.spread.landscape')
    },
    {
        value: "portrait",
        label: t('dialog.setting.spread.portrait')
    },
    {
        value: "both",
        label: t('dialog.setting.spread.both')
    },
    {
        value: "none",
        label: t('dialog.setting.spread.none')
    }
])
const orientation = computed(()=>[
    {
        value: "auto",
        label: t('dialog.setting.orientation.auto')
    },
    {
        value: "landscape",
        label: t('dialog.setting.orientation.landscape')
    },
    {
        value: "portrait",
        label: t('dialog.setting.orientation.portrait')
    }
])
const proportions = computed(()=>[
    {
        value: "width=1200, height=600",
        label: "2:1"
    },
    {
        value: "width=1200, height=800",
        label: "3:2"
    },
    {
        value: "width=1200, height=900",
        label: "4:3"
    },
    {
        value: "width=1200, height=1200",
        label: "1:1"
    },
    {
        value: "width=device-width, height=device-height",
        label: t('dialog.setting.proportions.auto')
    }
])

// if the setting is changed
const changeFlag = ref(false);

const activeBlock = ref([]);

// load the timer
const timer = TimerContext.getInstance(t);

// init the config
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

    // get the config that store from the application's core
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
    let textDir = await GetConfig("Epub", "TextDir").then(res => {
        if(handleError(res)){
            return res.Data;
        }
    });
    let layout = await GetConfig("Epub", "Layout").then(res => {
        if(handleError(res)){
            return res.Data;
        }
    });
    let flow = await GetConfig("Epub", "Flow").then(res => {
        if(handleError(res)){
            return res.Data;
        }
    });
    let spread = await GetConfig("Epub", "Spread").then(res => {
        if(handleError(res)){
            return res.Data;
        }
    });
    let orientation = await GetConfig("Epub", "Orientation").then(res => {
        if(handleError(res)){
            return res.Data;
        }
    });
    let proportions = await GetConfig("Epub", "Proportions").then(res => {
        if(handleError(res)){
            return res.Data;
        }
    });
    autoSave.value.isAutoSave = await GetConfig("Core", "AutoSave").then(res => {
        if(handleError(res)){
            return res.Data == "true" ? true : false;
        }
    });
    autoSave.value.autoSaveTime = await GetConfig("Core", "AutoSaveInterval").then(res => {
        if(handleError(res)){
            return parseInt(parseInt(res.Data));
        }
    });
    generalSetting.windowSize = size;
    generalSetting.resPort = parseInt(port, 10);
    epubLayout.textDirection = textDir;
    epubLayout.layout = layout;
    epubLayout.flow = flow;
    epubLayout.spread = spread;
    epubLayout.orientation = orientation;
    epubLayout.proportions = proportions;
    windowSetting.GPU = windowGPU == "true" ? true : false;
    linuxSetting.GPUPolicy = linuxGPU;
    timer.Start();
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

const saveTimes = [
    {label: "30s", value: 30},
    {label: "1min", value: 60},
    {label: "5min", value: 300},
    {label: "10min", value: 600},
    {label: "30min", value: 1800},
    {label: "1h", value: 3600},
]

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
    document.body.className = lang;
    let flag = true;
    flag = await saveConfig("Appearance", "DefaultOpen", generalSetting.windowSize);
    flag = await saveConfig("StaticResource", "Port", generalSetting.resPort.toString());
    flag = await saveConfig("Window", "GPUAccelerate", windowSetting.GPU === true ? "true" : "false");
    flag = await saveConfig("Linux", "GPUStrategy", linuxSetting.GPUPolicy);
    flag = await saveConfig("Core", "AutoSaveInterval", autoSave.value.autoSaveTime.toString());
    flag = await saveConfig("Core", "AutoSave", autoSave.value.isAutoSave === true ? "true" : "false");
    flag = await saveConfig("Epub", "TextDir", epubLayout.textDirection);
    flag = await saveConfig("Epub", "Layout", epubLayout.layout);
    flag = await saveConfig("Epub", "Flow", epubLayout.flow);
    flag = await saveConfig("Epub", "Spread", epubLayout.spread);
    flag = await saveConfig("Epub", "Orientation", epubLayout.orientation);
    flag = await saveConfig("Epub", "Proportions", epubLayout.proportions);
    if (flag && changeFlag.value) {
        ElNotification({
            title: t('message.saveSuccess'),
            message: t('message.configSaveInfo'),
            duration: 2000,
        })
    }
    changeFlag.value = false;
    visio.settingVisible = false;
    timer.Reset();
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
                <el-form-item :label="t('dialog.setting.windowSize')" prop="windowSize">
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
                <el-form-item :label="t('dialog.setting.autoSave')">
                    <el-switch v-model="autoSave.isAutoSave"></el-switch>
                </el-form-item>
                <el-form-item :label="t('dialog.setting.autoSaveTime')">
                    <el-select
                    v-model="autoSave.autoSaveTime"
                    style="width: 240px"
                    :disabled="!autoSave.isAutoSave"
                    >
                    <el-option
                        v-for="item in saveTimes"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                    />
                    </el-select>
                </el-form-item>
                <el-form-item :label="t('dialog.setting.resPort')">
                    <el-input @change="changeFlag = true" v-model="generalSetting.resPort" style="width: 240px"></el-input>
                </el-form-item>
            </el-form>
        </el-collapse-item>
        <el-collapse-item :title="t('dialog.setting.epubSaving')" name="4">
            <el-form :model="epubLayout"
            label-position="left"
            label-width="auto"
            style="max-width: 400px">
                <el-form-item :label="t('dialog.setting.textDirec.title')">
                    <el-select
                        v-model="epubLayout.textDirection"
                        placeholder="Select"
                        style="width: 240px"
                    >
                        <el-option
                            v-for="item in textDirect"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                        />
                    </el-select>
                </el-form-item>
                <el-form-item :label="t('dialog.setting.layout.title')">
                    <el-select
                        v-model="epubLayout.layout"
                        placeholder="Select"
                        style="width: 240px"
                    >
                        <el-option
                            v-for="item in layout"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                        />
                    </el-select>
                </el-form-item>
                <el-form-item :label="t('dialog.setting.flow.title')">
                    <el-select
                        v-model="epubLayout.flow"
                        :disabled="epubLayout.layout !== 'reflowable'"
                        placeholder="Select"
                        style="width: 240px"
                    >
                        <el-option
                            v-for="item in flow"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                        />
                    </el-select>
                </el-form-item>
                <el-form-item :label="t('dialog.setting.spread.title')">
                    <el-select
                        v-model="epubLayout.spread"
                        placeholder="Select"
                        style="width: 240px"
                    >
                        <el-option
                            v-for="item in spread"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                        />
                    </el-select>
                </el-form-item>
                <el-form-item :label="t('dialog.setting.orientation.title')">
                    <el-select
                        v-model="epubLayout.orientation"
                        placeholder="Select"
                        style="width: 240px"
                    >
                        <el-option
                            v-for="item in orientation"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                        />
                    </el-select>
                </el-form-item>
                <el-form-item :label="t('dialog.setting.proportions.title')">
                    <el-select
                        v-model="epubLayout.proportions"
                        :disabled="epubLayout.layout !== 'pre-paginated'"
                        placeholder="Select"
                        style="width: 240px"
                    >
                        <el-option
                            v-for="item in proportions"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                        />
                    </el-select>
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
