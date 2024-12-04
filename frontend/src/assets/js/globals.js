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
        cover: {
            name: "",
            data: ""
        }
    },
    content: "",
    resources: [],
    toc: []
})

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

const autoSave = ref({
    isAutoSave: true,
    autoSaveTime: 60
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
    autoSave
}

export default{}