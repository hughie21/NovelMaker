<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-08-16
@Description: This is the editor component, powered by wangeditor
*/

import { onBeforeUnmount } from 'vue';
import * as utils from '../assets/js/utils'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import { i18nChangeLanguage } from '@wangeditor/editor'
import '@wangeditor/editor/dist/css/style.css'

const mode = 'default';
const theme = utils.editTheme;
const lang = utils.editLang;
const bookinfo = utils.bookInfo;
i18nChangeLanguage(lang.value);

const editorRef = utils.editorRef;
const valueHtml = utils.valueHtml;
const toolbarConfig = {
    excludeKeys: [
        "group-image",
        "group-video",
        "fullScreen"
    ]
};
const editorConfig = {
    MENU_CONF: {}
};

let loadingEditor = true;

const handleCreated = (editor) => {
    editorRef.value = editor
    const headerContainer = document.getElementById('header-container')
    headerContainer.addEventListener('mousedown', event => {
        if (event.target.tagName !== 'LI') return
        event.preventDefault()
        const id = event.target.dataset.id
        // console.log(id)
        editor.scrollToElem(id)
    })
}

onBeforeUnmount(() => {
    const editor = editorRef.value
    if (editor == null) return
    editor.destroy()
})

const handleChange = (editor) => {
    const headers = editor.getElemsByTypePrefix('header');
    const headerContainer = document.getElementById('header-container');
    headerContainer.innerHTML = headers.map(header => {
        const id = header.id;
        const type = header.type;
        return `<li data-id="${id}" type="${type}">${header.children[0].text}</li>`;
    }).join('');
    bookinfo.content = valueHtml.value;
    if(loadingEditor){ // skip the first change event when the editor created
        loadingEditor = false;
        return;
    }
    utils.change.value = true;
}
</script>

<template>
<div style="border: 1px solid #ccc; height:100%" :class="theme" id="editor-container">
    <Toolbar
        style="border-bottom: 1px solid #ccc"
        :editor="editorRef"
        :defaultConfig="toolbarConfig"
        :mode="mode"
    />
    <div class="editor-box">
        <ul id="header-container"></ul>
        <Editor
        style="height:100%; width: 80%; overflow-y: hidden;"
        v-model="valueHtml"
        :defaultConfig="editorConfig"
        :mode="mode"
        @onCreated="handleCreated"
        @onChange="handleChange"
        />
    </div>
    
</div>
</template>

<style>
.editor-box {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: row;
}

#header-container {
    list-style-type: none;
    width: 20%;
    height: 100%;
    background-color: var(--w-e-textarea-bg-color);
    border-right: 1px solid var(--w-e-textarea-border-color);
    overflow-y: auto;
    margin: 0;
    padding: 0 0 0 10px;
}

#header-container li {
    color: var(--w-e-textarea-color);
    margin: 10px 0;
    cursor: pointer;
}
#header-container li:hover {
    text-decoration: underline;
}

#header-container li[type="header1"] {
font-size: 20px;
font-weight: bold;
}

#header-container li[type="header2"] {
font-size: 16px;
padding-left: 15px;
font-weight: bold;
}

#header-container li[type="header3"] {
font-size: 14px;
padding-left: 30px;
}

#header-container li[type="header4"] {
font-size: 12px;
padding-left: 45px;
}

#header-container li[type="header5"] {
font-size: 12px;
padding-left: 60px;
}

.dark{
    --w-e-textarea-bg-color: #111111;
    --w-e-textarea-color: #ebe9e9;
    --w-e-textarea-border-color: #2b2b2b;
    --w-e-textarea-slight-border-color: #222222;
    --w-e-textarea-slight-color: #3d3d3d;
    --w-e-textarea-slight-bg-color: #202020;
    --w-e-textarea-selected-border-color: #3973bd; 
    --w-e-textarea-handler-bg-color: #10478e; 

    --w-e-toolbar-color: #dfdfdf;
    --w-e-toolbar-bg-color: #221f1f;
    --w-e-toolbar-active-color: #dfdfdf;
    --w-e-toolbar-active-bg-color: #555555;
    --w-e-toolbar-disabled-color: #4a4a4a;
    --w-e-toolbar-border-color: #626262;

    --w-e-modal-button-bg-color: #686363;
    --w-e-modal-button-border-color: #d7d7d7;
}
</style>