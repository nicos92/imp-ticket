# Plan de Migración

## Plan de Migración: `imp-ticket` Console → Wails

### 1. Backend - `app.go`

Reescribir `app.go` para absorber toda la lógica de negocio de `backend/main.go`:

**Estructura `App`:**

```
App {
    ctx         context.Context
    configPath  string    // ruta al archivo de persistencia
}
```

**En `startup()`:** inicializar `configPath` con `os.UserConfigDir()` + `Nicolas-Sandoval/default-imp-ticket/config` y crear directorio si no existe.

**Métodos expuestos al frontend:**

| Método | Firma | Descripción |
| --- | --- | --- |
| `GetImpresoraActual` | `() string` | Retorna nombre de impresora por defecto |
| `GetContadorActual` | `() int` | Retorna último número guardado |
| `Imprimir` | `(cantidad int) error` | Ejecuta `ImprimirPares` sincrónicamente |
| `TestImpresora` | `() error` | Imprime una etiqueta ZPL de prueba |

**Funciones internas a mover sin cambios:**

- `ObtenerImpresoraPredeterminada()`
- `ImprimirTextoPlano()`
- `GenerarZPL()`
- `LeerUltimoNumero()` / `GuardarUltimoNumero()`
- `ImprimirPares()`
- `confGuardarUltimoNumero()`

**Eliminar:** `leerTexto()` (ya no se necesita input por consola).

---

### 2. Frontend - `App.vue`

Reemplazar el `HelloWorld` con una interfaz simple:

```
┌─────────────────────────────────────┐
│  Imp-Ticket - Impresión de Etiquetas│
├─────────────────────────────────────┤
│  Impresora: HP LaserJet Pro         │
│  Contador actual: Z0000452          │
├─────────────────────────────────────┤
│  [100] [1.000] [3.000] [5.000] [10K]│
│  [Etiqueta de Prueba]               │
├─────────────────────────────────────┤
│  Estado: Imprimiendo... 500/1000    │
└─────────────────────────────────────┘
```

**Estados reactivos:**

- `impresora: string` - nombre de impresora
- `contador: number` - contador actual
- `imprimiendo: boolean` - flag de proceso activo
- `showModal: boolean` - controla modal de confirmación
- `cantidadSeleccionada: number` - cantidad elegida
- `mensajeEstado: string` - texto de estado/error

**Modal de confirmación:**

- Muestra: "¿Imprimir X etiquetas? (Y pares)"
- Botones: Confirmar / Cancelar
- Al confirmar → llama `Imprimir(cantidad)`, muestra spinner/loading, al terminar muestra resultado y actualiza contador

---

### 3. Archivos a eliminar/modificar

| Archivo | Acción |
| --- | --- |
| `backend/main.go` | Eliminar |
| `backend/inputs/` | Eliminar directorio |
| `backend/go.mod` | Eliminar |
| `backend/README.md` | Eliminar |
| `backend/.gitignore` | Eliminar |
| `backend/.editorconfig` | Eliminar |
| `backend/.gitattributes` | Eliminar |
| `backend/LICENSE` | Eliminar |
| `app.go` | Reescribir con lógica de negocio |
| `main.go` | Sin cambios |
| `go.mod` | Sin cambios (ya tiene wails v2) |
| `frontend/src/App.vue` | Reescribir con UI de impresión |
| `frontend/src/components/HelloWorld.vue` | Eliminar |

---

### 4. Dependencias

No se necesitan dependencias adicionales en Go (ya usa `syscall`, `unsafe`, `os`, `fmt` que son stdlib). No se necesitan dependencias npm extras (Vue 3 vanilla es suficiente para la UI simple).

---
