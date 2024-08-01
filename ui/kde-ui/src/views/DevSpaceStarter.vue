<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import DevSpaceCreation from '../components/dialog/DevSpaceCreation.vue'

const route = useRoute()
const phase = ref('')
const link = ref('')
const devSpaceCreationVisible = ref(false)

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
           devSpaceCreationVisible.value = true
        } else {
           devSpaceCreationVisible.value = false
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

  <DevSpaceCreation :visible="devSpaceCreationVisible" />
</template>
