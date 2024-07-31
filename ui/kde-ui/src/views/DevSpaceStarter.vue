<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const phase = ref('')
const link = ref('')

const loading = () => {
    fetch('/api/devspace/devspace-sample', {
        method: 'GET'
    }).then(res => {
        return res.json()
    }).then(res => {
        phase.value = res.Status.Phase
        link.value = res.Status.Link
    })
}

setInterval(() => {
    loading()
}, 1000)

watch(phase, (p) => {
    if (p === 'Ready') {
        window.location.href = '/'
    }
})
</script>

<template>
  <div class="about">
    <h1>This is an about page{{ route.query.id }} - {{  route.hash }}</h1>
    <h2>{{ phase }}</h2>
  </div>
</template>
