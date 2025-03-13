<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-18
@LastEditors: Hughie
@LastEditTime: 2024-11-1
*/
import { useI18n } from 'vue-i18n';
import { FileOpen, FileSave, Base64Decode } from '../../../wailsjs/go/core/App.js'
import { ElMessage, ElMessageBox, ElLoading, ElNotification } from 'element-plus'
import { ref, reactive, inject, onMounted } from 'vue';
import { editorRef, change, visio, bookInfo, currentSave, title, generalSetting, epubLayout } from '../../assets/js/globals.js';
import { TocGenerator, initCover, resetState, getImageFiles, updateCatalog, TimerContext, setImage, normalizeHTML } from '../../assets/js/utils.js';
import { lookupSession, searchKey, replaceKey, resultCount } from '../../assets/js/lookup.js';
import "../../assets/css/tab.css"

const { t } = useI18n();
const disableBtn = reactive({
    file: false,
    edit: true,
})
const btnNormalClass = inject("btnNormalClass");
const btnDisabledClass = ref("el-button btn func_btn-big is-disabled");
const timer = new TimerContext(t);

onMounted(()=> {
    // start the timer which will be used to save the file automatically
    timer.Start();
    timer.State()
})
 
const openFilePicker = () => {
    async function innerOpner(){
        // open the file select dialog
        var res = await FileOpen().then((res) => {
            return res
        })
        if(res.Code == 1){ // failed to open the file
            ElMessage.error(t('message.openError'))
            return;
        }else if (res.Code == -1){ // cancel
            return;
        }

        let loading = ElLoading.service({
            lock: true,
            text: t('message.loadingMessage'),
            background: 'rgba(0, 0, 0, 0.7)',
        })

        try{
            // Escape some special characters into characters that can be parsed into JSON.
            res.Data.replaceAll('\n', '\\n');
            res.Data.replaceAll('\r', '\\r');
            res.Data.replaceAll('\t', '\\t');
            res.Data.replaceAll('"', "'");
            var rawData = JSON.parse(res.Data);
        }catch(e){
            loading.close();
            ElMessage.error(t('message.openError'))
            return;
        }
        rawData.metadata.creator = rawData.metadata.creator.join(',');

        rawData.metadata.contributors = rawData.metadata.contributors.join(',');

        rawData.content = await Base64Decode(rawData.content).then((res)=> {
            try{
                return JSON.parse(res);
            }catch{
                // replace the image urls
                // the image will store in the local server
                let imageUrls = /..\/Images/g
                res = res.replace(imageUrls, 'http://127.0.0.1:' + generalSetting.resPort)
                return res;
            }
        });

        // get the names of the images from the local server
        await getImageFiles();

        const E = editorRef.value;

        bookInfo.metadata = rawData.metadata;
        E.commands.setContent(rawData.content, true);
        
        // update the catalog
        updateCatalog();

        // load the cover image to the frontend
        initCover();
        loading.close();
        change.value = false;
        currentSave.value = res.Msg;
        title.value = res.Msg;

        // reset the timer
        timer.Reset();
    }
    function innerCallBack(action){
        if(action == 'confirm'){
            // delay the open file dialog in case both save file dialog
            // and the open file dialog are opened at the same time
            saveFilePicker(true).then(()=> {
                setTimeout(()=>{
                    innerOpner();
                }, 3000)
            })
        }else if(action == 'close'){
            return;
        }else if (action == 'cancel'){
            innerOpner();
            return;
        }
    }

    ElNotification({
        title: t('message.warning'),
        message: t('message.openInfo'),
        duration: 5000,
        type: 'warning',
    })

    // check if the file has been saved
    // if not, ask the user to save the file
    // if the user choose to save the file, save the file and then open the file
    if(change.value){
        ElMessageBox.confirm(t('message.saveWarning'),t('message.warning'), {
            confirmButtonText: t('message.confirm'),
            cancelButtonText: t('message.notSave'),
            distinguishCancelAndClose: true,
            type: 'warning',
            callback: innerCallBack
        })
    }else {
        innerOpner();
    }
    
}

// param saveAs: whether to save the file as a new file
const saveFilePicker = async (saveAs) => {
    // Deep copy the data
    let tempData = JSON.parse(JSON.stringify(bookInfo));
    const editor = editorRef.value;
    tempData.content = editor.getHTML();

    // load the image file's name and data to the data.resources
    await setImage(tempData);

    // Conversion of headings into special data structures for 
    // subsequent processing
    const headers = [] 
    editor.$nodes('custom-heading').forEach(h => {
        headers.push({
            type: "header" + h.attributes.level,
            text: h.textContent
        })
    })

    // change it to the TOC xml format
    const Toc =  new TocGenerator(headers)
    const res = Toc.process()

    tempData.toc = res;
    tempData.metadata.creator = tempData.metadata.creator.split(',');
    tempData.metadata.contributors = tempData.metadata.contributors.split(',');

    // Update the ids of all headers so that the epub's table of contents 
    // corresponds to the headers within the article
    const doc = new DOMParser().parseFromString(editor.getHTML(), 'text/html')
    let titles = doc.querySelectorAll("h1, h2, h3, h4, h5, h6")
    titles.forEach((e, i) => {
        e.id = "guide_signal_" + (i+1);
    })

    // Normalize html tags to conform to the xhtml specification
    tempData.content = normalizeHTML(doc.body.innerHTML);
    tempData.content = tempData.content.replaceAll("<p></p>", "<br></br>"); // Handling empty paragraphs
    tempData.content = tempData.content.replaceAll(" ", "");
    tempData.content = tempData.content.replaceAll('"', "'");
    tempData.metadata.meta = JSON.parse(JSON.stringify(epubLayout));
    let name = bookInfo.metadata.title;

    // If the user has already opened the file, it will be saved directly 
    // to that path without asking the user
    if (currentSave.value !== "" && !saveAs) {
        FileSave(currentSave.value, (()=>{
            return JSON.stringify(tempData)
        }  
        )(),
        true
        ).then((res)=> {
            if(res.Code == 1){
                ElMessage.error(t('message.saveError') + ": " + res.Msg)
                return
            }else if(res.Code == 0){
                currentSave.value = res.Data;
                title.value = res.Data;
                change.value = false;
                timer.Reset();
                ElMessage.success(t('message.saveSuccess'))
                return
            }
        })
        return
    }

    FileSave(name, (()=>{
            return JSON.stringify(tempData)
        }  
        )(),
        false
    ).then((res)=> {
        if(res.Code == 1){
            ElMessage.error(t('message.saveError') + ": " + res.Msg)
            return
        }else if(res.Code == 0){
            currentSave.value = res.Data;
            title.value = res.Data;
            change.value = false;
            timer.Reset();
            ElMessage.success(t('message.saveSuccess'))
        }
        return
    })
}

// const importFilePicker = async () => {
//     let loading = ElLoading.service({
//         lock: true,
//         text: t('message.loadingMessage'),
//         background: 'rgba(0, 0, 0, 0.7)',
//     })
//     var res = await FileImport().then((res) => {
//         return res
//     })
//     if(res.Code == 1){
//         loading.close();
//         ElMessage.error(t('message.openError'))
//         return;
//     }else if (res.Code == -1){
//         loading.close();
//         return;
//     }
//     const E = editorRef.value;
//     try{
//         var content = JSON.parse(res.Data);
//         console.log(content)
//         E.commands.insertContent(content);
//     }catch(e){
//         loading.close();
//         ElMessage.error(t('message.importError'))
//         return;
//     }
//     ElMessage.success(t('message.importSuccess'));
// }

const handleNew = ()=> {
    function innerCallBack(action){
        if(action == 'confirm'){
            let name = bookInfo.metadata.title;
            FileSave(name, JSON.stringify(bookInfo.value)).then((res)=>{
                if(res.Code == 0) {
                    ElMessage.info(t('message.saveSuccess'))
                    resetState();
                }else if(res.Code == -1){
                    return;
                }else{
                    ElMessage.error(t('message.saveError') + ": " + res);
                }
            })
        }else if(action == 'cancel'){
            resetState();
        }else if(action == 'close'){
            return;
        }
    }
    if(change.value){
        ElMessageBox.confirm(t('message.saveWarning'),t('message.warning'), {
            confirmButtonText: t('message.confirm'),
            cancelButtonText: t('message.notSave'),
            distinguishCancelAndClose: true,
            type: 'warning',
            callback: innerCallBack
        })
        return;
    }else {
        resetState(); // clear the editor and the bookinfo data
    }
}

const handleOpenLookup = () => {
    const E = editorRef.value;
    const resetSearch = () => {
        searchKey.value = "";
        replaceKey.value = "";
        resultCount.value = 0;
        E.chain().focus().run();
        E.setEditable(true);
    };

    const performLookup = (text) => {
        searchKey.value = text;
        lookupSession.lookup(E, text);
        visio.searchBarVisible = true;
    };

    E.setEditable(false);
    lookupSession.injectEditor(E);
    E.commands.blur();

    if (visio.searchBarVisible) {
        resetSearch();
    } else {
        const { from, to } = E.state.selection;
        const text = E.state.doc.textBetween(from, to);
        if (from !== to) {
            performLookup(text);
            return;
        }
    }

    lookupSession.destory();
    visio.searchBarVisible = !visio.searchBarVisible;
};

const openFileInfo = ()=> {
    visio.bookInfoVisible = true;
}
</script>

<template>
    <div class="func_btn-list">
        <div class="func_btn-container list_row">
            <span style="width:35px; " class="btn-group">
                <button class="el-button btn func_btn-big" @click="handleNew" id="btn-new">
                    <i class="el-icon">
                        <svg t="1717397117328" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4260" width="200" height="200"><path d="M524.8 64c4.693333 0 8.533333 3.84 8.533333 8.533333V490.666667h418.133334c4.693333 0 8.533333 3.84 8.533333 8.533333v46.933333a8.533333 8.533333 0 0 1-8.533333 8.533334H533.333333v418.133333a8.533333 8.533333 0 0 1-8.533333 8.533333h-46.933333a8.533333 8.533333 0 0 1-8.533334-8.533333V554.666667H51.2a8.533333 8.533333 0 0 1-8.533333-8.533334v-46.933333c0-4.693333 3.84-8.533333 8.533333-8.533333H469.333333V72.533333c0-4.693333 3.84-8.533333 8.533334-8.533333h46.933333z" p-id="4261"></path></svg>
                    </i>
                </button>
                <span>{{$t('toolBar.file.new')}}</span>
            </span>
            <div class="division-border"></div>
            <!-- <span style=" width:50px; " class="btn-group">
                <button class="el-button btn func_btn-big" @click="importFilePicker" id="btn-open">
                    <i class="el-icon">
                        <svg t="1717816393836" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="8897" width="200" height="200"><path d="M921.6 450.133333c-6.4-8.533333-14.933333-12.8-25.6-12.8h-10.666667V341.333333c0-40.533333-34.133333-74.666667-74.666666-74.666666H514.133333c-4.266667 0-6.4-2.133333-8.533333-4.266667l-38.4-66.133333c-12.8-21.333333-38.4-36.266667-64-36.266667H170.666667c-40.533333 0-74.666667 34.133333-74.666667 74.666667v597.333333c0 6.4 2.133333 12.8 6.4 19.2 6.4 8.533333 14.933333 12.8 25.6 12.8h640c12.8 0 25.6-8.533333 29.866667-21.333333l128-362.666667c4.266667-10.666667 2.133333-21.333333-4.266667-29.866667zM170.666667 224h232.533333c4.266667 0 6.4 2.133333 8.533333 4.266667l38.4 66.133333c12.8 21.333333 38.4 36.266667 64 36.266667H810.666667c6.4 0 10.666667 4.266667 10.666666 10.666666v96H256c-12.8 0-25.6 8.533333-29.866667 21.333334l-66.133333 185.6V234.666667c0-6.4 4.266667-10.666667 10.666667-10.666667z m573.866666 576H172.8l104.533333-298.666667h571.733334l-104.533334 298.666667z" p-id="8898"></path></svg>
                    </i>
                </button>
                <span>{{$t('toolBar.file.import')}}</span>
            </span> -->
            <span style=" width:50px; " class="btn-group">
                <button class="el-button btn func_btn-big" @click="openFilePicker" id="btn-open">
                    <i class="el-icon">
                        <svg t="1717816393836" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="8897" width="200" height="200"><path d="M921.6 450.133333c-6.4-8.533333-14.933333-12.8-25.6-12.8h-10.666667V341.333333c0-40.533333-34.133333-74.666667-74.666666-74.666666H514.133333c-4.266667 0-6.4-2.133333-8.533333-4.266667l-38.4-66.133333c-12.8-21.333333-38.4-36.266667-64-36.266667H170.666667c-40.533333 0-74.666667 34.133333-74.666667 74.666667v597.333333c0 6.4 2.133333 12.8 6.4 19.2 6.4 8.533333 14.933333 12.8 25.6 12.8h640c12.8 0 25.6-8.533333 29.866667-21.333333l128-362.666667c4.266667-10.666667 2.133333-21.333333-4.266667-29.866667zM170.666667 224h232.533333c4.266667 0 6.4 2.133333 8.533333 4.266667l38.4 66.133333c12.8 21.333333 38.4 36.266667 64 36.266667H810.666667c6.4 0 10.666667 4.266667 10.666666 10.666666v96H256c-12.8 0-25.6 8.533333-29.866667 21.333334l-66.133333 185.6V234.666667c0-6.4 4.266667-10.666667 10.666667-10.666667z m573.866666 576H172.8l104.533333-298.666667h571.733334l-104.533334 298.666667z" p-id="8898"></path></svg>
                    </i>
                </button>
                <span>{{$t('toolBar.file.open')}}</span>
            </span>
            <span class="btn-group" style="width:50px;">
                <button @click="saveFilePicker(false)" :class="disableBtn.file?btnDisabledClass:btnNormalClass" :aria-disabled="disableBtn.file" :disabled="disableBtn.file" id="btn-save">
                    <i class="el-icon">
                        <svg t="1717405451499" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="21449" width="200" height="200"><path d="M708.388571 121.904762L902.095238 320.804571V828.952381a73.142857 73.142857 0 0 1-73.142857 73.142857H195.047619a73.142857 73.142857 0 0 1-73.142857-73.142857V195.047619a73.142857 73.142857 0 0 1 73.142857-73.142857h513.340952zM292.571429 195.023238L195.047619 195.047619v633.904762l97.52381-0.024381V536.380952h438.857142v292.547048L828.952381 828.952381V350.549333l-97.52381-100.156952V365.714286H292.571429V195.023238zM658.285714 609.52381H365.714286v219.40419h292.571428V609.52381z m-48.761904 73.142857v73.142857h-195.04762v-73.142857h195.04762z m48.761904-487.619048H365.714286v97.52381h292.571428V195.047619z" p-id="21450"></path></svg>
                    </i>
                </button>
                <span>{{$t('toolBar.file.save')}}</span>
            </span>
            <span class="btn-group" style="width:50px;">
                <button @click="saveFilePicker(true)" :class="disableBtn.file?btnDisabledClass:btnNormalClass" :aria-disabled="disableBtn.file" :disabled="disableBtn.file" id="btn-save-as">
                    <i class="el-icon">
                        <svg t="1726113120392" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4276" width="200" height="200"><path d="M251.733333 533.333333c0-12.8 8.533333-21.333333 21.333333-21.333333l384 0c12.8 0 21.333333 8.533333 21.333333 21.333333S669.866667 554.666667 657.066667 554.666667l-384 0C264.533333 554.666667 251.733333 546.133333 251.733333 533.333333zM785.066667 640l-512 0c-12.8 0-21.333333 8.533333-21.333333 21.333333S264.533333 682.666667 273.066667 682.666667l512 0c12.8 0 21.333333-8.533333 21.333333-21.333333S797.866667 640 785.066667 640zM785.066667 768l-512 0c-12.8 0-21.333333 8.533333-21.333333 21.333333S264.533333 810.666667 273.066667 810.666667l512 0c12.8 0 21.333333-8.533333 21.333333-21.333333S797.866667 768 785.066667 768zM955.733333 541.866667c-12.8 0-21.333333 8.533333-21.333333 21.333333L934.4 853.333333c0 12.8-29.866667 42.666667-64 42.666667l-682.666667 0c-34.133333 0-64-29.866667-64-64l0-640C123.733333 157.866667 153.6 128 187.733333 128l64 0 0 149.333333C251.733333 337.066667 302.933333 384 358.4 384l541.866667 0 0 0-81.066667 81.066667c-8.533333 8.533333-8.533333 21.333333 0 29.866667 4.266667 4.266667 8.533333 8.533333 17.066667 8.533333 4.266667 0 12.8 0 17.066667-8.533333l115.2-115.2c0 0 0 0 0 0 0 0 4.266667-4.266667 4.266667-8.533333 0-4.266667 0-4.266667 0-8.533333 0 0 0 0 0 0 0 0 0 0 0 0 0-4.266667 0-4.266667 0-8.533333 0-4.266667-4.266667-4.266667-4.266667-4.266667 0 0 0 0 0 0L853.333333 230.4c-8.533333-8.533333-21.333333-8.533333-29.866667 0-8.533333 8.533333-8.533333 21.333333 0 29.866667L904.533333 341.333333l0 0L358.4 341.333333C324.266667 341.333333 294.4 311.466667 294.4 277.333333L294.4 128l149.333333 0 106.666667 0 0 128 128 0c12.8 0 21.333333-8.533333 21.333333-21.333333L699.733333 170.666667c0-34.133333-42.666667-76.8-93.866667-85.333333l-12.8 0-149.333333 0-256 0C132.266667 85.333333 81.066667 132.266667 81.066667 192l0 640c0 59.733333 46.933333 106.666667 106.666667 106.666667l682.666667 0c59.733333 0 106.666667-46.933333 106.666667-85.333333l0 0 0-294.4C977.066667 550.4 968.533333 541.866667 955.733333 541.866667z" p-id="4277"></path></svg>
                    </i>
                </button>
                <span>{{$t('toolBar.file.saveAs')}}</span>
            </span>
            <!-- <span class="btn-group" style="width:50px;">
                <button @click="exportFile" :class="disableBtn.file?btnDisabledClass:btnNormalClass" :aria-disabled="disableBtn.file" :disabled="disableBtn.file" id="btn-export">
                    <i class="el-icon">
                        <svg t="1717398848121" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4277" width="200" height="200"><path d="M895.99896 1024H127.99992a63.99992 63.99992 0 0 1-63.99992-63.99992V384.0008a63.99992 63.99992 0 0 1 63.99992-63.99992h223.99972a31.99996 31.99996 0 0 1 0 63.99992H127.99992v575.99928h767.99904V384.0008h-223.99972a31.99996 31.99996 0 0 1 0-63.99992h223.99972a63.99992 63.99992 0 0 1 63.99992 63.99992v575.99928a63.99992 63.99992 0 0 1-63.99992 63.99992zM543.9994 100.481154V704.0004a31.99996 31.99996 0 0 1-63.99992 0V100.481154L340.479654 216.641009 299.519706 167.361071l191.99976-159.9998 0.255999 0.319999a30.463962 30.463962 0 0 1 40.44795 0l0.255999-0.319999 191.99976 159.9998-40.959948 49.279938z" p-id="4278"></path></svg>
                    </i>
                </button>
                <span>{{$t('toolBar.file.export')}}</span>
            </span> -->
            <div class="division-border"></div>
            <span class="btn-group" style="width:80px;">
                <button @click="openFileInfo" :class="disableBtn.file?btnDisabledClass:btnNormalClass" :aria-disabled="disableBtn.file" :disabled="disableBtn.file" id="btn-file-info">
                    <i class="el-icon">
                        <svg t="1717677749278" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="7239" width="200" height="200"><path d="M688 312v-48c0-4.4-3.6-8-8-8H296c-4.4 0-8 3.6-8 8v48c0 4.4 3.6 8 8 8h384c4.4 0 8-3.6 8-8zM296 400c-4.4 0-8 3.6-8 8v48c0 4.4 3.6 8 8 8h184c4.4 0 8-3.6 8-8v-48c0-4.4-3.6-8-8-8H296zM672 516c-119.3 0-216 96.7-216 216s96.7 216 216 216 216-96.7 216-216-96.7-216-216-216z m107.5 323.5C750.8 868.2 712.6 884 672 884s-78.8-15.8-107.5-44.5C535.8 810.8 520 772.6 520 732s15.8-78.8 44.5-107.5C593.2 595.8 631.4 580 672 580s78.8 15.8 107.5 44.5C808.2 653.2 824 691.4 824 732s-15.8 78.8-44.5 107.5z" p-id="7240"></path><path d="M672 812m-32 0a32 32 0 1 0 64 0 32 32 0 1 0-64 0Z" p-id="7241"></path><path d="M652 748h40c4.4 0 8-3.6 8-8V628c0-4.4-3.6-8-8-8h-40c-4.4 0-8 3.6-8 8v112c0 4.4 3.6 8 8 8z" p-id="7242"></path><path d="M440 852H208V148h560v344c0 4.4 3.6 8 8 8h56c4.4 0 8-3.6 8-8V108c0-17.7-14.3-32-32-32H168c-17.7 0-32 14.3-32 32v784c0 17.7 14.3 32 32 32h272c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8z" p-id="7243"></path></svg>
                    </i>
                </button>
                <span>{{$t('toolBar.file.fileInfo')}}</span>
            </span>
            <div class="division-border"></div>
            <span class="btn-group" style="width:50px;">
                <button id="btn-lookup" @click="handleOpenLookup"  :class="disableBtn.file?btnDisabledClass:btnNormalClass" :aria-disabled="disableBtn.file" :disabled="disableBtn.file">
                    <i class="el-icon">
                        <svg t="1727537161154" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="11161" width="200" height="200"><path d="M699.05 251.05h-512v64h512v-64zM443.05 475.05h-256v64h256v-64zM379.05 699.05h-192v64h192v-64zM650.02 531.62c70.58 0 128 57.42 128 128s-57.42 128-128 128-128-57.42-128-128 57.42-128 128-128m0-64c-106.04 0-192 85.96-192 192s85.96 192 192 192 192-85.96 192-192-85.96-192-192-192zM801.81 766.16l-45.25 45.26L919.7 974.55l45.25-45.25-163.14-163.14z" p-id="11162"></path><path d="M59.05 49.45v915.2h512v-64h-448v-787.2h768v467.2h64V49.45z" p-id="11163"></path></svg>
                    </i>
                </button>
                <span>{{$t('toolBar.file.lookUp')}}</span>
            </span>
            <div class="division-border"></div>
        </div>
    </div>
</template>