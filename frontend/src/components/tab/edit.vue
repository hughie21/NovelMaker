<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-18
@LastEditors: Hughie
@LastEditTime: 2024-11-1
*/

import { useI18n } from 'vue-i18n';
import "../../assets/css/tab.css"
import { ref, inject } from 'vue';
import { editorRef, fontSizeVal, fontVal, headerVal, visio, fonts, isBold, isItalic, isStrike } from '../../assets/js/globals';
import { rgbaToHex } from '../../assets/js/utils';

const { t } = useI18n();
const fontColors = ref('#000000');
const backColors = ref('rgba(0, 0, 0, 0)');
const E = editorRef;
const fontColorPopper = ref()
const backColorPopper = ref()
const rowNum = ref(3)
const columnNum = ref(3)
const sizes = []
const btnNormalClass = inject("btnNormalClass");
const programLang = ref('')

for(var i = 12; i < 72; i+=2) {
    sizes.push({
        value: i*(1/16) + "rem",
        label: i + "px"
    })
}

const headers = [
    {
        value: 1,
        label: "H1",
        size: "2em"
    },
    {
        value: 2,
        label: "H2",
        size: "1.8em"
    
    },
    {
        value: 3,
        label: "H3",
        size: "1.6em"
    },
    {
        value: 4,
        label: "H4",
        size: "1.4em"
    },
    {
        value: 5,
        label: "H5",
        size: "1.2em"
    },
    {
        value: 6,
        label: "H6",
        size: "1em"
    },
    {
        value: 0,
        label: "Paragraph",
        size: "0.8em"
    }
]

const handleFontSelChange = (val) => {
    const E = editorRef.value;
    E.chain().focus().setFontFamily(val).run();
}

const handleFontSize = (val) => {
    const E = editorRef.value;
    E.chain().focus().setFontSize(val).run();
}

const handleHeaderSelChange = (val) => {
    const E = editorRef.value;
    E.chain().focus().setHeader(val).run();
}

const handleInsertTable = () => {
    const E = editorRef.value;
    E.chain().focus().insertTable({ rows: rowNum.value, cols: columnNum.value, withHeaderRow: true }).run();
    visio.tableInsertVisible = false;
}

const openFontColorPicker = () => {
    fontColorPopper.value.show();
}

const openBackColorPicker = () => {
    backColorPopper.value.show();
}

const changeFontColor = (val) => {
    const E = editorRef.value;
    E.chain().focus().setColor(val).run();
}

const changeBackColor = (val) => {
    const E = editorRef.value;
    E.chain().focus().setTextBackColor(rgbaToHex(val)).run();
}

const handlePaste = () => {
    const E = editorRef.value;
    const clipboard = localStorage.getItem('clipboard');
    E.chain().focus().insertContent(clipboard).run();
}

const handleCopy = () => {
    const text = getSelectedText();
    if(text == ""){
        return;
    }
    localStorage.setItem('clipboard', text);
}

const handleCut = () => {
    const E = editorRef.value;
    const text = getSelectedText();
    if(text == ""){
        return;
    }
    localStorage.setItem('clipboard', text);
    E.chain().focus().deleteSelection().run();
}

const getSelectedText = () => {
    const E = editorRef.value;
    const { state, dispatch } = E;
    const doc = state.doc;
    const { from, to } = state.selection;
    const text = doc.textBetween(from, to);
    return text;
}

const handleSetBold = () => {
    const E = editorRef.value;
    E.chain().focus().toggleBold().run();
    isBold.value = !isBold.value;
}

const handleSetItalic = () => {
    const E = editorRef.value;
    E.chain().focus().toggleItalic().run();
    isItalic.value = !isItalic.value;
}

const handleSetStrike = () => {
    const E = editorRef.value;
    E.chain().focus().toggleStrike().run();
    isStrike.value = !isStrike.value;
}

const handleInsetCodeBlock = () => {
    const E = editorRef.value;
    console.log(programLang.value)
    E.chain().focus().setCodeBlock({language: programLang.value}).run();
    visio.codeInsertVisible = false;
}

const handleCleanMarks = () => {
    const E = editorRef.value;
    E.chain().focus().unsetAllMarks().run()
    isBold.value = false;
    isItalic.value = false;
    isStrike.value = false;
}

const handleAlign = (val) => {
    const E = editorRef.value;
    E.chain().focus().setTextAlign(val).run();
}
</script>

<template>
    <div class="func_btn-list">
        <div class="func_btn-container list_row">
            <span class="btn-group" style="width:40px;">
                <button class="el-button btn func_btn-big" @click="handlePaste" id="btn-paste">
                    <i class="el-icon">
                        <svg t="1717403088327" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="12105" width="200" height="200"><path d="M704 128l-128 0 0-64c0-35.2-28.8-64-64-64l-128 0c-35.2 0-64 28.8-64 64l0 64-128 0 0 128 512 0 0-128zM512 128l-128 0 0-63.872c0.032-0.032 0.064-0.064 0.128-0.128l127.776 0c0.032 0.032 0.096 0.064 0.128 0.128l0 63.872zM832 320l0-160c0-17.6-14.4-32-32-32l-64 0 0 64 32 0 0 128-192 0-192 192 0 256-256 0 0-576 32 0 0-64-64 0c-17.6 0-32 14.4-32 32l0 640c0 17.6 14.4 32 32 32l288 0 0 192 640 0 0-704-192 0zM576 410.496l0 101.504-101.504 0 101.504-101.504zM960 960l-512 0 0-384 192 0 0-192 320 0 0 576z" p-id="12106"></path></svg>
                    </i>
                </button>
                {{$t('toolBar.edit.paste')}}
            </span>
            <span class="btn-group">
                <span style="margin-bottom: 2px;" >
                    <el-tooltip :content="t('toolBar.tooltip.copy')" placement="top" effect='dark'>
                        <button class="el-button btn is-link" @click="handleCopy" id="btn-cpoy">
                            <i class="el-icon">
                                <svg t="1717403312181" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="13146" width="200" height="200"><path d="M629.39164 415.032528l-163.616428-163.616428c-7.992021-7.992021-20.947078-7.992021-28.939099 0-7.992021 8.002254-7.992021 20.957311 0 28.949332l128.680754 128.680754-175.548178 0L389.968689 184.082552c0-11.2973-9.168824-20.466124-20.466124-20.466124L21.813818 163.616428c-11.307533 0-20.466124 9.168824-20.466124 20.466124l0 818.08214c0 11.307533 9.15859 20.466124 20.466124 20.466124l593.108273 0c11.307533 0 20.466124-9.15859 20.466124-20.466124L635.388215 429.512311C635.388215 424.078555 633.229039 418.880159 629.39164 415.032528zM594.455967 981.698568l-552.176025 0L42.279942 204.548676l306.756499 0 0 224.963635c0 11.2973 9.15859 20.466124 20.466124 20.466124l224.953402 0L594.455967 981.698568z" p-id="13147"></path><path d="M1023.978511 265.895883l0 572.652382c0 11.307533-9.15859 20.466124-20.466124 20.466124l-307.86167 0c-11.2973 0-20.466124-9.15859-20.466124-20.466124 0-11.2973 9.168824-20.466124 20.466124-20.466124l287.395546 0L983.046263 286.362007l-224.953402 0c-11.307533 0-20.466124-9.168824-20.466124-20.466124L737.626737 40.932248l-306.756499 0 0 75.693959c0 11.307533-9.168824 20.466124-20.466124 20.466124-11.307533 0-20.466124-9.15859-20.466124-20.466124L389.93799 20.466124c0-11.2973 9.15859-20.466124 20.466124-20.466124l347.688747 0c11.2973 0 20.466124 9.168824 20.466124 20.466124l0 224.963635 175.548178 0-128.680754-128.680754c-7.992021-7.992021-7.992021-20.947078 0-28.949332 7.992021-7.992021 20.947078-7.992021 28.939099 0l163.616428 163.626661C1021.819334 255.263731 1023.978511 260.462127 1023.978511 265.895883z" p-id="13148"></path></svg>
                            </i>
                        </button>
                    </el-tooltip>
                </span>
                <span>
                    <el-tooltip :content="t('toolBar.tooltip.cut')" placement="bottom" effect='dark'>
                        <button class="el-button btn is-link" @click="handleCut" id="btn-cut">
                            <i class="el-icon">
                                <svg t="1717403432425" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="14228" width="200" height="200"><path d="M790.4 640c-64 0-120 30.4-153.6 81.6L560 507.2l166.4-448c8-20.8-4.8-46.4-25.6-56-20.8-8-46.4 4.8-56 25.6L512 382.4 384 28.8C376 8 350.4-4.8 328 3.2c-20.8 8-33.6 33.6-25.6 56l161.6 448L388.8 720c-33.6-51.2-89.6-81.6-153.6-81.6-107.2 0-192 84.8-192 192s84.8 192 192 192c94.4 0 171.2-64 188.8-153.6V864l89.6-240 84.8 240c0 4.8 4.8 4.8 4.8 8C622.4 960 696 1024 790.4 1024c107.2 0 192-84.8 192-192s-86.4-192-192-192zM233.6 939.2c-59.2 0-105.6-48-105.6-107.2s46.4-107.2 107.2-107.2 107.2 46.4 107.2 107.2c-1.6 59.2-48 107.2-108.8 107.2z m556.8 0c-59.2 0-107.2-46.4-107.2-107.2 0-59.2 46.4-107.2 107.2-107.2 59.2 0 107.2 46.4 107.2 107.2-1.6 59.2-48 107.2-107.2 107.2z" p-id="14229"></path></svg>
                            </i>
                        </button>
                    </el-tooltip>
                </span>
            </span>
            <div class="division-border"></div>
            <span class="btn-group">
                <el-select v-model="headerVal" size="small" style="width: 160px" @change="handleHeaderSelChange">
                    <el-option
                        v-for="item in headers"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                    >
                    <span :style="`font-size:${item.size}`">{{ item.label }}</span>
                    </el-option>
                </el-select>
                <div style="display: flex; flex-direction:row; align-items: center; margin-left: 3px; margin-top: 5px">
                    <el-tooltip :content="t('toolBar.tooltip.bold')" placement="bottom" effect='dark'>
                        <el-button :class="{'selected-icon':isBold}" text size="small" @click="handleSetBold">
                            <i class="el-icon">
                                <svg t="1724558936693" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4250" width="200" height="180"><path d="M385.692278 918.576916 390.80771 919.307658C421.012953 920.769221 443.423035 921.499963 458.03843 921.499963 545.73119 921.499963 609.551045 905.545019 649.499963 873.63458 689.448881 841.724219 709.423104 790.449073 709.423104 719.807724 709.423104 655.986846 690.66689 608.365568 653.153831 576.942316 615.640852 545.519065 559.372209 529.807675 484.346171 529.807675 462.422961 529.807675 444.153935 530.051308 429.538462 530.538496 414.922988 531.025605 400.307751 531.756426 385.692278 532.730801L385.692278 918.576916ZM386.423099 445.769255C390.320522 446.256443 395.192241 446.621775 401.038494 446.865408 406.884667 447.108962 415.410176 447.23074 426.615414 447.23074 504.5645 447.23074 561.929295 431.51935 598.711532 400.096177 635.49377 368.672926 653.884652 319.590085 653.884652 252.846159 653.884652 194.384345 636.833477 151.025822 602.730732 122.76925 568.628066 94.512679 516.256768 80.38463 445.615419 80.38463 433.922993 80.38463 415.410412 81.358927 390.076889 83.307678L385.692278 83.307678 386.423099 445.769255ZM78.769231 0 528.192276 0C638.295434 0 721.115136 19.608812 776.653824 58.82691 832.192591 98.045086 859.961502 156.62781 859.961502 234.576896 859.961502 293.525898 841.936108 342.852372 805.884613 382.557657 769.833118 422.26302 716.731156 451.371717 646.576916 469.884613 732.320926 471.346176 799.063828 493.268992 846.80767 535.653849 894.551513 578.038705 918.423079 636.499653 918.423079 711.038425 918.423079 805.551734 886.026004 878.018954 821.23075 928.442289 756.435574 978.865625 663.141612 1004.076898 541.346186 1004.076898L78.769231 1004.076898 78.769231 931.000005 167.19234 922.961526 191.307697 899.576911 191.307697 101.57694 167.19234 80.38463 78.769231 73.076894 78.769231 0Z" p-id="4251"></path></svg>
                            </i>
                        </el-button>
                    </el-tooltip>
                    
                    <el-tooltip :content="t('toolBar.tooltip.italic')" placement="bottom" effect='dark'>
                        <el-button :class="{'selected-icon':isItalic}" text size="small" @click="handleSetItalic">
                            <i class="el-icon">
                                <svg t="1724559086522" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5526" width="200" height="200"><path d="M896 64 896 128 768 128 448 896 576 896 576 960 128 960 128 896 256 896 576 128 448 128 448 64Z" p-id="5527"></path></svg>
                            </i>
                        </el-button>
                    </el-tooltip>
                    
                    <el-tooltip :content="t('toolBar.tooltip.strike')" placement="bottom" effect='dark'>
                        <el-button :class="{'selected-icon':isStrike}" text size="small" @click="handleSetStrike">
                            <i class="el-icon">
                                <svg t="1724559302042" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6660" width="200" height="200"><path d="M731.424 517.024c63.904 47.936 100.576 116.096 100.576 186.976s-36.672 139.04-100.576 186.976c-59.36 44.512-137.28 69.024-219.424 69.024s-160.064-24.512-219.424-69.024c-63.936-47.936-100.576-116.096-100.576-186.976l128 0c0 69.376 87.936 128 192 128s192-58.624 192-128c0-69.376-87.936-128-192-128-82.144 0-160.064-24.512-219.424-69.024-63.936-47.936-100.576-116.096-100.576-186.976s36.672-139.04 100.576-186.976c59.36-44.512 137.28-69.024 219.424-69.024s160.064 24.512 219.424 69.024c63.904 47.936 100.576 116.096 100.576 186.976l-128 0c0-69.376-87.936-128-192-128s-192 58.624-192 128c0 69.376 87.936 128 192 128 82.144 0 160.064 24.512 219.424 69.024zM0 512l1024 0 0 64-1024 0z" p-id="6661"></path></svg>
                            </i>
                        </el-button>
                    </el-tooltip>
                    
                    <el-tooltip :content="t('toolBar.tooltip.clean')" placement="bottom" effect='dark'>
                        <el-button text size="small" @click="handleCleanMarks">
                            <i class="el-icon">
                                <svg t="1724559411257" class="icon" viewBox="0 0 1152 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="8509" width="200" height="200"><path d="M581.4 114.8L114.8 581.4c-50 50-50 131 0 181l160 160c24 24 56.6 37.4 90.6 37.4H1024c35.4 0 64-28.6 64-64s-28.6-64-64-64H775.8l261.4-261.2c50-50 50-131 0-181L762.6 114.8c-50-50-131-50-181 0z m13.4 717.2H365.2l-160-160 249.4-249.4 274.8 274.8-134.6 134.6z" p-id="8510"></path></svg>
                            </i>
                        </el-button>
                    </el-tooltip>
                </div>
            </span>
            <div class="division-border"></div>
            <span class="btn-group">
                <div class="col-display">
                    <el-select v-model="fontVal" size="small" style="width: 140px" @change="handleFontSelChange">
                        <el-option
                            v-for="item in fonts"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                        >
                        <span :style="`font-family: ${item.value}`">{{ item.label }}</span>
                        </el-option>
                    </el-select>
                    <el-select v-model="fontSizeVal" size="small" style="width: 80px; margin-left: 3px" @change="handleFontSize">
                        <el-option
                            v-for="item in sizes"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                        >
                        <span>{{ item.label }}</span>
                        </el-option>
                    </el-select>
                </div>
                
                <div class="col-display" style="margin-top: 5px;">
                    <el-tooltip :content="t('toolBar.tooltip.textAlignLeft')" placement="bottom" effect='dark'>
                        <el-button text size="small" @click="handleAlign('left')">
                            <i class="el-icon">
                                <svg t="1742202517010" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4310" width="200" height="200"><path d="M0 73.142857A36.571429 36.571429 0 0 1 36.571429 36.571429h950.857142a36.571429 36.571429 0 0 1 0 73.142857H36.571429A36.571429 36.571429 0 0 1 0 73.142857zM0 292.571429a36.571429 36.571429 0 0 1 36.571429-36.571429h731.428571a36.571429 36.571429 0 0 1 0 73.142857H36.571429A36.571429 36.571429 0 0 1 0 292.571429zM0 512a36.571429 36.571429 0 0 1 36.571429-36.571429h512a36.571429 36.571429 0 0 1 0 73.142858h-512A36.571429 36.571429 0 0 1 0 512zM0 950.857143a36.571429 36.571429 0 0 1 36.571429-36.571429h950.857142a36.571429 36.571429 0 0 1 0 73.142857H36.571429A36.571429 36.571429 0 0 1 0 950.857143zM0 731.428571a36.571429 36.571429 0 0 1 36.571429-36.571428h731.428571a36.571429 36.571429 0 0 1 0 73.142857H36.571429A36.571429 36.571429 0 0 1 0 731.428571z" p-id="4311"></path></svg>
                            </i>
                        </el-button>
                    </el-tooltip>

                    <el-tooltip :content="t('toolBar.tooltip.textAlignCenter')" placement="bottom" effect='dark'>
                        <el-button text size="small" @click="handleAlign('center')">
                            <i class="el-icon">
                                <svg t="1742202712001" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4474" width="200" height="200"><path d="M0 73.142857A36.571429 36.571429 0 0 1 36.571429 36.571429h950.857142a36.571429 36.571429 0 0 1 0 73.142857H36.571429A36.571429 36.571429 0 0 1 0 73.142857zM109.714286 292.571429a36.571429 36.571429 0 0 1 36.571428-36.571429h731.428572a36.571429 36.571429 0 0 1 0 73.142857H146.285714a36.571429 36.571429 0 0 1-36.571428-36.571428zM219.428571 512a36.571429 36.571429 0 0 1 36.571429-36.571429h512a36.571429 36.571429 0 0 1 0 73.142858h-512A36.571429 36.571429 0 0 1 219.428571 512zM0 950.857143a36.571429 36.571429 0 0 1 36.571429-36.571429h950.857142a36.571429 36.571429 0 0 1 0 73.142857H36.571429A36.571429 36.571429 0 0 1 0 950.857143zM109.714286 731.428571a36.571429 36.571429 0 0 1 36.571428-36.571428h731.428572a36.571429 36.571429 0 0 1 0 73.142857H146.285714a36.571429 36.571429 0 0 1-36.571428-36.571429z" p-id="4475"></path></svg>
                            </i>
                        </el-button>
                    </el-tooltip>

                    <el-tooltip :content="t('toolBar.tooltip.textAlignRight')" placement="bottom" effect='dark'>
                        <el-button text size="small" @click="handleAlign('right')">
                            <i class="el-icon">
                                <svg t="1742202730600" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4638" width="200" height="200"><path d="M0 73.142857A36.571429 36.571429 0 0 1 36.571429 36.571429h950.857142a36.571429 36.571429 0 0 1 0 73.142857H36.571429A36.571429 36.571429 0 0 1 0 73.142857zM219.428571 292.571429a36.571429 36.571429 0 0 1 36.571429-36.571429h731.428571a36.571429 36.571429 0 0 1 0 73.142857h-731.428571A36.571429 36.571429 0 0 1 219.428571 292.571429zM438.857143 512a36.571429 36.571429 0 0 1 36.571428-36.571429h512a36.571429 36.571429 0 0 1 0 73.142858h-512A36.571429 36.571429 0 0 1 438.857143 512zM0 950.857143a36.571429 36.571429 0 0 1 36.571429-36.571429h950.857142a36.571429 36.571429 0 0 1 0 73.142857H36.571429A36.571429 36.571429 0 0 1 0 950.857143zM219.428571 731.428571a36.571429 36.571429 0 0 1 36.571429-36.571428h731.428571a36.571429 36.571429 0 0 1 0 73.142857h-731.428571A36.571429 36.571429 0 0 1 219.428571 731.428571z" p-id="4639"></path></svg>
                            </i>
                        </el-button>
                    </el-tooltip>

                    <el-tooltip :content="t('toolBar.tooltip.textAlignJustify')" placement="bottom" effect='dark'>
                        <el-button text size="small" @click="handleAlign('justify')">
                            <i class="el-icon">
                                <svg t="1742203018007" class="icon" viewBox="0 0 1433 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="12223" width="200" height="200"><path d="M0 153.6h1433.9072V0H0zM0 588.8h1433.9072v-153.6H0zM0 1024h1433.9072V870.4H0z" p-id="12224"></path></svg>
                            </i>
                        </el-button>
                    </el-tooltip>
                </div>
            </span>
            <div class="division-border"></div>
            <span class="btn-group">
                <div class="col-display">
                    <el-popover
                        placement="bottom"
                        :width="200"
                        trigger="click"
                        :visible="visio.tableInsertVisible"
                    >
                        <template #reference>
                                <el-button text size="small"  @click="visio.tableInsertVisible = !visio.tableInsertVisible">
                                    <el-tooltip :content="t('toolBar.tooltip.table')" placement="top" effect='dark'>
                                    <i class="el-icon">
                                        <svg t="1742449322678" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="2651" width="200" height="200"><path d="M0 64l0 896 1024 0 0-896-1024 0zM384 640l0-192 256 0 0 192-256 0zM640 704l0 192-256 0 0-192 256 0zM640 192l0 192-256 0 0-192 256 0zM320 192l0 192-256 0 0-192 256 0zM64 448l256 0 0 192-256 0 0-192zM704 448l256 0 0 192-256 0 0-192zM704 384l0-192 256 0 0 192-256 0zM64 704l256 0 0 192-256 0 0-192zM704 896l0-192 256 0 0 192-256 0z" p-id="2652"></path></svg>
                                    </i>
                                    </el-tooltip>
                                </el-button>
                        </template>
                        <div>
                            <div style="display: flex; justify-content: space-between; align-items: center">
                                <span>{{t("toolBar.edit.row")}}</span>
                                <el-input-number v-model="rowNum" size="small" :min="1" controls-position="right"/>
                            </div>
                            <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 10px">
                                <span>{{t("toolBar.edit.column")}}</span>
                                <el-input-number v-model="columnNum" size="small" :min="2" controls-position="right"/>
                            </div>
                            <div style="display: flex; justify-content: center; align-items: center; margin-top: 5px">
                                <el-button size="small" text @click="handleInsertTable">{{t("toolBar.edit.insertTable")}}</el-button>
                                <el-button size="small" text @click="visio.tableInsertVisible = false">{{t("dialog.cancel")}}</el-button>
                            </div>
                        </div>
                    </el-popover>
                    <el-popover placement="bottom" :width="200" trigger="click" :visible="visio.codeInsertVisible">
                        <template #reference>
                            <el-button text size="small" @click="visio.codeInsertVisible = !visio.codeInsertVisible">
                                <el-tooltip :content="t('toolBar.tooltip.code')" placement="top" effect='dark'>
                                    <i class="el-icon">
                                        <svg t="1742449374319" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="3653" width="200" height="200"><path d="M322.133333 296.533333c-12.8-12.8-32-12.8-44.8 0l-192 192c-12.8 12.8-12.8 32 0 44.8l192 192c6.4 6.4 14.933333 8.533333 23.466667 8.533334s17.066667-2.133333 23.466667-8.533334c12.8-12.8 12.8-32 0-44.8L151.466667 512l168.533333-168.533333c12.8-12.8 12.8-34.133333 2.133333-46.933334zM940.8 488.533333l-192-192c-12.8-12.8-32-12.8-44.8 0-12.8 12.8-12.8 32 0 44.8l168.533333 168.533334-168.533333 168.533333c-12.8 12.8-12.8 32 0 44.8 6.4 6.4 14.933333 8.533333 23.466667 8.533333s17.066667-2.133333 23.466666-8.533333l192-192c8.533333-8.533333 8.533333-29.866667-2.133333-42.666667zM622.933333 76.8c-17.066667-4.266667-34.133333 6.4-38.4 23.466667L366.933333 902.4c-4.266667 17.066667 6.4 34.133333 23.466667 38.4 2.133333 0 6.4 2.133333 8.533333 2.133333 14.933333 0 27.733333-8.533333 29.866667-23.466666L644.266667 115.2c4.266667-17.066667-4.266667-34.133333-21.333334-38.4z" p-id="3654"></path></svg>
                                    </i>
                                </el-tooltip>
                            </el-button>
                        </template>
                        <div style="display: flex; justify-content: space-between; align-items: center;">
                            <el-input
                            style="width: 150px;border:none;"
                            :placeholder="t('toolBar.edit.code')"
                            clearable
                            size="small"
                            v-model="programLang"
                            ></el-input>
                            <el-button size="small" text @click="handleInsetCodeBlock">{{t("dialog.media.insert")}}</el-button>
                        </div>
                        
                    </el-popover>

                    <el-popover placement="bottom" :width="200" trigger="click" :visible="visio.aboveTextVisible">
                        <template #reference>
                            <el-button text size="small" @click="visio.aboveTextVisible = !visio.aboveTextVisible">
                                <el-tooltip :content="t('toolBar.tooltip.aboveText')" placement="top" effect='dark'>
                                    <i class="el-icon">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="icon" viewBox="0 0 1024 1024" style="font-style: normal;height:2em"><text data-v-19c4b0ab="" y="1200" x="0" font-size="100em">R</text><text data-v-19c4b0ab="" y="100" x="0" font-size="50em">abc</text></svg>
                                    </i>
                                </el-tooltip>
                            </el-button>
                        </template>
                        <div style="display: flex; justify-content: space-between; align-items: center;">
                            <el-input
                            style="width: 150px;border:none;"
                            :placeholder="t('toolBar.edit.abovePlaceHolder')"
                            clearable
                            size="small"
                            v-model="programLang"
                            ></el-input>
                            <el-button size="small" text>{{t("message.confirm")}}</el-button>
                        </div>
                    </el-popover>
                    
                </div>
                <div class="col-display" style="margin-top: 5px;">
                    <el-tooltip :content="t('toolBar.tooltip.color')" placement="bottom" effect='dark'>
                        <el-button text size="small" @click="openFontColorPicker">
                            <el-color-picker v-model="fontColors" size="small" ref="fontColorPopper" @change="changeFontColor" style="opacity: 0;width:0;height:0;padding:0"/>
                            <i class="el-icon">
                                <svg t="1725421210952" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4750" width="200" height="200"><path d="M1013.763272 901.117134H10.236925c-5.630255 0-10.236827 4.606572-10.236827 10.247063v102.398976c0 5.630255 4.606572 10.236827 10.236827 10.236827H1013.763272c5.630255 0 10.236827-4.606572 10.236826-10.236827V911.35396c0-5.630255-4.606572-10.236827-10.236826-10.236826zM181.376192 798.707921h108.796993c5.374334 0 10.247063-3.460047 11.905429-8.701303l68.740291-212.485809h280.570943l68.09537 212.485809c1.658366 5.118413 6.398017 8.701303 11.905429 8.701303h114.038248c1.412682 0 2.815127-0.255921 4.104968-0.634684 6.520859-2.303286 9.980906-9.346223 7.67762-15.877318L590.849255 8.445382c-1.791445-4.995571-6.531095-8.445382-11.77235-8.445382H448.127419c-5.374334 0-10.103748 3.326969-11.772351 8.445382L169.603841 782.195919c-0.511841 1.279603-0.634683 2.692285-0.634683 4.094731-0.133079 6.787016 5.497176 12.417271 12.407034 12.417271z m327.54774-660.602894h5.251492l107.261469 337.661725H400.515938l108.407994-337.661725z m0 0" p-id="4751"></path></svg>
                            </i>
                        </el-button>
                    </el-tooltip>
                    
                    <el-tooltip :content="t('toolBar.tooltip.background')" placement="bottom" effect='dark'>
                        <el-button text size="small" @click="openBackColorPicker">
                            <el-color-picker v-model="backColors" show-alpha size="small" @change="changeBackColor" ref="backColorPopper" style="opacity: 0;width:0;height:0;padding:0"/>
                            <i class="el-icon">
                                <svg t="1725440391145" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4705" width="300" height="300"><path d="M436.056925 570.188815l150.141511 0-75.943075-243.428342L436.056925 570.188815zM895.409902 896.025656 128.590098 896.025656 128.590098 127.974344l766.819804 0L895.409902 896.025656zM787.652836 793.707757l-202.17278-564.544398L444.369613 229.163359 238.297054 793.707757 368.426538 793.707757l41.460814-129.411104 204.738425 0L656.804971 793.707757 787.652836 793.707757z" p-id="4706"></path></svg>
                            </i>
                        </el-button>
                    </el-tooltip>
                    
                    <el-tooltip :content="t('toolBar.tooltip.orderlist')" placement="bottom" effect='dark'>
                        <el-button text size="small" @click="E.chain().focus().toggleOrderedList().run()">
                            <i class="el-icon">
                                <svg t="1725442807611" class="icon" viewBox="0 0 1030 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="12101" width="200" height="200"><path d="M971.950987 487.29438a50.629671 50.629671 0 0 1 1.419342 101.239898l-1.425823 0.025925H394.739038a50.629671 50.629671 0 0 1-1.419342-101.24638l1.425823-0.019443h577.211949z m0-420.254785a50.629671 50.629671 0 0 1 1.419342 101.24638l-1.425823 0.019443H394.739038a50.629671 50.629671 0 0 1-1.419342-101.24638l1.425823-0.019443h577.211949z m0 840.503089a50.629671 50.629671 0 0 1 1.419342 101.246379l-1.425823 0.019443H394.739038a50.629671 50.629671 0 0 1-1.419342-101.246379l1.425823-0.019443h577.211949zM217.457418 852.155949l-3.739545 3.998785-9.520607 10.434431-22.249317 24.738025-32.58005 36.48162h70.137519a45.568 45.568 0 0 1 45.548557 44.232912l0.019443 1.341569a45.568 45.568 0 0 1-44.226431 45.548557l-1.341569 0.019443H47.894684c-39.300861 0-60.169722-46.417013-34.090127-75.814886l79.327595-89.21762 32.631899-36.442734 16.824708-18.55514 7.213368-7.796658 3.156253-3.305316 2.060962-2.060962 0.213873-0.194431c6.092152-5.606076 9.397468-12.845367 9.397469-20.512405 0-15.833114-14.536911-29.644152-33.610532-29.644152-18.639392 0-32.949468 13.201823-33.578127 28.594228l-0.019443 1.049924a45.568 45.568 0 0 1-91.142481 0c0-67.240506 56.358886-120.780152 124.740051-120.780152 68.394127 0 124.746532 53.533165 124.746532 120.780152 0 32.89762-13.687899 63.656506-37.239899 86.093772l-1.075848 1.011038zM102.108354 18.49681c28.827544-26.844354 75.613975-6.902278 76.60557 32.081013l0.012962 1.270278v369.618633a45.568 45.568 0 0 1-91.116557 1.34157l-0.012962-1.34157v-265.397468a45.580962 45.580962 0 0 1-59.003139-4.86076l-1.004557-1.049924a45.568 45.568 0 0 1 1.250835-63.403747l1.036962-1.004557L102.114835 18.49681z" p-id="12102"></path></svg>
                            </i>
                        </el-button>
                    </el-tooltip>
                    
                    <el-tooltip :content="t('toolBar.tooltip.unorderlist')" placement="bottom" effect='dark'>
                        <el-button text size="small" @click="E.chain().focus().toggleBulletList().run()">
                            <i class="el-icon">
                                <svg t="1725442628444" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="9076" width="200" height="200"><path d="M274.976 96.32h748.56v113.52H274.976V96.336z m0 354.96h748.56V564.8H274.976v-113.52z m0 378.88h748.56v113.504H274.976v-113.504zM0.48 152.896c0-48.848 39.12-88.448 87.376-88.448 48.256 0 87.392 39.6 87.392 88.448 0 48.848-39.136 88.464-87.392 88.464-48.256 0-87.376-39.616-87.376-88.464z m0 378.88c0-48.864 39.12-88.464 87.376-88.464 48.256 0 87.392 39.6 87.392 88.464 0 48.848-39.136 88.448-87.392 88.448-48.256 0-87.376-39.6-87.376-88.448z m0 355.328c0-48.848 39.12-88.464 87.376-88.464 48.256 0 87.392 39.616 87.392 88.464s-39.136 88.448-87.392 88.448c-48.256 0-87.376-39.6-87.376-88.448z" p-id="9077"></path></svg>
                            </i>
                        </el-button>
                    </el-tooltip>
                </div> 
            </span>
            <div class="division-border"></div>
            <span class="btn-group">
                <button :class="btnNormalClass" @click="visio.mediaVisible = true" id="btn-media">
                    <i class="el-icon">
                        <svg t="1717816046552" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6855" width="200" height="200"><path d="M298.666667 277.333333L469.333333 384l-170.666666 106.666667v-213.333334z" p-id="6856"></path><path d="M213.333333 128a85.333333 85.333333 0 0 0-85.333333 85.333333v597.333334a85.333333 85.333333 0 0 0 85.333333 85.333333h597.333334a85.333333 85.333333 0 0 0 85.333333-85.333333V213.333333a85.333333 85.333333 0 0 0-85.333333-85.333333H213.333333z m597.333334 64H213.333333a21.333333 21.333333 0 0 0-21.333333 21.333333v456.832l154.88-112.896a32 32 0 0 1 36.266667-1.024l124.842666 80.938667 152.362667-148.138667a32 32 0 0 1 44.629333 0l127.018667 123.52V213.333333a21.333333 21.333333 0 0 0-21.333333-21.333333zM192 810.666667v-61.269334l174.762667-127.445333 127.829333 82.816a32 32 0 0 0 39.68-3.882667L682.666667 556.629333l148.352 144.213334 0.981333-0.981334V810.666667a21.333333 21.333333 0 0 1-21.333333 21.333333H213.333333a21.333333 21.333333 0 0 1-21.333333-21.333333z" p-id="6857"></path></svg>
                    </i>
                </button>
                {{$t('toolBar.edit.media')}}
            </span>
            <div class="division-border"></div>
            <span class="btn-group">
                <button :class="btnNormalClass" @click="visio.linkInsertVisible = true">
                    <i class="el-icon">
                        <svg t="1733985395557" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4272" width="200" height="200"><path d="M593.94368 715.648a10.688 10.688 0 0 0-14.976 0L424.21568 870.4c-71.68 71.68-192.576 79.232-271.68 0-79.232-79.232-71.616-200 0-271.616l154.752-154.752a10.688 10.688 0 0 0 0-15.04l-52.992-52.992a10.688 10.688 0 0 0-15.04 0L84.50368 530.688a287.872 287.872 0 0 0 0 407.488 288 288 0 0 0 407.488 0l154.752-154.752a10.688 10.688 0 0 0 0-15.04l-52.736-52.736z m344.384-631.168a288.256 288.256 0 0 1 0 407.616l-154.752 154.752a10.688 10.688 0 0 1-15.04 0l-52.992-52.992a10.688 10.688 0 0 1 0-15.104l154.752-154.688c71.68-71.68 79.232-192.448 0-271.68-79.104-79.232-200-71.68-271.68 0L443.92768 307.2a10.688 10.688 0 0 1-15.04 0l-52.864-52.864a10.688 10.688 0 0 1 0-15.04l154.88-154.752a287.872 287.872 0 0 1 407.424 0z m-296.32 240.896l52.672 52.736a10.688 10.688 0 0 1 0 15.04l-301.504 301.44a10.688 10.688 0 0 1-15.04 0l-52.736-52.672a10.688 10.688 0 0 1 0-15.04l301.632-301.504a10.688 10.688 0 0 1 15.04 0z" p-id="4273"></path></svg>
                    </i>
                </button>
                {{$t('toolBar.edit.link')}}
            </span>
            <div class="division-border"></div>
            <span class="btn-group" style="width:40px; ">
                <button class="el-button btn func_btn-big" @click="visio.settingVisible = true" id="btn-setting">
                    <i class="el-icon">
                        <svg t="1717404325809" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="19304" width="200" height="200"><path d="M885.196958 871.528145c-14.753141-5.464126-25.134981-8.742602-34.970408-12.56749-30.052695-11.474665-59.012564-24.588568-89.611671-34.970408-8.742602-2.732063-21.310093-1.639238-29.506282 2.732063-21.310093 10.38184-42.073773 22.402918-61.744627 36.063234-6.556952 4.371301-12.021078 14.206728-13.660316 22.402918-6.556952 36.063234-11.474665 72.67288-15.845966 109.282526-2.732063 20.76368-12.56749 30.052695-33.877583 29.506282-62.837452-0.546413-125.674905-0.546413-188.512357 0-20.76368 0-31.691933-8.196189-34.423996-28.959869-4.917714-36.609646-10.928253-72.67288-15.299554-109.282526-1.639238-12.56749-6.010539-20.217267-18.031617-26.227806-20.217267-9.835427-38.248884-23.495743-58.466151-33.33117-7.103364-3.824888-18.031617-5.464126-25.681394-2.732063-35.516821 12.56749-69.940817 26.774219-104.911225 40.434535-25.134981 9.835427-32.238345 7.649777-46.445074-16.392379-31.14552-53.002025-61.744627-106.00405-92.343735-159.006075-14.206728-24.588568-13.660316-30.599107 8.742602-48.084311 29.506282-23.495743 59.558977-45.898661 88.518846-69.940817 5.464126-4.371301 9.289015-14.206728 9.289015-21.310093 1.092825-20.217267-3.278476-41.52736 0-61.198215 3.824888-21.310093-4.917714-32.238345-20.217267-43.71301-27.320632-19.670855-52.455613-41.52736-79.229831-61.198215-15.845966-12.021078-19.670855-24.588568-9.289015-42.073773 33.33117-56.280501 66.115928-113.107414 98.900686-169.934328 9.289015-15.845966 21.310093-19.124442 38.248884-12.56749 36.063234 14.206728 72.126467 28.959869 108.736113 42.073773 7.649777 2.732063 19.670855 1.092825 27.320632-2.732063 20.76368-10.38184 39.341709-24.042156 60.105389-34.423996 10.928253-5.464126 14.753141-12.021078 15.845966-22.94933 4.917714-37.156059 10.928253-74.312118 15.845966-111.468177C386.322226 7.649777 397.796892 0 418.014159 0c62.837452 0.546413 125.674905 0.546413 188.512357 0 19.670855 0 30.599107 8.196189 32.784758 27.867044 4.917714 36.609646 10.928253 72.67288 15.299554 109.282526 1.639238 13.660316 6.556952 21.856505 19.670855 27.867044 19.124442 8.742602 37.156059 19.670855 54.641263 31.691933 11.474665 8.196189 21.310093 8.196189 33.877583 2.732063 33.877583-14.206728 68.301579-26.774219 102.725575-40.980947 19.124442-7.649777 32.238345-3.824888 42.620185 14.753141 31.14552 55.187676 63.930278 109.828939 95.62221 164.470202 12.56749 21.856505 10.928253 28.959869-9.289015 44.805836-28.959869 22.402918-57.919739 45.352248-87.426021 67.755166-9.835427 7.649777-12.56749 14.753141-11.474665 27.320632 2.185651 22.94933-0.546413 46.445074 0.546413 69.394404 0.546413 8.196189 4.371301 18.031617 10.38184 22.94933 28.959869 24.042156 59.012564 46.445074 88.518846 69.394404 20.217267 15.845966 21.310093 22.94933 8.742602 44.805836-32.238345 55.734088-64.47669 110.921764-97.261448 166.10944C901.042924 860.05348 891.753909 865.517606 885.196958 871.528145zM510.357893 692.304803c98.900686 0.546413 180.862581-79.776244 181.408993-177.584105 0.546413-100.539924-79.229831-181.955406-178.67693-182.501819-99.447099-0.546413-181.408993 80.322657-181.408993 179.769755C331.134551 610.342908 412.00362 691.75839 510.357893 692.304803z" p-id="19305"></path></svg>
                    </i>
                </button>
                {{$t('toolBar.edit.setting')}}
            </span>
            <div class="division-border"></div>
        </div>
    </div>
</template>

<style scoped>
.selected-icon{
    background: radial-gradient(circle, rgb(47 137 231 / 30%) 55%, rgba(174, 177, 238, 0) 50%, rgba(148, 187, 233, 0) 60%);
}
</style>