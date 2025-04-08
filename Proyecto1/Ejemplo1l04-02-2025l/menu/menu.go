package menu

import (
	"ejemplo/opciones"
	"fmt"
)

func IniciarMenu() {
	for {
		fmt.Println("\nMenú de Opciones")
		fmt.Println("1. Ingresar Usuario")
		fmt.Println("2. Mostrar Todos los Usuarios")
		fmt.Println("3. Buscar Usuario")
		fmt.Println("0. Salir")
		fmt.Print("Seleccione una opción: ")

		var opcion int
		fmt.Scanln(&opcion)

		switch opcion {
		case 1:
			opciones.IngresarUsuario()
		case 2:
			opciones.MostrarUsuarios()
		case 3:
			var nombre string
			fmt.Print("Ingrese el nombre del usuario a buscar: ")
			fmt.Scanln(&nombre)
			opciones.BuscarUsuario(nombre) // Pasamos el nombre como parámetro
		case 0:
			opciones.OpcionSalida() // Termina el programa
		default:
			fmt.Println("Opción inválida, intente de nuevo.")
		}
	}
}
