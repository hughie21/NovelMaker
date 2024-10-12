<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-9
@LastEditors: Hughie
@LastEditTime: 2024-10-10
@Description: This is the dialog allow user to look up or replace.
*/

import { ref } from 'vue'
import { lookupSession,searchKey, replaceKey, resultCount, currentPointer, searchOption} from "../../assets/js/lookup"
import { useI18n } from 'vue-i18n';
import { visio, editorRef } from '../../assets/js/utils';

const { t } = useI18n();
const switchRef = ref(false);

const handleLookup = (searchTerm) => {
    if (searchTerm == "") {
        lookupSession.destory()
        resultCount.value = 0
        return
    }
    const E = editorRef.value;
    lookupSession.lookup(E, searchTerm)
}

const handleNextResult = () => {
    const cuurentIndex = lookupSession.moveNext()
    currentPointer.value = cuurentIndex + 1;
}

const handlePrevResult = () => {
    const currentIndex = lookupSession.movePrev();
    currentPointer.value = currentIndex + 1;
}

const handleLookupClose = () => {
    const E = editorRef.value;
    visio.searchBarVisible = false;
    lookupSession.destory()
    E.setEditable(true);
    E.commands.focus();
    searchKey.value = "";
    replaceKey.value = "";
    resultCount.value = 0
}

const handleReplace = () => {
    if(replaceKey.value == "") {
        return
    }
    const E = editorRef.value;
    lookupSession.replace(E, replaceKey.value)
}

const handleReplaceAll = () => {
    if(replaceKey.value == "") {
        return
    }
    const E = editorRef.value;
    while(true) {
        if(lookupSession.replace(E, replaceKey.value)) {
            break;
        }
    }
}

const handleOptionChange = (option) => {
    searchOption[option] = !searchOption[option];
    lookupSession.lookup(editorRef.value, searchKey.value);
}

setTimeout(()=>{
    document.getElementById('lookupBar').addEventListener('keydown', (e)=>{
        if(e.key == "Enter") {
            handleNextResult();
        }
    })
}, 1000)

</script>

<template>
<div class="search-bar-mask" v-show="visio.searchBarVisible">
    <div id="btn-replace" style="display: none;" @click="visio.searchBarVisible = true; switchRef = !switchRef"></div> <!-- Shortcut keys -->
    <div class="search-bar">
        <div style="width: 5%;">
            <button class="search-bar-button" @click="switchRef = !switchRef">
                <i class="el-icon" style="width: 1.4em" v-if="!switchRef">
                    <svg t="1727513826212" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4251" width="200" height="200"><path d="M728.223744 520.22784a42.467328 42.467328 0 0 1-11.393024 20.503552L374.90688 882.65728c-16.662528 16.662528-43.677696 16.662528-60.340224 0s-16.662528-43.677696 0-60.340224L626.449408 510.43328 314.614784 198.598656c-16.662528-16.662528-16.662528-43.677696 0-60.340224 16.661504-16.662528 43.676672-16.662528 60.3392 0L716.879872 480.18432c10.860544 10.860544 14.642176 26.120192 11.343872 40.04352z" p-id="4252"></path></svg>
                </i>
                <i class="el-icon" style="width: 1.4em" v-else>
                    <svg t="1727514408248" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="3679" width="200" height="200"><path d="M505.952256 751.028224a42.467328 42.467328 0 0 1-20.503552-11.393024L143.52384 397.709312c-16.662528-16.661504-16.662528-43.676672 0-60.3392 16.661504-16.662528 43.676672-16.662528 60.3392 0L515.74784 649.253888 827.582464 337.41824c16.661504-16.662528 43.676672-16.662528 60.3392 0s16.662528 43.677696 0 60.340224L545.9968 739.683328c-10.861568 10.861568-26.120192 14.6432-40.044544 11.34592z" p-id="3680"></path></svg>
                </i>
            </button>
            
        </div>
        <div style="width: 94%; display: flex; flex-direction: column; justify-content: space-between">
            <div class="search-bar-input">
                <el-input style="width: 180px" id="lookupBar" size="small" v-model="searchKey" @input="handleLookup" :placeholder="t('dialog.search.lookUp')" />
                <div class="bar-inset">
                    <el-tooltip :content="t('dialog.search.caseSensitive')"  placement="bottom" effect='dark'>
                        <button :class="'search-bar-button ' + (searchOption.caseSensitive ? 'search-bar-button-active' : '')" @click="searchOption.caseSensitive = !searchOption.caseSensitive; lookupSession.lookup(editorRef ,searchKey)">
                            <i class="el-icon" style="width: 1.4em">
                                <svg t="1727515146209" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4756" width="200" height="200"><path d="M240.8704 250.31168H338.9952l217.6 552.41216H464.01024l-53.01248-142.35648H168.07424l-53.01248 142.35648H23.2704l217.6-552.41216z m-46.68416 340.42368h190.69952l-93.37344-249.9072h-3.16416l-94.1568 249.9072zM813.66016 391.89504c60.13952 0 104.448 14.69952 132.93568 44.10368 24.5248 25.52832 37.18656 62.6688 37.18656 111.4112v255.31392H906.24v-56.4736a169.58976 169.58976 0 0 1-59.3408 47.96416c-26.9056 12.38016-58.55744 19.34336-94.95552 19.34336-42.73152 0-75.96544-10.83392-99.70176-31.72352-26.112-20.8896-38.77376-47.96928-38.77376-81.23392 0-44.8768 18.2016-79.6928 54.59968-103.68 33.2288-23.20896 80.70656-34.816 140.84608-36.36224l91.78624-2.31936v-16.24576c0-55.7056-30.85824-83.5584-92.57984-83.5584-26.112 0-47.47264 4.63872-63.29856 13.9264-18.9952 10.8288-30.86336 27.07968-35.6096 49.5104l-83.08224-6.9632c8.704-43.32032 30.85824-75.81696 65.67424-95.93344 30.06976-18.57024 71.2192-27.07968 121.856-27.07968z m87.04 225.92l-86.2464 2.31936c-76.75392 1.54624-114.7392 27.8528-114.7392 77.36832 0 15.47264 6.33344 27.8528 19.78368 37.9136 12.66176 10.05568 30.06976 15.47264 51.43552 15.47264 35.60448 0 65.67424-10.83392 90.99264-31.72352 25.32352-20.8896 38.77376-47.19616 38.77376-78.14144v-23.20896z" p-id="4757"></path></svg>
                            </i>
                        </button>
                    </el-tooltip>
                    <el-tooltip :content="t('dialog.search.wholeWord')" placement="bottom" effect='dark'>
                        <button :class="'search-bar-button ' + (searchOption.wholeWord ? 'search-bar-button-active' : '')" @click="searchOption.wholeWord = !searchOption.wholeWord; lookupSession.lookup(editorRef, searchKey)">
                            <i class="el-icon" style="width: 1.4em">
                                <svg t="1727534219133" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5408" width="200" height="200"><path d="M870.4 814.08H153.6a128 128 0 0 1-128-128V500.736a25.6 25.6 0 0 1 51.2 0v185.344A76.8 76.8 0 0 0 153.6 762.88h716.8a76.8 76.8 0 0 0 76.8-76.8V500.736a25.6 25.6 0 0 1 51.2 0v185.344a128 128 0 0 1-128 128z" p-id="5409"></path><path d="M466.944 625.152H409.6v-36.352a102.4 102.4 0 0 1-91.648 43.52 127.488 127.488 0 0 1-102.4-45.056 166.912 166.912 0 0 1-36.864-111.616 175.616 175.616 0 0 1 40.448-120.32 125.952 125.952 0 0 1 98.816-44.544A128.512 128.512 0 0 1 409.6 339.456v-22.528h57.344z m-196.608-68.608a70.656 70.656 0 0 0 58.368 34.304 81.408 81.408 0 0 0 57.856-22.528 51.2 51.2 0 0 0 18.944-37.888V406.528a51.2 51.2 0 0 0-19.968-37.888 81.408 81.408 0 0 0-55.296-16.896 71.168 71.168 0 0 0-58.88 34.816 144.384 144.384 0 0 0-24.576 84.48 148.48 148.48 0 0 0 23.552 85.504zM614.4 348.672a128.512 128.512 0 0 1 97.28-33.792 118.272 118.272 0 0 1 90.624 43.52 166.4 166.4 0 0 1 40.448 116.736 168.448 168.448 0 0 1-40.96 119.296 120.32 120.32 0 0 1-91.136 39.424 153.6 153.6 0 0 1-102.4-32.768v27.136h-61.44V209.92H614.4v138.752z m0 70.656a422.4 422.4 0 0 0 0 113.664 67.072 67.072 0 0 0 17.408 36.352 85.504 85.504 0 0 0 64.512 23.04 71.68 71.68 0 0 0 59.392-35.328 145.92 145.92 0 0 0 22.528-82.944 132.608 132.608 0 0 0-24.576-83.968A74.752 74.752 0 0 0 690.688 358.4a77.824 77.824 0 0 0-58.88 23.04 69.632 69.632 0 0 0-17.408 37.888z" p-id="5410"></path></svg>
                            </i>
                        </button>
                    </el-tooltip>
                </div>
                <div class="search-result">
                    <span v-if="resultCount == 0">{{t('dialog.search.defaultResult')}}</span>
                    <span v-else>{{`${currentPointer}/${resultCount}`}}</span>
                </div>
                <div class="search-options">
                    <el-tooltip :content="t('dialog.search.nextResult')" placement="bottom" effect='dark'>
                        <button class="search-bar-button" @click="handleNextResult">
                            <i class="el-icon" style="width: 1.4em">
                                <svg t="1727535059856" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="2928" width="200" height="200"><path d="M528.252661 950.942077l295.703769-289.705148c9.624195-9.380648 9.624195-24.707728 0-34.138518-9.576099-9.381671-25.222451-9.381671-34.848693 0L535.495622 875.592972 535.495622 88.901773c0-13.3214-11.048637-24.142863-24.642237-24.142863-13.597693 0-24.645306 10.821463-24.645306 24.142863l0 786.690176L232.574475 627.098411c-9.603729-9.381671-25.248034-9.381671-34.847669 0-4.836145 4.76349-7.194866 10.939143-7.194866 17.116843 0 6.1777 2.405793 12.355399 7.194866 17.116843l295.723212 289.704125c9.602705 9.383718 25.244964 9.383718 34.850739 0L528.252661 950.942077z" p-id="2929"></path></svg>
                            </i>
                        </button>
                    </el-tooltip>
                    <el-tooltip :content="t('dialog.search.previousResult')" placement="bottom" effect='dark'>
                        <button class="search-bar-button" @click="handlePrevResult">
                            <i class="el-icon" style="width: 1.4em">
                                <svg t="1727535346660" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="3086" width="200" height="200"><path d="M493.45411 71.89033 197.751365 361.596502c-9.624195 9.380648-9.624195 24.707728 0 34.138518 9.576099 9.381671 25.222451 9.381671 34.848693 0l253.612115-248.494561 0 786.690176c0 13.3214 11.047614 24.142863 24.642237 24.142863 13.597693 0 24.645306-10.821463 24.645306-24.142863L535.499715 147.240459l253.632581 248.494561c9.603729 9.381671 25.248034 9.381671 34.847669 0 4.836145-4.76349 7.194866-10.939143 7.194866-17.116843 0-6.1777-2.405793-12.355399-7.194866-17.116843L528.255731 71.797209c-9.602705-9.383718-25.244964-9.383718-34.849716 0L493.45411 71.89033z" p-id="3087"></path></svg>
                            </i>
                        </button>
                    </el-tooltip>
                    <el-tooltip :content="t('dialog.search.close')" placement="bottom" effect='dark'>
                        <button @click="handleLookupClose" class="search-bar-button">
                            <i class="el-icon" style="width: 1.4em">
                                <svg t="1727535459487" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4636" width="200" height="200"><path d="M557.312 513.248l265.28-263.904c12.544-12.48 12.608-32.704 0.128-45.248-12.512-12.576-32.704-12.608-45.248-0.128l-265.344 263.936-263.04-263.84C236.64 191.584 216.384 191.52 203.84 204 191.328 216.48 191.296 236.736 203.776 249.28l262.976 263.776L201.6 776.8c-12.544 12.48-12.608 32.704-0.128 45.248 6.24 6.272 14.464 9.44 22.688 9.44 8.16 0 16.32-3.104 22.56-9.312l265.216-263.808 265.44 266.24c6.24 6.272 14.432 9.408 22.656 9.408 8.192 0 16.352-3.136 22.592-9.344 12.512-12.48 12.544-32.704 0.064-45.248L557.312 513.248z" p-id="4637"></path></svg>
                            </i>
                        </button>
                    </el-tooltip>
                </div>
            </div>
            <div class="replace-bar-input" v-show="switchRef">
                <el-input style="width: 180px" size="small" v-model="replaceKey" :placeholder="t('dialog.search.replace')" />
                <div class="replace-options">
                    <el-tooltip :content="t('dialog.search.replace')" placement="bottom" effect='dark'>
                        <button class="search-bar-button" @click="handleReplace">
                            <i class="el-icon" style="width: 1.4em; height: 1.4em">
                                <svg t="1727535959856" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="9153" width="200" height="200"><path d="M552.9088 430.9504L594.8928 102.4h51.4048l-11.6224 91.1872h61.184c30.976 0 54.272 6.8096 69.7344 20.48 15.4624 13.6192 23.2448 33.792 23.2448 60.4672 0 6.8096-0.512 13.9264-1.3824 21.2992l-3.9936 32.4096c-8.6528 68.4544-48.2816 102.7072-118.8864 102.7072H552.96z m117.504-45.568c35.7376 0 56.0128-18.6368 60.7744-55.8592l4.4544-34.9696c0.6144-5.632 0.8704-9.6256 0.8704-11.9296 0-28.928-15.4624-43.4176-46.4384-43.4176h-61.2352l-18.7392 146.1248h60.3136z m-142.08 253.2864c-18.4832 0-32.6144 4.6592-42.496 14.0288-9.5232 9.4208-15.4624 23.3472-17.8176 41.7792l-4.5056 36.6592a133.9392 133.9392 0 0 0-0.8704 12.3392c0 28.672 15.36 43.008 46.0288 43.008h17.408c13.4144 0 35.8912-1.2288 67.4816-3.7888l-4.9152 45.568c-26.2144 2.56-48.6912 3.84-67.4816 3.84h-18.3296c-61.6448 0-92.4672-26.6752-92.4672-80.0768 0-7.168 0.4096-14.4896 1.3312-22.1696l4.4544-34.0992c4.4544-34.0992 16.384-59.648 35.7376-76.6976 19.4048-17.3056 46.7968-26.0096 82.2272-26.0096h16.5376l62.5664 0.8704-5.8368 45.1584-79.0528-0.4096zM344.4224 495.4624C267.3152 495.4624 204.8 555.1104 204.8 628.6336V788.48C204.8 862.0032 267.3152 921.6 344.4224 921.6h335.1552C756.6848 921.6 819.2 862.0032 819.2 788.48v-159.8464c0-73.5232-62.5152-133.12-139.6224-133.12H344.4224z m-83.7632 133.1712c0-44.1344 37.5296-79.872 83.7632-79.872h335.1552c46.2848 0 83.7632 35.7376 83.7632 79.872V788.48c0 44.1344-37.4784 79.872-83.7632 79.872H344.4224c-46.2336 0-83.7632-35.7376-83.7632-79.872v-159.744z m83.7632-297.472V255.744c0-14.6944 12.4928-26.624 27.9552-26.624h111.7184a27.2896 27.2896 0 0 0 27.904-26.624 27.2896 27.2896 0 0 0-27.904-26.624H372.3776c-46.2848 0-83.8144 35.7376-83.8144 79.872v75.3152l-25.4976-24.32a28.8768 28.8768 0 0 0-39.5264 0 25.7536 25.7536 0 0 0 0 37.632l73.216 69.8368a28.8768 28.8768 0 0 0 39.4752 0l73.216-69.8368a25.7536 25.7536 0 0 0 0-37.632 28.8768 28.8768 0 0 0-39.4752 0l-25.5488 24.32z" p-id="9154"></path></svg>
                            </i>
                        </button>
                    </el-tooltip>
                    <el-tooltip :content="t('dialog.search.replaceAll')" placement="bottom" effect='dark'>
                        <button class="search-bar-button" @click="handleReplaceAll">
                            <i class="el-icon" style="width: 1.4em; height: 1.4em; margin-left:5px">
                                <svg t="1727535999726" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="9313" width="200" height="200"><path d="M673.792 107.008l-28.8768 236.8512h76.8c48.5376 0 75.776-24.6784 81.7152-74.0352l2.7648-23.3472a134.144 134.144 0 0 0 0.9216-15.36c0-19.2512-5.3248-33.792-15.9744-43.6224-10.6496-9.8304-26.624-14.7456-47.9232-14.7456h-42.0864l7.9872-65.7408h-35.328z m93.696 163.7376c-3.2768 26.8288-17.2032 40.2432-41.7792 40.2432h-41.472l12.9024-105.3696h42.0864c21.2992 0 31.9488 10.4448 31.9488 31.3344 0 1.6384-0.2048 4.5056-0.6144 8.6016l-3.072 25.1904z m-286.72 60.2112c9.216 8.6016 21.504 12.9024 36.864 12.9024h89.088l12.5952-106.2912c0.4096-4.5056 0.6144-7.68 0.6144-9.5232 0-16.7936-5.4272-30.208-16.2816-40.2432-10.8544-10.0352-26.8288-15.0528-47.9232-15.0528h-8.2944c-6.144 0-13.4144 0.2048-21.8112 0.6144l-17.2032 0.9216-8.9088 0.6144-4.608 31.3344c14.5408-1.024 30.9248-1.536 49.152-1.536h7.9872c11.264 0 19.5584 2.4576 24.8832 7.3728a25.088 25.088 0 0 1 7.9872 19.3536 23.1936 23.1936 0 0 1-0.3072 4.3008l-0.3072 2.4576h-51.6096c-14.5408 0-26.3168 1.7408-35.328 5.2224a39.168 39.168 0 0 0-20.2752 16.2816 76.4416 76.4416 0 0 0-10.1376 37.1712c0 13.9264 4.608 25.2928 13.824 34.0992z m26.4192-23.6544a18.688 18.688 0 0 1-4.9152-13.824l0.3072-3.6864c1.024-7.9872 3.7888-13.5168 8.2944-16.5888a30.72 30.72 0 0 1 18.7392-4.608l50.9952 0.3072-5.2224 43.008h-53.4528a21.5552 21.5552 0 0 1-14.7456-4.608zM311.808 791.7568c9.216 8.6016 21.504 12.9024 36.864 12.9024h89.088l12.5952-106.2912c0.4096-4.5056 0.6144-7.68 0.6144-9.5232 0-16.7936-5.4272-30.208-16.2816-40.2432-10.8544-10.0352-26.8288-15.0528-47.9232-15.0528h-8.2944c-6.144 0-13.4144 0.2048-21.8112 0.6144l-17.2032 0.9216-8.9088 0.6144-4.608 31.3344c14.5408-1.024 30.9248-1.536 49.152-1.536h7.9872c11.264 0 19.5584 2.4576 24.8832 7.3728a25.088 25.088 0 0 1 7.9872 19.3536 23.1936 23.1936 0 0 1-0.3072 4.3008l-0.3072 2.4576H363.7248c-14.5408 0-26.3168 1.7408-35.328 5.2224a39.168 39.168 0 0 0-20.2752 16.2816 76.4416 76.4416 0 0 0-10.1376 37.1712c0 13.9264 4.608 25.2928 13.824 34.0992z m26.4192-23.6544a18.688 18.688 0 0 1-4.9152-13.824l0.3072-3.6864c1.024-7.9872 3.7888-13.5168 8.2944-16.5888a30.72 30.72 0 0 1 18.7392-4.608l50.9952 0.3072-5.2224 43.008H352.9728a21.5552 21.5552 0 0 1-14.7456-4.608z m191.488-91.5456c6.7584-6.7584 16.4864-10.1376 29.184-10.1376l54.3744 0.3072 3.9936-32.5632-43.008-0.6144h-11.3664c-24.3712 0-43.2128 6.2464-56.5248 18.7392-13.312 12.288-21.504 30.72-24.576 55.296l-3.072 24.576a144.7424 144.7424 0 0 0-0.9216 15.9744c0 38.5024 21.1968 57.7536 63.5904 57.7536h12.5952c12.9024 0 28.3648-0.9216 46.3872-2.7648l3.3792-32.8704c-21.7088 1.8432-37.1712 2.7648-46.3872 2.7648h-11.9808c-21.0944 0-31.6416-10.3424-31.6416-31.0272 0-1.8432 0.2048-4.8128 0.6144-8.9088l3.072-26.4192c1.6384-13.312 5.7344-23.3472 12.288-30.1056zM665.6 563.2a51.2 51.2 0 0 1 51.2 51.2v204.8a51.2 51.2 0 0 1-51.2 51.2H256a51.2 51.2 0 0 1-51.2-51.2v-204.8a51.2 51.2 0 0 1 51.2-51.2h409.6z m-409.6-51.2a102.4 102.4 0 0 0-102.4 102.4v204.8a102.4 102.4 0 0 0 102.4 102.4h409.6a102.4 102.4 0 0 0 102.4-102.4v-204.8a102.4 102.4 0 0 0-102.4-102.4H256z m179.2-102.4a25.6 25.6 0 0 0 0 51.2h307.2a76.8 76.8 0 0 1 76.8 76.8v153.6a25.6 25.6 0 0 0 51.2 0v-153.6A128 128 0 0 0 742.4 409.6h-307.2z m-220.16-30.8736L174.1824 337.92a15.36 15.36 0 1 0-21.7088 21.7088l67.072 67.072a15.36 15.36 0 0 0 21.7088 0l67.072-67.072a15.36 15.36 0 1 0-21.6576-21.7088l-40.9088 40.8576V281.6a35.84 35.84 0 0 1 35.84-35.84h102.4a15.36 15.36 0 0 0 0-30.72h-102.4a66.56 66.56 0 0 0-66.56 66.56v97.1264z" p-id="9314"></path></svg>
                            </i>
                        </button>
                    </el-tooltip> 
                </div>
            </div>
        </div>
    </div>
</div>
</template>

<style scoped>
.el-icon svg {
    width: 100% !important;
    height: 100% !important;
}
.search-bar-mask{
    position: fixed;
    top: 17%;
    right:10px;
    height: 55px;
    width: 400px;
    overflow: hidden;
    z-index: 999;
    display: flex;
    justify-content: center;
}
.search-options{
    display: flex;
    flex-direction: row;
    align-items: center;
    height: 25px;
    position: absolute;
    right: 5px;
}

.search-bar-button-active{
    background: #456ec4cb !important;
}

.replace-options{
    display: flex;
    flex-direction: row;
    align-items: center;
    height: 25px;
    margin-left: 10px;
}
.search-result{
    margin-left: -28px;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
}
.search-result span{
    color: #fff;
    font-size: 12px;
}
.bar-inset{
    display: flex;
    align-items: center;
    position: relative;
    right: 40px;
}
.search-bar {
    background-color: var(--el-bg-color-overlay);
    width: 380px;
    height: fit-content;
    border: 1px solid rgb(51, 51, 51);
    border-left: 3px solid rgb(136, 136, 136);
    border-radius:2px;
    box-shadow: 5px 5px 5px #3939392e;
    display: flex;
    flex-direction: row;
    justify-content: start;
    position: relative;
    transition: 1s;
}
.search-bar-input {
    display: flex;
    flex-direction: row;
    align-items: center;
    height: 25px;
    width: 100%;
}
.replace-bar-input {
    display: flex;
    align-items: center;
    height: 28px;
    width: 100%;
}
.search-bar-button {
    border-radius: 2px;
    background: none;
    padding: 0;
    border: none;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
}
.search-bar-button:hover {
    background: #4a4a4a8a;
}
</style>