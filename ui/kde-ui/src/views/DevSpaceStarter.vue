<script setup lang="ts">
import { ref, watch } from 'vue'
import DevSpaceCreation from '../components/dialog/DevSpaceCreation.vue'

interface DevSpace {
    status: {
        phase: string
        deployStatus: string
        link: string
    }
}

const devSpace = ref({} as DevSpace)
const devSpaceCreationVisible = ref(false)

const loading = async () => {
    const name = 'sample'

    fetch(`/api/devspace/${name}`, {
        method: 'GET'
    }).then(res => {
        return res.json()
    }).then(res => {
        devSpace.value = res
        devSpaceCreationVisible.value = res?.ErrStatus?.code === 404
    }).catch(err => {
        console.log(err)
    })
}

setInterval(() => {
    loading()
}, 2000)

watch(() => devSpace.value?.status?.deployStatus, (p) => {
    if (p === 'Ready') {
        window.location.href = `http://${devSpace.value?.status?.link}`
    }
})
</script>

<template>
    <div class="about">
        <h2>Request Status: {{ devSpace.status?.phase }}</h2>
        <h2>Env Status: {{ devSpace.status?.deployStatus }}</h2>
    </div>

    <DevSpaceCreation :visible="devSpaceCreationVisible"/>
</template>
