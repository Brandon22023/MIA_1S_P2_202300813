package commands

import (
    "fmt"
    "strings"
    global "terminal/global"
    "terminal/structures"
    stores "terminal/stores"
)

// ParseMove procesa el comando move
func ParseMove(args []string) (string, error) {
    // 1. Parsear argumentos
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
    fmt.Printf("[move] srcPath: %s\n", srcPath)
    fmt.Printf("[move] dstPath: %s\n", dstPath)

    // 2. Obtener superbloque y path de partición activa
    partitionID, err := stores.GetActivePartitionID()
    if err != nil {
        return "", fmt.Errorf("no hay partición activa montada")
    }
    sb, _, partitionPath, err := stores.GetMountedPartitionSuperblock(partitionID)
    if err != nil {
        return "", fmt.Errorf("no se pudo obtener el superbloque de la partición activa")
    }

    // 3. Normalizar paths
    if !strings.HasPrefix(srcPath, "/") {
        srcPath = "/" + srcPath
    }
    if !strings.HasPrefix(dstPath, "/") {
        dstPath = "/" + dstPath
    }
    fmt.Printf("[move] srcPath normalizado: %s\n", srcPath)
    fmt.Printf("[move] dstPath normalizado: %s\n", dstPath)

    // 4. Validar que el archivo a mover esté en la lista global de mkfile y sea .txt
    validTxtPaths := global.GetValidFilePathsMkfile()
    found := false
    for _, p := range validTxtPaths {
        if p == srcPath && strings.HasSuffix(p, ".txt") {
            found = true
            break
        }
    }
    if !found {
        return "", fmt.Errorf("el archivo a mover no es válido o no existe en la lista global")
    }
    fmt.Println("[move] El archivo existe en la lista global y es .txt")

    // 5. Obtener el nombre del archivo y armar el path destino final
    parts := strings.Split(srcPath, "/")
    fileName := parts[len(parts)-1]
    dstFilePath := dstPath
    if !strings.HasSuffix(dstPath, ".txt") {
        if !strings.HasSuffix(dstPath, "/") {
            dstFilePath = dstPath + "/" + fileName
        } else {
            dstFilePath = dstPath + fileName
        }
    }
    fmt.Printf("[move] dstFilePath final: %s\n", dstFilePath)

    // Mostrar la lista global ANTES de actualizar
    fmt.Println("[move] Lista global de mkfile ANTES de actualizar:")
    for _, p := range global.GetValidFilePathsMkfile() {
        fmt.Println("  ", p)
    }

    // Actualizar la lista global: reemplazar srcPath por dstFilePath
    newPaths := []string{}
    for _, p := range global.GetValidFilePathsMkfile() {
        if p == srcPath {
            fmt.Printf("[move] Actualizando path en lista global: %s -> %s\n", p, dstFilePath)
            newPaths = append(newPaths, dstFilePath)
        } else {
            newPaths = append(newPaths, p)
        }
    }
    global.SetValidFilePathsMkfile(newPaths)

    // Mostrar la lista global DESPUÉS de actualizar
    fmt.Println("[move] Lista global de mkfile DESPUÉS de actualizar:")
    for _, p := range global.GetValidFilePathsMkfile() {
        fmt.Println("  ", p)
    }

    // 6. Leer el contenido del archivo origen
    inodeIndex, err := sb.FindInodeByPath(partitionPath, srcPath)
    if err != nil || inodeIndex < 0 {
        return "", fmt.Errorf("archivo origen no encontrado: %s", srcPath)
    }
    inode := &structures.Inode{}
    err = inode.Deserialize(partitionPath, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
    if err != nil {
        return "", fmt.Errorf("error deserializando inodo origen: %v", err)
    }
    var contenido string
    for _, bIdx := range inode.I_block {
        if bIdx == -1 {
            break
        }
        fileBlock := &structures.FileBlock{}
        err := fileBlock.Deserialize(partitionPath, int64(sb.S_block_start+(bIdx*sb.S_block_size)))
        if err != nil {
            continue
        }
        bloqueContenido := strings.TrimRight(string(fileBlock.B_content[:]), "\x00")
        contenido += bloqueContenido
    }
    fmt.Printf("[move] Contenido leído: '%s'\n", contenido)

    // 7. Crear el archivo en el destino
    parents, destFile := getParentDirs_w(dstFilePath)
    fmt.Printf("[move] Crear carpeta destino si no existe: %v, destFile: %s\n", parents, destFile)
    err = sb.CreateFolder(partitionPath, parents, "")
    if err != nil {
        return "", fmt.Errorf("no se pudo crear la carpeta destino: %v", err)
    }
    err = sb.CreateFile(partitionPath, parents, destFile, int(inode.I_size), []string{contenido})
    if err != nil {
        return "", fmt.Errorf("no se pudo crear el archivo en el destino: %v", err)
    }
    fmt.Println("[move] Archivo creado en el destino")

    // 8. Eliminar el archivo original
    srcParents, srcFile := getParentDirs_w(srcPath)
    err = sb.RemoveFile(partitionPath, srcParents, srcFile)
    if err != nil {
        return "", fmt.Errorf("no se pudo eliminar el archivo original: %v", err)
    }
    fmt.Println("[move] Archivo original eliminado")

    // 9. Imprimir inodos y bloques para depuración
    fmt.Println("\n--- INODOS DESPUÉS DEL MOVE ---")
    sb.PrintInodes(partitionPath)
    fmt.Println("\n--- BLOQUES DESPUÉS DEL MOVE ---")
    sb.PrintBlocks(partitionPath)

    return fmt.Sprintf("MOVE: Archivo %s movido a %s", srcPath, dstFilePath), nil
}

// getParentDirs obtiene las carpetas padres y el nombre del archivo
func getParentDirs_w(path string) ([]string, string) {
    parts := strings.Split(strings.Trim(path, "/"), "/")
    if len(parts) == 0 {
        return []string{}, ""
    }
    return parts[:len(parts)-1], parts[len(parts)-1]
}