package commands

import (
    "fmt"
    "strings"
    global "terminal/global"
    "terminal/structures"
    stores "terminal/stores"
)

// ParseRename procesa el comando rename
func ParseRename(args []string) (string, error) {
    var path, newName string
    for i := 0; i < len(args); i++ {
        if strings.HasPrefix(args[i], "-path=") {
            path = strings.TrimPrefix(args[i], "-path=")
        } else if strings.HasPrefix(args[i], "-name=") {
            newName = strings.TrimPrefix(args[i], "-name=")
        }
    }
    if path == "" || newName == "" {
        return "", fmt.Errorf("debe especificar los parámetros -path y -name")
    }

    partitionID, err := stores.GetActivePartitionID()
    if err != nil {
        return "", fmt.Errorf("no hay partición activa montada")
    }
    sb, _, partitionPath, err := stores.GetMountedPartitionSuperblock(partitionID)
    if err != nil {
        return "", fmt.Errorf("no se pudo obtener el superbloque de la partición activa")
    }

    normalizedPath := path
    if !strings.HasPrefix(normalizedPath, "/") {
        normalizedPath = "/" + normalizedPath
    }
    normalizedPath = strings.ReplaceAll(normalizedPath, "//", "/")

    parts := strings.Split(strings.Trim(normalizedPath, "/"), "/")
    for i := 0; i < len(parts); i++ {
        tryPath := "/" + strings.Join(parts[i:], "/")
        err = RenameFile(sb, partitionPath, tryPath, newName)
        if err == nil {
            newPath := strings.TrimSuffix(tryPath, "/"+parts[len(parts)-1]) + "/" + newName

            // Imprimir lista antes de cambiar
            fmt.Println("Lista antes del cambio:", global.GetValidFilePathsMkfile())

            // Actualizar la lista global usando get y set (comparando solo el nombre del archivo)
            fmt.Println("Renombrando archivo en la lista global:", tryPath, "a", newPath)
            paths := global.GetValidFilePathsMkfile()
            for i, p := range paths {
                // Extrae el nombre del archivo del path
                pathParts := strings.Split(p, "/")
                if len(pathParts) > 0 && pathParts[len(pathParts)-1] == parts[len(parts)-1] {
                    // Reconstruye el nuevo path manteniendo la ruta original
                    newPathFull := strings.TrimSuffix(p, pathParts[len(pathParts)-1]) + newName
                    paths[i] = newPathFull
                    break
                }
            }
            global.SetValidFilePathsMkfile(paths)

            // Imprimir lista después de cambiar
            fmt.Println("Lista después del cambio:", global.GetValidFilePathsMkfile())

            // También puedes actualizar la lista de archivos extraídos si lo necesitas
            oldName := parts[len(parts)-1]
            structures.RenameTxtFileInExtractedByName(oldName, newName)

            return fmt.Sprintf("RENAME: Archivo renombrado a: %s", newName), nil
        }
    }
    return "", fmt.Errorf("no se encontró el archivo en ninguna ruta posible a partir de: %s", path)
}

// RenameFile cambia el nombre del archivo en el bloque de carpeta
func RenameFile(sb *structures.SuperBlock, partitionPath string, filePath string, newName string) error {
    normalizedPath := filePath
    if !strings.HasPrefix(normalizedPath, "/") {
        normalizedPath = "/" + normalizedPath
    }
    normalizedPath = strings.ReplaceAll(normalizedPath, "//", "/")

    inodeIndex, err := sb.FindInodeByPath(partitionPath, normalizedPath)
    if err != nil || inodeIndex < 0 {
        return fmt.Errorf("archivo no encontrado: %s", normalizedPath)
    }

    parts := strings.Split(strings.Trim(normalizedPath, "/"), "/")
    parentPath := "/" + strings.Join(parts[:len(parts)-1], "/")
    oldName := parts[len(parts)-1]

    parentInodeIndex, err := sb.FindInodeByPath(partitionPath, parentPath)
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
            continue
        }
        for i, content := range block.B_content {
            name := strings.Trim(string(content.B_name[:]), "\x00 ")
            cmpName := oldName
            if len(cmpName) > 12 {
                cmpName = cmpName[:12]
            }
            if len(name) > 12 {
                name = name[:12]
            }
            if name == cmpName && content.B_inodo == inodeIndex {
                // Cambiar el nombre
                for j := range block.B_content[i].B_name {
                    block.B_content[i].B_name[j] = 0
                }
                copy(block.B_content[i].B_name[:], []byte(newName))
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
        return fmt.Errorf("no se encontró la entrada del archivo en la carpeta padre")
    }
    // Imprimir inodos y bloques para validar el cambio
    fmt.Println("\n--- INODOS DESPUÉS DEL RENAME ---")
    sb.PrintInodes(partitionPath)
    fmt.Println("\n--- BLOQUES DESPUÉS DEL RENAME ---")
    sb.PrintBlocks(partitionPath)
    return nil
}