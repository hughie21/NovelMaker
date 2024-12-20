<script setup>
/* 
@Author: Hughie
@CreateTime: 2024-10-9
@LastEditors: Hughie
@LastEditTime: 2024-12-20
@Description: Wandering guides that provide users with guidance on how to use the software
*/
import { visio } from '../../assets/js/globals';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const handleChange = (val) => {
    console.log(val)
    if(val == 5) {
        document.getElementById('tab-1').click();
    }
}
</script>

<template>
    <el-tour v-model="visio.tourVsisible" type="default" @change="handleChange" :mask="{
        style: {
          boxShadow: 'inset 0 0 15px #333',
        },
        color: 'var(--mask-color)',
      }">
        <el-tour-step
        target="#editor"
        :title="t('tour.step1.title')",
        :description="t('tour.step1.content')"
        :next-button-props="{
            children: t('tour.next'),
        }"
        :prev-button-props="{
            children: t('tour.previous'),
        }"
        />
        <el-tour-step
        target="#btn-new"
        :title="t('tour.step2.title')",
        :description="t('tour.step2.content')"
        :next-button-props="{
            children: t('tour.next'),
        }"
        :prev-button-props="{
            children: t('tour.previous'),
        }"
        />
        <el-tour-step
        target="#btn-open"
        :title="t('tour.step3.title')",
        :description="t('tour.step3.content')"
        :next-button-props="{
            children: t('tour.next'),
        }"
        :prev-button-props="{
            children: t('tour.previous'),
        }"
        />
        <el-tour-step
        target="#btn-save"
        :title="t('tour.step4.title')",
        :description="t('tour.step4.content')"
        :next-button-props="{
            children: t('tour.next'),
        }"
        :prev-button-props="{
            children: t('tour.previous'),
        }"
        />
        <el-tour-step
        target="#btn-file-info"
        :title="t('tour.step5.title')",
        :description="t('tour.step5.content')"
        :next-button-props="{
            children: t('tour.next'),
        }"
        :prev-button-props="{
            children: t('tour.previous'),
        }"
        />
        <el-tour-step
        target="#tab-1"
        :title="t('tour.step6.title')",
        :description="t('tour.step6.content')"
        :next-button-props="{
            children: t('tour.next'),
        }"
        :prev-button-props="{
            children: t('tour.previous'),
        }"
        />
        <el-tour-step
        target="#btn-media"
        :title="t('tour.step7.title')",
        :description="t('tour.step7.content')"
        :next-button-props="{
            children: t('tour.next'),
        }"
        :prev-button-props="{
            children: t('tour.previous'),
        }"
        />
        <el-tour-step
        target="#btn-setting"
        :title="t('tour.step8.title')",
        :description="t('tour.step8.content')"
        :next-button-props="{
            children: t('tour.next'),
        }"
        :prev-button-props="{
            children: t('tour.previous'),
        }"
        />
        <el-tour-step
        :title="t('tour.step9.title')",
        :description="t('tour.step9.content')"
        :next-button-props="{
            children: t('tour.next'),
        }"
        :prev-button-props="{
            children: t('tour.previous'),
        }"
        />
    </el-tour>
</template>