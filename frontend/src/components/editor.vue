<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-11-1
@Description: This is the editor component, powered by tiptap
*/
import { onMounted } from 'vue';
import * as constant from '../assets/js/globals';
import { updateCatalog } from '../assets/js/utils';
import { Editor, EditorContent } from "@tiptap/vue-3";
import StarterKit from "@tiptap/starter-kit";
import { CustomImage, CustomHeading, TextBackground, TextFontSize, SearchSelBackground } from '../assets/js/extension';
import BubbleMenu from "@tiptap/extension-bubble-menu";
import TextStyle from '@tiptap/extension-text-style'
import FontFamily from '@tiptap/extension-font-family'
import { Color } from '@tiptap/extension-color'
import Table from '@tiptap/extension-table'
import TableCell from '@tiptap/extension-table-cell'
import TableHeader from '@tiptap/extension-table-header'
import TableRow from '@tiptap/extension-table-row'
import '../assets/css/editor.css';

const tableMenu = document.getElementById("tableMenu").children[0];
 
function throttle(func, wait) {
    let timeout;
    return function(...args) {
        if (!timeout) {
            timeout = setTimeout(() => {
                timeout = null;
                func.apply(this, args);
            }, wait);
        }
    };
}

const editor = new Editor({
    content: ``,
    extensions: [
        StarterKit,
        CustomImage,
        // ImageResize,
        CustomHeading,
        // BubbleMenu.configure({
        //     pluginKey: "bubbleMenuImage",
        //     element: bubbleMenu,
        //     shouldShow: (({ editor, view, state, oldState, from, to }) => {
        //         return editor.isActive("image")
        //     })
        // }),
        BubbleMenu.configure({
            pluginKey: "bubbleMenuTable",
            element: tableMenu,
            shouldShow: (({ editor, view, state, oldState, from, to }) => {
                return editor.isActive("table")
            })
        }),
        TextStyle,
        FontFamily.configure({
            types: ['textStyle'],
        }),
        TextFontSize.configure({
            types: ['textStyle'],
        }),
        Color.configure({
            types: ["textStyle"],
        }),
        TextBackground.configure({
            types: ["textStyle"],
        }),
        Table.configure({
            resizable: true,
            defaultColumns: 3,
            defaultRows: 3,
        }),
        TableRow,
        TableHeader,
        TableCell,
        SearchSelBackground.configure({
            types: ['textStyle'],
        })
    ],
    onUpdate: throttle(({ editor } ) => { // Synchronising editor header to the catelogue
        updateCatalog();
        bookinfo.content = editor.getJSON();
        constant.change.value = true;
    }, 100),
    onSelectionUpdate: ({ editor }) => { // Synchronising editor content properties to tab option values
        for(var i = 1; i < 7; i++){
            if(editor.isActive('custom-heading', { level: i})){
                constant.headerVal.value = i;
                return;
            }
        }
        if(editor.isActive('paragraph')){
            constant.headerVal.value = 0;
            constant.fontVal.value = "Arial";
            constant.fontSizeVal.value = "16px";
        }
        constant.fonts.value.forEach((v)=>{
            if(editor.getAttributes("textStyle").fontFamily == v.value) {
                constant.fontVal.value = v.value;
            }
        })
        const fontSize = editor.getAttributes("font-size").fontSize;
        if(fontSize){
            constant.fontSizeVal.value = fontSize;
        }
        return;
    }
})

const theme = constant.editTheme;
const bookinfo = constant.bookInfo;

const editorRef = constant.editorRef;
editorRef.value = editor;

onMounted(()=>{ // Initialise the catelogue and add eventListener to the elemnt of it
    const headerContainer = document.getElementById('header-container');
    headerContainer.addEventListener('mousedown', event  => {
        if (event.target.tagName !== 'LI') return
        event.preventDefault()
        const id = event.target.dataset.id;
        const $doc = editor.$doc;
        $doc.querySelectorAll("custom-heading").map((header)=>{
            if(id == header.element.id){
                editor.chain().focus(header.pos).run();
            }
        })
        
    })

})
</script>

<template>
<div style="border: 1px solid #ccc; height:100%" :class="theme" id="editor-container">
    <div class="editor-box">
        <ul id="header-container"></ul>
        <editor-content class="editor-content" id="editor" :editor="editorRef" />
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
    background-color: var(--el-fill-color-blank);
    border-right: 1px solid #9b9b9b;
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

#header-container li a {
    text-decoration: none;
}

#header-container li[type="H1"] {
font-size: 20px;
font-weight: bold;
}

#header-container li[type="H2"] {
font-size: 16px;
padding-left: 15px;
font-weight: bold;
}

#header-container li[type="H3"] {
font-size: 14px;
padding-left: 30px;
}

#header-container li[type="H4"] {
font-size: 12px;
padding-left: 45px;
}

#header-container li[type="H5"] {
font-size: 12px;
padding-left: 60px;
}

#header-container li[type="H6"] {
font-size: 12px;
padding-left: 75px;
}
</style>