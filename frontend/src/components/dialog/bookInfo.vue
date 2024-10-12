<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-9
@LastEditors: Hughie
@LastEditTime: 2024-10-10
@Description: This the diaglog that allow user to editor the bookinfo.
*/
import { visio, change, cover, arrayRemove, language_list as languages, initCover } from '../../assets/js/utils';
import { ref, inject, nextTick } from 'vue';
import { useI18n } from 'vue-i18n';
import { Delete, Plus, ZoomIn } from '@element-plus/icons-vue'
import { OpenImage } from '../../../wailsjs/go/main/App.js'

const { t } = useI18n();
const bookinfo = inject("bookinfo");
const dialogImageUrl = ref('');
const dialogVisible = ref(false);
const inputValue = ref('');
const inputVisible = ref(false);
const InputRef = ref();
initCover();

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

const handleRemove = () => {
    arrayRemove(bookinfo.metadata.meta, "cover");
    arrayRemove(bookinfo.resources, "cover.jpg");
    initCover();
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
const handleBookInfo = () => {
    visio.bookInfoVisible = false;
}
</script>

<template>
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
        @input="()=>{change = true}"
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