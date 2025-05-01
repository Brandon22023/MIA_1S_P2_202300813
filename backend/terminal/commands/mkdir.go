package commands

import (
	stores "terminal/stores"
	structures "terminal/structures"
	utils "terminal/utils"
	"errors"
	"fmt"
	"regexp"
	"strings"


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
    re := regexp.MustCompile(`-path=[^\s]+|-p`)
    matches := re.FindAllString(args, -1)

    if len(matches) != len(tokens) {
        for _, token := range tokens {
            if !re.MatchString(token) {
                return "", fmt.Errorf("parámetro inválido: %s", token)
            }
        }
    }

    for _, match := range matches {
        kv := strings.SplitN(match, "=", 2)
        key := strings.ToLower(kv[0])

        switch key {
        case "-path":
            if len(kv) != 2 {
                return "", fmt.Errorf("formato de parámetro inválido: %s", match)
            }
            value := kv[1]
            if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
                value = strings.Trim(value, "\"")
            }
            cmd.path = value

            // Validar y agregar el path
            if cmd.path != "" {
                segments := strings.Split(cmd.path, "/")
                currentPath := ""
                valid := true

                if len(segments) == 2 && segments[1] != "" {
                    stores.ValidPaths = append(stores.ValidPaths, cmd.path)
                    break
                }

                for _, segment := range segments {
                    if segment == "" {
                        continue
                    }
                    currentPath += "/" + segment
                    if !contains_m(stores.ValidPaths, currentPath) {
                        valid = false
                        break
                    }
                }

                if valid || cmd.p {
                    stores.ValidPaths = append(stores.ValidPaths, cmd.path)
                } else {
                    return "", fmt.Errorf("el path '%s' no es válido porque faltan directorios intermedios", cmd.path)
                }
            }
        case "-p":
            cmd.p = true
        default:
            return "", fmt.Errorf("parámetro desconocido: %s", key)
        }
    }

    if cmd.path == "" {
        return "", errors.New("faltan parámetros requeridos: -path")
    }

    err := commandMkdir(cmd)
    if err != nil {
        fmt.Printf("[ERROR] %v\n", err)
        return "", err
    }

    return fmt.Sprintf("MKDIR: Directorio %s creado correctamente.", cmd.path), nil
}

// Función auxiliar para verificar si un path ya está en la lista
func contains_m(paths []string, path string) bool {
    for _, p := range paths {
        if p == path {
            return true
        }
    }
    return false
}

// Aquí debería de estar logeado un usuario, por lo cual el usuario debería tener consigo el id de la partición

func commandMkdir(mkdir *MKDIR) error {
    // Obtener el ID de la partición activa
    partitionID, err := stores.GetActivePartitionID()
    fmt.Println("ID de la partición activa:", partitionID)
    if err != nil {
        return fmt.Errorf("error al obtener el ID de la partición activa: %w", err)
    }

    // Verificar si hay un usuario autenticado
    if !stores.Auth.IsAuthenticated() {
        return fmt.Errorf("error: no hay un usuario autenticado")
    }

    // Obtener el usuario autenticado
    currentUser := stores.Auth.Username

    // Obtener la partición montada
    partitionSuperblock, mountedPartition, partitionPath, err := stores.GetMountedPartitionSuperblock(partitionID)
    if err != nil {
        return fmt.Errorf("error al obtener la partición montada: %w", err)
    }

    // Validar permisos de escritura en la carpeta padre
    parentDirs, destDir := utils.GetParentDirectories(mkdir.path)
    // Solo valida permisos sobre el padre inmediato si existe
    if len(parentDirs) > 0 {
        parentPath := "/" + strings.Join(parentDirs, "/")
        exists, err := partitionSuperblock.FolderExists(partitionPath, parentPath)
        if err != nil {
            return fmt.Errorf("error al verificar la existencia de la carpeta '%s': %w", parentPath, err)
        }
        if exists {
            inode, err := partitionSuperblock.FindInode(partitionPath, parentDirs[:len(parentDirs)-1], parentDirs[len(parentDirs)-1])
            if err != nil {
                return fmt.Errorf("error al obtener el inodo de la carpeta '%s': %w", parentPath, err)
            }
            if !hasWritePermission(inode, currentUser) {
                return fmt.Errorf("error: el usuario '%s' no tiene permisos de escritura en la carpeta '%s'", currentUser, parentPath)
            }
        }
    }

    // Crear la carpeta dentro de la partición (esto ya es recursivo)
    err = partitionSuperblock.CreateFolder(partitionPath, parentDirs, destDir, int64(mountedPartition.Part_start))
    if err != nil {
        return fmt.Errorf("error al crear la carpeta '%s': %w", mkdir.path, err)
    }

    // Serializar el superbloque para guardar los cambios
    err = partitionSuperblock.Serialize(partitionPath, int64(mountedPartition.Part_start))
    if err != nil {
        return fmt.Errorf("error al serializar el superbloque: %w", err)
    }

    // Imprimir inodos y bloques (opcional)
    partitionSuperblock.PrintInodes(partitionPath)
    partitionSuperblock.PrintBlocks(partitionPath)

    fmt.Printf("Directorio '%s' creado correctamente con propietario '%s'\n", mkdir.path, currentUser)
    return nil
}

// Verificar si el usuario tiene permisos de escritura
func hasWritePermission(inode *structures.Inode, user string) bool {
    // Verificar permisos de escritura (bit 2 de los permisos)
    permissions := string(inode.I_perm[:])
    return permissions[1] == '6' || permissions[1] == '2' // Permiso de escritura
}

func createDirectory(dirPath string, sb *structures.SuperBlock, partitionPath string, mountedPartition *structures.PARTITION, allowParents bool) error {
    fmt.Println("\nCreando directorio:", dirPath)

    parentDirs, destDir := utils.GetParentDirectories(dirPath)
    fmt.Println("\nDirectorios padres:", parentDirs)
    fmt.Println("Directorio destino:", destDir)

    // Solo llama a CreateFolder, que ya es recursiva
    //err := sb.CreateFolder(partitionPath, parentDirs, destDir)
    /*if err != nil {
        return fmt.Errorf("error al crear el directorio: %w", err)
    }*/

    // Imprimir inodos y bloques
    sb.PrintInodes(partitionPath)
    sb.PrintBlocks(partitionPath)

    // Serializar el superbloque
    err := sb.Serialize(partitionPath, int64(mountedPartition.Part_start))
    if err != nil {
        return fmt.Errorf("error al serializar el superbloque: %w", err)
    }

    return nil
}

