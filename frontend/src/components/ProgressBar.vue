<script lang="ts" setup>
defineProps<{
  actual: number
  total: number
  visible: boolean
}>()
</script>

<template>
  <div v-if="visible" class="progress-container">
    <div class="progress-header">
      <span class="progress-label">Imprimiendo pares...</span>
      <span class="progress-text">{{ actual.toLocaleString('es-AR') }} / {{ total.toLocaleString('es-AR') }}</span>
    </div>
    <div class="progress-track">
      <div class="progress-fill" :style="{ width: `${total ? (actual / total) * 100 : 0}%` }" />
    </div>
  </div>
</template>

<style scoped>
.progress-container {
  margin-top: 1.5rem;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.progress-label {
  color: #8a97a8;
  font-size: 0.9rem;
}

.progress-text {
  color: #e2e8f0;
  font-size: 0.9rem;
  font-variant-numeric: tabular-nums;
}

.progress-track {
  height: 8px;
  background: #2d3a4d;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #3b82f6;
  border-radius: 4px;
  transition: width 0.3s ease;
}
</style>
