<template>
    <el-form>
        <table>
            <tr>
                <td>Namespace/Name:</td>
                <td>{{ devspace.metadata.namespace }}/{{ devspace.metadata.name }}</td>
            </tr>
            <tr>
                <td>CPU</td>
                <td>
                    <el-input v-model="devspace.spec.cpu" />
                </td>
            </tr>
            <tr>
                <td>Memory</td>
                <td>
                    <el-input v-model="devspace.spec.memory" />
                </td>
            </tr>
            <tr>
                <td>
                    Services:
                </td>
                <td>
                    <el-checkbox v-model="devspace.spec.services.docker.enabled">
                        Docker
                    </el-checkbox>
                    <el-checkbox v-model="devspace.spec.services.mysql.enabled">
                        MySQL
                    </el-checkbox>
                    <el-checkbox v-model="devspace.spec.services.redis.enabled">
                        Redis
                    </el-checkbox>
                </td>
            </tr>
        </table>
        <el-button type="primary" @click="submitForm">Submit</el-button>
    </el-form>
</template>

<script lang="ts" setup>
import { ElMessage } from 'element-plus';
import { ref } from 'vue';
import { useRoute } from 'vue-router';
import type { DevSpace } from './types';

const route = useRoute();
const devspace = ref({} as DevSpace)

fetch(`/api/devspace/${route.params.name}?namespace=${route.params.namespace}`, {}).
    then(res => res.json()).
    then((d) => {
        if (!d.spec.services.docker) {
            d.spec.services.docker = {
                enabled: false
            }
        }
        if (!d.spec.services.mysql) {
            d.spec.services.mysql = {
                enabled: false
            }
        }
        if (!d.spec.services.redis) {
            d.spec.services.redis = {
                enabled: false
            }
        }
        devspace.value = d
    }).catch((e) => {
        ElMessage({
            message: e,
            type: 'error',
            plain: true,
        })
    })

const submitForm = () => {
    console.log(devspace.value.spec)
    fetch(`/api/devspace/${route.params.name}?namespace=${route.params.namespace}`, {
        method: 'PUT',
        body: JSON.stringify(devspace.value)
    }).then(res => {
        if (res.ok) {
            ElMessage({
                message: 'Updated sucessfully',
                type: 'success',
                plain: true,
            })
        }
    }).catch((e) => {
        ElMessage({
            message: e,
            type: 'error',
            plain: true,
        })
    })
}
</script>
