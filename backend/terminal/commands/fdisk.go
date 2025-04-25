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

// CommandFdisk parsea el comando fdisk y devuelve una instancia de FDISK
func ParseFdisk(tokens []string) (string, error) {
	cmd := &FDISK{} // Crea una nueva instancia de FDISK

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando fdisk
	re := regexp.MustCompile(`-size=\d+|-unit=[kKmMbB]|-fit=[bBfF]{2}|-path="[^"]+"|-path=[^\s]+|-type=[pPeElL]|-name="[^"]+"|-name=[^\s]+`)
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

		// Switch para manejar diferentes parámetros
		switch key {
		case "-size":
			// Convierte el valor del tamaño a un entero
			size, err := strconv.Atoi(value)
			if err != nil || size <= 0 {
				return "", errors.New("el tamaño debe ser un número entero positivo")
			}
			cmd.size = size
		case "-unit":
			// Verifica que la unidad sea "K", "M" o "B"
			if value != "K" && value != "M" && value != "B" {
				return "", errors.New("la unidad debe ser K, M o B")
			}

			cmd.unit = strings.ToUpper(value)
			fmt.Printf("Unidad procesada: %s\n", cmd.unit) // Verificar la unidad
		case "-fit":
			// Verifica que el ajuste sea "BF", "FF" o "WF"
			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return "", errors.New("el ajuste debe ser BF, FF o WF")
			}
			cmd.fit = value
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return "", errors.New("el path no puede estar vacío")
			}
			// Validar si el archivo existe
            if _, err := os.Stat(value); os.IsNotExist(err) {
                return "", fmt.Errorf("path no encontrado: %s", value)
            }
			cmd.path = value
		case "-type":
			// Verifica que el tipo sea "P", "E" o "L"
			value = strings.ToUpper(value)
			if value != "P" && value != "E" && value != "L" {
				return "", errors.New("el tipo debe ser P, E o L")
			}
			cmd.typ = value
		case "-name":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return "", errors.New("el nombre no puede estar vacío")
			}
			cmd.name = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -size, -path y -name hayan sido proporcionados
	if cmd.size == 0 {
		return "", errors.New("faltan parámetros requeridos: -size")
	}
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
		fmt.Println("Error:", err)
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
			fmt.Println("Error creando partición primaria:", err)
			return err
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
func isMBRFull(mbr *structures.MBR) bool {
    for _, partition := range mbr.Mbr_partitions {
        if partition.Part_status[0] != 1 { // Si hay al menos una partición inactiva
            return false
        }
    }
    return true // Todas las particiones están activas
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
	if availablePartition == nil {
		fmt.Println("No hay particiones disponibles.")
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
/*
func printLogicalPartitions(path string, start int32) error {
    var ebr structures.EBR
    offset := int64(start)

    fmt.Println("\nParticiones lógicas dentro de la partición extendida:")
	if ebr.Part_size == 0 && ebr.Part_next == -1 {
		fmt.Println("No hay particiones lógicas creadas aún.")
		return nil
	}

    for {
        // Leer el EBR desde el archivo
        err := ebr.DeserializeEBR(path, offset)
        if err != nil {
            return fmt.Errorf("error leyendo el EBR en offset %d: %v", offset, err)
        }

        // Imprimir la información del EBR
        ebr.PrintEBR()

        // Verificar si hay un siguiente EBR
        if ebr.Part_next == -1 {
            break // No hay más particiones lógicas
        }

        // Mover al siguiente EBR
        offset = int64(ebr.Part_next)
    }

    return nil
}
*/