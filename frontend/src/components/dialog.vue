<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-09-10
@Description: This is the dialog component, used element-plus
*/

import { reactive, ref, nextTick, inject } from 'vue';
import $ from 'jquery'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as utils from '../assets/js/utils'
import { useI18n } from 'vue-i18n';
import { languages as langs } from '../assets/js/i18n';
import { OpenImage, ImageUpload, GetStaticResources, FileDelete } from '../../wailsjs/go/main/App.js'
import { Delete, Plus, ZoomIn } from '@element-plus/icons-vue'
// import * as runtime from '../../wailsjs/runtime/runtime.js'

const { t, locale } = useI18n();
const bookinfo = inject("bookinfo");
const arrayRemove = utils.arrayRemove;
const editor = utils.editorRef;
const staticFiles = ref([]);
const cover = utils.cover;
const initCover = utils.initCover;
const dialogImageUrl = ref('');
const dialogVisible = ref(false);
const imageInfo = utils.imageInfo;
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
const languages = utils.language_list;
const visio = utils.visio;
const inputValue = ref('');
const inputVisible = ref(false);
const InputRef = ref();

const getImageFiles = async () => {
    const _ = await GetStaticResources().then((res)=>{
        var data = JSON.parse(res);
        // console.log(data)
        if(data.Code == 0) {
            staticFiles.value = data.FileList;
        }else {
            console.log(data.msg);
        }
        return true;
    })
}
getImageFiles();
initCover();

var lastClick = ref("");
const handleSelected = (event) => {
    var elem = $(event.target);
    if(lastClick.value != "") {
        lastClick.value.parent().css("border", "none");
        lastClick.value = elem;
        lastClick.value.parent().css("border", "1px solid #409EFF");
    }
    lastClick.value = elem;
    lastClick.value.parent().css("border", "1px solid #409EFF");
}

const uploadImage = async () => {
    var imageData = await ImageUpload().then((res)=>{
       return res;
    })
    if(imageData.Code == 0) {
        bookinfo.resources.push({
            id: imageData.Id + ".jpg",
            name: imageData.Id + ".jpg",
            type: "image/jpeg",
            data: ""
        });
        getImageFiles();
        ElMessage({
            message: t("dialog.media.uploadSuccess"),
            type: "success"
        });
    }else if(imageData.Code == 1) {
        ElMessage({
            message: t("dialog.media.uploadFail"),
            type: "error"
        });
    }
}

const uploadCover = async () => {
    var coverData = await OpenImage().then((res)=>{
       return res.Data;
    });
    bookinfo.metadata.meta.push({
        id: 'cover',
        content: "cover.jpg"
    });
    bookinfo.resources.push({
        id: "cover.jpg",
        name: "cover.jpg",
        type: "image/jpeg",
        data: coverData
    });
    initCover();
}

const handleRemove = () => {
    arrayRemove(bookinfo.metadata.meta, "cover");
    arrayRemove(bookinfo.resources, "cover.jpg");
    initCover();
}

const handlePictureCardPreview = (file) => {
  dialogImageUrl.value = "data:image/jpg;base64," + file;
  dialogVisible.value = true;
}

const handleClose = (tag)=>{
    bookinfo.metadata.subject.splice(bookinfo.metadata.subject.indexOf(tag), 1);
}
const showInput = ()=>{
    inputVisible.value = true;
    nextTick(() => {
        InputRef.value.input.focus();
    })
}

const handleInputConfirm = () => {
  if (inputValue.value) {
    bookinfo.metadata.subject.push(inputValue.value);
  }
  inputVisible.value = false;
  inputValue.value = '';
}

const handleSaveConfig = ()=> {
    localStorage.setItem('lang', setting.language);
    locale.value = setting.language;
    localStorage.setItem('theme', setting.theme);
    document.querySelector('html').setAttribute('class', setting.theme);
    var theme = setting.theme;
    var lang = setting.language;
    utils.editTheme.value = theme;
    utils.editLang.value = lang;
    visio.settingVisible = false;
}

const handleBookInfo = () => {
    console.log(bookinfo);
    visio.bookInfoVisible = false;
}

const handleInsert = () => {
    var elem = lastClick.value;
    const url = elem.attr("src");
    const id = elem.attr("id");
    const imgNode = {src: url, alt: id, title: id, zoom: 100, pos: "left"};
    const editor = utils.editorRef.value;
    editor.chain().focus().InsertImage(imgNode).createParagraphNear().run();
    imageInfo.elem = "";
    imageInfo.postition = "left";
    imageInfo.zoom = 100;
    visio.mediaVisible = false;
}

const handleDelImg = () => {
    var elem = lastClick.value;
    if(elem == ""){
        return;
    }
    ElMessageBox.confirm(t("dialog.media.deleteConfirm"), t("dialog.media.deleteTitle"), {
        confirmButtonText: t("dialog.confirm"),
        cancelButtonText: t("dialog.cancel"),
        type: 'warning',
    }).then(()=> {
        const name = elem.attr("id");     
        FileDelete(name).then((res) => {
            if(res.Code == 0){
                getImageFiles();
                ElMessage({
                    message: t("dialog.media.deleteSuccess"),
                    type: 'success'
                });
                visio.mediaVisible = false;
            }else{
                ElMessage({
                    message: t("dialog.media.deleteFail"),
                    type: 'error'
                });
            }
        })
    })
}

const handleImagePos = (justifyContent) => {
    imageInfo.postition = justifyContent;
    const E = utils.editorRef.value;
    console.log(justifyContent)
    E.chain().focus().updateAttributes("image", {pos: justifyContent}).run();
}

const handleImgSize = (val) => {
    imageInfo.zoom = val;
    const E = utils.editorRef.value;
    E.chain().focus().updateAttributes("image", {zoom: val}).run();
}

</script>

<template>
    <button style="display: none;" id="testing" @click="console.log(utils.metadata)"></button>

    <!--
    Setting dialog
    -->
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

    <!--
    Book info dialog
    -->
    <el-dialog
    v-model="visio.bookInfoVisible"
    :title="t('dialog.info.title')"
    width="500"
    >
        <el-form
        label-position="top"
        label-width="auto"
        :model="bookinfo.metadata"
        style="max-width: 600px"
        @input="()=>{utils.change.value = true}"
        >
            <el-form-item :label="t('dialog.info.name')">
                <el-input v-model="bookinfo.metadata.title" />
            </el-form-item>
            <el-form-item :label="t('dialog.info.identifier')">
                <el-input v-model="bookinfo.metadata.identifier" />
            </el-form-item>
            <el-form-item :label="t('dialog.info.publisher')">
                <el-input v-model="bookinfo.metadata.publisher" />
            </el-form-item>
            <el-form-item :label="t('dialog.info.author')">
                <el-input v-model="bookinfo.metadata.creator" />
            </el-form-item>
            <el-form-item :label="t('dialog.info.cover')">
                <div v-if="cover.isExist" class="image-preview_form">
                    <img :src="'data:image/jpg;base64,'+cover.data" alt="" />
                    <span class="list__item-actions">
                      <span
                        @click="handlePictureCardPreview(cover.data)"
                      >
                        <el-icon><zoom-in /></el-icon>
                      </span>
                      <span
                        @click="handleRemove()"
                      >
                        <el-icon><Delete /></el-icon>
                      </span>
                    </span>
                </div>
                <a class="form-cover-upload" @click="uploadCover" v-if="!cover.isExist">
                    <el-icon><Plus /></el-icon>
                </a>
                <el-dialog v-model="dialogVisible">
                    <img w-full :src="dialogImageUrl" alt="Preview Image" />
                </el-dialog>
            </el-form-item>
            <el-form-item :label="t('dialog.info.contributer')">
                <el-input v-model="bookinfo.metadata.contributors" />
            </el-form-item>
            <el-form-item :label="t('dialog.info.description')">
                <el-input 
                v-model="bookinfo.metadata.description" 
                type="textarea"
                :rows="3"/>
            </el-form-item>
            <el-form-item :label="t('dialog.info.subject')">
                <div class="flex gap-2">
                    <el-tag
                    v-for="tag in bookinfo.metadata.subject"
                    :key="tag"
                    closable
                    :disable-transitions="false"
                    @close="handleClose(tag)"
                    >
                    {{ tag }}
                    </el-tag>
                    <el-input
                    v-if="inputVisible"
                    ref="InputRef"
                    v-model="inputValue"
                    style="width: 50px;"
                    size="small"
                    @keyup.enter="handleInputConfirm"
                    @blur="handleInputConfirm"
                    />
                    <el-button v-else class="button-new-tag" size="small" @click="showInput">
                    + New Tag
                    </el-button>
                </div>
            </el-form-item>
            <el-form-item :label="t('dialog.info.language')">
                <el-select
                v-model="bookinfo.metadata.language"
                :placeholder="t('dialog.info.language')"
                size="large"
                style="width: 240px"
                >
                <el-option
                    v-for="item in languages"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                />
                </el-select>
            </el-form-item>
        </el-form>
        <template #footer>
            <div class="dialog-footer">
              <!-- <el-button @click="visio.bookInfoVisible = false">{{$t('dialog.cancel')}}</el-button> -->
              <el-button type="primary" @click="handleBookInfo">
                {{$t('dialog.confirm')}}
              </el-button>
            </div>
        </template>
    </el-dialog>

    <!--
    Media dialog
    -->
    <el-dialog
    v-model="visio.mediaVisible"
    :title="t('dialog.media.title')"
    width="60%"
    >
    <div class="media_container">
        <div class="media-toolbar">
            <el-button type="primary" plain @click="uploadImage">{{$t("dialog.media.upload")}}</el-button>
            <el-button type="danger" plain @click="handleDelImg">{{$t("dialog.media.delete")}}</el-button>
        </div>
        <div class="media-display">
                <el-image
                @click="handleSelected"
                v-for="item in staticFiles"
                :id="item"
                style="width: 160px; height: 160px;"
                :src="'http://127.0.0.1:7288/' + item"
                fit="contain"
                :z-index="-1"
                lazy
                ></el-image>
            </div>
        </div>
        <template #footer>
            <div class="dialog-footer">
            <el-button @click="visio.mediaVisible = false">{{$t('dialog.media.cancel')}}</el-button>
            <el-button type="primary" @click="handleInsert">
                {{$t('dialog.media.insert')}}
            </el-button>
            </div>
        </template>
    </el-dialog>

    <!-- help dialog -->
    <el-dialog
    v-model="visio.helpVisible"
    :title="t('dialog.help.title')"
    width="280px"
    >
    <div style="display: flex;justify-content: center;align-items: center;">
        <el-form style="width: 80%;">
            <el-form-item :label="t('dialog.help.version') + ':'">
                1.1.20
            </el-form-item>
            <el-form-item :label="t('dialog.help.date') + ':'">
                2024.9.25
            </el-form-item>
            <el-form-item :label="t('dialog.help.copyright') + ':'">
                2024 Hughie
            </el-form-item>
            <el-form-item label="Wails:">
                v2.8.1
            </el-form-item>
            <el-form-item label="Golang:">
                v1.20.3
            </el-form-item>
            <el-form-item label="Nodejs:">
                v16.15.1
            </el-form-item>
        </el-form>
    </div>
    </el-dialog>

    <!-- bubble menu -->
     <!-- image -->
    <div id="bubbleMenu" style="display: none;">
        <div class="bubble-menu">
            <div class="position-btn">
                <button :class="'bubble-btn '+ (imageInfo.postition == 'left' ? 'btn-active' : '')" @click="handleImagePos('left')">
                    <i class="el-icon">
                        <svg t="1725182706562" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="8984" width="200" height="200"><path d="M952 792H72c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8h880c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8zM952 160H72c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8h880c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8zM608 660c8.8 0 16-7.2 16-16V380c0-8.8-7.2-16-16-16H96c-8.8 0-16 7.2-16 16v264c0 8.8 7.2 16 16 16h512zM152 436h400v152H152V436zM704 646c0 4.4 3.6 8 8 8h224c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8H712c-4.4 0-8 3.6-8 8v56zM712 442h224c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8H712c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8z" p-id="8985"></path></svg>
                    </i>
                </button>
                <button :class="'bubble-btn '+ (imageInfo.postition == 'center' ? 'btn-active' : '')" @click="handleImagePos('center')">
                    <i class="el-icon">
                        <svg t="1725182692917" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="8826" width="200" height="200"><path d="M952 792H72c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8h880c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8zM952 160H72c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8h880c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8zM848 660c8.8 0 16-7.2 16-16V380c0-8.8-7.2-16-16-16H176c-8.8 0-16 7.2-16 16v264c0 8.8 7.2 16 16 16h672zM232 436h560v152H232V436z" p-id="8827"></path></svg>
                    </i>
                </button>
                <button :class="'bubble-btn '+ (imageInfo.postition == 'right' ? 'btn-active' : '')" @click="handleImagePos('right')">
                    <i class="el-icon">
                        <svg t="1725182718546" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="9142" width="200" height="200"><path d="M952 792H72c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8h880c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8zM952 160H72c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8h880c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8zM928 660c8.8 0 16-7.2 16-16V380c0-8.8-7.2-16-16-16H416c-8.8 0-16 7.2-16 16v264c0 8.8 7.2 16 16 16h512zM472 436h400v152H472V436zM80 646c0 4.4 3.6 8 8 8h224c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8H88c-4.4 0-8 3.6-8 8v56zM88 442h224c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8H88c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8z" p-id="9143"></path></svg>
                    </i>
                </button>
            </div>
            <div class="size-bar">
                <el-slider v-model="imageInfo.zoom" size="small" placement="bottom" @input="handleImgSize" :min="10"/>
            </div>
        </div>
    </div>

    <!-- table -->
    <div id="tableMenu" style="display: none;">
        <div class="bubble-menu" style="flex-direction: row; width: 390px">
            <el-tooltip :content="t('table.header')" placement="bottom" effect='dark'>
                <el-button text size="small" @click="editor.chain().focus().toggleHeaderRow().run()">
                    <i class="el-icon">
                        <svg t="1725548320998" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5117" width="200" height="200"><path d="M373.028571 365.714286v585.142857h-58.514285V365.714286zM709.485714 365.714286v585.142857h-58.514285V365.714286z" fill="#515151" p-id="5118"></path><path d="M950.857143 0H73.142857C29.257143 0 0 29.257143 0 73.142857v877.714286c0 43.885714 29.257143 73.142857 73.142857 73.142857h877.714286c43.885714 0 73.142857-29.257143 73.142857-73.142857V73.142857c0-43.885714-29.257143-73.142857-73.142857-73.142857z m0 950.857143H73.142857V438.857143h877.714286v512z" p-id="5119"></path></svg>
                    </i>
                </el-button>
            </el-tooltip>

            <el-tooltip :content="t('table.merge')" placement="bottom" effect='dark'>
                <el-button text size="small" @click="editor.chain().focus().mergeCells().run()">
                    <i class="el-icon" style="width: 2em">
                        <svg style="width: 2em; height: 2em" t="1725548409188" class="icon" viewBox="0 0 1365 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6488" width="200" height="200"><path d="M848.896 572.074667c0 6.656-2.730667 12.970667-7.68 17.92-4.949333 4.949333-11.264 7.68-17.92 7.68-6.656 0-12.970667-2.730667-17.92-7.68l-80.384-80.384c-3.413333-3.413333-3.413333-8.874667 0-12.117334l80.384-80.384c4.949333-4.949333 11.264-7.509333 17.92-7.509333 6.656 0 12.970667 2.730667 17.92 7.509333 4.949333 4.949333 7.68 11.264 7.68 17.92 0 6.656-2.730667 12.970667-7.68 17.92L816.128 477.866667h127.317333c14.165333 0 25.429333 11.434667 25.429334 25.429333 0 14.165333-11.434667 25.429333-25.429334 25.429333h-127.317333l25.258667 25.258667c4.778667 5.12 7.509333 11.434667 7.509333 18.090667zM516.437333 434.858667c0-6.656 2.730667-12.970667 7.68-17.92 4.949333-4.949333 11.264-7.509333 17.92-7.509334 6.656 0 12.970667 2.730667 17.92 7.509334l80.384 80.384c3.413333 3.413333 3.413333 8.874667 0 12.117333l-80.384 80.384c-4.949333 4.949333-11.264 7.68-17.92 7.68-6.656 0-12.970667-2.730667-17.92-7.509333-4.949333-4.949333-7.68-11.264-7.68-17.92s2.730667-12.970667 7.68-17.92l25.258667-25.258667h-127.317333c-14.165333 0-25.429333-11.434667-25.429334-25.429333 0-14.165333 11.434667-25.429333 25.429334-25.429334h127.317333l-25.258667-25.258666c-4.949333-4.949333-7.68-11.264-7.68-17.92z" p-id="6489"></path><path d="M792.234667 777.898667c0 5.461333 4.608 8.874667 8.874666 8.874666H989.866667c4.949333 0 8.874667-3.925333 8.874666-8.874666V246.101333c0-5.461333-4.608-8.874667-8.874666-8.874666H801.109333c-5.461333 0-8.874667 4.608-8.874666 8.874666v34.304c0 14.165333-11.434667 25.429333-25.429334 25.429334-14.165333 0-25.429333-11.434667-25.429333-25.429334V246.101333c0-33.109333 26.794667-59.904 59.904-59.904h188.757333c33.109333 0 59.904 26.794667 59.904 59.904v531.968c0 33.109333-26.794667 59.904-59.904 59.904H801.109333c-33.109333 0-59.904-26.794667-59.904-59.904v-34.304c0-14.165333 11.434667-25.429333 25.429334-25.429333 14.165333 0 25.429333 11.434667 25.429333 25.429333v34.133334zM375.637333 186.197333h188.757334c33.109333 0 59.904 26.794667 59.904 59.904v36.010667c0 14.165333-11.434667 25.429333-25.429334 25.429333-14.165333 0-25.429333-11.434667-25.429333-25.429333V246.101333c0-4.778667-3.925333-8.874667-8.874667-8.874666H375.637333c-4.949333 0-8.874667 3.925333-8.874666 8.874666v531.968c0 4.778667 3.925333 8.874667 8.874666 8.874667h188.757334c4.778667 0 8.874667-3.925333 8.874666-8.874667v-34.304c0-14.165333 11.434667-25.429333 25.429334-25.429333 14.165333 0 25.429333 11.434667 25.429333 25.429333v34.304c0 33.109333-26.794667 59.904-59.904 59.904H375.637333c-33.109333 0-59.904-26.794667-59.904-59.904V246.101333c0-33.109333 26.794667-59.904 59.904-59.904z" p-id="6490"></path></svg>
                    </i>
                </el-button>
            </el-tooltip>
            
            <el-tooltip :content="t('table.split')" placement="bottom" effect='dark'>
                <el-button text size="small" @click="editor.chain().focus().splitCell().run()">
                    <i class="el-icon" style="width: 1.4em">
                        <svg style="width: 1.4em; height: 1.4em" t="1725548432351" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6643" width="200" height="200"><path d="M882.688 522.24h-196.608v-40.96h196.608l-47.104-47.104c-8.192-8.192-8.192-20.48 0-28.672 8.192-8.192 20.48-8.192 28.672 0l96.256 96.256-96.256 96.256c-8.192 8.192-20.48 8.192-28.672 0-8.192-8.192-8.192-20.48 0-28.672l47.104-47.104z m-741.376-40.96h196.608v40.96h-196.608l47.104 47.104c8.192 8.192 8.192 20.48 0 28.672-8.192 8.192-20.48 8.192-28.672 0l-96.256-96.256 96.256-96.256c8.192-8.192 20.48-8.192 28.672 0 8.192 8.192 8.192 20.48 0 28.672l-47.104 47.104z m718.848 286.72h40.96v61.44c0 34.816-26.624 61.44-61.44 61.44h-225.28c-34.816 0-61.44-26.624-61.44-61.44v-634.88c0-34.816 26.624-61.44 61.44-61.44h225.28c34.816 0 61.44 26.624 61.44 61.44v61.44h-40.96v-61.44c0-10.24-8.192-20.48-20.48-20.48h-225.28c-10.24 0-20.48 8.192-20.48 20.48v634.88c0 12.288 8.192 20.48 20.48 20.48h225.28c10.24 0 20.48-8.192 20.48-20.48v-61.44z m-737.28 0h40.96v61.44c0 12.288 8.192 20.48 20.48 20.48h225.28c10.24 0 20.48-10.24 20.48-20.48v-634.88c0-12.288-8.192-20.48-20.48-20.48h-225.28c-10.24 0-20.48 10.24-20.48 20.48v63.488h-40.96V194.56c0-34.816 26.624-61.44 61.44-61.44h225.28c34.816 0 61.44 26.624 61.44 61.44v634.88c0 34.816-26.624 61.44-61.44 61.44h-225.28c-34.816 0-61.44-26.624-61.44-61.44v-61.44z" p-id="6644"></path></svg>
                    </i>
                </el-button>
            </el-tooltip>
            
            <el-tooltip :content="t('table.insertRow')" placement="bottom" effect='dark'>
                <el-button text size="small" @click="editor.chain().focus().addRowAfter().run()">
                    <i class="el-icon" style="width: 1.4em">
                        <svg style="width: 1.4em; height: 1.4em" t="1725549065412" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="2208" width="200" height="200"><path d="M552.96 286.72v-122.88h-409.6v-40.96h450.56v204.8h-450.56v-40.96h409.6z m0 573.44v-122.88h-409.6v-40.96h450.56v204.8h-450.56v-40.96h409.6z m-292.864-415.744l-67.584 67.584 67.584 67.584c8.192 8.192 8.192 20.48 0 28.672-8.192 8.192-20.48 8.192-28.672 0l-96.256-96.256 96.256-96.256c8.192-8.192 20.48-8.192 28.672 0 8.192 8.192 8.192 20.48 0 28.672z m129.024-34.816h512v204.8h-512v-204.8z m40.96 163.84h430.08v-122.88h-430.08v122.88z" p-id="2209"></path></svg>
                    </i>
                </el-button>
            </el-tooltip>
            
            <el-tooltip :content="t('table.deleteRow')" placement="bottom" effect='dark'>
                <el-button text size="small" @click="editor.chain().focus().deleteRow().run()">
                    <i class="el-icon" style="width: 1.4em">
                        <svg style="width: 1.4em; height: 1.4em" t="1725549090273" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="2366" width="200" height="200"><path d="M235.52 540.672l67.584 67.584c8.192 8.192 20.48 8.192 28.672 0 8.192-8.192 8.192-20.48 0-28.672l-67.584-67.584 67.584-67.584c8.192-8.192 8.192-20.48 0-28.672-8.192-8.192-20.48-8.192-28.672 0l-67.584 67.584-67.584-67.584c-8.192-8.192-20.48-8.192-28.672 0-8.192 8.192-8.192 20.48 0 28.672l67.584 67.584-67.584 67.584c-8.192 8.192-8.192 20.48 0 28.672 8.192 8.192 20.48 8.192 28.672 0l67.584-67.584z m307.2-253.952v-122.88h-409.6v-40.96h450.56v204.8h-450.56v-40.96h409.6z m0 573.44v-122.88h-409.6v-40.96h450.56v204.8h-450.56v-40.96h409.6z m-61.44-450.56h409.6v204.8h-409.6v-204.8z m40.96 163.84h327.68v-122.88h-327.68v122.88z" p-id="2367"></path></svg>
                    </i>
                </el-button>
            </el-tooltip>
            
            <el-tooltip :content="t('table.insertColumn')" placement="bottom" effect='dark'>
                <el-button text size="small" @click="editor.chain().focus().addColumnAfter().run()">
                    <i class="el-icon" style="width: 1.4em">
                        <svg style="width: 1.4em; height: 1.4em" t="1725549149331" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="2524" width="200" height="200"><path d="M286.72 481.28h-122.88v409.6h-40.96v-450.56h204.8v450.56h-40.96v-409.6z m573.44 0h-122.88v409.6h-40.96v-450.56h204.8v450.56h-40.96v-409.6z m-280.576 403.456l-67.584-67.584-67.584 67.584c-8.192 8.192-20.48 8.192-28.672 0-8.192-8.192-8.192-20.48 0-28.672l96.256-96.256 96.256 96.256c8.192 8.192 8.192 20.48 0 28.672-8.192 8.192-20.48 8.192-28.672 0zM409.6 645.12v-512h204.8v512h-204.8z m163.84-40.96v-430.08h-122.88v430.08h122.88z" p-id="2525"></path></svg>
                    </i>
                </el-button>
            </el-tooltip>
            
            <el-tooltip :content="t('table.deleteColumn')" placement="bottom" effect='dark'>
                <el-button text size="small" @click="editor.chain().focus().deleteColumn().run()">
                    <i class="el-icon" style="width: 1.4em">
                        <svg style="width: 1.4em; height: 1.4em" t="1725549164667" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="2682" width="200" height="200"><path d="M540.672 788.48l67.584-67.584c8.192-8.192 8.192-20.48 0-28.672-8.192-8.192-20.48-8.192-28.672 0l-67.584 67.584-67.584-67.584c-8.192-8.192-20.48-8.192-28.672 0-8.192 8.192-8.192 20.48 0 28.672l67.584 67.584-67.584 67.584c-8.192 8.192-8.192 20.48 0 28.672 8.192 8.192 20.48 8.192 28.672 0l67.584-67.584 67.584 67.584c8.192 8.192 20.48 8.192 28.672 0 8.192-8.192 8.192-20.48 0-28.672l-67.584-67.584z m-253.952-307.2h-122.88v409.6h-40.96v-450.56h204.8v450.56h-40.96v-409.6z m573.44 0h-122.88v409.6h-40.96v-450.56h204.8v450.56h-40.96v-409.6z m-450.56 61.44v-409.6h204.8v409.6h-204.8z m163.84-40.96v-327.68h-122.88v327.68h122.88z" p-id="2683"></path></svg>
                    </i>
                </el-button>
            </el-tooltip>

            <el-tooltip :content="t('table.deleteTable')" placement="bottom" effect='dark'>
                <el-button text size="small" @click="editor.chain().focus().deleteTable().run()">
                    <i class="el-icon" style="width: 1.4em">
                        <svg t="1725549411702" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="10786" width="200" height="200"><path d="M984.987 690.066a41.235 41.235 0 0 0-58.313 0l-83.42 83.42-83.42-83.42a41.235 41.235 0 0 0-58.314 0 41.235 41.235 0 0 0 0 58.313l83.42 83.42-83.42 83.387a41.235 41.235 0 0 0 0 58.313 41.235 41.235 0 0 0 58.313 0l83.42-83.42 83.387 83.454a41.235 41.235 0 0 0 58.313 0 41.235 41.235 0 0 0 0-58.313l-83.386-83.42 83.42-83.42a41.235 41.235 0 0 0 0-58.314z" fill="#666666" p-id="10787"></path><path d="M0 0v1024h672.988a256.813 256.813 0 0 1-62.036-82.47H390.208V698.702h233.745a257.306 257.306 0 0 1 83.157-83.867V392.877h234.42v203.15A256.4 256.4 0 0 1 1024 651.66V0z m316.902 941.53H82.47V698.702h234.432z m0-316.134H82.47V392.877h234.432z m0-305.826H82.47V82.47h234.432z m316.89 305.826H390.208V392.877h243.584z m0-305.826H390.208V82.47h243.584z m307.738 0H707.098V82.47H941.53z" p-id="10788"></path></svg>
                    </i>
                </el-button>
            </el-tooltip>
        </div>
    </div>
</template>

<style scoped>
.bubble-btn {
    background-color: transparent;
    cursor: pointer;
    border: none;
    padding: 8px 15px;
    display: inline-flex;
    justify-content: center;
    align-items: center;
    line-height: 1;
    text-align: center;
    box-sizing: border-box;
    outline: none;
    transition: .1s;
    font-weight: var(--el-button-font-weight);
    white-space: nowrap;
    font-size: var(--el-font-size-base);
    border-radius: var(--el-border-radius-base);
}

.bubble-btn:hover {
    background-color: #9a9a9c;
}

.btn-active {
    background-color: #9a9a9c;
}

.position-btn {
    display: flex;
    align-items: center;
    justify-content: center;
}

.size-bar {
    display: flex;
    align-items: center;
    width: 300px;
}

.bubble-menu {
    background-color: var(--el-bg-color-overlay);
    border: 1px solid rgb(39, 39, 39);
    border-radius: 0.7rem;
    box-shadow: #1111112e;
    display: flex;
    padding: 0.2rem 0.9rem;
    flex-direction: column;
}

.inner-image{
    display: flex;
    flex-direction: column;
    align-items: center;
    height: 420px;
    width: 100%;
}

.display-box {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 90%;
    width: 100%;
}

.display-box img {
    object-fit: contain;
}

.func_bar {
    width: 70%;
    height: 10%;
}

.media_container{
    height: 400px;
    overflow: hidden;
}
.media-toolbar{
    width: 100%;
    height: 8%;
    display: flex;
    justify-content:end;
    align-items: center;
    margin-left: -10px;
}

.dialog {
    position: fixed;
    background-color: white;
    width: 300px;
    padding: 20px;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 1000;
}

.el-upload {
    border: 1px dashed var(--el-border-color) !important;
    border-radius: 6px !important;
    cursor: pointer !important;
    position: relative !important;
    overflow: hidden !important;
    transition: var(--el-transition-duration-fast) !important;
  }
  
.el-upload:hover {
border-color: var(--el-color-primary) !important;
}
  
.el-icon.avatar-uploader-icon {
font-size: 28px !important;
color: #8c939d !important;
width: 178px !important;
height: 178px !important;
text-align: center !important;
}

.form-cover-upload {
    display: inline-flex;
    flex-wrap: wrap;
    background-color: transparent;
    height: 174px;
    width: 174px;
    border: 1px dashed var(--el-border-color);
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    justify-content: center;
    align-items: center;
    transition: var(--el-transition-duration-fast)
}

.form-cover-upload:hover {
    border-color: var(--el-color-primary);
}

.form-cover-upload .el-icon {
    font-size: 28px;
    color: var(--el-text-color-secondary);
}

.image-preview_form {
    height: 174px;
    width: 174px;
    display: flex;
    justify-content: center;
    align-items: center;
    overflow: hidden;
}

.image-preview_form img {
    height: 100%;
    width: 100%;
    object-fit: contain;
}

.list__item-actions {
    position: absolute;
    width: 174px;
    height: 100%;
    top: 0;
    left: 0;
    display: inline-flex;
    justify-content: center;
    align-items: center;
    color: #fff;
    font-size: 20px;
    cursor: pointer;
    opacity: 0;
    background-color: var(--el-overlay-color-lighter);
    transition: opacity var(--el-transition-duration);
}

.list__item-actions:hover {
    opacity: 1;
}

.list__item-actions span:last-child {
    margin-left: 40px;
}

.media-display{
    margin-top: 10px;
    height: 90%;
    width: 100%;
    overflow-y: auto !important;
}

.media-display .el-image {
    min-height: 200px;
    margin-bottom: 10px;
    margin-left: 12px;
}
.media-display .el-image:last-child {
    margin-bottom: 0;
}

</style>