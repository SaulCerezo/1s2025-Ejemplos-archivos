package opciones

import "fmt"

var usuarios []string // Slice para almacenar varios usuarios

func IngresarUsuario() {
	var usuario string
	fmt.Print("Ingrese su nombre de usuario: ")
	fmt.Scanln(&usuario)

	usuarios = append(usuarios, usuario) // Agrega el usuario a la lista
	fmt.Println("Usuario guardado correctamente.")
}
