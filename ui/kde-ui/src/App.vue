<script setup lang="ts">
import { ref } from 'vue';
import { RouterView } from 'vue-router';

const defaultMenu = ref(window.location.pathname)
</script>

<template>
    <el-menu mode="horizontal" :default-active="defaultMenu" router>
        <el-menu-item index="/">Home</el-menu-item>
        <el-menu-item index="/dev">DevSpace</el-menu-item>
        <el-menu-item index="/dashboard">Dashboard</el-menu-item>
        <el-menu-item index="/install">Installation</el-menu-item>
    </el-menu>

    <RouterView />
</template>
