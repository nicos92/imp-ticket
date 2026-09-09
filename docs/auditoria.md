
# AUDITORÍA: imp-ticket — Go + Wails v2 + Vue 3

## 1. [P0 - CRÍTICO] Backend Go / Wails

### 1.1 Bloqueo del hilo principal — `app.go:38-58`

**`Imprimir()` y `TestImpresora()` congelan la UI.** Wails ejecuta los métodos bindeados en la misma goroutine que el frontend. `imprimirPares()` hace un loop de `cantidadCodigos / 2` iteraciones, cada una ejecutando 3 operaciones pesadas: lectura/escritura de archivo (`counter`), generación ZPL, y syscalls `winspool.drv`. Para 10.000 etiquetas = 5.000 iteraciones de I/O de disco + impresión — la ventana queda congelada completamente.

**Solución:** envolver en goroutine + `runtime.EventsEmit` para progreso:

```go
func (a *App) Imprimir(cantidad int) error {
    go a.imprimirPares(cantidad)  // Background
    return nil
}
```

### 1.2 Goroutine sin contexto — `app.go:38`

`imprimirPares` no recibe `context.Context`. No hay forma de cancelar una impresión de 10.000 etiquetas una vez iniciada. Wails provee `a.ctx` que ya está en el struct pero no se usa.

### 1.3 Errores silenciados

| Ubicación | Problema |
| --- | --- |
| `app.go:18` (`startup`) | `os.UserConfigDir()` falla → `return` silencioso, `configPath` queda `""` |
| `app.go:21` (`startup`) | `os.MkdirAll()` — error ignorado completamente |
| `printer.go:89-90` (`iniciarDocumento`) | `syscall.UTF16PtrFromString()` — error ignorado con `_`, puede causar panic si el string es inválido |
| `printer.go:110` (`iniciarPagina`) | `procStartPagePrinter.Call` — 3er valor de retorno (err) ignorado |

### 1.4 Resource leak en printer — `printer.go:50-67`

`defer procClosePrinter.Call(hPrinter)` usa `uintptr` como handle. Si `iniciarDocumento` o `iniciarPagina` fallan, los defers se ejecutan en orden LIFO correcto, pero el `hPrinter` como `uintptr` no tiene garantía de que el runtime no lo reutilice entre llamadas. Además, `procEndDocPrinter.Call` en el defer se ejecuta aunque `StartDocPrinter` fallara (por el LIFO).

### 1.5 Race condition en counter — `counter.go:22-43`

`LeerUltimo` + `GuardarUltimo` no son atómicos. Si dos goroutines (o dos llamadas rápidas) lean y escriben el counter, se pierde incremento. Actualmente no es problema porque solo hay 1 goroutine (la UI congelada), pero será crítico cuando se corrija el P0 de la sección 1.1.

### 1.6 Validación de entrada

`Imprimir(cantidad int)` no valida:

- `cantidad < 0` → loop no se ejecuta pero no retorna error
- `cantidad == 1` → `totalPares = 0`, no imprime nada silenciosamente
- `cantidad ímpar` → se descarta el código impar sin aviso

---

## 2. [P1 - ARQUITECTURA] Clean Architecture

### 2.1 God Object — `app.go`

`App` struct mezcla 4 responsabilidades:

1. Lifecycle Wails (`ctx`, `startup`)
2. Configuración de directorios (`configPath`, `os.MkdirAll`)
3. Lógica de negocio (`imprimirPares`, contador)
4. Comunicación con hardware (`printer.ImprimirTextoPlano`)

### 2.2 Frontend conoce la infraestructura

El frontend llama directamente a `GetImpresoraActual()` que es un wrapper delgado sobre `printer.ObtenerImpresoraPredeterminada()`. No hay capa de UseCase. El frontend no debería saber que existe `winspool.drv`.

### 2.3 Propuesta de arquitectura

```
internal/
├── domain/
│   ├── counter.go          # Avanzar, WrapAround (sin I/O)
│   ├── ticket.go           # Ticket, Label (entidades)
│   └── errors.go           # ErrNoPrinter, ErrCounterOverflow
├── usecase/
│   ├── print_labels.go     # PrintLabelsUseCase (orchestration)
│   ├── get_status.go       # GetStatusUseCase
│   └── test_printer.go     # TestPrinterUseCase
├── infrastructure/
│   ├── persistence/
│   │   └── file_counter.go # FileCounter (LeeGuarda archivos)
│   └── printer/
│       └── winspool.go     # ZebraPrinter (syscalls)
app.go                      # Solo wiring + Wails bindings
```

### 2.4 Middleware Pattern propuesto

```go
type MethodHandler func(ctx context.Context, args ...any) (any, error)
type Middleware func(next MethodHandler) MethodHandler

func Chain(middlewares []Middleware, final MethodHandler) MethodHandler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        final = middlewares[i](final)
    }
    return final
}
```

Middlewares útiles:

- **Recovery** — `recover()` + log en vez de crash de la app
- **Logging** — timestamp + método + args + duración
- **Metrics** — conteo de llamadas por método, latencia
- **CircuitBreaker** — para la impresora (si falla 3 veces, no reintenta por 30s)

---

## 3. [P2 - SEGURIDAD] Desktop Security

### 3.1 Archivo de configuración sin permisos

`counter.go:43` usa `os.WriteFile` con permisos `0644` (world-readable). En Windows esto se traduce a ACLs por defecto, pero en un multi-usuario puede ser Problemático. Mejor: `0600`.

### 3.2 Sin sanitización de paths

`configPath` se construye con `os.UserConfigDir()` que es seguro, pero `counter.LeerUltimo(configPath)` hace `os.ReadFile(configPath)` directamente. Si el path viniera de userInput (no es el caso actual), sería path traversal. A futuro, validar que el path está dentro del directorio esperado.

### 3.3 Impresora como vector de ataque

`printer.ImprimirTextoPlano` acepta cualquier string y lo envía en crudo a la impresora via `RAW`. Un string ZPL malicioso podría:

- Cambiar configuraciones de la impresora
- Sobreescribir firmware (en impresoras con soporte)
- Ejecutar comandos de diagnóstico

**Mitigación:** validar que el ZPL generado solo contiene los caracteres del template (`^XA`, `^XZ`, alfanuméricos). Rechazar input externo.

### 3.4 Sin logging estructurado

No hay ningún `log` en el backend. Los errores se retornan como strings. Para una app de escritorio que maneja hardware, se necesita al menos un log a archivo para debugging post-mortem.

---

## 4. [P2 - PERFORMANCE] Optimizaciones

### 4.1 Apertura/cierre de impresora por par

`imprimirPares` llama `printer.ImprimirTextoPlano` por cada par, lo que ejecuta `OpenPrinter` → `StartDoc` → `Write` → `EndDoc` → `ClosePrinter` × N veces. Para 5.000 pares = 5.000 ciclos de apertura de handle.

**Solución:** abrir la impresora una vez, escribir todo el ZPL concatenado, cerrar una vez:

```go
func (a *App) imprimirPares(ctx context.Context, cantidad int) error {
    h, err := printer.Abrir(impresora)
    if err != nil { return err }
    defer printer.Cerrar(h)

    printer.IniciarDoc(h)
    for i := range totalPares {
        // generar + escribir
    }
    printer.FinDoc(h)
}
```

### 4.2 Alokaciones de strings en el loop

`zpl.Generar` usa `fmt.Sprintf` con template de ~300 bytes × 5.000 iteraciones. Considerar `strings.Builder` o pre-computar el template.

### 4.3 `os.ReadFile` + `os.WriteFile` en cada `GuardarUltimo`

Para un archivo de ~7 bytes, se abre el archivo, se lee/escribe, se cierra. Usar `os.WriteFile` es aceptable pero sin buffer. No es crítico pero sí innecesariamente pesado para operación tan simple.

### 4.4 Frontend: sin debounce en `cargarEstado()`

`confirmarImpresion` en `useImpresion.ts:58` llama `cargarEstado()` después de imprimir, que hace 2 llamadas Go consecutivas sin await paralelo:

```typescript
// Actual (secuencial)
impresora.value = await GetImpresoraActual()
contador.value = await GetContadorActual()

// Mejor (paralelo)
const [imp, cont] = await Promise.all([
    GetImpresoraActual(),
    GetContadorActual()
])
```

---

## 5. [P3 - UX / FRONTEND] Experiencia de Usuario

### 5.1 Sin feedback visual de éxito/error

`StatusMessage.vue:10` muestra todos los mensajes con el mismo color `#8a97a8`. No hay distinción entre éxito (verde) y error (rojo). El usuario no sabe si la impresión funcionó.

### 5.2 Sin barra de progreso para impresiones grandes

10.000 etiquetas = posible espera de minutos. No hay indicación de progreso. Wails soporta `runtime.EventsEmit` para enviar actualizaciones de progreso desde Go al frontend.

### 5.3 Botones sin accesibilidad

`PrintOptions.vue` — los `<button>` no tienen `aria-label` ni `role`. Para una app de escritorio esto es menor pero recomendado.

### 5.4 Sin validación de impresora desconectada

Si la impresora no está disponible al inicio, `GetImpresoraActual()` retorna `""` y el frontend muestra "No disponible". Pero si el usuario hace clic en "Etiqueta de prueba", se obtiene un error genérico. Mejor: deshabilitar los botones de impresión cuando no hay impresora.

### 5.5 Sin atajos de teclado

Para una app de uso repetitivo (un operador imprimiendo etiquetas todo el día), faltan atajos: Enter para confirmar, Escape para cancelar, números para selección rápida.

---

## Resumen de Prioridades

| Prioridad | Items | Esfuerzo |
| --- | --- | --- |
| **P0** | Goroutine para impresión + cancelación con context + validación input | Medio |
| **P1** | Separar God Object en clean architecture + middleware | Alto |
| **P2** | Seguridad (ZPL validation, logging, permisos) + Performance (abrir impresora 1 vez) | Bajo-Medio |
| **P3** | UX (colores éxito/error, progreso, accesibilidad, atajos) | Bajo |

**Impacto de NO corregir P0:** Al imprimir 10.000 etiquetas, la ventana queda congelada por varios minutos. Si el usuario intenta cerrarla, Windows la marca como "No responde" y puede matar el proceso, dejando el counter en un estado inválido (mitad de impresión).

¿Querés que proceda con la implementación de alguna de estas mejoras?
