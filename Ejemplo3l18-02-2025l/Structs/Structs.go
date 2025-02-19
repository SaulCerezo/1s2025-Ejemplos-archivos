package Structs

import (
	"fmt"
)

type MRB struct {
	MbrSize      int32    // 4 bytes //int32 va desde -2,147,483,648 hasta 2,147,483,647.
	CreationDate [10]byte // 10 bytes
	Signature    int32    // 4 bytes
	Fit          [1]byte  // 1 byte
	Partitions   [4]Partition
}

func PrintMBR(data MRB) {
	// Imprime los datos del MBR, formateando la fecha de creación, el ajuste (fit) y el tamaño total del disco
	fmt.Println(fmt.Sprintf("CreationDate: %s, fit: %s, size: %d", string(data.CreationDate[:]), string(data.Fit[:]), data.MbrSize))

	// Recorre las 4 particiones y las imprime
	for i := 0; i < 4; i++ {
		PrintPartition(data.Partitions[i])
	}
}

type Partition struct {
	Status      [1]byte
	Type        [1]byte
	Fit         [1]byte
	Start       int32
	Size        int32
	Name        [16]byte
	Correlative int32
	Id          [4]byte
}

func PrintPartition(data Partition) {
	fmt.Println(
		//Usa fmt.Sprintf para formatear la información de la partición en un solo string.
		fmt.Sprintf("Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s",
			//string(data.Name[:]) convierte el array de bytes [16]byte a una cadena de texto.
			string(data.Name[:]), string(data.Type[:]), data.Start, data.Size,
			//data.Start y data.Size se imprimen como enteros (int32).
			string(data.Status[:]), string(data.Id[:])))
}

type EBR struct {
	PartMount byte
	PartFit   byte
	PartStart int32
	PartSize  int32
	PartNext  int32
	PartName  [16]byte
}

func PrintEBR(data EBR) {
	fmt.Println(fmt.Sprintf("Name: %s, fit: %c, start: %d, size: %d, next: %d, mount: %c",
		string(data.PartName[:]),
		data.PartFit,
		data.PartStart,
		data.PartSize,
		data.PartNext,
		data.PartMount))
}

/*
   MbrSize (4 bytes):
       Hex:

   CreationDate (10 bytes):
       Hex:

	Los siguientes 10 bytes son: 32 30 32 34 2D 30 38 2D 30 39,
	que en texto plano representan la fecha 2024-08-09.


   Signature (4 bytes):
       Hex:

   Fit (1 byte):
       Hex:


*/
