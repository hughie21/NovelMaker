<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-18
@LastEditors: Hughie
@LastEditTime: 2024-11-1
*/
import { visio } from '../../assets/js/globals';
import { GetVersion } from '../../../wailsjs/go/core/App';
import "../../assets/css/tab.css"

const handleOpenHelp = () => {
    document.getElementById('tab-0').click();
    visio.tourVsisible = true
}

const handleOpenAbout = () => {
    GetVersion()
}
</script>

<template>
    <div class="func_btn-list">
        <div class="func_btn-container list_row">
            <span class="btn-group" style="width:40px; ">
                <button class="el-button btn func_btn-big" @click="handleOpenHelp" id="btn-help">
                    <i class="el-icon">
                        <svg t="1734074011736" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="1475" width="200" height="200"><path d="M512 74.666667C270.933333 74.666667 74.666667 270.933333 74.666667 512S270.933333 949.333333 512 949.333333 949.333333 753.066667 949.333333 512 753.066667 74.666667 512 74.666667z m0 810.666666c-204.8 0-373.333333-168.533333-373.333333-373.333333S307.2 138.666667 512 138.666667 885.333333 307.2 885.333333 512 716.8 885.333333 512 885.333333z" p-id="1476"></path><path d="M512 746.666667m-42.666667 0a42.666667 42.666667 0 1 0 85.333334 0 42.666667 42.666667 0 1 0-85.333334 0Z" p-id="1477"></path><path d="M512 245.333333c-76.8 0-138.666667 61.866667-138.666667 138.666667 0 17.066667 14.933333 32 32 32s32-14.933333 32-32c0-40.533333 34.133333-74.666667 74.666667-74.666667s74.666667 34.133333 74.666667 74.666667c0 27.733333-53.333333 76.8-91.733334 100.266667-8.533333 6.4-14.933333 17.066667-14.933333 27.733333v106.666667c0 17.066667 14.933333 32 32 32s32-14.933333 32-32v-89.6c34.133333-25.6 106.666667-83.2 106.666667-145.066667 0-76.8-61.866667-138.666667-138.666667-138.666667z" p-id="1478"></path></svg>
                    </i>
                </button>
                {{$t('toolBar.help.help')}}
            </span>
            <div class="division-border"></div>
            <span class="btn-group" style="width:40px; ">
                <button class="el-button btn func_btn-big" @click="handleOpenAbout" id="btn-about">
                    <i class="el-icon">
                        <svg t="1717404400946" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="20342" width="200" height="200"><path d="M512 97.52381c228.912762 0 414.47619 185.563429 414.47619 414.47619s-185.563429 414.47619-414.47619 414.47619S97.52381 740.912762 97.52381 512 283.087238 97.52381 512 97.52381z m0 73.142857C323.486476 170.666667 170.666667 323.486476 170.666667 512s152.81981 341.333333 341.333333 341.333333 341.333333-152.81981 341.333333-341.333333S700.513524 170.666667 512 170.666667z m36.571429 268.190476v292.571428h-73.142858V438.857143h73.142858z m0-121.904762v73.142857h-73.142858v-73.142857h73.142858z" p-id="20343"></path></svg>
                    </i>
                </button>
                {{$t('toolBar.help.about')}}
            </span>
            <div class="division-border"></div>
        </div>
    </div>
</template>