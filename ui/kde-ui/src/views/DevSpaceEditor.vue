<template>
    <el-form>
        <table style="width: 100%;">
            <tr>
                <td style="width: 50%;">Namespace/Name:</td>
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
                <td>Storage</td>
                <td>
                    <el-input v-model="devspace.spec.storage" />
                </td>
            </tr>
            <tr>
                <td>Host</td>
                <td>
                    <el-input v-model="devspace.spec.host" />
                </td>
            </tr>
            <tr>
                <td>Image</td>
                <td>
                    <el-input v-model="devspace.spec.image" />
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
                    <el-checkbox v-model="devspace.spec.services.postgres.enabled">
                        Postgres
                    </el-checkbox>
                    <el-checkbox v-model="devspace.spec.services.tDEngine.enabled">
                        TDEngine
                    </el-checkbox>
                    <el-checkbox v-model="devspace.spec.services.rabbitMQ.enabled">
                        RabbitMQ
                    </el-checkbox>
                </td>
            </tr>
            <tr v-if="devspace.spec.services.mysql.enabled">
                <td>
                    MySQL Password: <el-input v-model="devspace.spec.services.mysql.password" type="password"
                        style="width: 240px" />
                </td>
                <td>
                    MySQL Database: <el-input v-model="devspace.spec.services.mysql.database" style="width: 240px" />
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
import { NewEmptyDevSpace } from './types';

const route = useRoute();
const devspace = ref(NewEmptyDevSpace())

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
        if (!d.spec.services.postgres) {
            d.spec.services.postgres = {
                enabled: false
            }
        }
        if (!d.spec.services.tDEngine) {
            d.spec.services.tDEngine = {
                enabled: false
            }
        }
        if (!d.spec.services.rabbitMQ) {
            d.spec.services.rabbitMQ = {
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
