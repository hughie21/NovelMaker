<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-9
@LastEditors: Hughie
@LastEditTime: 2024-10-10
@Description: This is the dialog allow user to manage the media resources.
*/

import {  getImageFiles } from '../../assets/js/utils';
import { visio, imageInfo, editorRef, staticFiles } from '../../assets/js/globals.js';
import { ref } from 'vue';
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'
import $ from 'jquery'
import { ImageUpload, FileDelete, ImageDownload } from '../../../wailsjs/go/core/App.js'
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
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

getImageFiles();

const uploadImage = async () => {
    var imageData = await ImageUpload().then((res)=>{
       return res;
    })
    if(imageData.Code == 0) {
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

const downloadImage = () => {
    ElMessageBox.prompt(t('dialog.media.downloadImagePromt'), t('dialog.media.downloadImageTitle'), {
        confirmButtonText: t('message.confirm'),
        cancelButtonText: t('message.cancel'),
        inputPattern: /(?<!@)\b((http(s)?:\/\/)?((\d+)\.)?[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})(\/[^\s]*)?(.jpg|.png|.jpeg|.bmp|.gif|.webp)\b/,
        inputErrorMessage: t('message.invalidLink')
    }).then(async ({ value }) => {
        ElNotification({
            title: t('message.downloadTitle'),
            message: t('message.downloadNotice'),
            type: 'info'
        })
        await ImageDownload(value).then(res => {
            if(res.Code == 0) {
                getImageFiles();
                ElMessage({
                    message: t("message.downloadSuccess"),
                    type: "success"
                });
            }else if(res.Code == 1) {
                ElMessage({
                    message: t("message.downloadError") + res.Msg,
                    type: "error"
                });
            }
        })
    })
}


const handleInsert = () => {
    var elem = lastClick.value;
    const url = elem.attr("src");
    const id = elem.attr("id");
    const imgNode = {src: url, alt: id, title: id};
    const editor = editorRef.value;
    editor.chain().focus().InsertImage(imgNode).createParagraphNear().run();

    visio.mediaVisible = false;
}

const handleDelImg = () => {
    var elem = lastClick.value;
    const E = editorRef.value;
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
                const state = E.state;
                const doc = state.doc;
                const positionsToDelete = [];
                doc.descendants((node, pos) => {
                    if (node.type.name == "image" && node.attrs.alt == name) {
                        // record the positions to delete, if we delete them now, it will 
                        // change the structure of the document and the next position will be wrong
                        positionsToDelete.push({ from: pos, to: pos + node.nodeSize });
                    }
                });
                positionsToDelete.reverse().forEach(({ from, to }) => {
                    E.chain().focus().deleteRange({ from, to }).run();
                });
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
    <el-dialog
    v-model="visio.mediaVisible"
    :title="t('dialog.media.title')"
    width="60%"
    >
    <div class="media_container"> 
        <div class="media-toolbar">
            <el-button type="primary" plain @click="uploadImage">{{$t("dialog.media.upload")}}</el-button>
            <el-button type="primary" plain @click="downloadImage">{{$t("dialog.media.link")}}</el-button>
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