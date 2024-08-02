<script setup lang="ts">
import { ref } from 'vue'

const active = ref('Install')
</script>

<template>
    <div style="height: 300px; max-width: 600px">
        <el-tabs type="border-card" v-model="active">
            <el-tab-pane name="Install">
                <template #label>
                    <span>Install</span>
                </template>

                <span>Cluster</span>
            </el-tab-pane>
            <el-tab-pane name="Configuration">
                <template #label>
                    <span class="custom-tabs-label">
                      <el-icon><calendar /></el-icon>
                      <span>Configuration</span>
                    </span>
                </template>
                Configuration
            </el-tab-pane>
        </el-tabs>
    </div>
</template>
