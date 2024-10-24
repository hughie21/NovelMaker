<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-10-10
@Description: This is the main dialog component that contains all the sub dialog components.
*/

import Lookup from './dialog/lookup.vue'
import Setting from './dialog/setting.vue';
import BookInfo from './dialog/bookinfo.vue';
import Media from './dialog/media.vue';
import Help from './dialog/help.vue';
import Bubble from './dialog/bubble.vue';

</script>

<template>
    <!--
    Setting dialog
    -->
    <Setting></Setting>

    <!--
    Book info dialog
    -->
    <BookInfo></BookInfo>

    <!--
    Media dialog
    -->
    <Media></Media>

    <!-- help dialog -->
    <Help></Help>

    <!-- bubble menu -->
     <Bubble></Bubble>

    <!-- search bar -->
    <Lookup></Lookup>
</template>

<style>
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



</style>