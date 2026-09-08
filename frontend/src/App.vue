<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import {
  GetImpresoraActual,
  GetContadorActual,
  Imprimir,
  TestImpresora,
} from '../wailsjs/go/main/App'

const impresora = ref('')
const contador = ref(0)
const imprimiendo = ref(false)
const showModal = ref(false)
const cantidadSeleccionada = ref(0)
const mensajeEstado = ref('')

const opciones = [
  {valor: 100, descripcion: '100 etiquetas (50 pares)'},
  {valor: 1000, descripcion: '1.000 etiquetas (500 pares)'},
  {valor: 3000, descripcion: '3.000 etiquetas (1.500 pares)'},
  {valor: 5000, descripcion: '5.000 etiquetas (2.500 pares)'},
  {valor: 10000, descripcion: '10.000 etiquetas (5.000 pares)'},
]

async function cargarEstado() {
  impresora.value = await GetImpresoraActual()
  contador.value = await GetContadorActual()
}

function pedirConfirmacion(valor: number) {
  cantidadSeleccionada.value = valor
  showModal.value = true
}

async function confirmarImpresion() {
  showModal.value = false
  imprimiendo.value = true
  mensajeEstado.value = 'Imprimiendo, espere por favor...'
  try {
    await Imprimir(cantidadSeleccionada.value)
    mensajeEstado.value = 'Proceso finalizado correctamente.'
  } catch (e) {
    mensajeEstado.value = `Error: ${e}`
  } finally {
    imprimiendo.value = false
    cargarEstado()
  }
}

async function pruebaImpresora() {
  imprimiendo.value = true
  mensajeEstado.value = 'Enviando etiqueta de prueba...'
  try {
    await TestImpresora()
    mensajeEstado.value = 'Etiqueta de prueba enviada.'
  } catch (e) {
    mensajeEstado.value = `Error: ${e}`
  } finally {
    imprimiendo.value = false
  }
}

onMounted(cargarEstado)
</script>

<template>
  <main class="container">
    <h1>Imp-Ticket</h1>
    <p class="subtitle">Impresión de etiquetas</p>

    <div class="info">
      <p>Impresora: <strong>{{ impresora || 'No disponible' }}</strong></p>
      <p>Contador actual: <strong>Z{{ String(contador).padStart(7, '0') }}</strong></p>
    </div>

    <div class="botones">
      <button
        v-for="op in opciones"
        :key="op.valor"
        class="btn"
        :disabled="imprimiendo"
        @click="pedirConfirmacion(op.valor)"
      >
        {{ op.descripcion }}
      </button>

      <button class="btn btn-test" :disabled="imprimiendo" @click="pruebaImpresora">
        Etiqueta de prueba
      </button>
    </div>

    <div v-if="imprimiendo" class="estado">Imprimiendo, espere por favor...</div>
    <div v-else-if="mensajeEstado" class="estado">{{ mensajeEstado }}</div>

    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <h2>Confirmar impresión</h2>
        <p>
          ¿Imprimir <strong>{{ cantidadSeleccionada }}</strong> etiquetas
          ({{ cantidadSeleccionada / 2 }} pares)?
        </p>
        <div class="modal-acciones">
          <button class="btn" @click="showModal = false">Cancelar</button>
          <button class="btn btn-confirm" @click="confirmarImpresion">Confirmar</button>
        </div>
      </div>
    </div>
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

.botones {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
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

.btn-confirm {
  background: #2f855a;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal {
  background: #1b2636;
  border: 1px solid #2d3a4d;
  border-radius: 10px;
  padding: 2rem;
  max-width: 420px;
  width: 90%;
}

.modal-acciones {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}

.estado {
  margin-top: 1.5rem;
  color: #8a97a8;
}
</style>
