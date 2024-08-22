<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-08-16
@Description: This is the main frame of the application.
*/

import { ref, reactive, watch, onMounted, provide, h} from 'vue';
import { ElMessage, ElMessageBox, ElNotification  } from 'element-plus'
import { useI18n } from 'vue-i18n';
import editor from './editor.vue';
import { FileOpen, FileSave, Publish, GetImageData } from '../../wailsjs/go/main/App.js'
import dialogs from './dialog.vue'
import * as utils from '../assets/js/utils'

const editorRef = utils.editorRef;
const bookinfo = utils.bookInfo;
const { t } = useI18n();
const resetState = utils.resetState;


provide('bookinfo', bookinfo);

const visio = utils.visio;

const openFilePicker = async () => {
    var rawData = await FileOpen().then((res) => {
        var data = JSON.parse(res);
        return data
    })
    rawData.metadata.creator = rawData.metadata.creator.join(',');
    rawData.metadata.contributors = rawData.metadata.contributors.join(',');
    const editor = editorRef.value;
    editor.setHtml(rawData.content);
    bookinfo.metadata = rawData.metadata;
    bookinfo.content = rawData.content;
    bookinfo.toc = rawData.toc;
    bookinfo.resources = rawData.resources;
    utils.initCover();
}

const saveFilePicker = () => {
    let tempData = JSON.parse(JSON.stringify(bookinfo));
    tempData.metadata.creator = tempData.metadata.creator.split(',');
    tempData.metadata.contributors = tempData.metadata.contributors.split(',');
    tempData.content = tempData.content.replaceAll('"', "'");
    let name = bookinfo.metadata.title;
    FileSave(name, JSON.stringify(tempData)).then((res)=>{
        if(res.Code === 0) {
            ElMessage.info(t('message.saveSuccess'));
            utils.change.value = false;
        }else if(res.Code === -1){
            return;
        }else{
            ElMessage.error(t('message.saveError') + ": " + res);
        }
    })
}

const exportFile = () => {
    let tempData = JSON.parse(JSON.stringify(bookinfo));
    const editor = editorRef.value;
    const allImages = editor.getElemsByType("image");
    const imageIDs = Array.from(new Set(allImages.map(e => e.alt)))
    imageIDs.map(async (v)=>{
        let data = await GetImageData(v).then((res)=> {
            return res;
        })
        if(data.Code == 1){
            ElMessage.error(t('message.exportError') + ": " + data.Msg)
            return
        }
        tempData.resources.push({
            id: v,
            name: v,
            data: data.Data,
            type: "image/jpeg"
        })
        
    })
    const headers = editor.getElemsByTypePrefix('header');
    const Toc =  new utils.TocGenerator(headers)
    const res = Toc.process()
    tempData.toc = res;
    tempData.metadata.creator = tempData.metadata.creator.split(',');
    tempData.metadata.contributors = tempData.metadata.contributors.split(',');
    const doc = new DOMParser().parseFromString(tempData.content, 'text/html')
    let titles = doc.querySelectorAll("h1, h2, h3, h4, h5")
    titles.forEach((e, i) => {
        e.id = "guide_signal_" + (i+1);
    })
    tempData.content = doc.body.innerHTML;
    const regex = /<img(.*?)>/g;
    tempData.content = tempData.content.replaceAll("<br>", "<br></br>");
    tempData.content = tempData.content.replace(regex, (match, p1)=> {
        return `<img${p1}/>`;
    })
    tempData.content = tempData.content.replaceAll('"', "'");
    let name = bookinfo.metadata.title;
    Publish(name, (()=>{
            ElNotification({
                title: t('message.exportInfoTitle'),
                message: h('p', { style: 'color: teal' }, t('message.exportInfo')),
                type: 'warning',
            })
            return JSON.stringify(tempData)
        }  
        )()
    ).then((res)=> {
        if(res.Code == 1){
            ElMessage.error(t('message.exportError') + ": " + res.Msg)
            return
        }else if(res.Code == 0){
            ElMessage.success(t('message.exportSuccess'))
        }
        return
    })
}

const activeNames = ref("1")
const disableBtn = reactive({
    file: false,
    edit: true,
})

const btnDisabledClass = ref("el-button btn func_btn-big is-disabled");
const btnNormalClass = ref('el-button btn func_btn-big');

const judgeHeight = () => {
    var tabHeight = document.querySelector('.tab-container').offsetHeight;
    var totalHeight = window.innerHeight;
    var elem = document.querySelector('#editor-container');
    elem.setAttribute('style', 'height:' + (totalHeight - tabHeight - 40) + 'px');
}


onMounted(() => {
    setTimeout(()=>{
        judgeHeight();
    }, 500)
})

const handleChange = () => {
    setTimeout(()=>{
        judgeHeight();
    }, 300)
}

const openFileInfo = ()=> {
    visio.bookInfoVisible = true;
}

const handleNew = ()=> {
    function innerCallBack(action){
        console.log(action)
        if(action == 'confirm'){
            let name = bookinfo.metadata.title;
            FileSave(name, JSON.stringify(bookinfo.value)).then((res)=>{
                if(res.Code == 0) {
                    ElMessage.info(t('message.saveSuccess'))
                    resetState();
                }else if(res.Code == -1){
                    return;
                }else{
                    ElMessage.error(t('message.saveError') + ": " + res);
                }
            })
        }else if(action == 'cancel'){
            resetState();
        }else if(action == 'close'){
            return;
        }
    }
    if(utils.change.value){
        ElMessageBox.confirm(t('message.saveWarning'),t('message.warning'), {
            confirmButtonText: t('message.confirm'),
            cancelButtonText: t('message.notSave'),
            distinguishCancelAndClose: true,
            type: 'warning',
            callback: innerCallBack
        })
        return;
    }
}

const test = () => {
    // const editor = editorRef.value;
    // const allImage = editor.getElemsByType("image")
    // const id = allImage.map(e => e.alt)
    console.log(utils.change.value)
}

</script>

<template>
<button id="testing" @click="test" style="display: none;"></button>
<div class="tab-container">
    <dialogs></dialogs>
    <el-tabs type="border-card">
        <el-collapse v-model="activeNames" @change="handleChange">
            <el-collapse-item name="1">
                <el-tab-pane :label="$t('nav.file')">
                    <div class="func_btn-list">
                        <div class="func_btn-container list_row">
                            <span style="width:35px; " class="btn-group">
                                <button class="el-button btn func_btn-big" @click="handleNew" id="btn-new">
                                    <i class="el-icon">
                                        <svg t="1717397117328" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4260" width="200" height="200"><path d="M524.8 64c4.693333 0 8.533333 3.84 8.533333 8.533333V490.666667h418.133334c4.693333 0 8.533333 3.84 8.533333 8.533333v46.933333a8.533333 8.533333 0 0 1-8.533333 8.533334H533.333333v418.133333a8.533333 8.533333 0 0 1-8.533333 8.533333h-46.933333a8.533333 8.533333 0 0 1-8.533334-8.533333V554.666667H51.2a8.533333 8.533333 0 0 1-8.533333-8.533334v-46.933333c0-4.693333 3.84-8.533333 8.533333-8.533333H469.333333V72.533333c0-4.693333 3.84-8.533333 8.533334-8.533333h46.933333z" p-id="4261"></path></svg>
                                    </i>
                                </button>
                                <span>{{$t('toolBar.file.new')}}</span>
                            </span>
                            <span style=" width:70px; " class="btn-group">
                                <button class="el-button btn func_btn-big" @click="openFilePicker" id="btn-open">
                                    <i class="el-icon">
                                        <svg t="1717816393836" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="8897" width="200" height="200"><path d="M921.6 450.133333c-6.4-8.533333-14.933333-12.8-25.6-12.8h-10.666667V341.333333c0-40.533333-34.133333-74.666667-74.666666-74.666666H514.133333c-4.266667 0-6.4-2.133333-8.533333-4.266667l-38.4-66.133333c-12.8-21.333333-38.4-36.266667-64-36.266667H170.666667c-40.533333 0-74.666667 34.133333-74.666667 74.666667v597.333333c0 6.4 2.133333 12.8 6.4 19.2 6.4 8.533333 14.933333 12.8 25.6 12.8h640c12.8 0 25.6-8.533333 29.866667-21.333333l128-362.666667c4.266667-10.666667 2.133333-21.333333-4.266667-29.866667zM170.666667 224h232.533333c4.266667 0 6.4 2.133333 8.533333 4.266667l38.4 66.133333c12.8 21.333333 38.4 36.266667 64 36.266667H810.666667c6.4 0 10.666667 4.266667 10.666666 10.666666v96H256c-12.8 0-25.6 8.533333-29.866667 21.333334l-66.133333 185.6V234.666667c0-6.4 4.266667-10.666667 10.666667-10.666667z m573.866666 576H172.8l104.533333-298.666667h571.733334l-104.533334 298.666667z" p-id="8898"></path></svg>
                                    </i>
                                </button>
                                <span>{{$t('toolBar.file.open')}}</span>
                            </span>
                            <div class="division-border"></div>
                            <span class="btn-group" style="width:50px;">
                                <button @click="saveFilePicker" :class="disableBtn.file?btnDisabledClass:btnNormalClass" :aria-disabled="disableBtn.file" :disabled="disableBtn.file" id="btn-save">
                                    <i class="el-icon">
                                        <svg t="1717405451499" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="21449" width="200" height="200"><path d="M708.388571 121.904762L902.095238 320.804571V828.952381a73.142857 73.142857 0 0 1-73.142857 73.142857H195.047619a73.142857 73.142857 0 0 1-73.142857-73.142857V195.047619a73.142857 73.142857 0 0 1 73.142857-73.142857h513.340952zM292.571429 195.023238L195.047619 195.047619v633.904762l97.52381-0.024381V536.380952h438.857142v292.547048L828.952381 828.952381V350.549333l-97.52381-100.156952V365.714286H292.571429V195.023238zM658.285714 609.52381H365.714286v219.40419h292.571428V609.52381z m-48.761904 73.142857v73.142857h-195.04762v-73.142857h195.04762z m48.761904-487.619048H365.714286v97.52381h292.571428V195.047619z" p-id="21450"></path></svg>
                                    </i>
                                </button>
                                <span>{{$t('toolBar.file.save')}}</span>
                            </span>
                            <span class="btn-group" style="width:50px;">
                                <button @click="exportFile" :class="disableBtn.file?btnDisabledClass:btnNormalClass" :aria-disabled="disableBtn.file" :disabled="disableBtn.file" id="btn-export">
                                    <i class="el-icon">
                                        <svg t="1717398848121" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4277" width="200" height="200"><path d="M895.99896 1024H127.99992a63.99992 63.99992 0 0 1-63.99992-63.99992V384.0008a63.99992 63.99992 0 0 1 63.99992-63.99992h223.99972a31.99996 31.99996 0 0 1 0 63.99992H127.99992v575.99928h767.99904V384.0008h-223.99972a31.99996 31.99996 0 0 1 0-63.99992h223.99972a63.99992 63.99992 0 0 1 63.99992 63.99992v575.99928a63.99992 63.99992 0 0 1-63.99992 63.99992zM543.9994 100.481154V704.0004a31.99996 31.99996 0 0 1-63.99992 0V100.481154L340.479654 216.641009 299.519706 167.361071l191.99976-159.9998 0.255999 0.319999a30.463962 30.463962 0 0 1 40.44795 0l0.255999-0.319999 191.99976 159.9998-40.959948 49.279938z" p-id="4278"></path></svg>
                                    </i>
                                </button>
                                <span>{{$t('toolBar.file.export')}}</span>
                            </span>
                            <div class="division-border"></div>
                            <span class="btn-group" style="width:80px;">
                                <button @click="openFileInfo" :class="disableBtn.file?btnDisabledClass:btnNormalClass" :aria-disabled="disableBtn.file" :disabled="disableBtn.file">
                                    <i class="el-icon">
                                        <svg t="1717677749278" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="7239" width="200" height="200"><path d="M688 312v-48c0-4.4-3.6-8-8-8H296c-4.4 0-8 3.6-8 8v48c0 4.4 3.6 8 8 8h384c4.4 0 8-3.6 8-8zM296 400c-4.4 0-8 3.6-8 8v48c0 4.4 3.6 8 8 8h184c4.4 0 8-3.6 8-8v-48c0-4.4-3.6-8-8-8H296zM672 516c-119.3 0-216 96.7-216 216s96.7 216 216 216 216-96.7 216-216-96.7-216-216-216z m107.5 323.5C750.8 868.2 712.6 884 672 884s-78.8-15.8-107.5-44.5C535.8 810.8 520 772.6 520 732s15.8-78.8 44.5-107.5C593.2 595.8 631.4 580 672 580s78.8 15.8 107.5 44.5C808.2 653.2 824 691.4 824 732s-15.8 78.8-44.5 107.5z" p-id="7240"></path><path d="M672 812m-32 0a32 32 0 1 0 64 0 32 32 0 1 0-64 0Z" p-id="7241"></path><path d="M652 748h40c4.4 0 8-3.6 8-8V628c0-4.4-3.6-8-8-8h-40c-4.4 0-8 3.6-8 8v112c0 4.4 3.6 8 8 8z" p-id="7242"></path><path d="M440 852H208V148h560v344c0 4.4 3.6 8 8 8h56c4.4 0 8-3.6 8-8V108c0-17.7-14.3-32-32-32H168c-17.7 0-32 14.3-32 32v784c0 17.7 14.3 32 32 32h272c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8z" p-id="7243"></path></svg>
                                    </i>
                                </button>
                                <span>{{$t('toolBar.file.fileInfo')}}</span>
                            </span>
                            <div class="division-border"></div>
                        </div>
                    </div>
                </el-tab-pane>
                <el-tab-pane :label="$t('nav.edit')">
                    <div class="func_btn-list">
                        <div class="func_btn-container list_row">
                            <span class="btn-group" style="width:40px;">
                                <button class="el-button btn func_btn-big">
                                    <i class="el-icon">
                                        <svg t="1717403088327" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="12105" width="200" height="200"><path d="M704 128l-128 0 0-64c0-35.2-28.8-64-64-64l-128 0c-35.2 0-64 28.8-64 64l0 64-128 0 0 128 512 0 0-128zM512 128l-128 0 0-63.872c0.032-0.032 0.064-0.064 0.128-0.128l127.776 0c0.032 0.032 0.096 0.064 0.128 0.128l0 63.872zM832 320l0-160c0-17.6-14.4-32-32-32l-64 0 0 64 32 0 0 128-192 0-192 192 0 256-256 0 0-576 32 0 0-64-64 0c-17.6 0-32 14.4-32 32l0 640c0 17.6 14.4 32 32 32l288 0 0 192 640 0 0-704-192 0zM576 410.496l0 101.504-101.504 0 101.504-101.504zM960 960l-512 0 0-384 192 0 0-192 320 0 0 576z" p-id="12106"></path></svg>
                                    </i>
                                </button>
                                {{$t('toolBar.edit.paste')}}
                            </span>
                            <span class="btn-group">
                                <span style="margin-bottom: 2px;" >
                                    <button class="el-button btn is-link">
                                        <i class="el-icon">
                                            <svg t="1717403312181" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="13146" width="200" height="200"><path d="M629.39164 415.032528l-163.616428-163.616428c-7.992021-7.992021-20.947078-7.992021-28.939099 0-7.992021 8.002254-7.992021 20.957311 0 28.949332l128.680754 128.680754-175.548178 0L389.968689 184.082552c0-11.2973-9.168824-20.466124-20.466124-20.466124L21.813818 163.616428c-11.307533 0-20.466124 9.168824-20.466124 20.466124l0 818.08214c0 11.307533 9.15859 20.466124 20.466124 20.466124l593.108273 0c11.307533 0 20.466124-9.15859 20.466124-20.466124L635.388215 429.512311C635.388215 424.078555 633.229039 418.880159 629.39164 415.032528zM594.455967 981.698568l-552.176025 0L42.279942 204.548676l306.756499 0 0 224.963635c0 11.2973 9.15859 20.466124 20.466124 20.466124l224.953402 0L594.455967 981.698568z" p-id="13147"></path><path d="M1023.978511 265.895883l0 572.652382c0 11.307533-9.15859 20.466124-20.466124 20.466124l-307.86167 0c-11.2973 0-20.466124-9.15859-20.466124-20.466124 0-11.2973 9.168824-20.466124 20.466124-20.466124l287.395546 0L983.046263 286.362007l-224.953402 0c-11.307533 0-20.466124-9.168824-20.466124-20.466124L737.626737 40.932248l-306.756499 0 0 75.693959c0 11.307533-9.168824 20.466124-20.466124 20.466124-11.307533 0-20.466124-9.15859-20.466124-20.466124L389.93799 20.466124c0-11.2973 9.15859-20.466124 20.466124-20.466124l347.688747 0c11.2973 0 20.466124 9.168824 20.466124 20.466124l0 224.963635 175.548178 0-128.680754-128.680754c-7.992021-7.992021-7.992021-20.947078 0-28.949332 7.992021-7.992021 20.947078-7.992021 28.939099 0l163.616428 163.626661C1021.819334 255.263731 1023.978511 260.462127 1023.978511 265.895883z" p-id="13148"></path></svg>
                                        </i>
                                    </button>
                                </span>
                                <span>
                                    <button class="el-button btn is-link">
                                        <i class="el-icon">
                                            <svg t="1717403432425" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="14228" width="200" height="200"><path d="M790.4 640c-64 0-120 30.4-153.6 81.6L560 507.2l166.4-448c8-20.8-4.8-46.4-25.6-56-20.8-8-46.4 4.8-56 25.6L512 382.4 384 28.8C376 8 350.4-4.8 328 3.2c-20.8 8-33.6 33.6-25.6 56l161.6 448L388.8 720c-33.6-51.2-89.6-81.6-153.6-81.6-107.2 0-192 84.8-192 192s84.8 192 192 192c94.4 0 171.2-64 188.8-153.6V864l89.6-240 84.8 240c0 4.8 4.8 4.8 4.8 8C622.4 960 696 1024 790.4 1024c107.2 0 192-84.8 192-192s-86.4-192-192-192zM233.6 939.2c-59.2 0-105.6-48-105.6-107.2s46.4-107.2 107.2-107.2 107.2 46.4 107.2 107.2c-1.6 59.2-48 107.2-108.8 107.2z m556.8 0c-59.2 0-107.2-46.4-107.2-107.2 0-59.2 46.4-107.2 107.2-107.2 59.2 0 107.2 46.4 107.2 107.2-1.6 59.2-48 107.2-107.2 107.2z" p-id="14229"></path></svg>
                                        </i>
                                    </button>
                                </span>
                            </span>
                            <div class="division-border"></div>
                            <span class="btn-group">
                                <button :class="btnNormalClass" @click="visio.mediaVisible = true">
                                    <i class="el-icon">
                                        <svg t="1717816046552" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6855" width="200" height="200"><path d="M298.666667 277.333333L469.333333 384l-170.666666 106.666667v-213.333334z" p-id="6856"></path><path d="M213.333333 128a85.333333 85.333333 0 0 0-85.333333 85.333333v597.333334a85.333333 85.333333 0 0 0 85.333333 85.333333h597.333334a85.333333 85.333333 0 0 0 85.333333-85.333333V213.333333a85.333333 85.333333 0 0 0-85.333333-85.333333H213.333333z m597.333334 64H213.333333a21.333333 21.333333 0 0 0-21.333333 21.333333v456.832l154.88-112.896a32 32 0 0 1 36.266667-1.024l124.842666 80.938667 152.362667-148.138667a32 32 0 0 1 44.629333 0l127.018667 123.52V213.333333a21.333333 21.333333 0 0 0-21.333333-21.333333zM192 810.666667v-61.269334l174.762667-127.445333 127.829333 82.816a32 32 0 0 0 39.68-3.882667L682.666667 556.629333l148.352 144.213334 0.981333-0.981334V810.666667a21.333333 21.333333 0 0 1-21.333333 21.333333H213.333333a21.333333 21.333333 0 0 1-21.333333-21.333333z" p-id="6857"></path></svg>
                                    </i>
                                </button>
                                {{$t('toolBar.edit.media')}}
                            </span>
                            <div class="division-border"></div>
                            <span class="btn-group" style="width:40px; ">
                                <button class="el-button btn func_btn-big" @click="visio.settingVisible = true">
                                    <i class="el-icon">
                                        <svg t="1717404325809" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="19304" width="200" height="200"><path d="M885.196958 871.528145c-14.753141-5.464126-25.134981-8.742602-34.970408-12.56749-30.052695-11.474665-59.012564-24.588568-89.611671-34.970408-8.742602-2.732063-21.310093-1.639238-29.506282 2.732063-21.310093 10.38184-42.073773 22.402918-61.744627 36.063234-6.556952 4.371301-12.021078 14.206728-13.660316 22.402918-6.556952 36.063234-11.474665 72.67288-15.845966 109.282526-2.732063 20.76368-12.56749 30.052695-33.877583 29.506282-62.837452-0.546413-125.674905-0.546413-188.512357 0-20.76368 0-31.691933-8.196189-34.423996-28.959869-4.917714-36.609646-10.928253-72.67288-15.299554-109.282526-1.639238-12.56749-6.010539-20.217267-18.031617-26.227806-20.217267-9.835427-38.248884-23.495743-58.466151-33.33117-7.103364-3.824888-18.031617-5.464126-25.681394-2.732063-35.516821 12.56749-69.940817 26.774219-104.911225 40.434535-25.134981 9.835427-32.238345 7.649777-46.445074-16.392379-31.14552-53.002025-61.744627-106.00405-92.343735-159.006075-14.206728-24.588568-13.660316-30.599107 8.742602-48.084311 29.506282-23.495743 59.558977-45.898661 88.518846-69.940817 5.464126-4.371301 9.289015-14.206728 9.289015-21.310093 1.092825-20.217267-3.278476-41.52736 0-61.198215 3.824888-21.310093-4.917714-32.238345-20.217267-43.71301-27.320632-19.670855-52.455613-41.52736-79.229831-61.198215-15.845966-12.021078-19.670855-24.588568-9.289015-42.073773 33.33117-56.280501 66.115928-113.107414 98.900686-169.934328 9.289015-15.845966 21.310093-19.124442 38.248884-12.56749 36.063234 14.206728 72.126467 28.959869 108.736113 42.073773 7.649777 2.732063 19.670855 1.092825 27.320632-2.732063 20.76368-10.38184 39.341709-24.042156 60.105389-34.423996 10.928253-5.464126 14.753141-12.021078 15.845966-22.94933 4.917714-37.156059 10.928253-74.312118 15.845966-111.468177C386.322226 7.649777 397.796892 0 418.014159 0c62.837452 0.546413 125.674905 0.546413 188.512357 0 19.670855 0 30.599107 8.196189 32.784758 27.867044 4.917714 36.609646 10.928253 72.67288 15.299554 109.282526 1.639238 13.660316 6.556952 21.856505 19.670855 27.867044 19.124442 8.742602 37.156059 19.670855 54.641263 31.691933 11.474665 8.196189 21.310093 8.196189 33.877583 2.732063 33.877583-14.206728 68.301579-26.774219 102.725575-40.980947 19.124442-7.649777 32.238345-3.824888 42.620185 14.753141 31.14552 55.187676 63.930278 109.828939 95.62221 164.470202 12.56749 21.856505 10.928253 28.959869-9.289015 44.805836-28.959869 22.402918-57.919739 45.352248-87.426021 67.755166-9.835427 7.649777-12.56749 14.753141-11.474665 27.320632 2.185651 22.94933-0.546413 46.445074 0.546413 69.394404 0.546413 8.196189 4.371301 18.031617 10.38184 22.94933 28.959869 24.042156 59.012564 46.445074 88.518846 69.394404 20.217267 15.845966 21.310093 22.94933 8.742602 44.805836-32.238345 55.734088-64.47669 110.921764-97.261448 166.10944C901.042924 860.05348 891.753909 865.517606 885.196958 871.528145zM510.357893 692.304803c98.900686 0.546413 180.862581-79.776244 181.408993-177.584105 0.546413-100.539924-79.229831-181.955406-178.67693-182.501819-99.447099-0.546413-181.408993 80.322657-181.408993 179.769755C331.134551 610.342908 412.00362 691.75839 510.357893 692.304803z" p-id="19305"></path></svg>
                                    </i>
                                </button>
                                {{$t('toolBar.edit.setting')}}
                            </span>
                            <div class="division-border"></div>
                        </div>
                    </div>
                </el-tab-pane>
                <el-tab-pane :label="$t('nav.help')">
                    <div class="func_btn-list">
                        <div class="func_btn-container list_row">
                            <span class="btn-group" style="width:40px; ">
                                <button class="el-button btn func_btn-big">
                                    <i class="el-icon">
                                        <svg t="1717404400946" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="20342" width="200" height="200"><path d="M512 97.52381c228.912762 0 414.47619 185.563429 414.47619 414.47619s-185.563429 414.47619-414.47619 414.47619S97.52381 740.912762 97.52381 512 283.087238 97.52381 512 97.52381z m0 73.142857C323.486476 170.666667 170.666667 323.486476 170.666667 512s152.81981 341.333333 341.333333 341.333333 341.333333-152.81981 341.333333-341.333333S700.513524 170.666667 512 170.666667z m36.571429 268.190476v292.571428h-73.142858V438.857143h73.142858z m0-121.904762v73.142857h-73.142858v-73.142857h73.142858z" p-id="20343"></path></svg>
                                    </i>
                                </button>
                                {{$t('toolBar.help.help')}}
                            </span>
                        </div>
                    </div>
                </el-tab-pane>
            </el-collapse-item>
        </el-collapse>
    </el-tabs>
</div>
<div class="layout">
    <el-container>
        <div class="editor-container">
            <editor></editor>
        </div>
    </el-container>
</div>
</template>

<style>
.tab-container{
    width: 100%;
    margin: 0;
    padding: 0;
}

.layout {
    height: 100%;
}

.btn-group{
    display: flex; 
    flex-direction:column;
    height:100%;
    justify-content:center;
    align-items: center;
}

.el-tabs--border-card>.el-tabs__content {
    padding: 0 !important;
}
.el-collapse-item__header {
    height: 1em !important;
}

.el-collapse-item__content{
    padding-bottom: 5px !important;
}

.btn {
    border: none;
    background-color: transparent;
}

.btn:hover{
    cursor: pointer;
}

.icon{
    height: 100%;
    width: 100%;
    fill: var(--icon-color);
}

.func_btn-big{
    width: 32px;
    height: 32px;
    border: 1px solid #444444;
    border-radius: 5px;
}

.func_btn-list{
    display: flex;
    flex-direction: row;
    width: 100%;
    height: 55px;
    align-items: center;
}

.func_btn-container{
    margin-left: 3px;
}

.list_row{
    display: flex;
    flex-direction: row;
    justify-content: center;
}

.division-border{
    margin-left: 1.2em;
    margin-right: 1.2em;
    width: 1px;
    border-right: 1px solid #9b9b9b;
}

.list_col{
    display: flex;
    flex-direction: column;
    justify-content: center;
}

.el-main{
    padding: 0 !important;
    height: 100% !important;
}

.el-container {
    height: 100% !important;
}

.el-aside{
    overflow-x: hidden !important;
}

.novel_tree{
    overflow-y: auto;
    scrollbar-width: 2px;
    height: 100%;
    width: 100%;
    border-right: 1px solid #e6e6e6;
    background-color: var(--el-bg-color);
}

.func_toolbar{
    display: flex;
    flex-direction: row;
    justify-content: space-around;
    align-items: center;
    margin-top: 6px;
}

.func_toolbar a {
    cursor: pointer;
}

.editor-container {
    height: 100%;
    width: 100%;
}

.custom-tree-node {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 14px;
    padding-right: 8px;
}
.node-label{
    width: 90%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    background-color: transparent;
    border: none;
}
.novel_tree-title{
    display: flex;
    justify-content: center;
    margin-top: 5px;
}
.novel_tree-title input{
    border: none;
    background-color: transparent;
    font-size: 20px;
    font-weight: bold;
}

.novel_tree-title input:focus{
    outline: none;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.novel_tree-title input::placeholder{
    color: #9b9b9b;
}
</style>