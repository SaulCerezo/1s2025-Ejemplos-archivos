package opciones

import "fmt"

func MostrarUsuarios() {
	if len(usuarios) == 0 {
		fmt.Println("No hay usuarios ingresados.")
	} else {
		fmt.Println("Lista de usuarios ingresados:")
		for i, usuario := range usuarios {
			fmt.Printf("%d. %s\n", i+1, usuario)
		}
	}
}
