<script setup lang="ts">
import WelcomeItem from './WelcomeItem.vue'
import DocumentationIcon from './icons/IconDocumentation.vue'
import ToolingIcon from './icons/IconTooling.vue'
</script>

<template>
  <WelcomeItem>
    <template #icon>
      <DocumentationIcon />
    </template>
    <template #heading>Environment</template>

    Popullar
    <a href="https://github.com/LinuxSuRen/openvscode-server-images" target="_blank" rel="noopener">development</a>
    environments.
  </WelcomeItem>

  <WelcomeItem>
    <template #icon>
      <ToolingIcon />
    </template>
    <template #heading>Services</template>

    Easy to get services as you need: Docker, Redis, MySQL, and more.
  </WelcomeItem>
</template>
