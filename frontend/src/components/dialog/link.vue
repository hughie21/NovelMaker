<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-9
@LastEditors: Hughie
@LastEditTime: 2024-12-20
@Description: Hyperlink insertion popup
*/
import { editorRef, visio } from '../../assets/js/globals.js';
import { useI18n } from 'vue-i18n';
import { ref } from 'vue';
import { ElMessageBox } from 'element-plus'

const { t } = useI18n();
const url = ref('');
const prefix = ref('0');

const checkUrl = (text) => {
    if (text === '') {
        return false;
    }
    const urlReg = /(?<!@)\b((http(s)?:\/\/)?((\w+)\.)?[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})(\/[^\s]*)?\b/g
    if (!urlReg.test(text)) {
        return false;
    }
    return true;
}

const checkEmail = (text) => {
    if (text === '') {
        return false;
    }
    const emailReg = /\b([a-zA-Z0-9._%-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})\b/g
    if (!emailReg.test(text)) {
        return false;
    }
    return true;
}

const handleInsertLink = () => {
    const E = editorRef.value;
    var urlText = '';
    switch (prefix.value) {
        case "0":
            urlText = `http://${url.value}`
            if(!checkUrl(urlText)) {
                ElMessageBox.alert(t('dialog.link.error'), t('message.error'), {
                    confirmButtonText: t('message.confirm'),
                })
                return;
            }
            E.commands.setLink({ href: urlText });
            break;
        case "1":
            urlText = `https://${url.value}`
            if(!checkUrl(urlText)) {
                ElMessageBox.alert(t('dialog.link.error'), t('message.error'), {
                    confirmButtonText: t('message.confirm'),
                })
                return;
            }
            E.commands.setLink({ href: urlText });
            break;
        case "2":
            urlText = `ftp://${url.value}`
            if(!checkUrl(urlText)) {
                ElMessageBox.alert(t('dialog.link.error'), t('message.error'), {
                    confirmButtonText: t('message.confirm'),
                })
                return;
            }
            E.commands.setLink({ href: urlText });
            break;
        case "3":
            urlText = `mailto:${url.value}`
            if(!checkEmail(url.value)) {
                ElMessageBox.alert(t('dialog.link.error'), t('message.error'), {
                    confirmButtonText: t('message.confirm'),
                })
                return;
            }
            E.commands.setLink({ href: urlText });
            break;
    }
    visio.linkInsertVisible = false
}
</script>

<template>
    <el-dialog
    v-model="visio.linkInsertVisible"
    :title="t('dialog.link.title')"
    width="500"
    >
        <el-input
        v-model="url"
        style="max-width: 600px"
        :placeholder="t('dialog.link.text')"
        >
            <template #prepend>
                <el-select v-model="prefix" :placeholder="t('dialog.link.prefix')" style="width: 115px">
                    <el-option label="Http://" value="0" />
                    <el-option label="Https://" value="1" />
                    <el-option label="Ftp://" value="2" />
                    <el-option label="Mailto://" value="3" />
                </el-select>
            </template>
        </el-input>
        <template #footer>
            <div class="dialog-footer">
            <el-button @click="visio.linkInsertVisible = false">{{t('dialog.link.cancel')}}</el-button>
            <el-button type="primary" @click="handleInsertLink">
                {{t('dialog.link.insert')}}
            </el-button>
            </div>
        </template>
    </el-dialog>
</template>