package User

import (
	"ejemplo2/Gestion"
	"ejemplo2/Structs"
	"ejemplo2/Utilities"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func Login(user string, pass string, id string) {
	fmt.Println("======Start LOGIN======")
	fmt.Println("User:", user)
	fmt.Println("Pass:", pass)
	fmt.Println("Id:", id)

	// Obtener las particiones montadas
	mountedPartitions := Gestion.GetMountedPartitions()
	var filepath string
	var partitionFound bool
	var login bool = false

	// Verificar si el usuario ya está logueado en alguna partición
	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.ID == id && partition.LoggedIn { // Si la partición ya tiene un usuario logueado
				fmt.Println("Ya existe un usuario logueado!")
				return
			}
			if partition.ID == id { // Encuentra la partición correcta
				filepath = partition.Path
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	// Si no se encontró la partición montada, se detiene el proceso
	if !partitionFound {
		fmt.Println("Error: No se encontró ninguna partición montada con el ID proporcionado")
		return
	}

	// Abrir el archivo del sistema de archivos binario
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close() // Cierra el archivo al final de la ejecución

	var TempMBR Structs.MRB
	// Leer el MBR (Master Boot Record) del archivo binario
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	// Imprimir información del MBR
	Structs.PrintMBR(TempMBR)
	fmt.Println("-------------")

	var index int = -1
	// Buscar la partición en el MBR por su ID
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 { // Verifica que la partición tenga tamaño
			if strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) { // Compara el ID
				fmt.Println("Partition found")
				if TempMBR.Partitions[i].Status[0] == '1' { // Verifica si está montada
					fmt.Println("Partition is mounted")
					index = i
				} else {
					fmt.Println("Partition is not mounted")
					return
				}
				break
			}
		}
	}

	// Si se encontró la partición, imprimir su información
	if index != -1 {
		Structs.PrintPartition(TempMBR.Partitions[index])
	} else {
		fmt.Println("Partition not found")
		return
	}

	var tempSuperblock Structs.Superblock
	// Leer el Superblock de la partición
	if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	// Buscar el archivo de usuarios "/users.txt" dentro del sistema de archivos
	indexInode := InitSearch("/users.txt", file, tempSuperblock)

	var crrInode Structs.Inode
	// Leer el Inodo del archivo "users.txt"
	if err := Utilities.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{})))); err != nil {
		fmt.Println("Error: No se pudo leer el Inodo:", err)
		return
	}

	// Obtener el contenido del archivo users.txt desde los bloques del inodo
	data := GetInodeFileData(crrInode, file, tempSuperblock)

	// Dividir el contenido del archivo en líneas
	lines := strings.Split(data, "\n")

	// Iterar a través de las líneas para verificar las credenciales
	for _, line := range lines {
		words := strings.Split(line, ",")

		// Si la línea tiene 5 elementos, verificar si el usuario y contraseña coinciden
		if len(words) == 5 {
			if (strings.Contains(words[3], user)) && (strings.Contains(words[4], pass)) {
				login = true
				break
			}
		}
	}

	// Imprimir la información del Inodo
	fmt.Println("Inode", crrInode.I_block)

	// Si las credenciales son correctas, marcar la partición como logueada
	if login {
		fmt.Println("Usuario logueado con éxito")
		Gestion.MarkPartitionAsLoggedIn(id) // Marca la partición como logueada
	}

	fmt.Println("======End LOGIN======")
}

func InitSearch(path string, file *os.File, tempSuperblock Structs.Superblock) int32 {
	fmt.Println("======Start BUSQUEDA INICIAL ======")
	fmt.Println("path:", path)

	// Dividir la ruta en partes usando "/" como separador
	TempStepsPath := strings.Split(path, "/")
	StepsPath := TempStepsPath[1:] // Omitir el primer elemento vacío si la ruta empieza con "/"

	fmt.Println("StepsPath:", StepsPath, "len(StepsPath):", len(StepsPath))
	for _, step := range StepsPath {
		fmt.Println("step:", step)
	}

	var Inode0 Structs.Inode
	// Leer el inodo raíz (primer inodo del sistema de archivos)
	if err := Utilities.ReadObject(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		return -1 // Retornar -1 si hubo un error al leer
	}

	fmt.Println("======End BUSQUEDA INICIAL======")

	// Llamar a la función que busca el inodo del archivo según la ruta
	return SarchInodeByPath(StepsPath, Inode0, file, tempSuperblock)
}

// stack
func pop(s *[]string) string {
	lastIndex := len(*s) - 1
	last := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return last
}

func SarchInodeByPath(StepsPath []string, Inode Structs.Inode, file *os.File, tempSuperblock Structs.Superblock) int32 {
	fmt.Println("======Start BUSQUEDA INODO POR PATH======")

	index := int32(0) // Contador de bloques procesados en el inodo actual

	// Extrae el primer elemento del path y elimina espacios en blanco
	SearchedName := strings.Replace(pop(&StepsPath), " ", "", -1)

	fmt.Println("========== SearchedName:", SearchedName)

	// Iterar sobre los bloques del inodo
	for _, block := range Inode.I_block {
		if block != -1 { // Si el bloque es válido (no está vacío)
			if index < 13 { // Manejo de bloques directos (0-12)
				var crrFolderBlock Structs.Folderblock

				// Leer el bloque de carpeta desde el archivo binario
				if err := Utilities.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Folderblock{})))); err != nil {
					return -1
				}

				// Buscar el archivo/directorio dentro del bloque de carpeta
				for _, folder := range crrFolderBlock.B_content {
					fmt.Println("Folder === Name:", string(folder.B_name[:]), "B_inodo", folder.B_inodo)

					// Si el nombre del archivo o directorio coincide
					if strings.Contains(string(folder.B_name[:]), SearchedName) {
						fmt.Println("len(StepsPath)", len(StepsPath), "StepsPath", StepsPath)

						if len(StepsPath) == 0 { // Si llegamos al final de la ruta
							fmt.Println("Folder found======")
							return folder.B_inodo // Retornar índice del inodo encontrado
						} else { // Si aún hay más niveles en la ruta, seguir buscando
							fmt.Println("NextInode======")
							var NextInode Structs.Inode

							// Leer el siguiente inodo desde el archivo binario
							if err := Utilities.ReadObject(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Structs.Inode{})))); err != nil {
								return -1
							}

							// Llamada recursiva para seguir con la búsqueda
							return SarchInodeByPath(StepsPath, NextInode, file, tempSuperblock)
						}
					}
				}
			} else {
				fmt.Println("Manejo de bloques indirectos no implementado") // Falta implementar acceso a bloques indirectos
			}
		}
		index++ // Incrementar índice para saber si son bloques directos o indirectos
	}

	fmt.Println("======End BUSQUEDA INODO POR PATH======")
	return 0 // Si no se encontró, retornar 0
}

func GetInodeFileData(Inode Structs.Inode, file *os.File, tempSuperblock Structs.Superblock) string {
	fmt.Println("======Start CONTENIDO DEL BLOQUE======")
	index := int32(0)
	// define content as a string
	var content string

	// Iterate over i_blocks from Inode
	for _, block := range Inode.I_block {
		if block != -1 {
			//Dentro de los directos
			if index < 13 {
				var crrFileBlock Structs.Fileblock
				// Read object from bin file
				if err := Utilities.ReadObject(file, &crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Fileblock{})))); err != nil {
					return ""
				}

				content += string(crrFileBlock.B_content[:])

			} else {
				fmt.Print("indirectos")
			}
		}
		index++
	}

	fmt.Println("======End CONTENIDO DEL BLOQUE======")
	return content
}

// MKUSER
func AppendToFileBlock(inode *Structs.Inode, newData string, file *os.File, superblock Structs.Superblock) error {
	// Leer el contenido existente del archivo utilizando la función GetInodeFileData
	existingData := GetInodeFileData(*inode, file, superblock)

	// Concatenar el nuevo contenido
	fullData := existingData + newData

	// Asegurarse de que el contenido no exceda el tamaño del bloque
	if len(fullData) > len(inode.I_block)*binary.Size(Structs.Fileblock{}) {
		// Si el contenido excede, necesitas manejar bloques adicionales
		return fmt.Errorf("el tamaño del archivo excede la capacidad del bloque actual y no se ha implementado la creación de bloques adicionales")
	}

	// Escribir el contenido actualizado en el bloque existente
	var updatedFileBlock Structs.Fileblock
	copy(updatedFileBlock.B_content[:], fullData)
	if err := Utilities.WriteObject(file, updatedFileBlock, int64(superblock.S_block_start+inode.I_block[0]*int32(binary.Size(Structs.Fileblock{})))); err != nil {
		return fmt.Errorf("error al escribir el bloque actualizado: %v", err)
	}

	// Actualizar el tamaño del inodo
	inode.I_size = int32(len(fullData))
	if err := Utilities.WriteObject(file, *inode, int64(superblock.S_inode_start+inode.I_block[0]*int32(binary.Size(Structs.Inode{})))); err != nil {
		return fmt.Errorf("error al actualizar el inodo: %v", err)
	}

	return nil
}
