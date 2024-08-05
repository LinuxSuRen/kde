// This is a dialog for devspace creation
<script setup lang="ts">
import type { FormInstance } from 'element-plus';
import { reactive, ref } from 'vue';

const props = defineProps({
    visible: Boolean,
})

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
}

const createDevSpace = async (devSpace: DevSpace) =>  {
    fetch(`/api/devspace`, {
        method: 'POST',
        body: `{
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
              "cpu": "100m",
              "memory": "100Mi",
              "host": "dev-center.jenkins-zh.cn",
              "image": "${devSpace.image}"
            }
            }`
    }).then(res => {
        return res.json()
    })
}

const ruleFormRef = ref<FormInstance>()
const ruleForm = reactive({
  name: "sample"
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
}

const resetForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return
  formEl.resetFields()
}
</script>

<template>
  <el-dialog
    :modelValue="visible"
    title="Create DevSpace"
    width="500"
    destroy-on-close
    center
  >
    <el-form
      ref="ruleFormRef"
      style="max-width: 600px"
      :model="ruleForm"
      status-icon
      label-width="auto"
    >
      <el-form-item label="Name" prop="name">
        <el-input v-model="ruleForm.name" />
      </el-form-item>
      <el-form-item label="Language" prop="image">
        <el-select
          v-model="ruleForm.image"
          clearable
          placeholder="Select"
          style="width: 240px"
        >
          <el-option
            v-for="item in devLanguages"
            :key="item.name"
            :label="item.name"
            :value="item.image"
          />
        </el-select>
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
