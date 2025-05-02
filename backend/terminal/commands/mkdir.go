package commands

import (
    "errors"
    "fmt"
    "os"
    "regexp"
    "strings"
    "terminal/global"
    stores "terminal/stores"
    structures "terminal/structures"
    utils "terminal/utils"
)

// MKDIR estructura que representa el comando mkdir con sus parámetros
type MKDIR struct {
    path string // Path del directorio
    p    bool   // Opción -p (crea directorios padres si no existen)
}

func ParseMkdir(tokens []string) (string, error) {
    cmd := &MKDIR{} // Crea una nueva instancia de MKDIR

    // Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
    args := strings.Join(tokens, " ")
    re := regexp.MustCompile(`-path=("[^"]+"|[^\s]+)|-p`)
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
    parentDirs, _ := utils.GetParentDirectories(mkdir.path)
    currentParents := []string{}
    for _, parent := range parentDirs {
        currentParents = append(currentParents, parent)
        exists, err := partitionSuperblock.FolderExists(partitionPath, strings.Join(currentParents, "/"))
        if err != nil {
            return fmt.Errorf("error al verificar la existencia de la carpeta '%s': %w", parent, err)
        }
        if exists {
            inode, err := partitionSuperblock.FindInode(partitionPath, currentParents[:len(currentParents)-1], parent)
            if err != nil {
                return fmt.Errorf("error al obtener el inodo de la carpeta '%s': %w", parent, err)
            }
            if !hasWritePermission(inode, currentUser) {
                return fmt.Errorf("error: el usuario '%s' no tiene permisos de escritura en la carpeta '%s'", currentUser, parent)
            }
        }
    }

    // Crear el directorio (y padres si aplica)
    err = createDirectory(mkdir.path, partitionSuperblock, partitionPath, mountedPartition, mkdir.p)
    if err != nil {
        return fmt.Errorf("error al crear el directorio: %w", err)
    }

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

    // Si el dirPath es "/", requerir permisos de superusuario
    if dirPath == "/" {
        if os.Geteuid() != 0 {
            return fmt.Errorf("error: se requieren permisos de superusuario para crear el directorio raíz")
        }
    }

    // Validar si las carpetas padres existen y crearlas correctamente
    currentPath := ""
    for i, parent := range parentDirs {
        if i == 0 && parent == "" {
            continue
        }
        if currentPath == "" {
            currentPath = parent
        } else {
            currentPath = currentPath + "/" + parent
        }
        exists, err := sb.FolderExists(partitionPath, currentPath)
        fmt.Printf("Verificando existencia de la carpeta: %s\n", currentPath)
        if err != nil {
            return fmt.Errorf("error al verificar la existencia de la carpeta '%s': %w", currentPath, err)
        }
        if !exists {
            if !allowParents {
                return fmt.Errorf("error: no existen las carpetas padres para el directorio '%s'", dirPath)
            }
            // Crear las carpetas padres si la opción -p está habilitada
            fmt.Printf("Creando carpeta padre: %s\n", currentPath)
            err = sb.CreateFolder(partitionPath, parentDirs[:i], parent)
            if err != nil {
                return fmt.Errorf("error al crear la carpeta padre '%s': %w", currentPath, err)
            }
            // Guardar el path creado
            if !contains_m(global.ValidPaths, currentPath) {
                global.ValidPaths = append(global.ValidPaths, currentPath)
            }
        }
    }

    // Crear el directorio según el path proporcionado
    err := sb.CreateFolder(partitionPath, parentDirs, destDir)
    if err != nil {
        return fmt.Errorf("error al crear el directorio: %w", err)
    }
    normalizedPath := dirPath
    if !strings.HasPrefix(normalizedPath, "/") {
        normalizedPath = "/" + normalizedPath
    }
    normalizedPath = strings.ReplaceAll(normalizedPath, "//", "/")
    if !contains_m(global.ValidPaths, normalizedPath) {
        global.ValidPaths = append(global.ValidPaths, normalizedPath)
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