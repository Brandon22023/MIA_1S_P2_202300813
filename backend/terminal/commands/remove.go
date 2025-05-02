package commands

import (
    "fmt"
    "strings"
    global "terminal/global"
    "terminal/structures"
	 stores "terminal/stores"
)

func ParseRemove(args []string) (string, error) {
    var path string
    for i := 0; i < len(args); i++ {
        if strings.HasPrefix(args[i], "-path=") {
            path = strings.TrimPrefix(args[i], "-path=")
        }
    }
    if path == "" {
        return "", fmt.Errorf("debe especificar el parámetro -path")
    }

    partitionID, err := stores.GetActivePartitionID()
    if err != nil {
        return "", fmt.Errorf("no hay partición activa montada")
    }
    sb, _, partitionPath, err := stores.GetMountedPartitionSuperblock(partitionID)
    if err != nil {
        return "", fmt.Errorf("no se pudo obtener el superbloque de la partición activa")
    }

    // Intentar eliminar en la ruta original y luego ir subiendo niveles
    normalizedPath := path
    if !strings.HasPrefix(normalizedPath, "/") {
        normalizedPath = "/" + normalizedPath
    }
    normalizedPath = strings.ReplaceAll(normalizedPath, "//", "/")

    parts := strings.Split(strings.Trim(normalizedPath, "/"), "/")
    for i := 0; i < len(parts); i++ {
		tryPath := "/" + strings.Join(parts[i:], "/")
		err = RemoveFile(sb, partitionPath, tryPath)
		if err == nil {
			// Extraer el nombre del archivo del tryPath
			pathParts := strings.Split(tryPath, "/")
			fileName := pathParts[len(pathParts)-1]
			return fmt.Sprintf("REMOVE: Archivo eliminado: %s", fileName), nil
		}
	}
    return "", fmt.Errorf("no se encontró el archivo en ninguna ruta posible a partir de: %s", path)
}

func RemoveFile(sb *structures.SuperBlock, partitionPath string, filePath string) error {
    fmt.Println("========== INICIO RemoveFile ==========")
    normalizedPath := filePath
    if !strings.HasPrefix(normalizedPath, "/") {
        normalizedPath = "/" + normalizedPath
    }
    normalizedPath = strings.ReplaceAll(normalizedPath, "//", "/")
    fmt.Println("[Depuración] Path normalizado:", normalizedPath)

    // 1. Buscar el inodo del archivo
    inodeIndex, err := sb.FindInodeByPath(partitionPath, normalizedPath)
    fmt.Printf("[Depuración] InodeIndex para archivo '%s': %d (err: %v)\n", normalizedPath, inodeIndex, err)
    if err != nil || inodeIndex < 0 {
        return fmt.Errorf("archivo no encontrado: %s", normalizedPath)
    }

    // 2. Buscar el bloque de carpeta y la posición donde está el archivo
    parts := strings.Split(strings.Trim(normalizedPath, "/"), "/")
    parentPath := "/" + strings.Join(parts[:len(parts)-1], "/")
    fileName := parts[len(parts)-1]
    fmt.Printf("[Depuración] ParentPath: '%s', fileName: '%s'\n", parentPath, fileName)

    parentInodeIndex, err := sb.FindInodeByPath(partitionPath, parentPath)
    fmt.Printf("[Depuración] InodeIndex para carpeta padre '%s': %d (err: %v)\n", parentPath, parentInodeIndex, err)
    if err != nil || parentInodeIndex < 0 {
        return fmt.Errorf("carpeta padre no encontrada: %s", parentPath)
    }

    parentInode := &structures.Inode{}
    err = parentInode.Deserialize(partitionPath, int64(sb.S_inode_start+(parentInodeIndex*sb.S_inode_size)))
    if err != nil {
        return fmt.Errorf("error deserializando inodo padre: %v", err)
    }

    found := false
    for _, blockIndex := range parentInode.I_block {
        if blockIndex == -1 {
            continue
        }
        block := &structures.FolderBlock{}
        err := block.Deserialize(partitionPath, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
        if err != nil {
            fmt.Printf("[Depuración] Error deserializando bloque de carpeta %d: %v\n", blockIndex, err)
            continue
        }
        for i, content := range block.B_content {
            name := strings.Trim(string(content.B_name[:]), "\x00 ")
            cmpName := fileName
            if len(cmpName) > 12 {
                cmpName = cmpName[:12]
            }
            if len(name) > 12 {
                name = name[:12]
            }
            fmt.Printf("[Depuración] Comparando bloque %d, content %d: name='%s' (inodo=%d) vs fileName='%s' (inodeIndex=%d)\n", blockIndex, i, name, content.B_inodo, cmpName, inodeIndex)
            if name == cmpName && content.B_inodo == inodeIndex {
                fmt.Printf("[Depuración] ¡Coincidencia encontrada! Eliminando entrada...\n")
                // Eliminar la entrada del archivo en el bloque de carpeta
                block.B_content[i].B_inodo = -1
                for j := range block.B_content[i].B_name {
                    block.B_content[i].B_name[j] = 0
                }
                // Serializar el bloque actualizado
                block.Serialize(partitionPath, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
                found = true
                break
            }
        }
        if found {
            break
        }
    }
    if !found {
        fmt.Println("[Depuración] No se encontró la entrada del archivo en la carpeta padre")
        return fmt.Errorf("no se encontró la entrada del archivo en la carpeta padre")
    }

    // 3. Liberar el inodo y los bloques de datos del archivo
    inode := &structures.Inode{}
    err = inode.Deserialize(partitionPath, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
    if err == nil {
        fmt.Printf("[Depuración] Liberando inodo %d y sus bloques de datos\n", inodeIndex)
        // Liberar bloques de datos
        for _, blockIndex := range inode.I_block {
            if blockIndex == -1 {
                break
            }
            fmt.Printf("[Depuración] (Opcional) Liberar bloque de datos %d\n", blockIndex)
            // Aquí podrías limpiar el bloque si lo deseas
        }
        // Aquí podrías limpiar el inodo si lo deseas
    }

    // 4. Eliminar el path de la lista global
    for i, p := range global.ValidFilePaths_mkfile {
        fmt.Printf("[Depuración] Comparando path global: '%s' con '%s'\n", p, normalizedPath)
        if p == normalizedPath {
            fmt.Printf("[Depuración] Eliminando path de la lista global\n")
            global.ValidFilePaths_mkfile = append(global.ValidFilePaths_mkfile[:i], global.ValidFilePaths_mkfile[i+1:]...)
            break
        }
    }

    fmt.Printf("Archivo eliminado correctamente: %s\n", normalizedPath)
    fmt.Println("========== FIN RemoveFile ==========")
    return nil
}