package commands

import (
	stores "terminal/stores"
	structures "terminal/structures"
	utils "terminal/utils"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"os"
	"path/filepath"
)

// MKDIR estructura que representa el comando mkdir con sus parámetros
type MKDIR struct {
	path string // Path del directorio
	p    bool   // Opción -p (crea directorios padres si no existen)
}

/*
   mkdir -p -path=/home/user/docs/usac
   mkdir -path="/home/mis documentos/archivos clases"
*/

func ParseMkdir(tokens []string) (string, error) {
	cmd := &MKDIR{} // Crea una nueva instancia de MKDIR

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mkdir
	re := regexp.MustCompile(`-path=[^\s]+|-p`)
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

		// Switch para manejar diferentes parámetros
		switch key {
		case "-path":
			if len(kv) != 2 {
				return "", fmt.Errorf("formato de parámetro inválido: %s", match)
			}
			value := kv[1]
			// Remove quotes from value if present
			if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
				value = strings.Trim(value, "\"")
			}
			cmd.path = value
		case "-p":
			cmd.p = true
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -path haya sido proporcionado
	if cmd.path == "" {
		return "", errors.New("faltan parámetros requeridos: -path")
	}

	// Aquí se puede agregar la lógica para ejecutar el comando mkdir con los parámetros proporcionados
	err := commandMkdir(cmd)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("MKDIR: Directorio %s creado correctamente.", cmd.path), nil // Devuelve el comando MKDIR creado
}

// Aquí debería de estar logeado un usuario, por lo cual el usuario debería tener consigo el id de la partición

func commandMkdir(mkdir *MKDIR) error {
	// Obtener el id de la partición montada que está logueada
	var partitionID string

	partitionID, err := stores.GetActivePartitionID()
    if err != nil {
        return err
    }

	// Obtener la partición montada
	partitionSuperblock, mountedPartition, partitionPath, err := stores.GetMountedPartitionSuperblock(partitionID)
	if err != nil {
		return fmt.Errorf("error al obtener la partición montada: %w", err)
	}

	// Crear el directorio
	err = createDirectory(mkdir.path, partitionSuperblock, partitionPath, mountedPartition, mkdir.p)
	if err != nil {
		return fmt.Errorf("error al crear el directorio: %w", err)
	}

	return err
}

func createDirectory(dirPath string, sb *structures.SuperBlock, partitionPath string, mountedPartition *structures.PARTITION, allowParents bool) error {
	fmt.Println("\nCreando directorio:", dirPath)

	parentDirs, destDir := utils.GetParentDirectories(dirPath)
	fmt.Println("\nDirectorios padres:", parentDirs)
	fmt.Println("Directorio destino:", destDir)

	// Si el dirPath es "/", requerir permisos de superusuario
	if dirPath == "/" {
		if os.Geteuid() != 0 {
			return fmt.Errorf("error: se requieren permisos de superusuario para crear el directorio raíz")
		}
	}

	physicalBasePath := dirPath

	// Validar si las carpetas padres existen
    for _, parent := range parentDirs {
        exists, err := sb.FolderExists(partitionPath, parent)
		fmt.Printf("Verificando existencia de la carpeta: %s\n", parent)
        if err != nil {
            return fmt.Errorf("error al verificar la existencia de la carpeta '%s': %w", parent, err)
        }
        if !exists {
            if !allowParents {
                return fmt.Errorf("error: no existen las carpetas padres para el directorio '%s'", dirPath)
            }
            // Crear las carpetas padres si la opción -p está habilitada
			fmt.Printf("Creando carpeta padre: %s\n", parent)
			err = sb.CreateFolder(partitionPath, parentDirs, parent)
			if err != nil {
				return fmt.Errorf("error al crear la carpeta padre '%s': %w", parent, err)
			}

			// Crear físicamente la carpeta en el sistema operativo dentro de la ruta definida
			physicalPath := filepath.Join(physicalBasePath, parent)
			err = os.MkdirAll(physicalPath, 0755)
			if err != nil {
				return fmt.Errorf("error al crear físicamente la carpeta '%s': %w", physicalPath, err)
			}
			fmt.Printf("Carpeta creada físicamente: %s\n", physicalPath)
		}
    }

	// Crear físicamente el directorio destino dentro de la ruta definida
	fullPath := filepath.Join(physicalBasePath, filepath.Join(filepath.Join(parentDirs...), destDir))
	fullPath = filepath.Clean(fullPath)
	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		return fmt.Errorf("error al crear físicamente el directorio '%s': %w", fullPath, err)
	}
	fmt.Printf("Directorio creado físicamente: %s\n", fullPath)

	// Crear el directorio segun el path proporcionado
	err = sb.CreateFolder(partitionPath, parentDirs, destDir)
	if err != nil {
		return fmt.Errorf("error al crear el directorio: %w", err)
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
