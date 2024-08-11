// This is a dialog for devspace creation
<script setup lang="ts">
import type { FormInstance } from 'element-plus';
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const props = defineProps({
    visible: Boolean,
})

const emit = defineEmits(['created', 'closed'])

interface DevLanguage {
    name: string
    image: string
}

const devLanguages = ref([] as DevLanguage[])
fetch('/api/languages', {}).then(res => res.json()).then(data => {
    devLanguages.value = data
})

interface DevSpace {
    name: string
    image: string
    cpu: string
    memory: string
    docker: boolean
    mysql: boolean
    redis: boolean
    postgres: boolean
    tdEngine: boolean
    rabbitMQ: boolean
}

const createDevSpace = async (devSpace: DevSpace) => {
    const body = `{
              "apiVersion": "linuxsuren.github.io/v1alpha1",
              "kind": "DevSpace",
              "metadata": {
                "name": "${devSpace.name}",
                "annotations": {
                  "storageTemporary1": "a",
                  "ingressMode": "path1",
                  "volumeMode": "Filesystem",
                  "storageClassName1": "rook-cephfs"
                }
              },
              "spec": {
                "cpu": "${devSpace.cpu}",
                "memory": "${devSpace.memory}",
                "host": "devspace.jenkins-zh.cn",
                "image": "${devSpace.image}",
                "services": {}
              }
            }`

    const bodyObj = JSON.parse(body)
    if (devSpace.docker) {
        bodyObj.spec.services.docker = {
            enabled: true
        }
    }
    if (devSpace.mysql) {
        bodyObj.spec.services.mysql = {
            enabled: true
        }
    }
    if (devSpace.redis) {
        bodyObj.spec.services.redis = {
            enabled: true
        }
    }
    if (devSpace.postgres) {
        bodyObj.spec.services.postgres = {
            enabled: true
        }
    }
    if (devSpace.tdEngine) {
        bodyObj.spec.services.tdEngine = {
            enabled: true
        }
    }
    if (devSpace.rabbitMQ) {
        bodyObj.spec.services.rabbitMQ = {
            enabled: true
        }
    }

    fetch(`/api/devspace`, {
        method: 'POST',
        body: JSON.stringify(bodyObj),
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(res => {
        router.push(`/dev?name=${ruleForm.name}`)
        return res.json()
    }).finally(() => {
        emit('created')
    })
}

const ruleFormRef = ref<FormInstance>()
const ruleForm = reactive({
    cpu: "2",
    memory: "4Gi"
} as DevSpace)

const submitForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.validate((valid) => {
        if (valid) {
            createDevSpace(ruleForm)
        } else {
            console.log('error submit!')
        }
    })
    emit('closed')
}

const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
}

const cancelDialog = () => {
    emit('closed')
}
</script>

<template>
    <el-dialog :modelValue="visible" title="Create DevSpace" width="500" :before-close="cancelDialog">
        <el-form ref="ruleFormRef" style="max-width: 600px" :model="ruleForm" status-icon label-width="auto">
            <el-form-item label="Name" prop="name">
                <el-input v-model="ruleForm.name" />
            </el-form-item>
            <el-form-item label="Language" prop="image">
                <el-select v-model="ruleForm.image" clearable placeholder="Select" style="width: 240px">
                    <el-option v-for="item in devLanguages" :key="item.name" :label="item.name" :value="item.image" />
                </el-select>
            </el-form-item>
            <el-form-item label="CPU" prop="cpu">
                <el-input v-model="ruleForm.cpu" />
            </el-form-item>
            <el-form-item label="Memory" prop="memory">
                <el-input v-model="ruleForm.memory" />
            </el-form-item>
            <el-form-item label="Service" prop="docker">
                <el-checkbox v-model="ruleForm.docker">
                    Docker
                </el-checkbox>
                <el-checkbox v-model="ruleForm.mysql">
                    MySQL
                </el-checkbox>
                <el-checkbox v-model="ruleForm.redis">
                    Redis
                </el-checkbox>
                <el-checkbox v-model="ruleForm.postgres">
                    Postgres
                </el-checkbox>
                <el-checkbox v-model="ruleForm.tdEngine">
                    TDEngine
                </el-checkbox>
                <el-checkbox v-model="ruleForm.rabbitMQ">
                    RabbitMQ
                </el-checkbox>
            </el-form-item>

            <el-form-item>
                <el-button type="primary" @click="submitForm(ruleFormRef)">
                    Submit
                </el-button>
                <el-button @click="resetForm(ruleFormRef)">Reset</el-button>
            </el-form-item>
        </el-form>

    </el-dialog>
</template>
