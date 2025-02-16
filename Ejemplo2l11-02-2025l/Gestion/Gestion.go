package Gestion

import (
	"ejemplo2/Structs"
	"ejemplo2/Utilities"
	"fmt"
	"math/rand"
	"time"
)

func Mkdisk(size int, fit string, unit string, path string) {
	fmt.Println("======INICIO MKDISK======")
	fmt.Printf("Size: %d\nFit: %s\nUnit: %s\nPath: %s\n", size, fit, unit, path)

	// Validaciones
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: Fit debe ser 'bf', 'wf' o 'ff'")
		return
	}
	if size <= 0 {
		fmt.Println("Error: Size debe ser mayor a 0")
		return
	}
	if unit != "k" && unit != "m" {
		fmt.Println("Error: Las unidades válidas son 'k' o 'm'")
		return
	}

	// Crear archivo
	if err := Utilities.CreateFile(path); err != nil {
		fmt.Println("Error al crear archivo:", err)
		return
	}

	// Convertir tamaño a bytes
	sizeInBytes := size * 1024
	if unit == "m" {
		sizeInBytes *= 1024
	}

	// Abrir archivo
	file, err := Utilities.OpenFile(path)
	if err != nil {
		fmt.Println("Error al abrir archivo:", err)
		return
	}
	defer file.Close() // Asegura el cierre del archivo al salir de la función

	// Escribir ceros en un solo bloque en lugar de un bucle
	zeroBlock := make([]byte, sizeInBytes) // Crea un slice de bytes lleno de ceros
	if _, err := file.Write(zeroBlock); err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		return
	}

	// Crear MBR
	var newMBR Structs.MRB
	newMBR.MbrSize = int32(sizeInBytes)
	newMBR.Signature = rand.Int31()
	copy(newMBR.Fit[:], fit)

	// Obtener fecha actual en formato YYYY-MM-DD
	formattedDate := time.Now().Format("2006-01-02")
	copy(newMBR.CreationDate[:], formattedDate)

	// Escribir el MBR en el archivo
	if err := Utilities.WriteObject(file, newMBR, 0); err != nil {
		fmt.Println("Error al escribir el MBR:", err)
		return
	}

	// Leer el MBR para verificar que se escribió correctamente
	var tempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &tempMBR, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err)
		return
	}

	// Imprimir el MBR leído
	Structs.PrintMBR(tempMBR)

	fmt.Println("======FIN MKDISK======")
}
