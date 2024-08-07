<script setup lang="ts">
import { ElMessageBox } from 'element-plus';
import { ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import type { DevSpace } from './types';

const router = useRouter();
const route = useRoute();

interface DevSpaceStatus {
    name: string
    status: string
}

const devSpaceStatus = ref([] as DevSpaceStatus[])
const devSpace = ref({} as DevSpace)

const loading = async () => {
    if (!route.query.name) {
        router.push({
            path: '/dashboard',
        })
        return
    }

    fetch(`/api/devspace/${route.query.name}`, {
        method: 'GET'
    }).then(res => {
        return res.json()
    }).then(res => {
        devSpace.value = res

        devSpaceStatus.value = [{
            name: 'Request',
            status: res?.status?.phase
        }, {
            name: 'Deployment',
            status: res?.status?.deployStatus
        }]
    }).catch(err => {
        console.log(err)
    })
}

const autoReload = setInterval(() => {
    loading()
}, 2000)

const openIDE = () => {
    // open in a new window
    window.open(`http://${devSpace.value?.status?.link}`)
}

watch(() => devSpace.value?.status?.deployStatus, (p) => {
    if (p === 'Running') {
        ElMessageBox.confirm(`Are you confirm to open your IDE?`).then(openIDE)
        clearInterval(autoReload)
    }
})
</script>

<template>
    <div class="about">
        <el-button type="primary" @click="openIDE" v-if="devSpace?.status?.deployStatus === 'Running'">Open</el-button>

        <el-table :data="devSpaceStatus">
            <el-table-column prop="name" label="Name" />
            <el-table-column prop="status" label="Status" />
        </el-table>
    </div>
</template>
