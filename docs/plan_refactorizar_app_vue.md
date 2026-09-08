# Plan: Refactorizar App.vue

## Plan: Refactorizar App.vue en componentes reutilizables

### Análisis

`App.vue` actualmente es un "mega componente" que mezcla:

- **Lógica de negocio** (llamadas a Wails: impresora, contador, imprimir)
- **Estado UI** (modal, imprimiendo, mensajes)
- **Presentación** (panel info, botones, modal, mensajes)

### Componentes a crear

| Componente | Responsabilidad | Props | Emits |
| --- | --- | --- | --- |
| `PrinterInfo.vue` | Muestra impresora y contador | `impresora: string`, `contador: number` | — |
| `PrintOptions.vue` | Botones de cantidad + prueba | `opciones: Array`, `deshabilitado: boolean` | `seleccionar(valor)`, `probar()` |
| `ConfirmModal.vue` | Modal de confirmación | `visible: boolean`, `cantidad: number` | `confirmar()`, `cancelar()` |
| `StatusMessage.vue` | Muestra estado/error | `imprimiendo: boolean`, `mensaje: string` | — |

### Composable

`composables/useImpresion.ts` — encapsula toda la lógica:

- `impresora`, `contador` (refs)
- `imprimiendo`, `mensajeEstado` (refs)
- `cargarEstado()`, `confirmarImpresion()`, `pruebaImpresora()`
- `cantidadSeleccionada`, `showModal` + `pedirConfirmacion()`

### Resultado

`App.vue` queda como superficie de composición delgada:

```vue
<script setup>
import { useImpresion } from './composables/useImpresion'
// solo compone los componentes
</script>
```

### Estructura final

```
src/
├── composables/
│   └── useImpresion.ts
├── components/
│   ├── PrinterInfo.vue
│   ├── PrintOptions.vue
│   ├── ConfirmModal.vue
│   └── StatusMessage.vue
└── App.vue
```
