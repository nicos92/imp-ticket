<script lang="ts" setup>
defineProps<{
  impresora: string
  contador: number
}>()
</script>

<template>
  <div class="info">
    <p>Impresora: <strong>{{ impresora || 'No disponible' }}</strong></p>
    <p>Contador actual: <strong>Z{{ String(contador).padStart(7, '0') }}</strong></p>
  </div>
</template>

<style scoped>
.info {
  background: rgba(255, 255, 255, 0.06);
  border-radius: 8px;
  padding: 1rem 1.5rem;
  margin-bottom: 2rem;
  text-align: left;
}

.info p {
  margin: 0.25rem 0;
}
</style>
