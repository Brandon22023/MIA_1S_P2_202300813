package commands

import (
	structures "terminal/structures"
	utils "terminal/utils"
	"errors"  // Paquete para manejar errores y crear nuevos errores con mensajes personalizados
	"fmt"     // Paquete para formatear cadenas y realizar operaciones de entrada/salida
	"regexp"  // Paquete para trabajar con expresiones regulares, útil para encontrar y manipular patrones en cadenas
	"strconv" // Paquete para convertir cadenas a otros tipos de datos, como enteros
	"strings" // Paquete para manipular cadenas, como unir, dividir, y modificar contenido de cadenas
	"encoding/binary" // Agregar esta línea
	"os" // Agregar esta línea

)

// FDISK estructura que representa el comando fdisk con sus parámetros
type FDISK struct {
	size int    // Tamaño de la partición
	unit string // Unidad de medida del tamaño (K o M o B)
	fit  string // Tipo de ajuste (BF, FF, WF)
	path string // Ruta del archivo del disco
	typ  string // Tipo de partición (P, E, L)
	name string // Nombre de la partición
}

/*
	fdisk -size=1 -type=L -unit=M -fit=BF -name="Particion3" -path="/home/keviin/University/PRACTICAS/MIA_LAB_S2_2024/CLASEEXTRA/disks/Disco1.mia"
	fdisk -size=300 -path=/home/Disco1.mia -name=Particion1
	fdisk -type=E -path=/home/Disco2.mia -Unit=K -name=Particion2 -size=300
*/
var addSet, sizeSet bool

func ParseFdisk(tokens []string) (string, error) {
	cmd := &FDISK{} // Crea una nueva instancia de FDISK
	sizeSet = false
	addSet = false
	deleteSet := false // Nueva bandera para el comando delete
    deleteType := ""   // Variable para almacenar el tipo de eliminación (fast o full)
	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando fdisk
	re := regexp.MustCompile(`-size=\d+|-add=-?\d+|-unit=[kKmMbB]|-fit=[bBfF]{2}|-path="[^"]+"|-path=[^\s]+|-type=[pPeElL]|-name="[^"]+"|-name=[^\s]+|-delete=(fast|full)`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return "", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]
	
		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}
	
		// Procesar los parámetros en el orden en que aparecen
		if key == "-size" {
			if !addSet { // Priorizar el primero que aparezca
				size, err := strconv.Atoi(value)
				if err != nil || size <= 0 {
					return "", errors.New("el tamaño debe ser un número entero positivo")
				}
				cmd.size = size
				sizeSet = true
			}
		} else if key == "-add" {
			if !sizeSet { // Priorizar el primero que aparezca
				add, err := strconv.Atoi(value)
				if err != nil {
					return "", errors.New("el valor de -add debe ser un número entero")
				}
				cmd.size = add // Usamos el mismo campo `size` para almacenar el valor de `add`
				addSet = true
			}
		} else if key == "-unit" {
			// Verifica que la unidad sea "K", "M" o "B"
			if value != "K" && value != "M" && value != "B" {
				return "", errors.New("la unidad debe ser K, M o B")
			}
			cmd.unit = strings.ToUpper(value)
			fmt.Printf("Unidad procesada: %s\n", cmd.unit) // Verificar la unidad
		} else if key == "-fit" {
			// Verifica que el ajuste sea "BF", "FF" o "WF"
			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return "", errors.New("el ajuste debe ser BF, FF o WF")
			}
			cmd.fit = value
		} else if key == "-path" {
			// Verifica que el path no esté vacío
			if value == "" {
				return "", errors.New("el path no puede estar vacío")
			}
			// Validar si el archivo existe
			if _, err := os.Stat(value); os.IsNotExist(err) {
				return "", fmt.Errorf("path no encontrado: %s", value)
			}
			cmd.path = value
		} else if key == "-type" {
			// Verifica que el tipo sea "P", "E" o "L"
			value = strings.ToUpper(value)
			if value != "P" && value != "E" && value != "L" {
				return "", errors.New("el tipo debe ser P, E o L")
			}
			cmd.typ = value
		} else if key == "-name" {
			// Verifica que el nombre no esté vacío
			if value == "" {
				return "", errors.New("el nombre no puede estar vacío")
			}
			cmd.name = value
		} else if key == "-delete" {
            // Procesar el comando delete
            if value != "fast" && value != "full" {
                return "", errors.New("el valor de -delete debe ser 'fast' o 'full'")
            }
            deleteType = value
            deleteSet = true
        } else {
			// Si el parámetro no es reconocido, devuelve un error
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -size, -path y -name hayan sido proporcionados
	// Verifica que al menos uno de los parámetros -size o -add haya sido proporcionado
	if cmd.path == "" {
		return "", errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return "", errors.New("faltan parámetros requeridos: -name")
	}

	// Si no se proporcionó la unidad, se establece por defecto a "M"
	if cmd.unit == "" {
		cmd.unit = "M"
	}

	// Si no se proporcionó el ajuste, se establece por defecto a "FF"
	if cmd.fit == "" {
		cmd.fit = "WF"
	}

	// Si no se proporcionó el tipo, se establece por defecto a "P"
	if cmd.typ == "" {
		cmd.typ = "P"
	}
	println("booleano de addSet:", addSet)
	if addSet {
		return handleAddSpace(cmd) // Llama a handleAddSpace si se utilizó -add
	}
	println("booleano de addSet:", addSet)
	if deleteSet {
        return NewDeletePartition(cmd, deleteType)
    }
	
		// Deserializar el MBR para verificar las particiones
		var mbr structures.MBR
		err := mbr.DeserializeMBR(cmd.path)
		if err != nil {
			return "", fmt.Errorf("error deserializando el MBR: %v", err)
		}
	
		// Validar si ya se alcanzó el límite de particiones
		if !mbr.HasAvailablePartition() {
			return "", errors.New("no se pueden agregar más particiones: las 4 particiones del MBR ya están ocupadas")
		}
	
		// Validar que el tamaño de la partición no exceda el tamaño del disco
		err = validatePartitionSize(cmd, &mbr)
		if err != nil {
			return "", err
		}
		// Crear la partición con los parámetros proporcionados
		err = commandFdisk(cmd)
		if err != nil {
			return "", fmt.Errorf("error: %v", err)
		}
			// Devuelve un mensaje de éxito con los detalles de la partición creada
		return fmt.Sprintf("FDISK: Partición creada exitosamente\n"+
		"-> Path: %s\n"+
		"-> Nombre: %s\n"+
		"-> Tamaño: %d %s\n"+
		"-> Tipo: %s\n"+
		"-> Fit: %s",
		cmd.path, cmd.name, cmd.size, cmd.unit, cmd.typ, cmd.fit), nil
	
}

func validatePartitionSize(fdisk *FDISK, mbr *structures.MBR) error {
    // Convertir el tamaño solicitado a bytes
    sizeBytes, err := utils.ConvertToBytes(fdisk.size, fdisk.unit)
    if err != nil {
        return fmt.Errorf("error convirtiendo el tamaño: %v", err)
    }

    // Calcular el espacio ya utilizado por las particiones existentes
    usedSpace := int32(0)
    for _, partition := range mbr.Mbr_partitions {
        if partition.Part_status[0] != 'N' { // Si la partición está activa o creada
            usedSpace += partition.Part_size
        }
    }

    // Verificar si el tamaño solicitado más el espacio utilizado excede el tamaño total del disco
    if usedSpace+int32(sizeBytes) > mbr.Mbr_size {
        return fmt.Errorf("el tamaño solicitado (%d bytes) excede el espacio disponible en el disco. Espacio disponible: %d bytes",
            sizeBytes, mbr.Mbr_size-usedSpace)
    }

    return nil
}
func commandFdisk(fdisk *FDISK) error {
	// Convertir el tamaño a bytes
	fmt.Printf("Unidad antes de la conversión: %s\n", fdisk.unit)
	sizeBytes, err := utils.ConvertToBytes(fdisk.size, fdisk.unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return err
	}
	fmt.Printf("Tamaño convertido a bytes: %d\n", sizeBytes)

	if fdisk.typ == "P" {
		// Crear partición primaria
		err = createPrimaryPartition(fdisk, sizeBytes)
		if err != nil {
			return fmt.Errorf("error creando partición primaria: %v", err)
		}
	} else if fdisk.typ == "E" {
        fmt.Println("Creando partición extendida...")
        // Llamar al método para crear la partición extendida
        err = createExtendedPartition(fdisk, sizeBytes)
        if err != nil {
            fmt.Println("Error creando partición extendida:", err)
            return err
        } // Les toca a ustedes implementar la partición extendida
	} else if fdisk.typ == "L" {
		fmt.Println("Creando partición lógica...")
        // Llamar al método CreateLogicalPartition del paquete structures
        err = CreateLogicalPartition(fdisk.path, int32(sizeBytes), fdisk.fit, fdisk.name)
        if err != nil {
            fmt.Println("Error creando partición lógica:", err)
            return err
        }
		
	}

	return nil
}

func createPrimaryPartition(fdisk *FDISK, sizeBytes int) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}


	/* SOLO PARA VERIFICACIÓN */
	// Imprimir MBR
	fmt.Println("\nMBR original:")
	mbr.PrintMBR()
	
	// Obtener la primera partición disponible
	availablePartition, startPartition, indexPartition := mbr.GetFirstAvailablePartition()
	
	
	// Validar si el espacio es inutilizable
	if availablePartition.Part_status[0] == 'I' {
		return errors.New("el espacio de la partición está marcado como inutilizable y no puede ser reutilizado")
	}

	/* SOLO PARA VERIFICACIÓN */
	// Print para verificar que la partición esté disponible
	fmt.Println("\nPartición disponible:")
	availablePartition.PrintPartition()

	// Crear la partición con los parámetros proporcionados
	availablePartition.CreatePartition(startPartition, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name)

	// Print para verificar que la partición se haya creado correctamente
	fmt.Println("\nPartición creada (modificada):")
	availablePartition.PrintPartition()

	// Colocar la partición en el MBR
	if availablePartition != nil {
		mbr.Mbr_partitions[indexPartition] = *availablePartition
	}

	// Imprimir las particiones del MBR
	fmt.Println("\nParticiones del MBR:")
	mbr.PrintPartitions()

	// Serializar el MBR en el archivo binario
	err = mbr.SerializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}


func createExtendedPartition(fdisk *FDISK, sizeBytes int) error {
    // Crear una instancia de MBR
    var mbr structures.MBR

    // Deserializar la estructura MBR desde un archivo binario
    err := mbr.DeserializeMBR(fdisk.path)
    if err != nil {
        return fmt.Errorf("error deserializando el MBR: %v", err)
    }

    // Verificar si ya existe una partición extendida
    for _, partition := range mbr.Mbr_partitions {
        if partition.Part_type == [1]byte{'E'} {
            return errors.New("ya existe una partición extendida")
			
        }
    }

    // Obtener la primera partición disponible
    availablePartition, startPartition, indexPartition := mbr.GetFirstAvailablePartition()
    if availablePartition == nil {
        return errors.New("no hay particiones disponibles")
    }

    // Crear la partición extendida con los parámetros proporcionados
    availablePartition.CreatePartition(startPartition, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name)
    availablePartition.Part_type = [1]byte{'E'}

    // Crear el primer EBR dentro de la partición extendida
    var ebr structures.EBR
    ebr.CreatePartition(int32(startPartition)+int32(binary.Size(mbr)), 30, "F", "") // El tamaño inicial del EBR es 30 bytes

    // Serializar el EBR en el archivo binario
    err = ebr.SerializeEBR(fdisk.path, int64(startPartition))
    if err != nil {
        return fmt.Errorf("error serializando el EBR: %v", err)
    }

    // Colocar la partición en el MBR
    mbr.Mbr_partitions[indexPartition] = *availablePartition

    // Imprimir las particiones del MBR
    fmt.Println("\nParticiones del MBR:")
    mbr.PrintPartitions()

    // Imprimir el EBR creado
    fmt.Println("\nEBR creado:")
    ebr.PrintEBR()

	


    // Serializar el MBR en el archivo binario
    err = mbr.SerializeMBR(fdisk.path)
    if err != nil {
        return fmt.Errorf("error serializando el MBR: %v", err)
    }

    return nil
}

func CreateLogicalPartition(path string, size int32, fit string, name string) error {
    file, err := os.OpenFile(path, os.O_RDWR, 0644)
    if err != nil {
        return fmt.Errorf("error abriendo el archivo: %v", err)
    }
    defer file.Close()

    var ebr structures.EBR
    offset := int64(0)

    // Leer el primer EBR de la partición extendida
    err = ebr.DeserializeEBR(path, offset)
    if err != nil {
        return fmt.Errorf("error leyendo el primer EBR: %v", err)
    }

    // Recorrer la lista enlazada de EBRs para encontrar el último
    for ebr.Part_next != -1 {
        offset = int64(ebr.Part_next)
        err = ebr.DeserializeEBR(path, offset)
        if err != nil {
            return fmt.Errorf("error leyendo el siguiente EBR: %v", err)
        }
    }

    // Calcular el inicio de la nueva partición lógica
    newPartitionStart := ebr.Part_start + ebr.Part_size + int32(binary.Size(ebr))

    // Verificar si hay suficiente espacio en la partición extendida
    if newPartitionStart+size > ebr.Part_start+ebr.Part_size {
        return errors.New("no hay suficiente espacio en la partición extendida para crear la partición lógica")
    }

    // Crear el nuevo EBR para la partición lógica
    newEBR := structures.EBR{}
    err = newEBR.CreatePartition(newPartitionStart, size, fit, name)
    if err != nil {
        return fmt.Errorf("error creando el nuevo EBR: %v", err)
    }

    // Actualizar el puntero `Part_next` del EBR actual
    ebr.Part_next = newPartitionStart
    err = ebr.SerializeEBR(path, offset)
    if err != nil {
        return fmt.Errorf("error actualizando el EBR actual: %v", err)
    }

    // Escribir el nuevo EBR en el archivo
    err = newEBR.SerializeEBR(path, int64(newPartitionStart))
    if err != nil {
        return fmt.Errorf("error escribiendo el nuevo EBR: %v", err)
    }

    fmt.Println("Partición lógica creada exitosamente.")
    return nil

}

func handleAddSpace(fdisk *FDISK) (string, error) {
    // Convertir el valor de `add` a bytes
    addBytes, err := utils.ConvertToBytes(fdisk.size, fdisk.unit)
    if err != nil {
        return "", fmt.Errorf("error convirtiendo el tamaño: %v", err)
    }

    // Deserializar el MBR
    var mbr structures.MBR
    err = mbr.DeserializeMBR(fdisk.path)
    if err != nil {
        return "", fmt.Errorf("error deserializando el MBR: %v", err)
    }

    // Buscar la partición por nombre
    var partition *structures.PARTITION
    for i := range mbr.Mbr_partitions {
        if strings.Trim(string(mbr.Mbr_partitions[i].Part_name[:]), "\x00") == fdisk.name {
            partition = &mbr.Mbr_partitions[i]
            break
        }
    }
    if partition == nil {
        return "", fmt.Errorf("la partición con nombre '%s' no existe", fdisk.name)
    }

    // Validar si es posible agregar o quitar espacio
    if addBytes > 0 {
        // Verificar que no exceda el tamaño del disco
        usedSpace := int32(0)
        for _, p := range mbr.Mbr_partitions {
            if p.Part_status[0] != 'N' {
                usedSpace += p.Part_size
            }
        }
        if usedSpace+int32(addBytes) > mbr.Mbr_size {
            return "", fmt.Errorf("no hay suficiente espacio en el disco para agregar %d bytes", addBytes)
        }
        partition.Part_size += int32(addBytes)
    }  else {
        // Verificar que no se reduzca más allá de 1 byte
        if partition.Part_size+int32(addBytes) < 1 {
            maxRemovable := partition.Part_size - 1
            return "", fmt.Errorf("no se puede reducir la partición '%s' más allá de su tamaño actual. "+
                "El máximo espacio que puede eliminar es %d bytes", fdisk.name, maxRemovable)
        }
        partition.Part_size += int32(addBytes)
    }

    // Serializar el MBR actualizado
    err = mbr.SerializeMBR(fdisk.path)
    if err != nil {
        return "", fmt.Errorf("error serializando el MBR: %v", err)
    }

    return fmt.Sprintf("FDISK: Espacio %s exitosamente a la partición '%s'\n"+
        "-> Tamaño actual: %d bytes",
        func() string {
            if addBytes > 0 {
                return "agregado"
            }
            return "reducido"
        }(), fdisk.name, partition.Part_size), nil
}

func NewDeletePartition(fdisk *FDISK, deleteType string) (string, error) {
    // Deserializar el MBR
    var mbr structures.MBR
    err := mbr.DeserializeMBR(fdisk.path)
    if err != nil {
        return "", fmt.Errorf("error deserializando el MBR: %v", err)
    }

    // Buscar la partición por nombre
    var partitionIndex int = -1
    for i := range mbr.Mbr_partitions {
        if strings.Trim(string(mbr.Mbr_partitions[i].Part_name[:]), "\x00") == fdisk.name {
            partitionIndex = i
            break
        }
    }

    if partitionIndex == -1 {
        return "", fmt.Errorf("la partición con nombre '%s' no existe", fdisk.name)
    }

    // Obtener la partición a eliminar
    partition := &mbr.Mbr_partitions[partitionIndex]

    if deleteType == "fast" {
        // Eliminar la partición de la lista (marcar como disponible)
        partition.Part_status[0] = 'N'
        partition.Part_name = [16]byte{}
        partition.Part_size = 0
        partition.Part_start = -1
        partition.Part_type = [1]byte{}
        partition.Part_fit = [1]byte{}
    } else if deleteType == "full" {
		fmt.Printf("llega el full") // Depuración
        // Sobrescribir el espacio de la partición con \0
        file, err := os.OpenFile(fdisk.path, os.O_RDWR, 0644)
        if err != nil {
            return "", fmt.Errorf("error abriendo el archivo: %v", err)
        }
        defer file.Close()

        // Sobrescribir el espacio con \0
        zeroBuffer := make([]byte, partition.Part_size)
        _, err = file.WriteAt(zeroBuffer, int64(partition.Part_start))
        if err != nil {
            return "", fmt.Errorf("error sobrescribiendo el espacio de la partición: %v", err)
        }

        // Eliminar la partición de la lista
        partition.Part_status[0] = 'I'
		fmt.Printf("Part_status después de asignar 'I': %c\n", partition.Part_status[0]) // Depuración
        partition.Part_name = [16]byte{}
        partition.Part_size = 0
        partition.Part_start = -1
        partition.Part_type = [1]byte{}
        partition.Part_fit = [1]byte{}
    } else {
        return "", fmt.Errorf("tipo de eliminación no válido: '%s'. Use 'fast' o 'full'", deleteType)
    }

	fmt.Println("\nEstado del MBR antes de la serialización:")
	mbr.PrintPartitions()
	

    // Serializar el MBR actualizado
    err = mbr.SerializeMBR(fdisk.path)
    if err != nil {
        return "", fmt.Errorf("error serializando el MBR: %v", err)
    }
	fmt.Println("\nEstado actual de las particiones del MBR:")
    mbr.PrintPartitions()

	// Imprimir el estado de la partición después de la asignación
    fmt.Printf("Estado de la partición después de eliminar (full): %+v\n", *partition)

    return fmt.Sprintf("FDISK: Partición '%s' eliminada exitosamente con el método '%s'.", fdisk.name, deleteType), nil
}