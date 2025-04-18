/* 
@Author: Hughie
@CreateTime: 2024-10-31
@LastEditors: Hughie
@LastEditTime: 2025-3-28
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
const isBold = ref(false);
const isItalic = ref(false);
const isStrike = ref(false);
const isTextAlign = ref(0b00) // 00: left, 01: center, 10: right, 11: justify
const cateWidth = ref(20);
const editorWidth = ref(80);

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
    aboutVisible: false,
    searchBarVisible: false,
    codeInsertVisible: false,
    linkInsertVisible: false,
    tourVsisible: false,
    saveConfigVisible: true,
    aboveTextVisible: false,
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
        cover: {
            name: "",
            data: ""
        },
        meta: {}
    },
    content: "",
    resources: [],
    toc: []
})

// pulic variables
const change = ref(false);
const title = ref("Untitle");

const currentSave = ref("");

const staticFiles = ref([]);

const fileSuffix = {
    "jpg": "image/jpeg",
    "png": "image/png",
    "gif": "image/gif",
    "svg": "image/svg",
    "tiff": "image/tiff",
    "bmp": "image/bmp",
    "webp": "image/webp",
    "ico": "image/x-icon",
    "mp3": "audio/x-ms-mp3",
    "ogg": "audio/ogg",
    "wav": "audio/x-ms-wav",
    "flac": "audio/x-ms-flac",
    "aac": "audio/x-ms-aac",
    "wma": "audio/x-ms-wma",
    "wmv": "audio/x-ms-wmv",
    "mp4": "audio/x-ms-mp4",
    "aiff": "audio/x-ms-aiff",
    "aif": "audio/x-ms-aif",
    "aifc": "audio/x-ms-aifc",
    "m4a": "audio/x-ms-m4a",
    "m4b": "audio/x-ms-m4b",
    "m4p": "audio/x-ms-m4p",
    "m4r": "audio/x-ms-m4r"
}

// setting
const autoSave = ref({
    isAutoSave: true,
    autoSaveTime: 60
})

let defaultLang = localStorage.getItem('lang');
let defaultTheme =localStorage.getItem('theme');
const generalSetting = reactive({
    language: defaultLang,
    theme: defaultTheme,
    windowSize: "normal",
    resPort: 7288
});
const windowSetting = reactive({
    GPU : true
})
const linuxSetting = reactive({
    GPUPolicy: "auto"
})
const epubLayout = reactive({
    textDirection: "",
    layout: "",
    flow: "",
    spread: "",
    orientation: "",
    proportions: ""
})

export {
    editorRef,
    valueHtml,
    headerVal,
    fontVal,
    fonts,
    fontSizeVal,
    language_list,
    editLang,
    editTheme,
    visio,
    cover,
    bookInfo,
    change,
    title,
    currentSave,
    staticFiles,
    fileSuffix,
    autoSave,
    generalSetting,
    windowSetting,
    linuxSetting,
    isBold,
    isItalic,
    isStrike,
    epubLayout,
    cateWidth,
    editorWidth,
    isTextAlign
}

export default{}