/* 
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-12-20
@Description: This is the public methods for the whole program.
*/

import {Base64Decode, GetStaticResources, DirectLoading, FileSave, GetImageData, NewFile } from "../../../wailsjs/go/core/App.js"
import { ElMessage, ElLoading } from 'element-plus'
import {
    editorRef,
    cover,
    bookInfo,
    change,
    staticFiles,
    title,
    currentSave,
    autoSave,
    fileSuffix,
    epubLayout
} from "./globals.js"


/*
public methods
*/

function rgbaToHex(rgba) {
    const [r, g, b, a] = rgba.match(/\d+(\.\d+)?/g).map(Number);
    const hexR = r.toString(16).padStart(2, '0');
    const hexG = g.toString(16).padStart(2, '0');
    const hexB = b.toString(16).padStart(2, '0');
    const hexA = Math.round(a * 255).toString(16).padStart(2, '0');
    return `#${hexR}${hexG}${hexB}${hexA}`;
}

const arrayRemove = (array, id) => {
    for(let i in array) {
        if(array[i].id == id){
            array.splice(i, 1);
            break;
        }
    }
}

const deepClone = (obj) => {
    return JSON.parse(JSON.stringify(obj));
}

const arrayEquel = (arr1, arr2) => {
    if(arr1.length === arr2.length && arr1.every((value, index) => value === arr2[index])){
        return true;
    }
    return false;
}

const normalizeHTML = (html) => {
    // change the tag name to lowercase
    html = html.replace(/<\/?([A-Z][A-Z0-9]*)\b[^>]*>/gi, function (match) {
        return match.toLowerCase();
    });

    // makesure the self-closing tags are closed properly
    html = html.replace(/<([a-z]+)([^>]*)\/?>/g, function (match, tagName, attributes) {
        const selfClosingTags = ['area', 'base', 'br', 'col', 'embed', 'hr', 'img', 'input', 'link', 'meta', 'param', 'source', 'track', 'wbr'];
        if (selfClosingTags.includes(tagName)) {
            return `<${tagName}${attributes}></${tagName}>`;
        } else {
            return match;
        }
    });

    return html;
}

const resetState = async () => {
    let ok = await NewFile().then((res)=> {
        if(res.Code == 1) {
            ElMessage.error(t("message.error"), res.Msg);
            return false;
        }else {
            return true;
        }
    })
    if(!ok) return;
    getImageFiles();
    bookInfo.metadata = {
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
    }
    bookInfo.content = "";
    bookInfo.resources = [];
    bookInfo.toc = [];
    initCover();
    const editor = editorRef.value;
    editor.chain().setContent("", true).run();
    change.value = false;
    title.value = "Untitle";
    currentSave.value = "";
}

const getImageFiles = async () => {
    const _ = await GetStaticResources().then((res)=>{
        if(res.Code == 0) {
            if(res.Data == '') {
                staticFiles.value = [];
                return;
            }
            let data = JSON.parse(res.Data);
            staticFiles.value = data.FileList;
        }else {
            ElMessage.error(t("message.error"));
        }
    })
}

// loading the cover data to the frontend 
const initCover = ()=>{
    if(bookInfo.metadata.cover.name == ""){
        cover.isExist = false;
        cover.data = "";
        return
    }
    bookInfo.resources.push({
        id: "cover",
        name: "cover",
        type: "image/jpeg",
        data: bookInfo.metadata.cover.data
    })
    cover.isExist = true;
    cover.data = bookInfo.metadata.cover.data;
}

// Check if open the epub suffix file directly and load the data
const checkIfOpenFileDirectly = async (t) => {
    let loading = ElLoading.service({
        lock: true,
        text: t("message.loadingMessage"),
        background: 'rgba(0, 0, 0, 0.7)',
    })
    let message = await DirectLoading().then((res)=> {
        return res;
    })
    if(message.Code == -1){
        loading.close();
        return;
    }
    if (message.Code == 1) {
        loading.close();
        ElMessage.error(t("message.openError") + message.Msg);
        return;
    }
    let rawData = JSON.parse(message.Data);
    if(rawData != null) {
        rawData.metadata.creator = rawData.metadata.creator.join(',');

        rawData.metadata.contributors = rawData.metadata.contributors.join(',');

        rawData.content = await Base64Decode(rawData.content).then((res)=> {
            return JSON.parse(res);
        });

        await getImageFiles();

        const E = editorRef.value;

        bookInfo.metadata = rawData.metadata;
        E.commands.setContent(rawData.content, false);

        updateCatalog();

        initCover();
        loading.close();
        change.value = false;
        currentSave.value = message.Msg;
        title.value = message.Msg;
    }
    return
}

const updateCatalog = () => {
    const headerContainer = document.getElementById('header-container');
    const editorContainer = document.getElementById('editor');
    const headers = editorContainer.querySelectorAll('h1, h2, h3, h4, h5, h6')
    headerContainer.innerHTML = [...headers].map((header)=>{
        let id = header.id;
        let type = header.tagName;
        let text = header.innerText;
        return `<li data-id="${id}" type="${type}">${text}</li>`;
    }).join('');
}

// Based on the headers array, generate the toc structure
class TocGenerator {
    constructor(headers) {
        this.count = 1;
        this.weight = {
            "header1": 6,
            "header2": 5,
            "header3": 4,
            "header4": 3,
            "header5": 2,
            "header6": 1
        }
        this.generateTempData(headers);
    }

    generateID(){
        let idPrefix = 'guide_signal_';
        return idPrefix + this.count++;
    }

    generateTempData(headers) { // make it more easy to deal with
        let temp = []
        headers.forEach((item) => {
            temp.push({
                weight: this.weight[item.type],
                type: item.type,
                text: item.text
            });
        })
        this.headers = temp;
    }

    process(){
        let that = this;
        let result = [];
        let segment = [];
        let maxVal = Math.max.apply(null, this.headers.map((item) => item.weight)); // the biggest weight is the root node
        let tops = [0];
        this.headers.forEach((item, i) => {
            if(item.weight == maxVal) {
                tops.push(i);
            }
        })
        tops.push(this.headers.length); // the range of slice
        tops = Array.from(new Set(tops));

        for(var i = 0; i < tops.length-1; i++) {
            segment.push(this.headers.slice(tops[i], tops[i+1]));
        }

        function formTree(seg) {
            // the first of each segment is the root node
            let newItem = {
                id: seg[0].type,
                label: seg[0].text,
                href: that.generateID(),
                children: []
            };
            if(seg.length == 1){
                return newItem;
            }else if(seg.length == 2){
                // the second one definitely is the child of the root node
                newItem.children.push({
                    id: seg[1].type,
                    label: seg[1].text,
                    href: that.generateID(),
                    children: []
                });
                return newItem;
            }else {
                newItem.children.push({
                    id: seg[1].type,
                    label: seg[1].text,
                    href: that.generateID(),
                    children: []
                });
            }
            // restore the current node
            let temp = newItem.children[0];
            for(var i = 2; i < seg.length; i++) {
                // compare to the next node with the cureent, if smaller, then is the child of the current
                if(seg[i] < temp.value){
                    temp.children.push({
                        id: seg[i].type,
                        label: seg[i].text,
                        href: that.generateID(),
                        children: []
                    });
                    temp = temp.children[temp.children.indexOf(temp) + 1];
                    continue;
                }
                // else, the next node is the sibling of the current
                newItem.children.push({
                    id: seg[i].type,
                    label: seg[i].text,
                    href: that.generateID(),
                    children: []
                });
                temp = newItem.children[newItem.children.indexOf(temp) + 1]; // whatever it goes, the current node should be changed to the next node
            }
            return newItem;
        }

        segment.forEach(v => {
            result.push(formTree(v));
        })
        return result;
    }

}

const setImage = async (book) => {
    if (staticFiles.value === null || staticFiles.value.length === 0) {
        return;
    }
    await Promise.all(staticFiles.value.map(async (v) => {
        let data = await GetImageData(v);
        if (data.Code == 1) {
            ElMessage.error(t('message.exportError') + ": " + data.Msg);
            return;
        }
        let [name,suffix] = v.split('\\')[2].split(".");
        book.resources.push({
            id: name,
            name: name,
            data: data.Data,
            type: fileSuffix[suffix]
        });
    }));
}

// Timer context for auto save
class TimerContext {
    constructor(t) {
        // singleton
        if (TimerContext.instance) {
            return TimerContext.instance;
        }
        this.timer = null;
        this.time = autoSave.value.autoSaveTime * 1000;
        this.t = t;
        this.state = false; // false: not running, true: running
        this.currentTime = null;
        TimerContext.instance = this;
    }
    // save the data
    async save() {
        if (autoSave.value.isAutoSave == false && currentSave.value == "") return;

        let tempData = JSON.parse(JSON.stringify(bookInfo));
        if(editorRef.value == null) {
            return;
        }
        const editor = editorRef.value;
        tempData.content = editor.getHTML();

        await setImage(tempData);

        const headers = [];
        editor.$nodes('custom-heading').forEach(h => {
            headers.push({
                type: "header" + h.attributes.level,
                text: h.textContent
            });
        });

        const Toc = new TocGenerator(headers);
        const res = Toc.process();

        tempData.toc = res;
        tempData.metadata.creator = tempData.metadata.creator.split(',');
        tempData.metadata.contributors = tempData.metadata.contributors.split(',');

        const doc = new DOMParser().parseFromString(editor.getHTML(), 'text/html')
        let titles = doc.querySelectorAll("h1, h2, h3, h4, h5, h6")
        titles.forEach((e, i) => {
            e.id = "guide_signal_" + (i+1);
        })

        tempData.content = normalizeHTML(doc.body.innerHTML);
        tempData.content = tempData.content.replaceAll("<p></p>", "<br></br>");
        tempData.content = tempData.content.replaceAll(" ", "");
        tempData.content = tempData.content.replaceAll('"', "'");
        tempData.metadata.meta = JSON.parse(JSON.stringify(epubLayout));

        FileSave(currentSave.value, (() => {
            return JSON.stringify(tempData);
        })(),
            true).then((res) => {
                if (res.Code == 1) {
                    ElMessage.error(this.t('message.saveError') + ": " + res.Msg);
                    return;
                } else if (res.Code == 0) {
                    currentSave.value = res.Data;
                    title.value = res.Data;
                    change.value = false;
                    ElMessage.success(this.t('message.saveSuccess'));
                    return;
                }
                return;
            });
    }

    Start() {
        if (autoSave.value.isAutoSave == false || currentSave.value == "") {
            this.stop();
            return;
        }
        this.state = true;
        this.currentTime = new Date().getTime();
        this.timer = setInterval(() => {
            this.save();
            this.currentTime = new Date().getTime();
        }, this.time);
    }

    Reset() {
        this.stop();
        this.time = autoSave.value.autoSaveTime * 1000;
        if (autoSave.value.isAutoSave == false || currentSave.value == "") return;
        this.Start();
    }

    stop() {
        if(this.timer == null) return;
        this.state = false;
        clearInterval(this.timer);
        this.timer = null;
        this.currentTime = null;
    }

    static getInstance(t) {
        if (!TimerContext.instance) {
            TimerContext.instance = new TimerContext(t);
        }
        return TimerContext.instance;
    }

    State() {
        let currentTime = new Date().getTime();
        console.log(`Timer state: ${this.state == true ? "running" : "stop"}
Time each round: ${this.time/1000} seconds`);
        if(this.currentTime != null) {
            console.log(`Have run: ${(currentTime - this.currentTime) / 1000} seconds
Time left: ${(this.time - (currentTime - this.currentTime)) / 1000} seconds`);
        }
    }
}

export {
    TocGenerator,
    checkIfOpenFileDirectly,
    initCover,
    arrayRemove,
    resetState,
    rgbaToHex,
    arrayEquel,
    getImageFiles,
    updateCatalog,
    TimerContext,
    setImage,
    normalizeHTML,
    deepClone
}

export default {
    
}