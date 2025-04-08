package opciones

import "fmt"

func BuscarUsuario(nombre string) {
	for _, usuario := range usuarios {
		if usuario == nombre {
			fmt.Println("El usuario", nombre, "está en el sistema.")
			return
		}
	}
	fmt.Println("El usuario", nombre, "NO está en el sistema.")
}
