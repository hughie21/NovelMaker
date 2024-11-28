/* 
@Author: Hughie
@CreateTime: 2024-7-5
@LastEditors: Hughie
@LastEditTime: 2024-11-1
@Description: This is the public methods for the whole program.
*/

import {Base64Decode, GetStaticResources} from "../../../wailsjs/go/core/App.js"

import { ElMessage } from 'element-plus'
import {
    editorRef,
    cover,
    bookInfo,
    change,
    staticFiles,
    title,
    currentSave
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

const arrayEquel = (arr1, arr2) => {
    if(arr1.length === arr2.length && arr1.every((value, index) => value === arr2[index])){
        return true;
    }
    return false;
}

const resetState = () => {
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
    const editor = editorRef.value;
    editor.chain().clearContent().run();
    change.value = false;
    title.value = "Untitle";
    currentSave.value = "";
}

const getImageFiles = async () => {
    const _ = await GetStaticResources().then((res)=>{
        if(res.Code == 0) {
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
const checkIfOpenFileDirectly = async () => { 
    // let message = await DirectLoading().then((res)=> {
    //     return res;
    // })
    // if(message.Code != 0){
    //     return;
    // }
    // let rawData = JSON.parse(message.Data);
    // if(rawData != null) {
    //     rawData.metadata.creator = rawData.metadata.creator.join(',');
    //     rawData.metadata.contributors = rawData.metadata.contributors.join(',');
    //     rawData.content = await Base64Decode(rawData.content).then((res)=> {
    //         return JSON.parse(res);
    //     });
    //     const E = editorRef.value;
    //     E.chain().setContent(rawData.content, true).run();
    //     E.chain().focus().insertContent().run();
    //     bookInfo.metadata = rawData.metadata;
    //     bookInfo.content = rawData.content;
    //     bookInfo.toc = rawData.toc;
    //     bookInfo.resources = rawData.resources;
    //     initCover();
    //     change.value = false;
    // }
    // return
}

// Based on the headers array, generate the toc structure
class TocGenerator {
    constructor(headers) {
        this.count = 1;
        this.weight = {
            "header1": 5,
            "header2": 4,
            "header3": 3,
            "header4": 2,
            "header5": 1
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

export {
    TocGenerator,
    checkIfOpenFileDirectly,
    initCover,
    arrayRemove,
    resetState,
    rgbaToHex,
    arrayEquel,
    getImageFiles,
}

export default {
    
}