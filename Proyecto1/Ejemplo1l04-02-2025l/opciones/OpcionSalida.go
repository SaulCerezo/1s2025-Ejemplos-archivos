package opciones

import (
	"fmt"
	"os"
)

func OpcionSalida() {
	fmt.Println("Terminando el programa ...")
	os.Exit(0) //para terminar la ejecucion del programa
}
