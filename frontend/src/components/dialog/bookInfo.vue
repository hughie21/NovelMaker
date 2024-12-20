<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-9
@LastEditors: Hughie
@LastEditTime: 2024-10-15
@Description: This the diaglog that allow user to editor the bookinfo.
*/
import {  arrayRemove, initCover, deepClone  } from '../../assets/js/utils.js';
import { visio, change, cover, language_list as languages, bookInfo } from '../../assets/js/globals.js';
import { ref, nextTick, reactive, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { Delete, Plus, ZoomIn } from '@element-plus/icons-vue'
import { OpenImage } from '../../../wailsjs/go/core/App.js'
import { ElMessage } from 'element-plus';
const { t } = useI18n();

// deep copy the metadata from the booinfo
const metadata = reactive(deepClone(bookInfo.metadata));
// the cover image
const dialogImageUrl = ref('');
const dialogVisible = ref(false);
const inputValue = ref('');
const inputVisible = ref(false);
const InputRef = ref();
const formRef = ref();
const illegalCharRegex = /[\"|\\]/;
// Calibration rules: the title can not be empty and can not contain illegal characters like " and \
const rules = reactive({
    title: [
        { required: true, message: t('dialog.info.requireMessage'), trigger: 'blur' },
        {validator: (rule, value, callback) => {
            if (illegalCharRegex.test(value)) {
                callback(new Error(t('dialog.info.illegalMessage')));
            } else {
                callback();
            }
        }, trigger: 'blur'},
    ],
    identifier: [
        {validator: (rule, value, callback)=> {
            if (illegalCharRegex.test(value)) {
                callback(new Error(t('dialog.info.illegalMessage')));
            } else {
                callback();
            }

        }, trigger: 'blur'}
    ],
    publisher: [
        {validator: (rule, value, callback)=> {
            if (illegalCharRegex.test(value)) {
                callback(new Error(t('dialog.info.illegalMessage')));
            } else {
                callback();
            }

        }, trigger: 'blur'}
    ],
    author: [
        {validator: (rule, value, callback)=> {
            if (illegalCharRegex.test(metadata.creator)) {
                callback(new Error(t('dialog.info.illegalMessage')));
            } else {
                callback();
            }

        }, trigger: 'blur'}
    ],
    contributors: [
        {validator: (rule, value, callback)=> {
            if (illegalCharRegex.test(value)) {
                callback(new Error(t('dialog.info.illegalMessage')));
            } else {
                callback();
            }

        }, trigger: 'blur'}
    ],
    description: [
        {validator: (rule, value, callback)=> {
            if (illegalCharRegex.test(value)) {
                callback(new Error(t('dialog.info.illegalMessage')));
            } else {
                callback();
            }

        }, trigger: 'blur'}
    ],
    subject: [
        {validator: (rule, value, callback)=> {
            metadata.subject.forEach((item) => {
                if (illegalCharRegex.test(item)) {
                    return callback(new Error(t('dialog.info.illegalMessage')));
                }
            })
            callback();
        }, trigger: 'blur'}
    ]
})

initCover();

// Monitor changes to raw data such as new or open files
watch(bookInfo, (newVal, oldVal) => {
    Object.assign(metadata, deepClone(newVal.metadata));
})

const handlePictureCardPreview = (file) => {
  dialogImageUrl.value = "data:image/jpg;base64," + file;
  dialogVisible.value = true;
}

const handleClose = (tag)=>{
    metadata.subject.splice(metadata.subject.indexOf(tag), 1);
}
const showInput = ()=>{
    inputVisible.value = true;
    nextTick(() => {
        InputRef.value.input.focus();
    })
}

const handleInputConfirm = () => {
  if (inputValue.value) {
    metadata.subject.push(inputValue.value);
  }
  inputVisible.value = false;
  inputValue.value = '';
}

const handleRemove = () => {
    bookInfo.metadata.cover.name = "";
    bookInfo.metadata.cover.data = "";
    arrayRemove(bookInfo.resources, "cover.jpg");
    initCover();
}

const uploadCover = async () => {
    var coverData = await OpenImage().then((res)=>{
       return res.Data;
    });
    bookInfo.metadata.cover.name = "cover.jpg";
    bookInfo.metadata.cover.data = coverData;
    initCover();
}
const handleBookInfo = (formEl) => {
    if (!formEl) return;
    formEl.validate((valid) => {
        if (valid) {
            let temp = deepClone(metadata);
            console.log(temp);
            bookInfo.metadata.title = temp.title;
            bookInfo.metadata.identifier = temp.identifier;
            bookInfo.metadata.contributors = temp.contributors;
            bookInfo.metadata.description = temp.description;
            bookInfo.metadata.subject = temp.subject;
            bookInfo.metadata.language = temp.language;
            bookInfo.metadata.publisher = temp.publisher;
            bookInfo.metadata.creator = temp.creator;
            visio.bookInfoVisible = false;
        }else {
            ElMessage.error(t('message.bookinfoSaveError'));
            return;
        }
    })
    console.log("pass")
}

</script>

<template>
    <el-dialog
    v-model="visio.bookInfoVisible"
    :title="t('dialog.info.title')"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    width="500"
    >
        <el-form
        label-position="top"
        label-width="auto"
        :model="metadata"
        ref="formRef"
        style="max-width: 600px"
        :rules="rules"
        @input="()=>{change = true}"
        >
            <el-form-item :label="t('dialog.info.name')" prop="title">
                <el-input v-model="metadata.title" />
            </el-form-item>
            <el-form-item :label="t('dialog.info.identifier')"  prop="identifier">
                <el-input v-model="metadata.identifier" />
            </el-form-item>
            <el-form-item :label="t('dialog.info.publisher')" prop="publisher">
                <el-input v-model="metadata.publisher" />
            </el-form-item>
            <el-form-item :label="t('dialog.info.author')" prop="author">
                <el-input v-model="metadata.creator" />
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
            <el-form-item :label="t('dialog.info.contributer')" prop="contributors">
                <el-input v-model="metadata.contributors" />
            </el-form-item>
            <el-form-item :label="t('dialog.info.description')" prop="description">
                <el-input 
                v-model="metadata.description" 
                type="textarea"
                :rows="3"/>
            </el-form-item>
            <el-form-item :label="t('dialog.info.subject')" prop="subject">
                <div class="flex gap-2">
                    <el-tag
                    v-for="tag in metadata.subject"
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
                v-model="metadata.language"
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
              <el-button type="primary" @click="handleBookInfo(formRef)">
                {{$t('dialog.confirm')}}
              </el-button>
            </div>
        </template>
    </el-dialog>
</template>

<style scoped>
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
</style>