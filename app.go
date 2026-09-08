package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	maxContador = 9999999
)

var (
	winspool              = syscall.NewLazyDLL("winspool.drv")
	procGetDefaultPrinter = winspool.NewProc("GetDefaultPrinterW")
	procOpenPrinter       = winspool.NewProc("OpenPrinterW")
	procClosePrinter      = winspool.NewProc("ClosePrinter")
	procStartDocPrinter   = winspool.NewProc("StartDocPrinterW")
	procEndDocPrinter     = winspool.NewProc("EndDocPrinter")
	procStartPagePrinter  = winspool.NewProc("StartPagePrinter")
	procEndPagePrinter    = winspool.NewProc("EndPagePrinter")
	procWritePrinter      = winspool.NewProc("WritePrinter")
)

type DOC_INFO_1 struct {
	pDocName    *uint16
	pOutputFile *uint16
	pDatatype   *uint16
}

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
	nombre, err := obtenerImpresoraPredeterminada()
	if err != nil {
		return ""
	}
	return nombre
}

func (a *App) GetContadorActual() int {
	return a.leerUltimoNumero()
}

func (a *App) Imprimir(cantidad int) error {
	return a.imprimirPares(cantidad)
}

func (a *App) TestImpresora() error {
	impresora, err := obtenerImpresoraPredeterminada()
	if err != nil {
		return err
	}

	zpl := generarZPL("0000001", "0000002", time.Now())
	return imprimirTextoPlano(impresora, zpl)
}

func (a *App) imprimirPares(cantidadCodigos int) error {
	impresora, err := obtenerImpresoraPredeterminada()
	if err != nil {
		return err
	}

	ultimoNumero := a.leerUltimoNumero()
	if err := a.guardarUltimoNumero(ultimoNumero + cantidadCodigos); err != nil {
		return err
	}

	totalPares := cantidadCodigos / 2
	contador := ultimoNumero

	for i := range totalPares {
		contador++
		if contador > maxContador {
			contador = 1
		}
		c1 := fmt.Sprintf("%07d", contador)

		contador++
		if contador > maxContador {
			contador = 1
		}
		c2 := fmt.Sprintf("%07d", contador)

		zpl := generarZPL(c1, c2, time.Now())
		if err := imprimirTextoPlano(impresora, zpl); err != nil {
			return fmt.Errorf("error en par %d: %v", i+1, err)
		}
	}

	return nil
}

func (a *App) leerUltimoNumero() int {
	if a.configPath == "" {
		return 0
	}

	data, err := os.ReadFile(a.configPath)
	if err != nil {
		return 0
	}

	num, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0
	}
	return num
}

func (a *App) guardarUltimoNumero(num int) error {
	if a.configPath == "" {
		return fmt.Errorf("configPath no inicializado")
	}
	return os.WriteFile(a.configPath, fmt.Appendf(nil, "%d", num), 0644)
}

func obtenerImpresoraPredeterminada() (string, error) {
	var size uint32
	procGetDefaultPrinter.Call(0, uintptr(unsafe.Pointer(&size)))
	if size == 0 {
		return "", fmt.Errorf("no se pudo determinar el tamaño")
	}

	buffer := make([]uint16, size)
	ret, _, err := procGetDefaultPrinter.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	if ret == 0 {
		return "", fmt.Errorf("error al obtener impresora: %v", err)
	}

	return syscall.UTF16ToString(buffer), nil
}

func imprimirTextoPlano(nombreImpresora string, texto string) error {
	var hPrinter uintptr

	pPrinterName, err := syscall.UTF16PtrFromString(nombreImpresora)
	if err != nil {
		return err
	}

	ret, _, err := procOpenPrinter.Call(
		uintptr(unsafe.Pointer(pPrinterName)),
		uintptr(unsafe.Pointer(&hPrinter)),
		0,
	)
	if ret == 0 {
		return fmt.Errorf("no se pudo abrir la impresora: %v", err)
	}
	defer procClosePrinter.Call(hPrinter)

	docName, _ := syscall.UTF16PtrFromString("NSS-Default-Imp-Ticket")
	dataType, _ := syscall.UTF16PtrFromString("RAW")

	di := DOC_INFO_1{
		pDocName:    docName,
		pOutputFile: nil,
		pDatatype:   dataType,
	}

	ret, _, err = procStartDocPrinter.Call(
		hPrinter,
		1,
		uintptr(unsafe.Pointer(&di)),
	)
	if ret == 0 {
		return fmt.Errorf("error en StartDocPrinter: %v", err)
	}
	defer procEndDocPrinter.Call(hPrinter)

	ret, _, _ = procStartPagePrinter.Call(hPrinter)
	if ret == 0 {
		return fmt.Errorf("error en StartPagePrinter")
	}
	defer procEndPagePrinter.Call(hPrinter)

	bytesTexto := []byte(texto)
	var bytesEscritos uint32

	ret, _, err = procWritePrinter.Call(
		hPrinter,
		uintptr(unsafe.Pointer(&bytesTexto[0])),
		uintptr(len(bytesTexto)),
		uintptr(unsafe.Pointer(&bytesEscritos)),
	)
	if ret == 0 {
		return fmt.Errorf("error al escribir datos: %v", err)
	}

	return nil
}

func generarZPL(codigo1, codigo2 string, fecha time.Time) string {
	fechaFormateada := fecha.Format("02/01/2006 15:04")

	plantillaZPL := `^XA
	^MMT
	^PW832
	^LL392
	^LS0
	^FO310,80^A0R,35,35^FD%s^FS
	^BY3,3,240^FT310,368^BCB,,N,N
	^FH\^FD>:Z>5123456>67^FS
	^PQ1,0,1,Y
	^FO010,030^A0R,50,90^FDZ%s^FS

	^FO720,80^A0R,35,35^FD%s^FS
	^BY3,3,240^FT720,368^BCB,,N,N
	^FH\^FD>:Z>5123456>67^FS
	^PQ1,0,1,Y
	^FO420,030^A0R,50,90^FDZ%s^FS
	^XZ`

	return fmt.Sprintf(plantillaZPL, fechaFormateada, codigo1, fechaFormateada, codigo2)
}
