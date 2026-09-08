<script lang="ts" setup>
defineProps<{
  imprimiendo: boolean
  mensaje: string
}>()
</script>

<template>
  <div v-if="imprimiendo" class="estado">Imprimiendo, espere por favor...</div>
  <div v-else-if="mensaje" class="estado">{{ mensaje }}</div>
</template>

<style scoped>
.estado {
  margin-top: 1.5rem;
  color: #8a97a8;
}
</style>
