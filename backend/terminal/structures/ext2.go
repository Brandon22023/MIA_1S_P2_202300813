package structures

import (
	"terminal/utils"
	"strings"
	"time"
	"fmt"
)
type TxtFile struct {
    Path      string `json:"path"`
    ID        string `json:"id"`
    Contenido string `json:"contenido"`
	Size      int32  `json:"size"`
}
// Crear users.txt en nuestro sistema de archivos
func (sb *SuperBlock) CreateUsersFileExt2(path string) error {
	// ----------- Creamos / -----------
	// Creamos el inodo raíz
	rootInode := &Inode{
		I_uid:   1,
		I_gid:   1,
		I_size:  0,
		I_atime: float32(time.Now().Unix()),
		I_ctime: float32(time.Now().Unix()),
		I_mtime: float32(time.Now().Unix()),
		I_block: [15]int32{sb.S_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  [1]byte{'0'},
		I_perm:  [3]byte{'7', '7', '7'},
	}

	// Serializar el inodo raíz
	err := rootInode.Serialize(path, int64(sb.S_first_ino))
	if err != nil {
		return err
	}

	// Actualizar el bitmap de inodos
	err = sb.UpdateBitmapInode(path)
	if err != nil {
		return err
	}

	// Actualizar el superbloque
	sb.S_inodes_count++
	sb.S_free_inodes_count--
	sb.S_first_ino += sb.S_inode_size

	// Creamos el bloque del Inodo Raíz
	rootBlock := &FolderBlock{
		B_content: [4]FolderContent{
			{B_name: [12]byte{'.'}, B_inodo: 0},
			{B_name: [12]byte{'.', '.'}, B_inodo: 0},
			{B_name: [12]byte{'-'}, B_inodo: -1},
			{B_name: [12]byte{'-'}, B_inodo: -1},
		},
	}

	// Actualizar el bitmap de bloques
	err = sb.UpdateBitmapBlock(path)
	if err != nil {
		return err
	}

	// Serializar el bloque de carpeta raíz
	err = rootBlock.Serialize(path, int64(sb.S_first_blo))
	if err != nil {
		return err
	}

	// Actualizar el superbloque
	sb.S_blocks_count++
	sb.S_free_blocks_count--
	sb.S_first_blo += sb.S_block_size

	// ----------- Creamos /users.txt -----------
	usersText := "1,G,root\n1,U,root,root,123\n"

	// Deserializar el inodo raíz
	err = rootInode.Deserialize(path, int64(sb.S_inode_start+0)) // 0 porque es el inodo raíz
	if err != nil {
		return err
	}

	// Actualizamos el inodo raíz
	rootInode.I_atime = float32(time.Now().Unix())

	// Serializar el inodo raíz
	err = rootInode.Serialize(path, int64(sb.S_inode_start+0)) // 0 porque es el inodo raíz
	if err != nil {
		return err
	}

	// Deserializar el bloque de carpeta raíz
	err = rootBlock.Deserialize(path, int64(sb.S_block_start+0)) // 0 porque es el bloque de carpeta raíz
	if err != nil {
		return err
	}

	// Actualizamos el bloque de carpeta raíz
	rootBlock.B_content[2] = FolderContent{B_name: [12]byte{'u', 's', 'e', 'r', 's', '.', 't', 'x', 't'}, B_inodo: sb.S_inodes_count}

	// Serializar el bloque de carpeta raíz
	err = rootBlock.Serialize(path, int64(sb.S_block_start+0)) // 0 porque es el bloque de carpeta raíz
	if err != nil {
		return err
	}

	// Creamos el inodo users.txt
	usersInode := &Inode{
		I_uid:   1,
		I_gid:   1,
		I_size:  int32(len(usersText)),
		I_atime: float32(time.Now().Unix()),
		I_ctime: float32(time.Now().Unix()),
		I_mtime: float32(time.Now().Unix()),
		I_block: [15]int32{sb.S_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  [1]byte{'1'},
		I_perm:  [3]byte{'7', '7', '7'},
	}

	// Actualizar el bitmap de inodos
	err = sb.UpdateBitmapInode(path)
	if err != nil {
		return err
	}

	// Serializar el inodo users.txt
	err = usersInode.Serialize(path, int64(sb.S_first_ino))
	if err != nil {
		return err
	}

	// Actualizamos el superbloque
	sb.S_inodes_count++
	sb.S_free_inodes_count--
	sb.S_first_ino += sb.S_inode_size

	// Creamos el bloque de users.txt
	usersBlock := &FileBlock{
		B_content: [64]byte{},
	}
	// Copiamos el texto de usuarios en el bloque
	copy(usersBlock.B_content[:], usersText)

	// Serializar el bloque de users.txt
	err = usersBlock.Serialize(path, int64(sb.S_first_blo))
	if err != nil {
		return err
	}

	// Actualizar el bitmap de bloques
	err = sb.UpdateBitmapBlock(path)
	if err != nil {
		return err
	}

	// Actualizamos el superbloque
	sb.S_blocks_count++
	sb.S_free_blocks_count--
	sb.S_first_blo += sb.S_block_size

	return nil
}

func (sb *SuperBlock) createFolderInInodeExt2(path string, inodeIndex int32, parentsDir []string, destDir string) error {
    inode := &Inode{}
    err := inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
    if err != nil {
        return err
    }
    if inode.I_type[0] == '1' {
        return nil
    }

    for bIdx, blockIndex := range inode.I_block {
        if blockIndex == -1 {
            fmt.Printf("[DEBUG] Creando nuevo bloque de carpeta para inodo %d en posición %d\n", inodeIndex, bIdx)
            newBlock := &FolderBlock{
                B_content: [4]FolderContent{
                    {B_name: [12]byte{'.'}, B_inodo: inodeIndex},
                    {B_name: [12]byte{'.', '.'}, B_inodo: inodeIndex},
                    {B_name: [12]byte{'-'}, B_inodo: -1},
                    {B_name: [12]byte{'-'}, B_inodo: -1},
                },
            }
            err := newBlock.Serialize(path, int64(sb.S_first_blo))
            if err != nil {
                return err
            }
            err = sb.UpdateBitmapBlock(path)
            if err != nil {
                return err
            }
            inode.I_block[bIdx] = sb.S_blocks_count
            err = inode.Serialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
            if err != nil {
                return err
            }
            sb.S_blocks_count++
            sb.S_free_blocks_count--
            sb.S_first_blo += sb.S_block_size
            // break aquí para que la próxima iteración procese el nuevo bloque
            break
        }

        block := &FolderBlock{}
        err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
        if err != nil {
            return err
        }

        for indexContent := 2; indexContent < len(block.B_content); indexContent++ {
            content := block.B_content[indexContent]

            if len(parentsDir) != 0 {
                parentDir, err := utils.First(parentsDir)
                if err != nil {
                    return err
                }
                contentName := strings.Trim(string(content.B_name[:]), "\x00 ")
                parentDirName := strings.Trim(parentDir, "\x00 ")
                if strings.EqualFold(contentName, parentDirName) {
                    return sb.createFolderInInodeExt2(path, content.B_inodo, utils.RemoveElement(parentsDir, 0), destDir)
                }
            }

            // Si ya no hay padres pendientes y hay espacio libre, crea la carpeta
            if len(parentsDir) == 0 && content.B_inodo == -1 {
                fmt.Printf("[DEBUG] Creando carpeta '%s' en inodo %d (nuevo inodo: %d, nuevo bloque: %d)\n", destDir, inodeIndex, sb.S_inodes_count, sb.S_blocks_count)
                fmt.Printf("[DEBUG] Estado antes de crear: S_first_ino=%d, S_first_blo=%d, S_inodes_count=%d, S_blocks_count=%d\n", sb.S_first_ino, sb.S_first_blo, sb.S_inodes_count, sb.S_blocks_count)
                newInode := sb.S_inodes_count
                copy(content.B_name[:], destDir)
                content.B_inodo = newInode
                block.B_content[indexContent] = content
                err = block.Serialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
                if err != nil {
                    fmt.Printf("[ERROR] Al serializar folderInode: %v\n", err)
                    return err
}

                folderInode := &Inode{
                    I_uid:   1,
                    I_gid:   1,
                    I_size:  0,
                    I_atime: float32(time.Now().Unix()),
                    I_ctime: float32(time.Now().Unix()),
                    I_mtime: float32(time.Now().Unix()),
                    I_block: [15]int32{sb.S_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
                    I_type:  [1]byte{'0'},
                    I_perm:  [3]byte{'6', '6', '4'},
                }
                err = folderInode.Serialize(path, int64(sb.S_first_ino))
                if err != nil {
                    fmt.Printf("[ERROR] Al serializar folderBlock: %v\n", err)
                    return err
                }
                err = sb.UpdateBitmapInode(path)
                if err != nil {
                    return err
                }
                sb.S_inodes_count++
                sb.S_free_inodes_count--
                sb.S_first_ino += sb.S_inode_size

                folderBlock := &FolderBlock{
                    B_content: [4]FolderContent{
                        {B_name: [12]byte{'.'}, B_inodo: newInode},
                        {B_name: [12]byte{'.', '.'}, B_inodo: inodeIndex},
                        {B_name: [12]byte{'-'}, B_inodo: -1},
                        {B_name: [12]byte{'-'}, B_inodo: -1},
                    },
                }
                err = folderBlock.Serialize(path, int64(sb.S_first_blo))
                if err != nil {
                    return err
                }
                err = sb.UpdateBitmapBlock(path)
                if err != nil {
                    return err
                }
                sb.S_blocks_count++
                sb.S_free_blocks_count--
                sb.S_first_blo += sb.S_block_size
                fmt.Printf("[DEBUG] Carpeta '%s' creada correctamente en inodo %d\n", destDir, newInode)
                if len(parentsDir) > 0 {
                    fmt.Printf("[DEBUG] Continúa con el siguiente padre: %v\n", utils.RemoveElement(parentsDir, 0))
                    return sb.createFolderInInodeExt2(path, content.B_inodo, utils.RemoveElement(parentsDir, 0), destDir)
                }
                return nil
            }
        }
        // Si llegaste aquí y no encontraste espacio, sigue con el siguiente bloque
    }

    // Si no se encontró el padre y aún quedan padres por crear, créalo y vuelve a intentar
    if len(parentsDir) != 0 {
        parentDir := parentsDir[0]
        fmt.Printf("[DEPURACION] Intentando crear padre '%s' en inodo %d. parentsDir: %v\n", parentDir, inodeIndex, parentsDir)
        err := sb.createFolderInInodeExt2(path, inodeIndex, []string{}, parentDir)
        if err != nil {
            fmt.Printf("[DEPURACION] Error creando padre '%s': %v\n", parentDir, err)
            return err
        }
        fmt.Printf("[DEPURACION] Padre '%s' creado. Buscando inodo para continuar la recursión...\n", parentDir)
        inode := &Inode{}
        fmt.Printf("[DEBUG] Deserializando inodo %d en offset %d\n", inodeIndex, sb.S_inode_start+(inodeIndex*sb.S_inode_size))
        err = inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
        if err != nil {
            fmt.Printf("[DEPURACION] Error deserializando inodo %d: %v\n", inodeIndex, err)
            return err
        }
        for _, blockIndex := range inode.I_block {
            if blockIndex == -1 {
                break
            }
            block := &FolderBlock{}
            fmt.Printf("[DEBUG] Deserializando inodo %d en offset %d\n", inodeIndex, sb.S_inode_start+(inodeIndex*sb.S_inode_size))

            err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
            if err != nil {
                fmt.Printf("[DEPURACION] Error deserializando bloque %d: %v\n", blockIndex, err)
                return err
            }
            for _, content := range block.B_content {
                contentName := strings.Trim(string(content.B_name[:]), "\x00 ")
                fmt.Printf("[DEPURACION] Revisando contenido: '%s' (inodo: %d) buscando '%s'\n", contentName, content.B_inodo, parentDir)
                if strings.EqualFold(contentName, parentDir) {
                    fmt.Printf("[DEPURACION] Encontrado padre '%s' en inodo %d. Continuando recursión...\n", parentDir, content.B_inodo)
                    return sb.createFolderInInodeExt2(path, content.B_inodo, utils.RemoveElement(parentsDir, 0), destDir)
                }
            }
        }
        fmt.Printf("[DEPURACION] No se pudo encontrar el inodo del padre recién creado: %s\n", parentDir)
        return fmt.Errorf("no se pudo encontrar el inodo del padre recién creado: %s", parentDir)
    }
    return nil
}

// createFolderinode crea una carpeta en un inodo específico
func (sb *SuperBlock) createFileInInode(path string, inodeIndex int32, parentsDir []string, destFile string, fileSize int, fileContent []string) error {
	// Crear un nuevo inodo
	inode := &Inode{}
	// Deserializar el inodo
	err := inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
	if err != nil {
		return err
	}
	// Verificar si el inodo es de tipo carpeta
	if inode.I_type[0] == '1' {
		fmt.Printf("[DEBUG] Inodo %d no es carpeta, es archivo. No se crea nada.\n", inodeIndex)
		return nil
	}

	// Iterar sobre cada bloque del inodo (apuntadores)
	for _, blockIndex := range inode.I_block {
		// Si el bloque no existe, salir
		if blockIndex == -1 {
			break
		}

		// Crear un nuevo bloque de carpeta
		block := &FolderBlock{}

		// Deserializar el bloque
		err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
		if err != nil {
			return err
		}

		// Iterar sobre cada contenido del bloque, desde el index 2 porque los primeros dos son . y ..
		for indexContent := 2; indexContent < len(block.B_content); indexContent++ {
			// Obtener el contenido del bloque
			content := block.B_content[indexContent]

			// Sí las carpetas padre no están vacías debereamos buscar la carpeta padre más cercana
			if len(parentsDir) != 0 {
				//fmt.Println("---------ESTOY  VISITANDO--------")

				// Si el contenido está vacío, salir
				if content.B_inodo == -1 {
					break
				}

				// Obtenemos la carpeta padre más cercana
				parentDir, err := utils.First(parentsDir)
				if err != nil {
					return err
				}

				// Convertir B_name a string y eliminar los caracteres nulos
				contentName := strings.Trim(string(content.B_name[:]), "\x00 ")
				// Convertir parentDir a string y eliminar los caracteres nulos
				parentDirName := strings.Trim(parentDir, "\x00 ")
				// Si el nombre del contenido coincide con el nombre de la carpeta padre
				if strings.EqualFold(contentName, parentDirName) {
					//fmt.Println("---------ESTOY  ENCONTRANDO--------")
					// Si son las mismas, entonces entramos al inodo que apunta el bloque
					err := sb.createFileInInode(path, content.B_inodo, utils.RemoveElement(parentsDir, 0), destFile, fileSize, fileContent)
					if err != nil {
						return err
					}
					return nil
				}
			} else {
				//fmt.Println("---------ESTOY  CREANDO--------")

				// Si el apuntador al inodo está ocupado, continuar con el siguiente
				if content.B_inodo != -1 {
					continue
				}

				// Actualizar el contenido del bloque
				copy(content.B_name[:], []byte(destFile))
				content.B_inodo = sb.S_inodes_count

				// Actualizar el bloque
				block.B_content[indexContent] = content

				// Serializar el bloque
				err = block.Serialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
				if err != nil {
					return err
				}

				// Crear el inodo del archivo
				fileInode := &Inode{
					I_uid:   1,
					I_gid:   1,
					I_size:  int32(fileSize),
					I_atime: float32(time.Now().Unix()),
					I_ctime: float32(time.Now().Unix()),
					I_mtime: float32(time.Now().Unix()),
					I_block: [15]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
					I_type:  [1]byte{'1'},
					I_perm:  [3]byte{'6', '6', '4'},
				}

				// Crear el bloques del archivo
				for i := 0; i < len(fileContent); i++ {
					// Actualizamos el inodo del archivo
					fileInode.I_block[i] = sb.S_blocks_count

					// Creamos el bloque del archivo
					fileBlock := &FileBlock{
						B_content: [64]byte{},
					}
					// Copiamos el texto de usuarios en el bloque
					copy(fileBlock.B_content[:], fileContent[i])

					// Serializar el bloque de users.txt
					err = fileBlock.Serialize(path, int64(sb.S_first_blo))
					if err != nil {
						return err
					}

					// Actualizar el bitmap de bloques
					err = sb.UpdateBitmapBlock(path)
					if err != nil {
						return err
					}

					// Actualizamos el superbloque
					sb.S_blocks_count++
					sb.S_free_blocks_count--
					sb.S_first_blo += sb.S_block_size
				}

				// Serializar el inodo de la carpeta
				err = fileInode.Serialize(path, int64(sb.S_first_ino))
				if err != nil {
					return err
				}

				// Actualizar el bitmap de inodos
				err = sb.UpdateBitmapInode(path)
				if err != nil {
					return err
				}

				// Actualizar el superbloque
				sb.S_inodes_count++
				sb.S_free_inodes_count--
				sb.S_first_ino += sb.S_inode_size

				return nil
			}
		}

	}
	return nil
}

// Extrae archivos .txt y su contenido del sistema de archivos ext2, mostrando la ruta completa
func (sb *SuperBlock) ExtractTxtFiles(path string) error {
    for i := int32(0); i < sb.S_inodes_count; i++ {
        inode := &Inode{}
        err := inode.Deserialize(path, int64(sb.S_inode_start+(i*sb.S_inode_size)))
        if err != nil {
            return err
        }
        // Solo archivos
        if inode.I_type[0] == '1' {
            for _, blockIndex := range inode.I_block {
                if blockIndex == -1 {
                    break
                }
                block := &FileBlock{}
                err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
                if err != nil {
                    return err
                }
                // Buscar el nombre y la ruta en los bloques de carpeta
                fileName := ""
                filePath := ""
                for j := int32(0); j < sb.S_inodes_count; j++ {
                    parentInode := &Inode{}
                    _ = parentInode.Deserialize(path, int64(sb.S_inode_start+(j*sb.S_inode_size)))
                    if parentInode.I_type[0] == '0' {
                        for _, parentBlockIndex := range parentInode.I_block {
                            if parentBlockIndex == -1 {
                                break
                            }
                            parentBlock := &FolderBlock{}
                            _ = parentBlock.Deserialize(path, int64(sb.S_block_start+(parentBlockIndex*sb.S_block_size)))
                            for _, content := range parentBlock.B_content {
                                if content.B_inodo == i {
                                    name := strings.TrimRight(string(content.B_name[:]), "\x00")
                                    if strings.HasSuffix(name, ".txt") {
                                        fileName = name
                                        // Reconstruir la ruta desde la raíz
                                        filePath = reconstructPath(sb, path, j, name)
                                    }
                                }
                            }
                        }
                    }
                }
                if fileName != "" {
					content := strings.TrimRight(string(block.B_content[:]), "\x00")
					fmt.Printf(
						"Archivo encontrado: %s\nRuta: %s\nTamaño: %d bytes\nContenido:\n%s\n\n",
						fileName, filePath, inode.I_size, content,
					)
				}
            }
        }
    }
    return nil
}

// Función auxiliar para reconstruir la ruta completa de un archivo dado el inodo padre y el nombre
func reconstructPath(sb *SuperBlock, path string, parentInodeIndex int32, fileName string) string {
    var dirs []string
    currentInodeIndex := parentInodeIndex

    // Recorrer hacia atrás hasta el inodo raíz (0)
    for currentInodeIndex != 0 {
        found := false
        for i := int32(0); i < sb.S_inodes_count; i++ {
            inode := &Inode{}
            _ = inode.Deserialize(path, int64(sb.S_inode_start+(i*sb.S_inode_size)))
            if inode.I_type[0] == '0' {
                for _, blockIndex := range inode.I_block {
                    if blockIndex == -1 {
                        break
                    }
                    block := &FolderBlock{}
                    _ = block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
                    for _, content := range block.B_content {
                        if content.B_inodo == currentInodeIndex && string(content.B_name[:1]) != "." {
                            dirName := strings.TrimRight(string(content.B_name[:]), "\x00")
                            dirs = append([]string{dirName}, dirs...)
                            currentInodeIndex = i
                            found = true
                            break
                        }
                    }
                    if found {
                        break
                    }
                }
            }
            if found {
                break
            }
        }
        if !found {
            break
        }
    }
    // Agregar la raíz
    dirs = append([]string{""}, dirs...)
    // Agregar el nombre del archivo
    dirs = append(dirs, fileName)
    return strings.Join(dirs, "/")
}
// Devuelve archivos .txt y su contenido del sistema de archivos ext2, mostrando la ruta completa
func (sb *SuperBlock) GetTxtFiles(path string, partitionID string) ([]TxtFile, error) {
    var files []TxtFile
    for i := int32(0); i < sb.S_inodes_count; i++ {
        inode := &Inode{}
        err := inode.Deserialize(path, int64(sb.S_inode_start+(i*sb.S_inode_size)))
        if err != nil {
            continue
        }
        // Solo archivos
        if inode.I_type[0] == '1' {
            for _, blockIndex := range inode.I_block {
                if blockIndex == -1 {
                    break
                }
                block := &FileBlock{}
                err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
                if err != nil {
                    continue
                }
                fileName := ""
                filePath := ""
                for j := int32(0); j < sb.S_inodes_count; j++ {
                    parentInode := &Inode{}
                    _ = parentInode.Deserialize(path, int64(sb.S_inode_start+(j*sb.S_inode_size)))
                    if parentInode.I_type[0] == '0' {
                        for _, parentBlockIndex := range parentInode.I_block {
                            if parentBlockIndex == -1 {
                                break
                            }
                            parentBlock := &FolderBlock{}
                            _ = parentBlock.Deserialize(path, int64(sb.S_block_start+(parentBlockIndex*sb.S_block_size)))
                            for _, content := range parentBlock.B_content {
                                if content.B_inodo == i {
                                    name := strings.TrimRight(string(content.B_name[:]), "\x00")
                                    if strings.HasSuffix(name, ".txt") {
                                        fileName = name
                                        filePath = reconstructPath(sb, path, j, name)
                                    }
                                }
                            }
                        }
                    }
                }
                if fileName != "" {
                    content := strings.TrimRight(string(block.B_content[:]), "\x00")
                    files = append(files, TxtFile{
                        Path:      filePath,
                        ID:        partitionID,
                        Contenido: content,
						Size:      inode.I_size,
                    })
                }
            }
        }
    }
    return files, nil
}