import { ref, readonly, onMounted, onUnmounted } from 'vue'
import {
  GetImpresoraActual,
  GetContadorActual,
  Imprimir,
  TestImpresora,
} from '../../wailsjs/go/main/App'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'

export interface OpcionImpresion {
  valor: number
  descripcion: string
}

const OPCIONES_DEFAULT: OpcionImpresion[] = [
  { valor: 10, descripcion: '10 etiquetas (5 pares)' },
  { valor: 100, descripcion: '100 etiquetas (50 pares)' },
  { valor: 1000, descripcion: '1.000 etiquetas (500 pares)' },
  { valor: 3000, descripcion: '3.000 etiquetas (1.500 pares)' },
  { valor: 5000, descripcion: '5.000 etiquetas (2.500 pares)' },
  { valor: 10000, descripcion: '10.000 etiquetas (5.000 pares)' },
]

export function useImpresion() {
  const impresora = ref('')
  const contador = ref(0)
  const imprimiendo = ref(false)
  const showModal = ref(false)
  const cantidadSeleccionada = ref(0)
  const mensajeEstado = ref('')

  const opciones = OPCIONES_DEFAULT

  async function cargarEstado() {
    const [imp, cont] = await Promise.all([
      GetImpresoraActual(),
      GetContadorActual(),
    ])
    impresora.value = imp
    contador.value = cont
  }

  onMounted(() => {
    EventsOn('print:done', () => {
      imprimiendo.value = false
      mensajeEstado.value = 'Proceso finalizado correctamente.'
      cargarEstado()
    })
    EventsOn('print:error', (msg: string) => {
      imprimiendo.value = false
      mensajeEstado.value = `Error: ${msg}`
    })
  })

  onUnmounted(() => {
    EventsOff('print:done')
    EventsOff('print:error')
  })

  function pedirConfirmacion(valor: number) {
    cantidadSeleccionada.value = valor
    showModal.value = true
  }

  function cancelarImpresion() {
    showModal.value = false
  }

  async function confirmarImpresion() {
    showModal.value = false
    imprimiendo.value = true
    mensajeEstado.value = 'Imprimiendo, espere por favor...'
    try {
      await Imprimir(cantidadSeleccionada.value)
    } catch (e) {
      imprimiendo.value = false
      mensajeEstado.value = `Error: ${e}`
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

  return {
    impresora: readonly(impresora),
    contador: readonly(contador),
    imprimiendo: readonly(imprimiendo),
    showModal,
    cantidadSeleccionada: readonly(cantidadSeleccionada),
    mensajeEstado: readonly(mensajeEstado),
    opciones,
    cargarEstado,
    pedirConfirmacion,
    cancelarImpresion,
    confirmarImpresion,
    pruebaImpresora,
  }
}
