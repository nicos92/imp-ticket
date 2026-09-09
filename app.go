package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"imp-ticket/internal/counter"
	"imp-ticket/internal/printer"
	"imp-ticket/internal/zpl"
)

const maxCantidadImpresion = 20000

type App struct {
	ctx        context.Context
	configPath string
	mu         sync.Mutex
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	baseDir, err := os.UserConfigDir()
	if err != nil {
		return
	}

	appDir := filepath.Join(baseDir, "Nicolas-Sandoval", "default-imp-ticket")
	os.MkdirAll(appDir, 0755)
	a.configPath = filepath.Join(appDir, "config")
}

func (a *App) GetImpresoraActual() string {
	nombre, err := printer.ObtenerImpresoraPredeterminada()
	if err != nil {
		return ""
	}
	return nombre
}

func (a *App) GetContadorActual() int {
	return counter.LeerUltimo(a.configPath)
}

func (a *App) Imprimir(cantidad int) error {
	if err := validarCantidad(cantidad); err != nil {
		return err
	}

	go a.imprimirPares(cantidad)
	return nil
}

func validarCantidad(cantidad int) error {
	if cantidad <= 0 {
		return fmt.Errorf("la cantidad debe ser mayor a cero")
	}
	if cantidad%2 != 0 {
		return fmt.Errorf("la cantidad debe ser un número par")
	}
	if cantidad > maxCantidadImpresion {
		return fmt.Errorf("la cantidad excede el máximo permitido de %d", maxCantidadImpresion)
	}
	return nil
}

func (a *App) TestImpresora() error {
	impresora, err := printer.ObtenerImpresoraPredeterminada()
	if err != nil {
		return err
	}

	zplData := zpl.Generar("0000001", "0000002", time.Now())
	return printer.ImprimirTextoPlano(impresora, zplData)
}

func (a *App) imprimirPares(cantidadCodigos int) {
	impresora, err := printer.ObtenerImpresoraPredeterminada()
	if err != nil {
		a.emitirError(err)
		return
	}

	a.mu.Lock()
	contador := counter.LeerUltimo(a.configPath)
	nuevoUltimoNumero := contador + cantidadCodigos
	errGuardar := counter.GuardarUltimo(a.configPath, nuevoUltimoNumero)
	a.mu.Unlock()
	if errGuardar != nil {
		a.emitirError(errGuardar)
		return
	}
	totalPares := cantidadCodigos / 2

	shouldReturn := recorrerPares(totalPares, contador, impresora, a)
	if shouldReturn {
		return
	}

	runtime.EventsEmit(a.ctx, "print:done")
}

func recorrerPares(totalPares int, contador int, impresora string, a *App) bool {
	for i := range totalPares {
		contador = counter.Avanzar(contador)
		c1 := fmt.Sprintf("%07d", contador)

		contador = counter.Avanzar(contador)
		c2 := fmt.Sprintf("%07d", contador)

		zplData := zpl.Generar(c1, c2, time.Now())
		if err := printer.ImprimirTextoPlano(impresora, zplData); err != nil {
			a.emitirError(fmt.Errorf("par %d: %w", i+1, err))
			return true
		}
	}
	return false
}

func (a *App) emitirError(err error) {
	runtime.EventsEmit(a.ctx, "print:error", err.Error())
}
