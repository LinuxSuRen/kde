<script setup lang="ts">
import { ElNotification } from 'element-plus';
import { reactive, ref, watch } from 'vue';
import type { Config, Cluster } from './types';

const tabActive = ref('Install')
watch(tabActive, () => {
    if (tabActive.value === 'Install') {
        loadInstanceStatusData()
    } else if (tabActive.value === 'Configuration') {
        loadConfig()
    } else if (tabActive.value === 'Cluster') {
        loadClusterInfo()
    }
})

const installForm = reactive({
    namespace: '',
    image: '',
})

const instanceStatusData = ref({
    status: [],
    namespace: '',
})
const installStatusLoading = ref(true)

const loadInstanceStatusData = () => {
    installStatusLoading.value = true
    fetch(`/api/instanceStatus?namespace=${installForm.namespace}`, {
        method: 'GET'
    }).then(res => res.json()).then(d => {
        instanceStatusData.value = d
        if (installForm.namespace === '') {
            installForm.namespace = d.namespace
        }
    }).finally(() => {
        installStatusLoading.value = false
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

const configLoading = ref(true)
const config = ref({
    languages: [{
        name: '',
        image: '',
    }]
} as Config)
const loadConfig = () => {
    configLoading.value = true
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
        if (!d.languages) {
            d.languages = [{
                name: '',
                image: '',
            }]
        }
        config.value = d
    }).finally(() => {
        configLoading.value = false
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

const clusterInfoLoading = ref(true)
const clusterInfo = ref({} as Cluster)
const loadClusterInfo = () => {
    clusterInfoLoading.value = true
    fetch(`/api/cluster/info`, {
        method: 'GET'
    }).then(res => res.json()).then(d => {
        clusterInfo.value = d
    }).finally(() => {
        clusterInfoLoading.value = false
    })
}
</script>

<template>
    <div>
        <el-tabs type="border-card" v-model="tabActive">
            <el-tab-pane name="Install" v-loading="installStatusLoading">
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
                    <el-table :data="instanceStatusData.status" style="width: 100%">
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
            <el-tab-pane name="Configuration" v-loading="configLoading">
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
                        <el-select v-model="config.volumeAccessMode">
                            <el-option label="ReadWriteMany" value="ReadWriteMany" />
                            <el-option label="ReadWriteOnce" value="ReadWriteOnce" />
                        </el-select>
                    </el-form-item>
                    <el-form-item label="VolumeMode">
                        <el-select v-model="config.volumeMode">
                            <el-option label="Filesystem" value="Filesystem" />
                        </el-select>
                    </el-form-item>
                    <el-form-item label="StorageClassName">
                        <el-input v-model="config.storageClassName" />
                    </el-form-item>
                    <el-form-item label="Custom Languages">
                        <div v-for="(lan, index) in config.languages" :key="lan.name">
                            Name:<el-input v-model="lan.name" style="width: 200px;" />
                            Image:<el-input v-model="lan.image" style="width: 400px;" />
                            <el-icon v-if="index == config.languages.length - 1">
                                <Plus @click="config.languages.push({ name: '', image: '' })" />
                            </el-icon>
                        </div>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="updateConfig">Submit</el-button>
                    </el-form-item>
                </el-form>
            </el-tab-pane>
            <el-tab-pane name="Cluster">
                <template #label>
                    <span class="custom-tabs-label">
                        <span>Cluster</span>
                    </span>
                </template>
                <el-table :data="clusterInfo.nodes" style="width: 100%">
                    <el-table-column prop="name" label="Name" />
                    <el-table-column prop="os" label="OS">
                        <template #default="scope">
                            {{ scope.row.os }}/{{ scope.row.osImage }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="arch" label="Arch" />
                    <el-table-column prop="containerRuntime" label="ContainerRuntime" />
                    <el-table-column prop="allocatable.cpu" label="CPU" />
                    <el-table-column prop="allocatable.memory" label="Memory" />
                    <el-table-column prop="allocatable.pods" label="Pods" />
                    <el-table-column prop="images" label="Images">
                        <template #default="scope">
                            {{ scope.row.images.length }}
                        </template>
                    </el-table-column>
                </el-table>
            </el-tab-pane>
        </el-tabs>
    </div>
</template>
