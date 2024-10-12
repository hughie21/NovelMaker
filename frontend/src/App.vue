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
