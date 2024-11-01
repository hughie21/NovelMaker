<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-18
@LastEditors: Hughie
@LastEditTime: 2024-11-1
*/
import { visio } from '../../assets/js/globals';
import "../../assets/css/tab.css"
</script>

<template>
    <div class="func_btn-list">
        <div class="func_btn-container list_row">
            <span class="btn-group" style="width:40px; ">
                <button class="el-button btn func_btn-big" @click="visio.helpVisible = true">
                    <i class="el-icon">
                        <svg t="1717404400946" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="20342" width="200" height="200"><path d="M512 97.52381c228.912762 0 414.47619 185.563429 414.47619 414.47619s-185.563429 414.47619-414.47619 414.47619S97.52381 740.912762 97.52381 512 283.087238 97.52381 512 97.52381z m0 73.142857C323.486476 170.666667 170.666667 323.486476 170.666667 512s152.81981 341.333333 341.333333 341.333333 341.333333-152.81981 341.333333-341.333333S700.513524 170.666667 512 170.666667z m36.571429 268.190476v292.571428h-73.142858V438.857143h73.142858z m0-121.904762v73.142857h-73.142858v-73.142857h73.142858z" p-id="20343"></path></svg>
                    </i>
                </button>
                {{$t('toolBar.help.help')}}
            </span>
        </div>
    </div>
</template>