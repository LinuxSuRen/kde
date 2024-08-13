<script setup lang="ts">
import { ElNotification } from 'element-plus';
import { reactive, ref, watch } from 'vue';
import type { Config } from './types';

const tabActive = ref('Install')
watch(tabActive, () => {
    if (tabActive.value === 'Install') {
        loadInstanceStatusData()
    } else if (tabActive.value === 'Configuration') {
        loadConfig()
    }
})

const installForm = reactive({
    namespace: '',
    image: '',
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
}).then(res => { return res.json() }).then(d => {
    namespaceList.value = d
})

const imageList = ref([])
fetch('/api/images', {}).then(res => { return res.json() }).then(d => {
    imageList.value = d
})

watch(() => installForm.namespace, loadInstanceStatusData)

const install = () => {
    fetch(`/api/install`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(installForm),
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

const config = ref({} as Config)
const loadConfig = () => {
    fetch(`/api/config?namespace=${installForm.namespace}`, {}).then(res => {
        if (res.status !== 200) {
            ElNotification({
                title: 'Failed to load config',
                type: 'error',
            })
        } else {
            return res.json()
        }
    }).then(d => {
        config.value = d
    })
}
const updateConfig = () => {
    fetch(`/api/config?namespace=${installForm.namespace}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(config.value)
    }).then(res => {
        if (res.status === 200) {
            ElNotification({
                title: 'Update successful',
                type: 'success',
            })
        } else {
            ElNotification({
                title: 'Update failed',
                type: 'error',
            })
        }
    })
}
</script>

<template>
    <div>
        <el-tabs type="border-card" v-model="tabActive">
            <el-tab-pane name="Install">
                <template #label>
                    <span>Install</span>
                </template>

                <el-form :model="installForm" label-width="auto" inline>
                    <el-form-item label="Namespace" prop="namespace">
                        <el-select v-model="installForm.namespace" clearable placeholder="Select" filterable
                            style="width: 240px">
                            <el-option v-for="item in namespaceList.items" :key="item.metadata.name"
                                :label="item.metadata.name" :value="item.metadata.name" />
                        </el-select>
                    </el-form-item>
                    <el-form-item label="Image" prop="image">
                        <el-select v-model="installForm.image" style="width: 400px">
                            <el-option v-for="item in imageList" :key="item" :label="item" :value="item" />
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
                    <el-button type="primary" @click="tabActive = 'Configuration'">Next</el-button>
                </div>
            </el-tab-pane>
            <el-tab-pane name="Configuration">
                <template #label>
                    <span class="custom-tabs-label">
                        <span>Configuration</span>
                    </span>
                </template>
                <el-form label-width="auto">
                    <el-form-item label="Host">
                        <el-input v-model="config.host" />
                    </el-form-item>
                    <el-form-item label="ImagePullPolicy">
                        <el-select v-model="config.imagePullPolicy">
                            <el-option label="Always" value="Always" />
                            <el-option label="Never" value="Never" />
                            <el-option label="IfNotPresent" value="IfNotPresent" />
                        </el-select>
                    </el-form-item>
                    <el-form-item label="IngressMode">
                        <el-input v-model="config.ingressMode" />
                    </el-form-item>
                    <el-form-item label="VolumeAccessMode">
                        <el-input v-model="config.volumeAccessMode" />
                    </el-form-item>
                    <el-form-item label="VolumeMode">
                        <el-input v-model="config.volumeMode" />
                    </el-form-item>
                    <el-form-item label="StorageClassName">
                        <el-input v-model="config.storageClassName" />
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="updateConfig">Submit</el-button>
                    </el-form-item>
                </el-form>
            </el-tab-pane>
        </el-tabs>
    </div>
</template>
