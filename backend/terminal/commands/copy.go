package commands

import (
    "fmt"
    "strings"
    global "terminal/global"
    "terminal/structures"
    stores "terminal/stores"
)

// ParseCopy procesa el comando copy
func ParseCopy(args []string) (string, error) {
    var srcPath, dstPath string
    for i := 0; i < len(args); i++ {
        if strings.HasPrefix(args[i], "-path=") {
            srcPath = strings.TrimPrefix(args[i], "-path=")
        } else if strings.HasPrefix(args[i], "-destino=") {
            dstPath = strings.TrimPrefix(args[i], "-destino=")
        }
    }
    if srcPath == "" || dstPath == "" {
        return "", fmt.Errorf("debe especificar los parámetros -path y -destino")
    }

    partitionID, err := stores.GetActivePartitionID()
    if err != nil {
        return "", fmt.Errorf("no hay partición activa montada")
    }
    sb, _, partitionPath, err := stores.GetMountedPartitionSuperblock(partitionID)
    if err != nil {
        return "", fmt.Errorf("no se pudo obtener el superbloque de la partición activa")
    }

    // Normalizar paths
    if !strings.HasPrefix(srcPath, "/") {
        srcPath = "/" + srcPath
    }
    if !strings.HasPrefix(dstPath, "/") {
        dstPath = "/" + dstPath
    }

    // Llama a la función recursiva para copiar (con overwrite=true)
    err = copyFolderRecursive(sb, partitionPath, srcPath, dstPath, true)
    if err != nil {
        return "", fmt.Errorf("error al copiar: %v", err)
    }
    // Imprimir inodos y bloques para validar el cambio
    fmt.Println("\n--- INODOS DESPUÉS DEL COPY ---")
    sb.PrintInodes(partitionPath)
    fmt.Println("\n--- BLOQUES DESPUÉS DEL COPY ---")
    sb.PrintBlocks(partitionPath)

    return fmt.Sprintf("COPY: Carpeta %s copiada a %s", srcPath, dstPath), nil
}

func copyFolderRecursive(sb *structures.SuperBlock, partitionPath, src, dst string, overwrite bool) error {
    fmt.Printf("\n[copyFolderRecursive] INICIO: src='%s', dst='%s', overwrite=%v\n", src, dst, overwrite)

    // Crear la carpeta destino si no existe
    parents, destDir := getParentDirs(dst)
    fmt.Printf("[copyFolderRecursive] Crear carpeta destino: parents=%v, destDir=%s\n", parents, destDir)
    err := sb.CreateFolder(partitionPath, parents, destDir)
    if err != nil {
        fmt.Printf("[copyFolderRecursive] Error creando carpeta destino: %v\n", err)
        return err
    }
    // Agregar la carpeta destino a la lista global de carpetas si no existe
    already := false
    for _, p := range global.ValidPaths {
        if p == dst {
            already = true
            break
        }
    }
    if !already {
        global.ValidPaths = append(global.ValidPaths, dst)
        fmt.Printf("[copyFolderRecursive] Carpeta agregada a la lista global: %s\n", dst)
    }

    // Obtener la lista global de archivos .txt válidos (ANTES)
    validTxtPaths := global.GetValidFilePathsMkfile()
    fmt.Println("[copyFolderRecursive] Lista global de archivos .txt ANTES del copy:")
    for _, p := range validTxtPaths {
        fmt.Println("  ", p)
    }

    // Obtener el inodo de la carpeta origen
    inodeIndex, err := sb.FindInodeByPath(partitionPath, src)
    fmt.Printf("[copyFolderRecursive] Buscar inodo origen para '%s' -> inodeIndex=%d, err=%v\n", src, inodeIndex, err)
    if err != nil || inodeIndex < 0 {
        return fmt.Errorf("carpeta origen no encontrada: %s", src)
    }
    inode := &structures.Inode{}
    err = inode.Deserialize(partitionPath, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
    if err != nil {
        fmt.Printf("[copyFolderRecursive] Error deserializando inodo origen: %v\n", err)
        return fmt.Errorf("error deserializando inodo origen: %v", err)
    }
    if inode.I_type[0] != '0' {
        fmt.Printf("[copyFolderRecursive] El path origen no es una carpeta: %s\n", src)
        return fmt.Errorf("el path origen no es una carpeta: %s", src)
    }

    // Recorrer los bloques de la carpeta origen
    for _, blockIndex := range inode.I_block {
        if blockIndex == -1 {
            continue
        }
        fmt.Printf("[copyFolderRecursive] Procesando blockIndex=%d de carpeta origen\n", blockIndex)
        block := &structures.FolderBlock{}
        err := block.Deserialize(partitionPath, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
        if err != nil {
            fmt.Printf("[copyFolderRecursive] Error deserializando FolderBlock: %v\n", err)
            continue
        }
        for _, content := range block.B_content {
            name := strings.Trim(string(content.B_name[:]), "\x00 ")
            fmt.Printf("[copyFolderRecursive] Analizando content: name='%s', inodo=%d\n", name, content.B_inodo)
            if name == "" || name == "." || name == ".." || content.B_inodo == -1 {
                continue
            }
            childInode := &structures.Inode{}
            err := childInode.Deserialize(partitionPath, int64(sb.S_inode_start+(content.B_inodo*sb.S_inode_size)))
            if err != nil {
                fmt.Printf("[copyFolderRecursive] Error deserializando childInode: %v\n", err)
                continue
            }
            srcChildPath := src + "/" + name
            dstChildPath := dst + "/" + name
            fmt.Printf("[copyFolderRecursive] srcChildPath='%s', dstChildPath='%s', tipo='%c'\n", srcChildPath, dstChildPath, childInode.I_type[0])
            if childInode.I_type[0] == '0' {
                // Es carpeta, copiar recursivamente
                fmt.Printf("[copyFolderRecursive] Es carpeta: copiando recursivamente '%s' -> '%s'\n", srcChildPath, dstChildPath)
                err := copyFolderRecursive(sb, partitionPath, srcChildPath, dstChildPath, overwrite)
                if err != nil {
                    fmt.Printf("[copyFolderRecursive] Error copiando subcarpeta: %v\n", err)
                }
            } else if childInode.I_type[0] == '1' {
                // Es archivo, pero solo copia si está en la lista global de mkfile y es .txt
                if strings.HasSuffix(name, ".txt") {
                    for _, validPath := range validTxtPaths {
                        if validPath == srcChildPath {
                            fmt.Printf("[copyFolderRecursive] Es .txt válido, copiando '%s' -> '%s'\n", srcChildPath, dstChildPath)
                            var contenido string
                            for _, bIdx := range childInode.I_block {
                                if bIdx == -1 {
                                    break
                                }
                                fileBlock := &structures.FileBlock{}
                                err := fileBlock.Deserialize(partitionPath, int64(sb.S_block_start+(bIdx*sb.S_block_size)))
                                if err != nil {
                                    fmt.Printf("[copyFolderRecursive] Error deserializando FileBlock: %v\n", err)
                                    continue
                                }
                                bloqueContenido := strings.TrimRight(string(fileBlock.B_content[:]), "\x00")
                                fmt.Printf("[copyFolderRecursive]   Bloque %d: contenido='%s'\n", bIdx, bloqueContenido)
                                contenido += bloqueContenido
                            }
                            // Crear el archivo en destino (sobrescribe si existe)
                            parents, destFile := getParentDirs(dstChildPath)
                            if overwrite {
                                fmt.Printf("[copyFolderRecursive] Overwrite activo: eliminando '%s' si existe\n", dstChildPath)
                                sb.RemoveFile(partitionPath, parents, destFile)
                            }
                            fmt.Printf("[copyFolderRecursive] Creando archivo en destino: parents=%v, destFile=%s, size=%d\n", parents, destFile, childInode.I_size)
                            err := sb.CreateFile(partitionPath, parents, destFile, int(childInode.I_size), []string{contenido})
                            if err != nil {
                                fmt.Printf("[copyFolderRecursive] Error creando archivo en destino: %v\n", err)
                            } else {
                                fmt.Printf("[copyFolderRecursive] Archivo creado exitosamente: %s\n", dstChildPath)
                            }
                            break // Solo copiar una vez por coincidencia
                        }
                    }
                }
            }
        }
    }

    // Al final, actualiza la lista global de archivos .txt agregando los nuevos paths en el destino
    nuevosPaths := []string{}
    for _, p := range validTxtPaths {
        if strings.HasPrefix(p, src) {
            // Solo los que están bajo el path origen
            nuevo := strings.Replace(p, src, dst, 1)
            nuevosPaths = append(nuevosPaths, nuevo)
        }
    }
    // Mantén los paths originales y agrega los nuevos (evita duplicados)
    finalPaths := global.GetValidFilePathsMkfile()
    for _, np := range nuevosPaths {
        existe := false
        for _, fp := range finalPaths {
            if fp == np {
                existe = true
                break
            }
        }
        if !existe {
            finalPaths = append(finalPaths, np)
        }
    }
    global.SetValidFilePathsMkfile(finalPaths)

    // Mostrar la lista global de archivos .txt DESPUÉS del copy
    fmt.Println("[copyFolderRecursive] Lista global de archivos .txt DESPUÉS del copy:")
    for _, p := range global.GetValidFilePathsMkfile() {
        fmt.Println("  ", p)
    }

    fmt.Printf("[copyFolderRecursive] FIN: src='%s', dst='%s'\n", src, dst)
    return nil
}

// getParentDirs obtiene las carpetas padres y el nombre del archivo
func getParentDirs(path string) ([]string, string) {
    parts := strings.Split(strings.Trim(path, "/"), "/")
    if len(parts) == 0 {
        return []string{}, ""
    }
    return parts[:len(parts)-1], parts[len(parts)-1]
}