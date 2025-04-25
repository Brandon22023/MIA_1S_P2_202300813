package commands

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"os"
	"path/filepath"
	stores "terminal/stores"
	structures "terminal/structures"
	utils "terminal/utils"
)

// MKFILE estructura que representa el comando mkfile con sus parámetros
type MKFILE struct {
	path string // Ruta del archivo
	r    bool   // Opción recursiva
	size int    // Tamaño del archivo
	cont string // Contenido del archivo
}

// ParserMkfile parsea el comando mkfile y devuelve una instancia de MKFILE
func ParserMkfile(tokens []string) (string, error) {
	cmd := &MKFILE{} // Crea una nueva instancia de MKFILE

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mkfile
	re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+|-r|-size=\d+|-cont="[^"]+"|-cont=[^\s]+`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Verificar que todos los tokens fueron reconocidos por la expresión regular
	if len(matches) != len(tokens) {
		// Identificar el parámetro inválido
		for _, token := range tokens {
			if !re.MatchString(token) {
				return "", fmt.Errorf("parámetro inválido: %s", token)
			}
		}
	}

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		key := strings.ToLower(kv[0])
		var value string
		if len(kv) == 2 {
			value = kv[1]
		}

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return "", errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		case "-r":
			// Establece el valor de r a true
			cmd.r = true
		case "-size":
			// Convierte el valor del tamaño a un entero
			size, err := strconv.Atoi(value)
			if err != nil || size < 0 {
				return "", errors.New("el tamaño debe ser un número entero no negativo")
			}
			cmd.size = size
		case "-cont":
			// Verifica que el contenido no esté vacío
			if value == "" {
				return "", errors.New("el contenido no puede estar vacío")
			}
			cmd.cont = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -path haya sido proporcionado
	if cmd.path == "" {
		return "", errors.New("faltan parámetros requeridos: -path")
	}

	// Si no se proporcionó el tamaño, se establece por defecto a 0
	if cmd.size == 0 {
		cmd.size = 0
	}

	// Si no se proporcionó el contenido, se establece por defecto a ""
	if cmd.cont == "" {
		cmd.cont = ""
	}

	// Crear el archivo con los parámetros proporcionados
	err := commandMkfile(cmd)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("MKFILE:Archivo creado exitosamente\n"+
	 "-> Archivo: %s\n", cmd.path), nil // Devuelve el comando MKFILE creado

	
}

// Función ficticia para crear el archivo (debe ser implementada)
func commandMkfile(mkfile *MKFILE, ) error {
	// Obtener el id de la partición montada que está logueada
    partitionID, err := stores.GetActivePartitionID()
    if err != nil {
        return err
    }
	// Obtener la partición montada
	partitionSuperblock, mountedPartition, partitionPath, err := stores.GetMountedPartitionSuperblock(partitionID)
	
	if err != nil {
		return fmt.Errorf("error al obtener la partición montada: %w", err)
	}

	// Generar el contenido del archivo si no se proporcionó
	if mkfile.cont == "" {
		mkfile.cont = generateContent(mkfile.size)
	}

	// Crear el archivo
	err = createFile(mkfile.path, mkfile.size, mkfile.cont, partitionSuperblock, partitionPath, mountedPartition, mkfile.r)
	if err != nil {
		err = fmt.Errorf("error al crear el archivo: %w", err)
	}

	return err
}

// generateContent genera una cadena de números del 0 al 9 hasta cumplir el tamaño ingresado
func generateContent(size int) string {
	content := ""
	for len(content) < size {
		content += "0123456789"
	}
	return content[:size] // Recorta la cadena al tamaño exacto
}

// Funcion para crear un archivo
func createFile(filePath string, size int, content string, sb *structures.SuperBlock, partitionPath string, mountedPartition *structures.PARTITION, recursive bool) error {
	fmt.Println("\nCreando archivo:", filePath)

	parentDirs, destDir := utils.GetParentDirectories(filePath)
	fmt.Println("\nDirectorios padres:", parentDirs)
	fmt.Println("Directorio destino:", destDir)

	// Validar si el parámetro `content` es un path a un archivo
    if content != "" {
        if _, err := os.Stat(content); err == nil {
            // Leer el contenido del archivo especificado en `content`
            fileData, err := os.ReadFile(content)
            if err != nil {
                return fmt.Errorf("error al leer el archivo '%s': %w", content, err)
            }
            content = string(fileData) // Usar el contenido del archivo como contenido
            fmt.Printf("Contenido leído del archivo '%s':\n%s\n", content, fileData)
        } else {
            fmt.Printf("El parámetro `content` no es un path válido, se usará como contenido directo.\n")
        }
    }

	// Verificar si el directorio base existe físicamente
    baseDir := filepath.Dir(filePath) // Obtiene el directorio base del filePath
    if _, err := os.Stat(baseDir); os.IsNotExist(err) {
        if recursive {
            // Crear los directorios faltantes si la opción -r está habilitada
            fmt.Printf("Creando directorios faltantes: %s\n", baseDir)
            err := os.MkdirAll(baseDir, 0755)
            if err != nil {
                return fmt.Errorf("error al crear los directorios faltantes '%s': %w", baseDir, err)
            }
            fmt.Printf("Directorios creados: %s\n", baseDir)
        } else {
            return fmt.Errorf("error: el directorio base '%s' no existe", baseDir)
        }
    } else {
        fmt.Printf("Directorio base verificado: %s\n", baseDir)
    }

	// Si el contenido proviene de un archivo, reemplazar los chunks con el contenido leído
    var chunks []string
    if content != "" {
        chunks = []string{content} // Usar el contenido leído como un único chunk
    } else {
        // Obtener contenido por chunks si no se proporcionó un archivo en `-cont`
        chunks = utils.SplitStringIntoChunks(generateContent(size))
    }
    fmt.Println("\nChunks del contenido:", chunks)

    // Crear el archivo físicamente en el sistema operativo
    file, err := os.Create(filePath) // Crea el archivo en la ruta especificada
    if err != nil {
        return fmt.Errorf("error al crear el archivo físicamente: %w", err)
    }
    defer file.Close()

    // Escribir los chunks en el archivo
    for _, chunk := range chunks {
        _, err := file.WriteString(chunk)
        if err != nil {
            return fmt.Errorf("error al escribir en el archivo: %w", err)
        }
    }
    fmt.Printf("Archivo '%s' creado físicamente con contenido.\n", filePath)



	// Crear el archivo
	err = sb.CreateFile(partitionPath, parentDirs, destDir, size, chunks)
	if err != nil {
		return fmt.Errorf("error al crear el archivo: %w", err)
	}

	// Imprimir inodos y bloques
	sb.PrintInodes(partitionPath)
	sb.PrintBlocks(partitionPath)

	// Serializar el superbloque
	err = sb.Serialize(partitionPath, int64(mountedPartition.Part_start))
	if err != nil {
		return fmt.Errorf("error al serializar el superbloque: %w", err)
	}

	return nil
}

