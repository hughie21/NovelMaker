<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-09-23
@Description: This is the tab bar of the application.
*/

import { ref, onMounted, provide} from 'vue';
import { editorRef, fonts } from '../assets/js/globals';
import fileTab from '../components/tab/file.vue'
import editTab from '../components/tab/edit.vue'
import helpTab from '../components/tab/help.vue'
import "../assets/css/tab.css"

async function getFonts() {
    if(!("queryLocalFonts" in window)){
        return [
            {
                value: "serif",
                label: "serif"
            },
            {
                value: "sans-serif",
                label: "sans-serif"
            },
            {
                value: "monospace",
                label: "monospace"
            },
            {
                value: "Inter",
                label: "Inter"
            },
            {
                value: "cursive",
                label: "cursive"
            },
            {
                value: "Comic Sans",
                label: "Comic Sans"
            }
        ];
    }
    const res = await window.queryLocalFonts();
    let familys = []
    let formattedFonts = []
    for(var i = 0; i < res.length; i++) {
        if(familys.includes(res[i].family)){
            continue;
        }
        familys.push(res[i].family);
        formattedFonts.push({
            value: res[i].family,
            label: res[i].fullName
        })
    }
    fonts.value = formattedFonts;
}

const activeNames = ref("1")
const btnNormalClass = ref('el-button btn func_btn-big');

provide("btnNormalClass", btnNormalClass);

const judgeHeight = () => {
    var tabHeight = document.querySelector('.tab-container').offsetHeight;
    var totalHeight = window.innerHeight;
    var elem = document.querySelector('#editor-container');
    elem.setAttribute('style', 'height:' + (totalHeight - tabHeight - 35) + 'px');
}
onMounted(() => {
    setTimeout(()=>{
        judgeHeight();
        getFonts();
    }, 500)
})

const handleChange = () => {
    setTimeout(()=>{
        judgeHeight();
    }, 300)
}

const test = () => {
    const E = editorRef.value;
    console.log(E.getHTML());
}
</script>

<template>
<button id="testing" @click="test" style="display: none;"></button>
<el-tabs type="border-card">
    <el-collapse v-model="activeNames" @change="handleChange">
        <el-collapse-item name="1" id="tab-collapse">
            <!-- File tab -->
            <el-tab-pane :label="$t('nav.file')">
                <fileTab></fileTab>
            </el-tab-pane>

            <!-- Edit tab -->
            <el-tab-pane :label="$t('nav.edit')">
                <editTab></editTab>
            </el-tab-pane>

            <!-- Help tab -->
            <el-tab-pane :label="$t('nav.help')">
                <helpTab></helpTab>
            </el-tab-pane>
        </el-collapse-item>
    </el-collapse>
</el-tabs>
</template>

<style>
#tab-collapse .el-collapse-item__header {
    height: 1em !important;
}

#tab-collapse .el-collapse-item__content{
    padding-bottom: 5px !important;
}

.custom-tree-node {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 14px;
    padding-right: 8px;
}
.node-label{
    width: 90%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    background-color: transparent;
    border: none;
}
.novel_tree-title{
    display: flex;
    justify-content: center;
    margin-top: 5px;
}
.novel_tree-title input{
    border: none;
    background-color: transparent;
    font-size: 20px;
    font-weight: bold;
}

.novel_tree-title input:focus{
    outline: none;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.novel_tree-title input::placeholder{
    color: #9b9b9b;
}
</style>