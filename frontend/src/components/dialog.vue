<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-08-16
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
import { i18nChangeLanguage } from '@wangeditor/editor'

const { t, locale } = useI18n();
const bookinfo = inject("bookinfo");
const arrayRemove = utils.arrayRemove;
const editorRef = utils.editorRef;
const staticFiles = ref([]);
const cover = utils.cover;
const initCover = utils.initCover;
const dialogImageUrl = ref('');
const dialogVisible = ref(false);
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
    i18nChangeLanguage(lang);
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
    const imgNode = {type: "image", children: [{ text: '' }], src: url, alt: id};
    const editor = editorRef.value;
    editor.insertNode(imgNode)
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
</template>

<style scoped>
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