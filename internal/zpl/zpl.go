package zpl

import (
	"fmt"
	"time"
)

func Generar(codigo1, codigo2 string, fecha time.Time) string {
	fechaFormateada := fecha.Format("02/01/2006 15:04")

	plantilla := `^XA
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

	return fmt.Sprintf(
		plantilla,
		fechaFormateada, codigo1,
		fechaFormateada, codigo2,
	)
}