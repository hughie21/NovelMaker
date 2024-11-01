<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-9
@LastEditors: Hughie
@LastEditTime: 2024-10-10
@Description: This is the dialog show the help information.
*/

import { visio } from '../../assets/js/globals.js';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>

<template>
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
</template>

<style scoped>

</style>