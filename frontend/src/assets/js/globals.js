/* 
@Author: Hughie
@CreateTime: 2024-10-31
@LastEditors: Hughie
@LastEditTime: 2024-11-1
@Description: This is the public variables and constant of the program.
*/

import lang from './lang.js'
import { reactive, ref, shallowRef } from "vue";

// editor
const editorRef = shallowRef();
const valueHtml = ref('');
const headerVal = ref(0);
const fontVal = ref("Arial");
const fonts = ref([])
const fontSizeVal = ref("16px")
const imageInfo = reactive({
    zoom: 100,
    postition: "left",
    elem: ""
})

// i18n
const language_list = ref(lang);

// setting
var storgeLang = localStorage.getItem("lang");
const editLang = ref(storgeLang);
const editTheme = ref(localStorage.getItem("theme"));

// dialog
const visio = reactive({
    bookInfoVisible: false,
    settingVisible: false,
    mediaVisible: false,
    tableInsertVisible: false,
    helpVisible: false,
    searchBarVisible: false,
})
const cover = reactive({
    isExist: false,
    data: ''
})

// program state
const bookInfo = reactive({
    metadata: {
        title: "Untitle",
        creator: "",
        identifier: "",
        language: "",
        contributors: "",
        description: "",
        publisher: "",
        subject: [],
        date: "",
        meta: []
    },
    content: "",
    resources: [],
    toc: []
})

const change = ref(false);
const title = ref("Untitle");

const currentSave = ref("");

const staticFiles = ref([]);

export {
    editorRef,
    valueHtml,
    headerVal,
    fontVal,
    fonts,
    fontSizeVal,
    imageInfo,
    language_list,
    editLang,
    editTheme,
    visio,
    cover,
    bookInfo,
    change,
    title,
    currentSave,
    staticFiles
}

export default{}