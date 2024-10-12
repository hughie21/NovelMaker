/* 
@Author: Hughie
@CreateTime: 2024-10-5
@LastEditors: Hughie
@LastEditTime: 2024-10-10
@Description: The main implementation of the look up and replace function
*/

import { editorRef } from "./utils.js";
import { ref } from 'vue'

const searchKey = ref('');
const replaceKey = ref('');
const resultCount = ref(0);
const currentPointer = ref(0);
const showMask = ref(false);
const searchOption = ref({
    caseSensitive: false,
    wholeWord: false,
})

class lookupState {
    constructor() {
        this.init()
        this.E = editorRef.value;
    }
    injectEditor(editor) {
        this.E = editor; // make sure the editor is loaded properly
    }
    init() {
        this.searchPosition = [];
        this.currentIndex = 0;
    }
    at(index) {
        if(this.searchPosition.length == 0){
            return
        }
        this.E.commands.setSearchSelection(this.searchPosition[index]);
        const element = this.E.view.domAtPos(this.searchPosition[index].from).node;
        element.scrollIntoView({ behavior: 'instant' });
    }
    moveNext() {
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
    replace(editor, replaceTerm) {
        if(this.searchPosition.length == 0){
            return true
        }
        this.E = editor;
        const range = this.searchPosition[this.currentIndex];
        editor.chain().unsetSearchSelection(range).deleteRange(range).insertContentAt(range.from, replaceTerm, {
            updateSelection: false,
        }).run();
        const prevIndex = this.currentIndex > 0 ? this.currentIndex - 1 : 0;
        this.lookup(editor, this.searchKey, this.searchOptions);
        return false
    }
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
    showMask,
    resultCount,
    currentPointer
}

export default {}