<script lang="ts" setup>
import { onMounted } from 'vue'
import { useImpresion } from './composables/useImpresion'
import PrinterInfo from './components/PrinterInfo.vue'
import PrintOptions from './components/PrintOptions.vue'
import ConfirmModal from './components/ConfirmModal.vue'
import StatusMessage from './components/StatusMessage.vue'

const {
  impresora,
  contador,
  imprimiendo,
  showModal,
  cantidadSeleccionada,
  mensajeEstado,
  opciones,
  cargarEstado,
  pedirConfirmacion,
  cancelarImpresion,
  confirmarImpresion,
  pruebaImpresora,
} = useImpresion()

onMounted(cargarEstado)
</script>

<template>
  <main class="container">
    <h1>Pre-Etiquetas</h1>

    <PrinterInfo :impresora="impresora" :contador="contador" />

    <PrintOptions
      :opciones="opciones"
      :deshabilitado="imprimiendo"
      @seleccionar="pedirConfirmacion"
      @probar="pruebaImpresora"
    />

    <StatusMessage :imprimiendo="imprimiendo" :mensaje="mensajeEstado" />

    <ConfirmModal
      :visible="showModal"
      :cantidad="cantidadSeleccionada"
      @confirmar="confirmarImpresion"
      @cancelar="cancelarImpresion"
    />
  </main>
</template>

<style scoped>
.container {
  max-width: 640px;
  margin: 0 auto;
  padding: 2rem;
}

h1 {
  margin: 0;
}

.subtitle {
  color: #8a97a8;
  margin: 0 0 1.5rem;
}
</style>
