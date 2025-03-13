package Gestion

import (
	"bytes"
	"ejemplo2/Structs"
	"ejemplo2/Utilities"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Estructura para representar una partición montada
type MountedPartition struct {
	Path     string
	Name     string
	ID       string
	Status   byte // 0: no montada, 1: montada
	LoggedIn bool // true: usuario ha iniciado sesión, false: no ha iniciado sesión
}

// Mapa para almacenar las particiones montadas, organizadas por disco
var mountedPartitions = make(map[string][]MountedPartition)

// Función para imprimir las particiones montadas ---------------------------------------  +Agregado
// PrintMountedPartitions imprime en la consola todas las particiones montadas.
func PrintMountedPartitions() {
	fmt.Println("Particiones montadas:")

	// Si no hay particiones montadas, muestra un mensaje y termina la función.
	if len(mountedPartitions) == 0 {
		fmt.Println("No hay particiones montadas.")
		return
	}

	// Itera sobre cada disco montado y sus particiones.
	for diskID, partitions := range mountedPartitions {
		fmt.Printf("Disco ID: %s\n", diskID)
		for _, partition := range partitions {
			// Determina si la partición está logueada o no.
			loginStatus := "No"
			if partition.LoggedIn {
				loginStatus = "Sí"
			}
			// Imprime los detalles de la partición.
			fmt.Printf(" - Partición Name: %s, ID: %s, Path: %s, Status: %c, LoggedIn: %s\n",
				partition.Name, partition.ID, partition.Path, partition.Status, loginStatus)
		}
	}
	fmt.Println("")
}

// GetMountedPartitions devuelve el mapa de particiones montadas.
// Retorna un mapa donde la clave es el ID del disco y el valor es una lista de particiones montadas en ese disco.
func GetMountedPartitions() map[string][]MountedPartition {
	return mountedPartitions
}

// --------------------------------------------------------------------- +Agregado

// MarkPartitionAsLoggedIn busca una partición por su ID y la marca como logueada (LoggedIn = true).
func MarkPartitionAsLoggedIn(id string) {
	// Recorre todas las particiones montadas en los discos.
	for diskID, partitions := range mountedPartitions {
		for i, partition := range partitions {
			// Si la partición coincide con el ID buscado, se marca como logueada.
			if partition.ID == id {
				mountedPartitions[diskID][i].LoggedIn = true
				fmt.Printf("Partición con ID %s marcada como logueada.\n", id)
				return
			}
		}
	}
	// Si no se encuentra la partición, se muestra un mensaje de error.
	fmt.Printf("No se encontró la partición con ID %s para marcarla como logueada.\n", id)
}

// --------------------------------------------------------------------- +

// ////////////////////////////////////////////////////////////////////////////
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

func Fdisk(size int, path string, name string, unit string, type_ string, fit string) {
	fmt.Println("======Start FDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Path:", path)
	fmt.Println("Name:", name)
	fmt.Println("Unit:", unit)
	fmt.Println("Type:", type_)
	fmt.Println("Fit:", fit)

	// Validar fit (b/w/f)
	if fit != "b" && fit != "f" && fit != "w" {
		fmt.Println("Error: Fit must be 'b', 'f', or 'w'")
		return
	}

	// Validar size > 0
	if size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	// Validar unit (b/k/m)
	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: Unit must be 'b', 'k', or 'm'")
		return
	}

	// Ajustar el tamaño en bytes
	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	}

	// Abrir el archivo binario en la ruta proporcionada
	file, err := Utilities.OpenFile(path)
	if err != nil {
		fmt.Println("Error: Could not open file at path:", path)
		return
	}

	//Leer el MBR (Master Boot Record), Se
	// lee el MBR desde el archivo para obtener la estructura del disco.
	var TempMBR Structs.MRB
	// Leer el objeto desde el archivo binario
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: Could not read MBR from file")
		return
	}

	// Imprimir el objeto MBR
	Structs.PrintMBR(TempMBR)

	fmt.Println("-------------")

	// Validaciones de las particiones
	var primaryCount, extendedCount, totalPartitions int
	var usedSpace int32 = 0

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			totalPartitions++
			usedSpace += TempMBR.Partitions[i].Size

			if TempMBR.Partitions[i].Type[0] == 'p' {
				primaryCount++
			} else if TempMBR.Partitions[i].Type[0] == 'e' {
				extendedCount++
			}
		}
	}

	// Validar que no se exceda el número máximo de particiones primarias y extendidas
	if totalPartitions >= 4 {
		fmt.Println("Error: No se pueden crear más de 4 particiones primarias o extendidas en total.")
		return
	}

	// Validar que solo haya una partición extendida
	if type_ == "e" && extendedCount > 0 {
		fmt.Println("Error: Solo se permite una partición extendida por disco.")
		return
	}

	// Validar que no se pueda crear una partición lógica sin una extendida
	if type_ == "l" && extendedCount == 0 {
		fmt.Println("Error: No se puede crear una partición lógica sin una partición extendida.")
		return
	}

	// Validar que el tamaño de la nueva partición no exceda el tamaño del disco
	if usedSpace+int32(size) > TempMBR.MbrSize {
		fmt.Println("Error: No hay suficiente espacio en el disco para crear esta partición.")
		return
	}

	// Determinar la posición de inicio de la nueva partición
	var gap int32 = int32(binary.Size(TempMBR))
	if totalPartitions > 0 {
		gap = TempMBR.Partitions[totalPartitions-1].Start + TempMBR.Partitions[totalPartitions-1].Size
	}

	// Encontrar una posición vacía para la nueva partición
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size == 0 {
			if type_ == "p" || type_ == "e" {
				// Crear partición primaria o extendida
				//copy se usa para copiar el contenido de name
				// (que es un string) a TempMBR.Partitions[i].Name, que es un array de bytes.
				TempMBR.Partitions[i].Size = int32(size)
				TempMBR.Partitions[i].Start = gap
				copy(TempMBR.Partitions[i].Name[:], name)
				copy(TempMBR.Partitions[i].Fit[:], fit)
				copy(TempMBR.Partitions[i].Status[:], "0")
				copy(TempMBR.Partitions[i].Type[:], type_)
				TempMBR.Partitions[i].Correlative = int32(totalPartitions + 1)

				if type_ == "e" {
					// Inicializar el primer EBR en la partición extendida
					ebrStart := gap // El primer EBR se coloca al inicio de la partición extendida
					ebr := Structs.EBR{
						PartFit:   fit[0],
						PartStart: ebrStart,
						PartSize:  0,
						PartNext:  -1,
					}
					copy(ebr.PartName[:], "")
					Utilities.WriteObject(file, ebr, int64(ebrStart))
				}

				break
			}
		}
	}

	// Manejar la creación de particiones lógicas dentro de una partición extendida
	if type_ == "l" {
		for i := 0; i < 4; i++ {
			if TempMBR.Partitions[i].Type[0] == 'e' {
				ebrPos := TempMBR.Partitions[i].Start
				var ebr Structs.EBR
				for {
					Utilities.ReadObject(file, &ebr, int64(ebrPos))
					if ebr.PartNext == -1 {
						break
					}
					ebrPos = ebr.PartNext
				}

				// Calcular la posición de inicio de la nueva partición lógica
				newEBRPos := ebr.PartStart + ebr.PartSize                    // El nuevo EBR se coloca después de la partición lógica anterior
				logicalPartitionStart := newEBRPos + int32(binary.Size(ebr)) // El inicio de la partición lógica es justo después del EBR

				// Ajustar el siguiente EBR
				ebr.PartNext = newEBRPos
				Utilities.WriteObject(file, ebr, int64(ebrPos))

				// Crear y escribir el nuevo EBR
				newEBR := Structs.EBR{
					PartFit:   fit[0],
					PartStart: logicalPartitionStart,
					PartSize:  int32(size),
					PartNext:  -1,
				}
				copy(newEBR.PartName[:], name)
				Utilities.WriteObject(file, newEBR, int64(newEBRPos))

				// Imprimir el nuevo EBR creado
				fmt.Println("Nuevo EBR creado:")
				Structs.PrintEBR(newEBR)
				break
			}
		}
	}

	// Sobrescribir el MBR
	if err := Utilities.WriteObject(file, TempMBR, 0); err != nil {
		fmt.Println("Error: Could not write MBR to file")
		return
	}

	var TempMBR2 Structs.MRB
	// Leer el objeto nuevamente para verificar
	if err := Utilities.ReadObject(file, &TempMBR2, 0); err != nil {
		fmt.Println("Error: Could not read MBR from file after writing")
		return
	}

	// Imprimir el objeto MBR actualizado
	Structs.PrintMBR(TempMBR2)

	// Cerrar el archivo binario
	defer file.Close()

	fmt.Println("======FIN FDISK======")
}

//////////////////////////////////////////////////////////////////////////////

// Función para montar particiones
func Mount(path string, name string) {
	file, err := Utilities.OpenFile(path)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo en la ruta:", path)
		return
	}
	defer file.Close()

	var TempMBR Structs.MRB
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR desde el archivo")
		return
	}

	fmt.Printf("Buscando partición con nombre: '%s'\n", name)

	partitionFound := false
	var partition Structs.Partition
	var partitionIndex int

	// Convertir el nombre a comparar a un arreglo de bytes de longitud fija
	nameBytes := [16]byte{}
	copy(nameBytes[:], []byte(name))

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Type[0] == 'p' && bytes.Equal(TempMBR.Partitions[i].Name[:], nameBytes[:]) {
			partition = TempMBR.Partitions[i]
			partitionIndex = i
			partitionFound = true
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: Partición no encontrada o no es una partición primaria")
		return
	}

	// Verificar si la partición ya está montada
	if partition.Status[0] == '1' {
		fmt.Println("Error: La partición ya está montada")
		return
	}

	//fmt.Printf("Partición encontrada: '%s' en posición %d\n", string(partition.Name[:]), partitionIndex+1)

	// Generar el ID de la partición utilizando la función `generateDiskID`
	// Esta función genera un identificador único para el disco basado en su ruta.
	diskID := generateDiskID(path)

	// Verificar si ya se ha montado alguna partición en este disco específico.
	// `mountedPartitions` es un mapa que guarda las particiones montadas por disco.
	mountedPartitionsInDisk := mountedPartitions[diskID]
	var letter byte

	// Si no hay particiones montadas en este disco (len(mountedPartitionsInDisk) == 0),
	// se considera que es un nuevo disco y asignamos una letra a la partición.
	if len(mountedPartitionsInDisk) == 0 {
		// Si no hay ningún disco montado (len(mountedPartitions) == 0),
		// asignamos la letra 'a' para la primera partición en el disco.
		if len(mountedPartitions) == 0 {
			letter = 'a'
		} else {
			// Si ya hay discos montados, obtenemos el último disco montado (`lastDiskID`),
			// luego obtenemos la última letra de la partición montada en ese disco.
			// A partir de esa letra, incrementamos para asignar la siguiente letra disponible.
			lastDiskID := getLastDiskID()
			lastLetter := mountedPartitions[lastDiskID][0].ID[len(mountedPartitions[lastDiskID][0].ID)-1]
			letter = lastLetter + 1 // Incrementamos la letra para la siguiente partición.
		}
	} else {
		// Si ya hay particiones montadas en este disco, utilizamos la misma letra
		// que la primera partición montada en este disco.
		letter = mountedPartitionsInDisk[0].ID[len(mountedPartitionsInDisk[0].ID)-1]
	}

	// Crear el ID de la partición utilizando el último par de dígitos de un carnet
	// (por ejemplo, "202501234"), el índice de la partición (`partitionIndex`),
	// y la letra que hemos asignado a la partición.
	carnet := "202501234"                   // Cambiar su carnet aquí
	lastTwoDigits := carnet[len(carnet)-2:] // Obtener los últimos dos dígitos del carnet
	partitionID := fmt.Sprintf("%s%d%c", lastTwoDigits, partitionIndex+1, letter)
	// Formateamos el ID como un string que contiene los dos últimos dígitos del carnet,
	// el índice de la partición incrementado, y la letra asignada.

	// Actualizar el estado de la partición a "montada" y asignar el ID generado a la partición.
	// `partition.Status[0]` se establece en '1' para indicar que la partición está montada.
	// `copy(partition.Id[:], partitionID)` asigna el ID generado a la partición.
	partition.Status[0] = '1'
	copy(partition.Id[:], partitionID)

	// Actualizamos el `TempMBR.Partitions[partitionIndex]` para reflejar los cambios en la partición.
	TempMBR.Partitions[partitionIndex] = partition

	// Agregamos la partición montada al mapa `mountedPartitions` para mantener un registro
	// de las particiones montadas en el disco correspondiente.
	mountedPartitions[diskID] = append(mountedPartitions[diskID], MountedPartition{
		Path:   path,        // Ruta del disco.
		Name:   name,        // Nombre de la partición.
		ID:     partitionID, // ID de la partición.
		Status: '1',         // Estado de la partición (montada).
	})

	// Escribir el MBR actualizado en el archivo utilizando la función `Utilities.WriteObject`.
	// Si la escritura falla, se imprime un mensaje de error.
	if err := Utilities.WriteObject(file, TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo sobrescribir el MBR en el archivo")
		return
	}

	// Imprimir el mensaje confirmando que la partición ha sido montada, junto con su ID.
	fmt.Printf("Partición montada con ID: %s\n", partitionID)

	fmt.Println("")
	// Imprimir el MBR actualizado (esto muestra el estado actual del MBR con las particiones montadas).
	fmt.Println("MBR actualizado:")
	Structs.PrintMBR(TempMBR)
	fmt.Println("")

	// Imprimir las particiones que están montadas actualmente en el sistema (solo se mantienen mientras dure la sesión de la consola).
	PrintMountedPartitions()

}

// Función para obtener el ID del último disco montado
func getLastDiskID() string {
	var lastDiskID string
	for diskID := range mountedPartitions {
		lastDiskID = diskID
	}
	return lastDiskID
}

func generateDiskID(path string) string {
	return strings.ToLower(path)
}
