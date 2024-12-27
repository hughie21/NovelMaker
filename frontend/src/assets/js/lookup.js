/* 
@Author: Hughie
@CreateTime: 2024-10-5
@LastEditors: Hughie
@LastEditTime: 2024-10-10
@Description: The main implementation of the look up and replace function
*/

import { editorRef } from "./globals.js";
import { ref } from 'vue'

// the keyword that user want to search
const searchKey = ref('');
// the keyword that user want to replace
const replaceKey = ref('');
// the total count of the search result
const resultCount = ref(0);
// where the current search result is
const currentPointer = ref(0);
// the search options
const searchOption = ref({
    caseSensitive: false,
    wholeWord: false,
})

// The main class of the lookup state
class lookupState {
    constructor() {
        this.init()
        this.E = editorRef.value;
    }
    injectEditor(editor) {
        this.E = editor; // make sure the editor is loaded properly
    }
    init() {
        // record all the search result in the editor
        this.searchPosition = [];
        // the index of the searchPosition array
        this.currentIndex = 0;
    }
    // set the current search result in the background color
    at(index) {
        if(this.searchPosition.length == 0){
            return
        }
        this.E.commands.setSearchSelection(this.searchPosition[index]);
        const element = this.E.view.domAtPos(this.searchPosition[index].from).node;
        element.scrollIntoView({ behavior: 'instant' });
    }
    // move the current search result to the next one
    moveNext() {
        // Enables users to loop the selection
        if(this.currentIndex == this.searchPosition.length - 1){
            this.E.commands.unsetSearchSelection(this.searchPosition[this.currentIndex]);
            this.currentIndex = 0;
            this.at(this.currentIndex)
            return this.currentIndex;
        }
        this.E.commands.unsetSearchSelection(this.searchPosition[this.currentIndex]);
        this.currentIndex++
        this.at(this.currentIndex)
        return this.currentIndex;
    }
    // move the current search result to the previous one
    movePrev() {
        if(this.currentIndex == 0){
            this.E.commands.unsetSearchSelection(this.searchPosition[this.currentIndex]);
            this.currentIndex = this.searchPosition.length - 1;
            this.at(this.currentIndex)
            return this.currentIndex;
        }
        this.E.commands.unsetSearchSelection(this.searchPosition[this.currentIndex]);
        this.currentIndex--
        this.at(this.currentIndex)
        return this.currentIndex;
    }
    // the key method that enables the search function
    lookup(editor, key) {
        this.E = editor;
        this.destory();
        const state = editor.state;
        const doc = state.doc;
        this.searchKey = key;
        const caseSensitive = searchOption.value.caseSensitive;
        const wholeWord = searchOption.value.wholeWord;
        let regex;
    
        if (caseSensitive && wholeWord) {
            regex = new RegExp("\\b" + key + "\\b", 'g');
        } else if (caseSensitive) {
            regex = new RegExp(key, 'g');
        } else if (wholeWord) {
            regex = new RegExp("\\b" + key + "\\b", 'ig');
        } else {
            regex = new RegExp(key, 'ig');
        }
    
        let array1;
        doc.descendants((node, pos) => {
            if (node.type.name == "text") {
                while ((array1 = regex.exec(node.text)) !== null) {
                    if (array1.length == 0) {
                        break;
                    }
                    var selectRange = {
                        from: pos + array1.index,
                        to: pos + array1.index + array1[0].length
                    }
                    this.searchPosition.push(selectRange);
                }
            }
        });
    
        // If color all the search result, the editor will be very slow and even break down
        // So I just set the each of it when it is selected
        this.at(0);
        resultCount.value = this.searchPosition.length;
        currentPointer.value = this.currentIndex + 1;
    }
    // the key method that enables the replace function
    replace(editor, replaceTerm) {
        if(this.searchPosition.length == 0){
            return
        }
        this.E = editor;
        const range = this.searchPosition[this.currentIndex];
        editor.chain().unsetSearchSelection(range).deleteRange(range).insertContentAt(range.from, replaceTerm, {
            updateSelection: false,
        }).run();

        // update the following search result
        for(var i = this.currentIndex+1; i < this.searchPosition.length; i++){
            this.searchPosition[i].from = this.searchPosition[i].from - (range.to - range.from) + replaceTerm.length;
            this.searchPosition[i].to = this.searchPosition[i].to - (range.to - range.from) + replaceTerm.length;
        }

        // remove the current search result
        this.searchPosition.splice(this.currentIndex, 1);
        if(this.currentIndex > this.searchPosition.length - 1){
            this.currentIndex = this.searchPosition.length - 1
        }
        // highlight the next search result
        this.at(this.currentIndex);
        resultCount.value = this.searchPosition.length;
    }
    // the key method that enables the replace all function
    replaceAll(editor, replaceTerm) {
        return new Promise((resolve, reject) => {
            if(this.searchPosition.length == 0){
                resolve();
                return
            }
            this.E = editor;
            const changes = [];
            for(var i = 0; i < this.searchPosition.length; i++){
                const range = this.searchPosition[i];
                changes.push({
                    from: range.from,
                    to: range.to
                })
                for(var j = i+1; j < this.searchPosition.length; j++){
                    this.searchPosition[j].from = this.searchPosition[j].from - (range.to - range.from) + replaceTerm.length;
                    this.searchPosition[j].to = this.searchPosition[j].to - (range.to - range.from) + replaceTerm.length;
                }
            }
            changes.forEach((change) => {
                editor.chain().unsetSearchSelection({from: change.from, to: change.to}).deleteRange({from: change.from, to: change.to}).insertContentAt(change.from, replaceTerm, {
                    updateSelection: false,
                }).run();
            })
            this.destory();
            resultCount.value = this.searchPosition.length;
            resolve();
        })
    }
    // reset the state of the lookup
    destory() {
        // clear all the search selection while close the dialog
        this.E.commands.unsetSearchSelection(this.searchPosition[this.currentIndex]);
        this.searchPosition = [];
        this.currentIndex = 0;
    }
}
const lookupSession = new lookupState();

export {
    lookupSession,
    searchKey,
    replaceKey,
    searchOption,
    resultCount,
    currentPointer
}

export default {}