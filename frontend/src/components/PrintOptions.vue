<script lang="ts" setup>
import type { OpcionImpresion } from '../composables/useImpresion'

defineProps<{
  opciones: OpcionImpresion[]
  deshabilitado: boolean
}>()

const emit = defineEmits<{
  seleccionar: [valor: number]
  probar: []
}>()
</script>

<template>
  <div class="botones">
    <button
      v-for="op in opciones"
      :key="op.valor"
      class="btn"
      :disabled="deshabilitado"
      @click="emit('seleccionar', op.valor)"
    >
      {{ op.descripcion }}
    </button>

    <button
      class="btn btn-test"
      :disabled="deshabilitado"
      @click="emit('probar')"
    >
      Etiqueta de prueba
    </button>
  </div>
</template>

<style scoped>
.botones {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.btn {
  background: #2d3a4d;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 0.9rem 1.5rem;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}

.btn:hover:not(:disabled) {
  background: #3b4c63;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-test {
  background: #4a5568;
}
</style>
