package zpl

import (
	"fmt"
	"time"
)

func Generar(codigo1, codigo2 string, fecha time.Time) string {
	fechaFormateada := fecha.Format("02/01/2006 15:04")

	codBarra1 := formatoCodigo(codigo1)
	codBarra2 := formatoCodigo(codigo2)

	plantilla := `^XA
	^MMT
	^PW832
	^LL392
	^LS0
	^FO310,80^A0R,35,35^FD%s^FS
	^BY3,3,240^FT310,368^BCB,,N,N
	^FH\^FD>:%s^FS
	^PQ1,0,1,Y
	^FO010,030^A0R,50,90^FDZ%s^FS

	^FO720,80^A0R,35,35^FD%s^FS
	^BY3,3,240^FT720,368^BCB,,N,N
	^FH\^FD>:%s^FS
	^PQ1,0,1,Y
	^FO420,030^A0R,50,90^FDZ%s^FS
	^XZ`

	return fmt.Sprintf(
		plantilla,
		fechaFormateada, codBarra1, codigo1,
		fechaFormateada, codBarra2, codigo2,
	)
}

func formatoCodigo(codigo string) string {
	cod1 := codigo[0:6]
	cod12 := codigo[6:7]
	return fmt.Sprintf("Z>5%s>6%s", cod1, cod12)
}
