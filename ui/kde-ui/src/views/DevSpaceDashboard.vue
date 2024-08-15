<template>
    <el-button type="primary" @click="devSpaceCreationVisible = true">New</el-button>
    <el-button type="primary" @click="loadData">Reload</el-button>
    <el-table :data="tableData" style="width: 100%" :row-class-name="tableRowClassName" v-loading="dashboardLoading">
        <el-table-column prop="metadata.name" label="Name" width="100">
            <template #default="scope">
                <el-button link type="primary" size="small"
                    @click.prevent="router.push({ path: '/dev', query: { name: scope.row.metadata.name } })">
                    {{ scope.row.metadata.name }}
                </el-button>
            </template>
        </el-table-column>
        <el-table-column prop="status.link" label="Address" />
        <el-table-column prop="status.deployStatus" label="Deployment" width="120" />
        <el-table-column prop="status.phase" label="Status" width="80" />
        <el-table-column fixed="right" label="Operations" min-width="80">
            <template #default="scope">
                <el-button link type="primary" size="small" @click.prevent="deleteRow(scope.$index)">
                    Remove
                </el-button>
                <el-button link type="primary" size="small"
                    @click.prevent="editDevSpace(scope.row.metadata.namespace, scope.row.metadata.name)">
                    Edit
                </el-button>
                <el-button link type="primary" size="small"
                    @click.prevent="restartDevSpace(scope.row.metadata.namespace, scope.row.metadata.name)">
                    Restart
                </el-button>
            </template>
        </el-table-column>
    </el-table>

    <DevSpaceCreation :visible="devSpaceCreationVisible" @closed="devSpaceCreationVisible = false"
        @created="loadData()" />
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import DevSpaceCreation from '../components/dialog/DevSpaceCreation.vue';
import type { DevSpace } from './types';

const router = useRouter();

const tableRowClassName = ({
    row,
    rowIndex,
}: {
    row: DevSpace
    rowIndex: number
}) => {
    if (rowIndex === 1) {
        return 'warning-row'
    } else if (row?.status?.phase === 'Ready') {
        return 'success-row'
    }
    return ''
}

const dashboardLoading = ref(true)
const tableData = ref([] as DevSpace[])
const loadData = () => {
    dashboardLoading.value = true
    fetch('/api/devspace', {
        method: 'GET'
    }).then(res => {
        return res.json()
    }).then(d => {
        tableData.value = d.items
    }).finally(() => {
        dashboardLoading.value = false
    })
}
loadData()

const deleteRow = (index: number) => {
    fetch(`/api/devspace/${tableData.value[index].metadata.name}`, {
        method: 'DELETE'
    }).finally(() => {
        loadData()
    })
}
const editDevSpace = (namespace: string, name: string) => {
    router.push({
        path: `/devspace/${namespace}/${name}`,
    })
}
const restartDevSpace = (namespace: string, name: string) => {
    fetch(`/api/devspace/${name}/restart?namespace=${namespace}`, {
        method: 'PUT'
    }).finally(() => {
        loadData()
    })
}

const devSpaceCreationVisible = ref(false)
</script>

<style>
.el-table .warning-row {
    --el-table-tr-bg-color: var(--el-color-warning-light-9);
}

.el-table .success-row {
    --el-table-tr-bg-color: var(--el-color-success-light-9);
}
</style>
