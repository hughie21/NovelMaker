<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-24
@LastEditors: Hughie
@LastEditTime: 2024-11-1
@Description: This is header of the application.
*/

import { ref, onMounted } from 'vue';
import * as runtime from "../wailsjs/runtime/runtime";
import { FileSave } from '../wailsjs/go/main/App';
import { useI18n } from 'vue-i18n';
import { change, title, bookInfo, editorRef } from './assets/js/globals';
import { ElMessage, ElMessageBox } from 'element-plus';

const { t } = useI18n();
const activeIndex = ref('1');
const maximised = ref(false);
const E = editorRef.value;

const shortCutState = {
    'CONTROL+N': () => {
        menu.new();
    },
    'CONTROL+O': () => {
        menu.open();
    },
    'CONTROL+SHIFT+O': () => {
        menu.import();
    },
    'CONTROL+S': () => {
        menu.save();
    },
    'CONTROL+SHIFT+S': () => {
        menu.saveAs();
    },
    'CONTROL+E': () => {
        menu.export();
    },
    'CONTROL+F': () => {
        menu.lookUp();
    },
    'CONTROL+H': () => {
        menu.replace();
    }
}

let keys = new Set();
let timer = null;

document.addEventListener('keydown', (e) => {
    keys.add(e.key.toUpperCase());
    clearTimeout(timer);
    timer = setTimeout(() => {
        keys.clear();
    }, 300);

    let shortCut = Array.from(keys).join("+");
    if (shortCutState[shortCut]) {
        shortCutState[shortCut]();
        keys.clear();
    }
});

const menu = {
    new : () => {
        document.getElementById('btn-new').click();
    },
    open : () => {
        document.getElementById('btn-open').click();
    },
    import : () => {
        document.getElementById('btn-import').click();
    },
    save : () => {
        document.getElementById('btn-save').click();
    },
    saveAs : () => {
        document.getElementById('btn-save-as').click();
    },
    export : () => {
        document.getElementById('btn-export').click();
    },
    lookUp : () => {
        document.getElementById('btn-lookup').click();
    },
    replace : () => {
        document.getElementById('btn-replace').click();
    },
    exit : () => {
        handleClose();
    },
    undo : () => {
        E.chain().undo().focus().run();
    },
    redo : () => {
        E.chain().redo().focus().run();
    },
    copy : () => {
        document.getElementById('btn-copy').click();
    },
    paste : () => {
        document.getElementById('btn-paste').click();
    },
    cut : () => {
        document.getElementById('btn-cut').click();
    },
    selectAll : () => {
        E.chain().focus().selectAll().run();
    },
    about : () => {
        document.getElementById('btn-about').click();
    }
}

onMounted(() => {
    runtime.WindowIsMaximised().then(res=> {
        maximised.value = res;
    })
})

const handleMaximise = () => {
    runtime.WindowMaximise();
    maximised.value = true;
}

const handleUnMaximise = () => {
    runtime.WindowUnmaximise();
    maximised.value = false;
}

const handleClose = () => {
    function innerCallBack(action){
        if(action == 'confirm'){
            let name = bookInfo.metadata.title;
            FileSave(name, JSON.stringify(bookInfo.value)).then((res)=>{
                if(res.Code == 0) {
                    ElMessage.info(t('message.saveSuccess'))
                    runtime.Quit();
                }else if(res.Code == -1){
                    return;
                }else{
                    ElMessage.error(t('message.saveError') + ": " + res);
                }
            })
        }else if(action == 'cancel'){
            runtime.Quit();
        }else if(action == 'close'){
            return;
        }
    }
    if(change.value){
        ElMessageBox.confirm(t('message.saveWarning'),t('message.warning'), {
            confirmButtonText: t('message.confirm'),
            cancelButtonText: t('message.notSave'),
            distinguishCancelAndClose: true,
            type: 'warning',
            callback: innerCallBack
        })
        return;
    }
    runtime.Quit();
}

const handleSelect = (index) => {
    // console.log(index);
    menu[index]();
}
</script>

<template>
<div class="header_container" style="--wails-draggable:drag">
    <div class="header-left">
        <div class="header-icon">
            <img src="./assets/images/appicon.png" />
        </div>
        <el-menu
            :default-active="activeIndex"
            class="header-menu"
            menu-trigger="click"
            mode="horizontal"
            :ellipsis="false"
            id="header_menu"
            :close-on-click-outside="true"
            :unique-opened="true"
            @select="handleSelect"
        >   <!-- File -->
            <el-sub-menu index="file">
                <template #title>{{t("header.file")}}</template>
                <el-menu-item index="new"><span class="sub-menu-item">{{t("header.new")}}<span class="sub-menu-shortcut">Ctrl+N</span></span></el-menu-item>
                <el-menu-item index="open"><span class="sub-menu-item">{{t("header.open")}}<span class="sub-menu-shortcut">Ctrl+O</span></span></el-menu-item>
                <el-menu-item index="import"><span class="sub-menu-item seperator">{{t("header.import")}}<span class="sub-menu-shortcut">Ctrl+Shift+O</span></span></el-menu-item>
                <el-menu-item index="save"><span class="sub-menu-item">{{t("header.save")}}<span class="sub-menu-shortcut">Ctrl+S</span></span></el-menu-item>
                <el-menu-item index="saveAs"><span class="sub-menu-item">{{t("header.saveAs")}}<span class="sub-menu-shortcut">Ctrl+Shift+S</span></span></el-menu-item>
                <el-menu-item index="export"><span class="sub-menu-item seperator">{{t("header.export")}}<span class="sub-menu-shortcut">Ctrl+E</span></span></el-menu-item>
                <el-menu-item index="lookUp"><span class="sub-menu-item">{{t("header.lookUp")}}<span class="sub-menu-shortcut">Ctrl+F</span></span></el-menu-item>
                <el-menu-item index="replace"><span class="sub-menu-item seperator">{{t("header.replace")}}<span class="sub-menu-shortcut">Ctrl+H</span></span></el-menu-item>
                <el-menu-item index="exit"><span class="sub-menu-item">{{t("header.exit")}}<span class="sub-menu-shortcut"></span></span></el-menu-item>
            </el-sub-menu>
            <!-- Edit -->
            <el-sub-menu index="edit">
                <template #title>{{t("header.edit")}}</template>
                <el-menu-item index="undo"><span class="sub-menu-item">{{t("header.undo")}}<span class="sub-menu-shortcut">Ctrl+Z</span></span></el-menu-item>
                <el-menu-item index="redo"><span class="sub-menu-item seperator">{{t("header.redo")}}<span class="sub-menu-shortcut">Ctrl+Shift+Z</span></span></el-menu-item>
                <el-menu-item index="copy"><span class="sub-menu-item">{{t("header.copy")}}<span class="sub-menu-shortcut">Ctrl+C</span></span></el-menu-item>
                <el-menu-item index="paste"><span class="sub-menu-item">{{t("header.paste")}}<span class="sub-menu-shortcut">Ctrl+V</span></span></el-menu-item>
                <el-menu-item index="cut"><span class="sub-menu-item seperator">{{t("header.cut")}}<span class="sub-menu-shortcut">Ctrl+X</span></span></el-menu-item>
                <el-menu-item index="selectAll"><span class="sub-menu-item">{{t("header.selectAll")}}<span class="sub-menu-shortcut">Ctrl+A</span></span></el-menu-item>
            </el-sub-menu>
            <!-- Help -->
            <el-sub-menu index="help">
                <template #title>{{t("header.help")}}</template>
                <el-menu-item index="about"><span class="sub-menu-item">{{t("header.about")}}<span class="sub-menu-shortcut"></span></span></el-menu-item>
            </el-sub-menu>
        </el-menu>
    </div>
    <div class="header-center">
        <span class="save-flag" v-show="change">*</span>
        <div class="header-title">
            {{title}}
        </div>
    </div>
    <div class="header-right">
        <button class="button-item" @click="runtime.WindowMinimise()">
            <span>
                <i class="el-icon">
                    <svg t="1729826233754" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4437" width="200" height="200"><path d="M153.6 470.3744h716.8v83.2512H153.6V470.3744z" fill-opacity=".9" p-id="4438"></path></svg>
                </i>
            </span>
        </button>
        <button v-if="maximised" class="button-item" @click="handleUnMaximise">
            <span>
                <i class="el-icon">
                    <svg t="1729826641136" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6677" width="200" height="200"><path d="M861.3 741.2v-70c17.5 0 31.7-14 31.7-31.2V164.4c0-17.2-14.2-31.2-31.7-31.2H386.7c-17.5 0-31.7 14-31.7 31.2h-70c0-55.8 45.6-101.2 101.7-101.2h474.6c56.1 0 101.7 45.4 101.7 101.2V640c0 55.8-45.6 101.2-101.7 101.2z" p-id="6678"></path><path d="M861.3 741.2h-93.2v-70h93.2c17.5 0 31.7-14 31.7-31.2V164.4c0-17.2-14.2-31.2-31.7-31.2H386.7c-17.5 0-31.7 14-31.7 31.2V256h-70v-91.5c0-55.8 45.6-101.2 101.7-101.2h474.6c56.1 0 101.7 45.4 101.7 101.2V640c0 55.8-45.6 101.2-101.7 101.2z" p-id="6679"></path><path d="M639.1 963.3H164.5c-56.1 0-101.7-45.4-101.7-101.2V386.5c0-55.8 45.6-101.2 101.7-101.2h474.6c56.1 0 101.7 45.4 101.7 101.2v475.6c0 55.8-45.6 101.2-101.7 101.2z m-474.6-608c-17.5 0-31.7 14-31.7 31.2v475.6c0 17.2 14.2 31.2 31.7 31.2h474.6c17.5 0 31.7-14 31.7-31.2V386.5c0-17.2-14.2-31.2-31.7-31.2H164.5z" p-id="6680"></path></svg>
                </i>
            </span>
        </button>
        <button v-else class="button-item" @click="handleMaximise">
            <span>
                <i class="el-icon">
                    <svg t="1729827242226" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="13190" width="200" height="200"><path d="M806.4 153.6H217.6C179.2 153.6 153.6 179.2 153.6 217.6v588.8C153.6 844.8 179.2 870.4 217.6 870.4h588.8c38.4 0 64-25.6 64-64V217.6C870.4 179.2 844.8 153.6 806.4 153.6z m0 652.8H217.6V217.6h588.8v588.8z" fill-opacity=".9" p-id="13191"></path></svg>
                </i>
            </span>
        </button>
        <button class="button-item" @click="handleClose">
            <span>
                <i class="el-icon">
                    <svg t="1729826719524" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="7976" width="200" height="200"><path d="M925.468404 822.294069 622.19831 512.00614l303.311027-310.331931c34.682917-27.842115 38.299281-75.80243 8.121981-107.216907-30.135344-31.369452-82.733283-34.259268-117.408013-6.463202L512.000512 399.25724 207.776695 87.993077c-34.675754-27.796066-87.272669-24.90625-117.408013 6.463202-30.178323 31.414477-26.560936 79.375815 8.121981 107.216907l303.311027 310.331931L98.531596 822.294069c-34.724873 27.820626-38.341237 75.846432-8.117888 107.195418 30.135344 31.43699 82.72919 34.326806 117.408013 6.485715l304.178791-311.219137 304.177767 311.219137c34.678824 27.841092 87.271646 24.951275 117.408013-6.485715C963.808618 898.140501 960.146205 850.113671 925.468404 822.294069z" p-id="7977"></path></svg>
                </i>
            </span>
        </button>
    </div>
</div>
</template>

<style scoped>
.header_container {
    height: 20px;
    width: 100%;
    background-color: var(--el-bg-color);
    padding: 5px 5px;
    display: flex;
    justify-content: start;
    align-items: center;
    
}
.header-left{
    width: 28%;
    height: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.header-icon{
    width: 22px;
    height: 20px;
}
.header-icon img{
    height: 100%;
    width: 100%;
}
.header-menu {
    height: 100%;
    display: flex;
    width: 95%;
    border: none !important;
}
.el-dropdown-link {
    display: flex;
    align-items: center;
    cursor: pointer;
    padding: 8px 5px;
    font-size: 0.9em;
}
.el-dropdown-link:hover {
    background-color: var(--el-fill-color-dark);
    border-radius: 2px;
}
.sub-menu-item {
    width: 100%;
    display: flex;
    justify-content: space-between;
}
.sub-menu-shortcut{
    margin-right: 10px;
}
.seperator {
    width: 100%;
    border-bottom: 1px solid var(--el-border-color);
}
.header-center {
    width: 62%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
}
.save-flag {
    color: var(--el-text-color-secondary);
    font-size: 1em;
}
.header-title {
    width: 300px;
    height: 100%;
    font-size: 0.8em;
    color: var(--el-text-color-secondary);
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
}
.header-right {
    width: 10%;
    height: 100%;
    display: flex;
    justify-content: flex-end;
    align-items: center;
}
.button-item {
    width: 28px;
    height: 25px;
    border: none;
    background-color: var(--el-bg-color);
    padding: 4px 5px;
    cursor: pointer;
    margin-right: 15px;
}

.button-item:last-child {
    margin-right: 10px;
}

.button-item:hover {
    background-color: var(--el-fill-color-dark);
    border-radius: 2px;
}

/* Cancel the hover style of the menu item */
:deep(li.el-sub-menu:hover div.el-sub-menu__title) {
    color: var(--el-menu-text-color) !important;
}

.el-menu-item:hover{
    color: var(--el-menu-text-color) !important;
}

/* Cancel the acitve style of the sub menu item */
:deep(li.is-active>div.el-sub-menu__title) {
    border-bottom: 2px solid transparent !important;
    color: var(--el-menu-text-color) !important;
}

ul.el-menu li.el-menu-item.is-active {
    color: var(--el-menu-text-color) !important;
}

</style>