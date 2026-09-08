package printer

import (
	"fmt"
	"syscall"
	"unsafe"
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

type docInfo1 struct {
	pDocName    *uint16
	pOutputFile *uint16
	pDatatype   *uint16
}

func ObtenerImpresoraPredeterminada() (string, error) {
	var size uint32
	procGetDefaultPrinter.Call(
		0,
		uintptr(unsafe.Pointer(&size)),
	)
	if size == 0 {
		return "", fmt.Errorf("no se pudo determinar el tamaño")
	}

	buffer := make([]uint16, size)
	ret, _, err := procGetDefaultPrinter.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	if ret == 0 {
		return "", fmt.Errorf("error al obtener impresora: %w", err)
	}

	return syscall.UTF16ToString(buffer), nil
}

func ImprimirTextoPlano(nombreImpresora, texto string) error {
	hPrinter, err := abrirImpresora(nombreImpresora)
	if err != nil {
		return err
	}
	defer procClosePrinter.Call(hPrinter)

	if err := iniciarDocumento(hPrinter); err != nil {
		return err
	}
	defer procEndDocPrinter.Call(hPrinter)

	if err := iniciarPagina(hPrinter); err != nil {
		return err
	}
	defer procEndPagePrinter.Call(hPrinter)

	return escribirDatos(hPrinter, texto)
}

func abrirImpresora(nombre string) (uintptr, error) {
	var hPrinter uintptr

	pPrinterName, err := syscall.UTF16PtrFromString(nombre)
	if err != nil {
		return 0, fmt.Errorf("nombre de impresora inválido: %w", err)
	}

	ret, _, err := procOpenPrinter.Call(
		uintptr(unsafe.Pointer(pPrinterName)),
		uintptr(unsafe.Pointer(&hPrinter)),
		0,
	)
	if ret == 0 {
		return 0, fmt.Errorf("no se pudo abrir la impresora: %w", err)
	}
	return hPrinter, nil
}

func iniciarDocumento(hPrinter uintptr) error {
	docName, _ := syscall.UTF16PtrFromString("NSS-Default-Imp-Ticket")
	dataType, _ := syscall.UTF16PtrFromString("RAW")

	di := docInfo1{
		pDocName:    docName,
		pOutputFile: nil,
		pDatatype:   dataType,
	}

	ret, _, err := procStartDocPrinter.Call(
		hPrinter,
		1,
		uintptr(unsafe.Pointer(&di)),
	)
	if ret == 0 {
		return fmt.Errorf("error en StartDocPrinter: %w", err)
	}
	return nil
}

func iniciarPagina(hPrinter uintptr) error {
	ret, _, _ := procStartPagePrinter.Call(hPrinter)
	if ret == 0 {
		return fmt.Errorf("error en StartPagePrinter")
	}
	return nil
}

func escribirDatos(hPrinter uintptr, texto string) error {
	bytesTexto := []byte(texto)
	var bytesEscritos uint32

	ret, _, err := procWritePrinter.Call(
		hPrinter,
		uintptr(unsafe.Pointer(&bytesTexto[0])),
		uintptr(len(bytesTexto)),
		uintptr(unsafe.Pointer(&bytesEscritos)),
	)
	if ret == 0 {
		return fmt.Errorf("error al escribir datos: %w", err)
	}
	return nil
}