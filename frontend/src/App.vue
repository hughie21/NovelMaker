<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-09-23
@Description: This is main frame of the application.
*/

import tab from './components/tab.vue'
import dialogs from './components/dialog.vue'
import editor from './components/editor.vue'
import { useI18n } from 'vue-i18n';
import {checkIfOpenFileDirectly} from './assets/js/utils'

const { t } = useI18n();

window.onerror = (message, source, lineno, colno, error) => {
  console.log("Error: " + message + " in " + source + " at line " + lineno + " and column " + colno);
  console.log(error);
}

checkIfOpenFileDirectly(t);
</script>

<template>
  <div class="tab-container">
    <dialogs></dialogs>
    <tab></tab>
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
.editor-container {
  height: 100%;
  width: 100%;
}
</style>
