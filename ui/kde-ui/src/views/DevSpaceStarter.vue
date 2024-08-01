<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const phase = ref('')
const link = ref('')

const createDevSpace = async (name: string) =>  {
    fetch(`/api/devspace`, {
        method: 'POST',
        body: `{
            "apiVersion": "linuxsuren.github.io/v1alpha1",
            "kind": "DevSpace",
            "metadata": {
              "name": "${name}",
              "annotations": {
                "storageTemporary1": "a",
                "ingressMode": "path1",
                "volumeMode": "Filesystem",
                "storageClassName": "rook-cephfs"
              }
            },
            "spec": {
              "cpu": "100m",
              "memory": "100Mi",
              "host": "dev-center.jenkins-zh.cn"
            }
            }`
    }).then(res => {
        return res.json()
    }).then(res => {
        phase.value = res?.Status?.Phase
        link.value
    })
}

const loading = async () => {
    const name = 'sample'

    fetch(`/api/devspace/${name}`, {
        method: 'GET'
    }).then(res => {
        return res.json()
    }).then(res => {
        phase.value = res?.status?.phase
        link.value = res?.status?.link

        if (res?.ErrStatus?.code === 404) {
           createDevSpace(name)
        }
    }).catch(err => {
        console.log(err)
    })
}

setInterval(() => {
    loading()
}, 1000)

watch(phase, (p) => {
    if (p === 'Ready') {
        window.location.href = `http://${link.value}`
    }
})
</script>

<template>
  <div class="about">
    <h1>This is an about page{{ route.query.id }} - {{  route.hash }}</h1>
    <h2>{{ phase }}</h2>
  </div>
</template>
