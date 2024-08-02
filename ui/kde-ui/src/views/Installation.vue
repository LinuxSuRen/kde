<script setup lang="ts">
import { ref, reactive } from 'vue'

const active = ref('Install')

const instanceStatusData = ref()
fetch('/api/instanceStatus', {
    method: 'GET'
}).then(res => res.json()).then(d => {
    instanceStatusData.value = d
})

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

const installForm = reactive({
  namespace: '',
})
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
                        <el-table-column prop="component" label="Component" />
                        <el-table-column prop="status" label="status" width="180" />
                    </el-table>
                </div>

                <div>
                    <el-button type="primary">Install</el-button>
                    <el-button type="primary">Next</el-button>
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
