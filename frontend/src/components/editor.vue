<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2025-3-28
@Description: This is the editor component, powered by tiptap
*/
import { onMounted, onUnmounted, ref } from 'vue';
import { change, headerVal, fontSizeVal, fontVal, editTheme, editorRef, isBold, isItalic, isStrike, fonts, editorWidth, cateWidth, isTextAlign } from '../assets/js/globals';
import { updateCatalog } from '../assets/js/utils';
import { Editor, EditorContent } from "@tiptap/vue-3";
import StarterKit from "@tiptap/starter-kit";
import CustomImage from '../assets/extension/image';
import CustomHeading from '../assets/extension/header';
import {SearchSelBackground, TextStyleExtends} from '../assets/extension/textStyle';
import bold from "@tiptap/extension-bold"
import BubbleMenu from "@tiptap/extension-bubble-menu";
import Table from '@tiptap/extension-table'
import TableCell from '@tiptap/extension-table-cell'
import TableHeader from '@tiptap/extension-table-header'
import TableRow from '@tiptap/extension-table-row'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import MathExtension from '@aarkue/tiptap-math-extension'
import TextAlign from '../assets/extension/textalign';
import Link from '@tiptap/extension-link'
import RubyExtension from '../assets/extension/ruby'
import { createLowlight, all } from 'lowlight'
import '../assets/css/editor.css';
import "katex/dist/katex.min.css"

const lowlight = createLowlight(all)
const curHeight = ref(1);
const activeIndex = ref(0);
const headerHeight = ref([]);
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
        }),
        MathExtension.configure({
            evaluation: false
        }),
        TextAlign,
        RubyExtension,
    ],
    onUpdate: throttle(({ editor } ) => { // Synchronising editor header to the catelogue
        updateCatalog();
        getHeadersHeight();
        change.value = true;
        activeScroll();
    }, 100),
    onSelectionUpdate: ({ editor }) => { // Synchronising editor content properties to tab option values
        for(var i = 1; i < 7; i++){
            if(editor.isActive('custom-heading', { level: i})){
                headerVal.value = i;
                return;
            }
        }

        if(editor.isActive('paragraph', {textAlign: 'left'})) {
            isTextAlign.value = 0b00;
        }

        if(editor.isActive('paragraph', {textAlign: 'center'})) {
            isTextAlign.value = 0b01;
        }

        if(editor.isActive('paragraph', {textAlign: 'right'})) {
            isTextAlign.value = 0b10;
        }

        if(editor.isActive('paragraph', {textAlign: 'justify'})) {
            isTextAlign.value = 0b11;
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
        const { state } = editor;
        // get the current cursor position
        const cursorPosition = state.selection.$from.pos;
        // get the editor scroll height
        const editorScrollContainer = document.getElementById('editor');
        const scrollTop = editorScrollContainer.scrollTop;
        if (cursorPosition !== undefined) {
            const coords = editor.view.coordsAtPos(cursorPosition, -1);
            curHeight.value = coords.top + scrollTop; // absolute position of the cursor
        }
        activeScroll();
        return;
    }
})

const getHeadersHeight = () => {
    let arr = [];
    const editorScrollContainer = document.getElementById('editor');
    const headers = editorScrollContainer.querySelectorAll('h1, h2, h3, h4, h5, h6');
    const scrollTop = editorScrollContainer.scrollTop;
    headers.forEach((header) => {
        const rect = header.getBoundingClientRect();
        // get each header's top position
        arr.push(rect.top + scrollTop);
    })
    headerHeight.value = arr;
}

const activeScroll = () => {
    // get the header's top position each time the catelogue list need to be updated
    getHeadersHeight();
    let arr = headerHeight.value;
    if (arr[0] > curHeight.value) {
        activateHeader(-1);
        return (activeIndex.value = 0);
    }
    if (arr[arr.length - 1] < curHeight.value) {
        activateHeader(arr.length - 1);
        return (activeIndex.value = arr.length - 1);
    }
    for(let i = 0; i < arr.length - 1; i++) {
        if (arr[i] < curHeight.value && arr[i + 1] > curHeight.value) {
            activateHeader(i);
            return (activeIndex.value = i);
        }
    }
}

const activateHeader = (index) => {
    const listContainer = document.getElementById('header-container');
    const lists = listContainer.querySelectorAll('li');
    if(lists.length === 0) {
        return;
    }
    if(index === -1) {
        lists[activeIndex.value].classList.remove('active');
        return;
    }
    if (activeIndex.value !== index) {
        lists[activeIndex.value].classList.remove('active');
    };
    lists[index].classList.add('active');
    lists[index].scrollIntoView({ behavior: 'smooth', block: 'center' });
}

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

onUnmounted(()=>{
    editor.destroy();
})
</script>

<template>
<div style="border: 1px solid #ccc; height:100%" :class="editTheme" id="editor-container">
    <div class="editor-box">
        <ul id="header-container" :style="`width:${cateWidth}%`"></ul>
        <editor-content class="editor-content" id="editor" :editor="editorRef" :style="`width:${editorWidth}%`"/>
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
    height: 100%;
    background-color: var(--el-fill-color-blank);
    border-right: 1px solid #9b9b9b;
    overflow-y: auto;
    margin: 0;
    padding: 0 0 0 0px;
    transition: all 0.3s;
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

#header-container li.active {
    color: var(--el-color-primary);
}
</style>