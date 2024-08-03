<script setup lang="ts">
import {ref, reactive, h, watch} from 'vue'
import { ElNotification } from 'element-plus'

const active = ref('Install')

const installForm = reactive({
    namespace: '',
})

const instanceStatusData = ref()
const loadInstanceStatusData = () => {
    fetch(`/api/instanceStatus?namespace=${installForm.namespace}`, {
        method: 'GET'
    }).then(res => res.json()).then(d => {
        instanceStatusData.value = d
    })
}
loadInstanceStatusData()

const namespaceList = ref({
    items: [{
        metadata: {
            name: ''
        }
    }]
})
fetch('/api/namespaces', {
    method: 'GET'
}).then(res => {return res.json()}).then(d => {
    namespaceList.value = d
})

watch(() => installForm.namespace, loadInstanceStatusData)

const install = () => {
    fetch(`/api/install?namespace=${installForm.namespace}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(res => {
        if (res.status === 200) {
            ElNotification({
                title: 'Installation successful',
                type: 'success',
            })
        } else {
            ElNotification({
                title: 'Installation failed',
                type: 'error',
            })
        }
    }).finally(loadInstanceStatusData)
}
const uninstall = () => {
    fetch(`/api/uninstall?namespace=${installForm.namespace}`, {
        method: 'DELETE',
    }).then(res => {
        if (res.status === 200) {
            ElNotification({
                title: 'Uninstallation successful',
                type: 'success',
            })
        } else {
            ElNotification({
                title: 'Uninstallation failed',
                type: 'error',
            })
        }
    }).finally(loadInstanceStatusData)
}
</script>

<template>
    <div style="height: 300px; max-width: 600px">
        <el-tabs type="border-card" v-model="active">
            <el-tab-pane name="Install">
                <template #label>
                    <span>Install</span>
                </template>

                <el-form
                    :model="installForm"
                    label-width="auto"
                >
                    <el-form-item label="Namespace" prop="namespace">
                        <el-select
                            v-model="installForm.namespace"
                            clearable
                            placeholder="Select"
                            style="width: 240px"
                            >
                            <el-option
                                v-for="item in namespaceList.items"
                                :key="item.metadata.name"
                                :label="item.metadata.name"
                                :value="item.metadata.name"
                            />
                            </el-select>
                    </el-form-item>
                </el-form>

                <div>
                    <div>Instance status:</div>
                    <el-table :data="instanceStatusData" style="width: 100%">
                        <el-table-column prop="component" label="Component" width="150" />
                        <el-table-column prop="name" label="Name" />
                        <el-table-column prop="status" label="status" width="150" />
                    </el-table>
                </div>

                <div>
                    <el-button type="primary" @click="install">Install</el-button>
                    <el-button type="danger" @click="uninstall">Uninstall</el-button>
                    <el-button type="primary" @click="active='Configuration'">Next</el-button>
                </div>
            </el-tab-pane>
            <el-tab-pane name="Configuration">
                <template #label>
                    <span class="custom-tabs-label">
                      <span>Configuration</span>
                    </span>
                </template>
                Configuration
            </el-tab-pane>
        </el-tabs>
    </div>
</template>
