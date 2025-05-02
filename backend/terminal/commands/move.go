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
    fmt.Println("========== INICIO MOVE ==========")
    // 1. Parsear argumentos
    var srcPath, dstPath string
    fmt.Println("[move] Argumentos recibidos:", args)
    for i := 0; i < len(args); i++ {
        if strings.HasPrefix(args[i], "-path=") {
            srcPath = strings.TrimPrefix(args[i], "-path=")
        } else if strings.HasPrefix(args[i], "-destino=") {
            dstPath = strings.TrimPrefix(args[i], "-destino=")
        }
    }
    fmt.Println("[move] srcPath:", srcPath)
    fmt.Println("[move] dstPath:", dstPath)
    if srcPath == "" || dstPath == "" {
        return "", fmt.Errorf("debe especificar los parámetros -path y -destino")
    }

    // 2. Obtener superbloque y path de partición activa
    partitionID, err := stores.GetActivePartitionID()
    fmt.Println("[move] partitionID:", partitionID, "err:", err)
    if err != nil {
        return "", fmt.Errorf("no hay partición activa montada")
    }
    sb, _, partitionPath, err := stores.GetMountedPartitionSuperblock(partitionID)
    fmt.Println("[move] partitionPath:", partitionPath, "err:", err)
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
    fmt.Println("[move] srcPath normalizado:", srcPath)
    fmt.Println("[move] dstPath normalizado:", dstPath)

    // 4. Validar que el archivo a mover esté en la lista global de mkfile y sea .txt
    validTxtPaths := global.GetValidFilePathsMkfile()
    fmt.Println("[move] Lista global de mkfile actual:", validTxtPaths)
    found := false
    for _, p := range validTxtPaths {
        if p == srcPath && strings.HasSuffix(p, ".txt") {
            found = true
            break
        }
    }
    fmt.Println("[move] ¿Archivo a mover está en la lista global y es .txt?:", found)
    if !found {
        return "", fmt.Errorf("el archivo a mover no es válido o no existe en la lista global")
    }

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
    fmt.Println("[move] dstFilePath final:", dstFilePath)

    // Mostrar la lista global ANTES de actualizar
    fmt.Println("[move] Lista global de mkfile ANTES de actualizar:")
    for _, p := range global.GetValidFilePathsMkfile() {
        fmt.Println("  ", p)
    }

    // Actualizar la lista global: reemplazar srcPath por dstFilePath
    newPaths := []string{}
    for _, p := range global.GetValidFilePathsMkfile() {
        if p == srcPath {
            fmt.Printf("[move] Reemplazando en lista global: %s -> %s\n", p, dstFilePath)
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

    // 6. Leer el contenido del archivo origen (buscando rutas intermedias)
    fmt.Println("[move] Buscando inodo del archivo origen...")
    var inodeIndex int32 = -1
    var foundPath string
    srcParts := strings.Split(strings.Trim(srcPath, "/"), "/")
    for i := 0; i < len(srcParts); i++ {
        tryPath := "/" + strings.Join(srcParts[i:], "/")
        fmt.Printf("[move] Intentando buscar inodo con path: %s\n", tryPath)
        idx, err := sb.FindInodeByPath(partitionPath, tryPath)
        fmt.Printf("[move] Resultado búsqueda: idx=%d, err=%v\n", idx, err)
        if err == nil && idx >= 0 {
            inodeIndex = idx
            foundPath = tryPath
            fmt.Printf("[move] Inodo encontrado en path: %s, idx: %d\n", foundPath, inodeIndex)
            break
        }
    }
    if inodeIndex < 0 {
        fmt.Printf("[move] Advertencia: el archivo origen %s no existe en el sistema de archivos. Se eliminará de la lista global.\n", srcPath)
        // Elimina el path de la lista global para mantenerla sincronizada
        newPaths := []string{}
        for _, p := range global.GetValidFilePathsMkfile() {
            if p != srcPath {
                newPaths = append(newPaths, p)
            }
        }
        global.SetValidFilePathsMkfile(newPaths)
        return fmt.Sprintf("MOVE: El archivo %s no existe físicamente, eliminado de la lista global.", srcPath), nil
    }

    fmt.Printf("[move] Deserializando inodo %d...\n", inodeIndex)
    inode := &structures.Inode{}
    err = inode.Deserialize(partitionPath, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
    if err != nil {
        fmt.Printf("[move] Error deserializando inodo: %v\n", err)
        return "", fmt.Errorf("error deserializando inodo origen: %v", err)
    }
    var contenido string
    fmt.Println("[move] Leyendo bloques del archivo origen...")
    for _, bIdx := range inode.I_block {
        if bIdx == -1 {
            break
        }
        fmt.Printf("[move] Leyendo bloque %d\n", bIdx)
        fileBlock := &structures.FileBlock{}
        err := fileBlock.Deserialize(partitionPath, int64(sb.S_block_start+(bIdx*sb.S_block_size)))
        if err != nil {
            fmt.Printf("[move] Error leyendo bloque %d: %v\n", bIdx, err)
            continue
        }
        bloqueContenido := strings.TrimRight(string(fileBlock.B_content[:]), "\x00")
        fmt.Printf("[move] Contenido bloque %d: '%s'\n", bIdx, bloqueContenido)
        contenido += bloqueContenido
    }

    // 7. Crear el archivo en el destino
    fmt.Println("[move] Preparando para crear archivo en destino...")
    parents, destFile := getParentDirs_w(dstFilePath)
    fmt.Printf("[move] parents: %v, destFile: %s\n", parents, destFile)

    // Validar existencia de las carpetas padres solo en la lista global
    for i := 1; i <= len(parents); i++ {
        subPath := "/" + strings.Join(parents[:i], "/")
        fmt.Printf("[move] Validando carpeta en lista global: %s\n", subPath)
        if !contains(global.ValidPaths, subPath) {
            return "", fmt.Errorf("no se encontró la carpeta destino: %s", subPath)
        }
    }

    // 9. Imprimir inodos y bloques para depuración
    fmt.Println("\n--- INODOS DESPUÉS DEL MOVE ---")
    sb.PrintInodes(partitionPath)
    fmt.Println("\n--- BLOQUES DESPUÉS DEL MOVE ---")
    sb.PrintBlocks(partitionPath)

    fmt.Println("========== FIN MOVE ==========")
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