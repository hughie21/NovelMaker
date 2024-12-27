<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-11-1
@Description: This is the editor component, powered by tiptap
*/
import { onMounted } from 'vue';
// import * as constant from '../assets/js/globals';
import { change, headerVal, fontSizeVal, fontVal, editTheme, editorRef, isBold, isItalic, isStrike, fonts  } from '../assets/js/globals';
import { updateCatalog } from '../assets/js/utils';
import { Editor, EditorContent } from "@tiptap/vue-3";
import StarterKit from "@tiptap/starter-kit";
import { CustomImage, CustomHeading, SearchSelBackground, TextStyleExtends } from '../assets/js/extension';
import bold from "@tiptap/extension-bold"
import BubbleMenu from "@tiptap/extension-bubble-menu";
import Table from '@tiptap/extension-table'
import TableCell from '@tiptap/extension-table-cell'
import TableHeader from '@tiptap/extension-table-header'
import TableRow from '@tiptap/extension-table-row'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import Link from '@tiptap/extension-link'
import { createLowlight, all } from 'lowlight'
import '../assets/css/editor.css';

const lowlight = createLowlight(all)

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
        CodeBlockLowlight.configure({
          lowlight,
        }),
        bold.configure({
            HTMLAttributes: {
                class: 'bold'     
            }
        }),
        StarterKit.configure({
            bold: false,
            codeBlock: false,
        }),
        CustomImage,
        CustomHeading,
        BubbleMenu.configure({
            pluginKey: "bubbleMenuTable",
            element: tableMenu,
            shouldShow: (({ editor, view, state, oldState, from, to }) => {
                return editor.isActive("table")
            })
        }),
        TextStyleExtends.configure({
            types: ['textStyle'],
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
        }),
        Link.configure({
            openOnClick: false,
            linkOnPaste: true,
        })
    ],
    onUpdate: throttle(({ editor } ) => { // Synchronising editor header to the catelogue
        updateCatalog();
        change.value = true;
    }, 100),
    onSelectionUpdate: ({ editor }) => { // Synchronising editor content properties to tab option values
        for(var i = 1; i < 7; i++){
            if(editor.isActive('custom-heading', { level: i})){
                headerVal.value = i;
                return;
            }
        }
        if(editor.isActive('bold') && editor.isActive('paragraph')) {
            isBold.value = true;
        }else {
            isBold.value = false;
        }


        if(editor.isActive('italic') && editor.isActive('paragraph')) {
            isItalic.value = true;
        }else {
            isItalic.value = false;
        }

        if(editor.isActive('strike') && editor.isActive('paragraph')) {
            isStrike.value = true;
        }else {
            isStrike.value = false;
        }

        if(editor.isActive('paragraph')){
            headerVal.value = 0;
            fontVal.value = "Arial";
            fontSizeVal.value = "16px";
        }
        fonts.value.forEach((v)=>{
            if(editor.getAttributes("textStyle").fontFamily == v.value) {
                fontVal.value = v.value;
            }
        })
        const fontSize = editor.getAttributes("textStyle").fontSize;
        if(fontSize){
            fontSizeVal.value = fontSize;
        }
        return;
    }
})

editorRef.value = editor;

// Initialise the catelogue and add eventListener to the elemnt of it
onMounted(()=>{ 
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
<div style="border: 1px solid #ccc; height:100%" :class="editTheme" id="editor-container">
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
    color: var(--el-text-color-regular);
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