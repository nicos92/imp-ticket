package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"imp-ticket/internal/counter"
	"imp-ticket/internal/printer"
	"imp-ticket/internal/zpl"
)

type App struct {
	ctx        context.Context
	configPath string
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
	return a.imprimirPares(cantidad)
}

func (a *App) TestImpresora() error {
	impresora, err := printer.ObtenerImpresoraPredeterminada()
	if err != nil {
		return err
	}

	zplData := zpl.Generar("0000001", "0000002", time.Now())
	return printer.ImprimirTextoPlano(impresora, zplData)
}

func (a *App) imprimirPares(cantidadCodigos int) error {
	impresora, err := printer.ObtenerImpresoraPredeterminada()
	if err != nil {
		return err
	}

	contador := counter.LeerUltimo(a.configPath)
	nuevoUltimoNumero := contador + cantidadCodigos
	errNuevoUltimoNumero := counter.GuardarUltimo(a.configPath, nuevoUltimoNumero)
	if errNuevoUltimoNumero != nil {
		return errNuevoUltimoNumero
	}
	totalPares := cantidadCodigos / 2

	for i := range totalPares {
		contador = counter.Avanzar(contador)
		c1 := fmt.Sprintf("%07d", contador)

		contador = counter.Avanzar(contador)
		c2 := fmt.Sprintf("%07d", contador)

		zplData := zpl.Generar(c1, c2, time.Now())
		if err := printer.ImprimirTextoPlano(impresora, zplData); err != nil {
			return fmt.Errorf("par %d: %w", i+1, err)
		}
	}

	return nil
}
