package structures

import (
	"fmt"
	"strings"
	"terminal/utils"
	"time"
    "terminal/global"

)
var TxtFilesExtracted []TxtFile
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

// createFolderInInode crea una carpeta en un inodo específico
func (sb *SuperBlock) createFolderInInodeExt2(path string, inodeIndex int32, parentsDir []string, destDir string) error {
    fmt.Printf("-> createFolderInInodeExt2: inodeIndex=%d, parentsDir=%v, destDir=%s\n", inodeIndex, parentsDir, destDir)
    inode := &Inode{}
    err := inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
    if err != nil {
        return err
    }
    // Si no es carpeta, salir
    if inode.I_type[0] == '1' {
        return nil
    }
    // 1. Buscar si ya existe la carpeta en los bloques del inodo padre
    for _, blockIndex := range inode.I_block {
        if blockIndex == -1 {
            break
        }
        block := &FolderBlock{}
        err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
        if err != nil {
            return err
        }
        for k := 2; k < 4; k++ {
            name := strings.TrimRight(string(block.B_content[k].B_name[:]), "\x00")
            if name == destDir && block.B_content[k].B_inodo != -1 {
                // Ya existe la carpeta, no crearla de nuevo
                return nil
            }
        }
    }

    // 1. Buscar un bloque de carpeta del padre con espacio para el nuevo folder
    for i, blockIndex := range inode.I_block {
        if blockIndex == -1 {
            // No hay bloque, creamos uno nuevo para la carpeta hija
            newBlockIndex := sb.S_blocks_count
            inode.I_block[i] = newBlockIndex

            // Crear bloque de la nueva carpeta
            newBlock := &FolderBlock{
                B_content: [4]FolderContent{
                    {B_name: [12]byte{'.'}, B_inodo: sb.S_inodes_count},
                    {B_name: [12]byte{'.', '.'}, B_inodo: inodeIndex},
                    {B_name: [12]byte{'-'}, B_inodo: -1},
                    {B_name: [12]byte{'-'}, B_inodo: -1},
                },
            }
            // Serializar el bloque de la nueva carpeta
            err := newBlock.Serialize(path, int64(sb.S_first_blo))
            if err != nil {
                return err
            }
            // Actualizar bitmap y superbloque de bloques
            err = sb.UpdateBitmapBlock(path)
            if err != nil {
                return err
            }
            sb.S_blocks_count++
            sb.S_free_blocks_count--
            sb.S_first_blo += sb.S_block_size

            // Crear el inodo de la nueva carpeta
            newInode := &Inode{
                I_uid:   1,
                I_gid:   1,
                I_size:  0,
                I_atime: float32(time.Now().Unix()),
                I_ctime: float32(time.Now().Unix()),
                I_mtime: float32(time.Now().Unix()),
                I_block: [15]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
                I_type:  [1]byte{'0'},
                I_perm:  [3]byte{'6', '6', '4'},
            }
            newInode.I_block[0] = newBlockIndex

            // Serializar el nuevo inodo
            err = newInode.Serialize(path, int64(sb.S_first_ino))
            if err != nil {
                return err
            }
            // Actualizar bitmap y superbloque de inodos
            err = sb.UpdateBitmapInode(path)
            if err != nil {
                return err
            }
            sb.S_inodes_count++
            sb.S_free_inodes_count--
            sb.S_first_ino += sb.S_inode_size

            // ENLAZAR la nueva carpeta en el bloque de su padre
            // Buscar el primer bloque del padre con espacio
            for j := 0; j < len(inode.I_block); j++ {
                parentBlockIndex := inode.I_block[j]
                if parentBlockIndex == -1 {
                    break
                }
                parentBlock := &FolderBlock{}
                err := parentBlock.Deserialize(path, int64(sb.S_block_start+(parentBlockIndex*sb.S_block_size)))
                if err != nil {
                    return err
                }
                for k := 2; k < 4; k++ {
                    if parentBlock.B_content[k].B_inodo == -1 {
                        copy(parentBlock.B_content[k].B_name[:], destDir)
                        parentBlock.B_content[k].B_inodo = sb.S_inodes_count - 1 // El último inodo creado
                        // Serializar el bloque actualizado
                        err = parentBlock.Serialize(path, int64(sb.S_block_start+(parentBlockIndex*sb.S_block_size)))
                        if err != nil {
                            return err
                        }
                        // Serializar el inodo padre actualizado
                        err = inode.Serialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
                        if err != nil {
                            return err
                        }
                        return nil
                    }
                }
            }
            // Si no hay espacio en ningún bloque del padre, deberías crear un nuevo bloque y enlazarlo (no implementado aquí)
            return nil
        } else {
            // Si el bloque existe, buscar espacio para enlazar la nueva carpeta
            block := &FolderBlock{}
            err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
            if err != nil {
                return err
            }
            for k := 2; k < 4; k++ {
                if block.B_content[k].B_inodo == -1 {
                    // Ya existe espacio, así que solo crear el inodo y el bloque de la nueva carpeta
                    // Crear el bloque de la nueva carpeta
                    newBlockIndex := sb.S_blocks_count
                    newBlock := &FolderBlock{
                        B_content: [4]FolderContent{
                            {B_name: [12]byte{'.'}, B_inodo: sb.S_inodes_count},
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
                    sb.S_blocks_count++
                    sb.S_free_blocks_count--
                    sb.S_first_blo += sb.S_block_size

                    // Crear el inodo de la nueva carpeta
                    newInode := &Inode{
                        I_uid:   1,
                        I_gid:   1,
                        I_size:  0,
                        I_atime: float32(time.Now().Unix()),
                        I_ctime: float32(time.Now().Unix()),
                        I_mtime: float32(time.Now().Unix()),
                        I_block: [15]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
                        I_type:  [1]byte{'0'},
                        I_perm:  [3]byte{'6', '6', '4'},
                    }
                    newInode.I_block[0] = newBlockIndex

                    err = newInode.Serialize(path, int64(sb.S_first_ino))
                    if err != nil {
                        return err
                    }
                    err = sb.UpdateBitmapInode(path)
                    if err != nil {
                        return err
                    }
                    sb.S_inodes_count++
                    sb.S_free_inodes_count--
                    sb.S_first_ino += sb.S_inode_size

                    // Enlazar la nueva carpeta en el bloque actual
                    copy(block.B_content[k].B_name[:], destDir)
                    block.B_content[k].B_inodo = sb.S_inodes_count - 1
                    err = block.Serialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
                    if err != nil {
                        return err
                    }
                    err = inode.Serialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
                    if err != nil {
                        return err
                    }
                    return nil
                }
            }
        }
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
func (sb *SuperBlock) ExtractTxtFiles(path string, partitionID string) error {
    TxtFilesExtracted = []TxtFile{}
    validPaths := global.GetValidFilePathsMkfile()
    for inodeIndex := int32(0); inodeIndex < sb.S_inodes_count; inodeIndex++ {
        inode := &Inode{}
        err := inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
        if err != nil {
            continue
        }
        if inode.I_type[0] == '1' {
            foundName := ""
            // Buscar el nombre en todos los bloques de carpeta
            for blockIdx := int32(0); blockIdx < sb.S_blocks_count; blockIdx++ {
                folderBlock := &FolderBlock{}
                err := folderBlock.Deserialize(path, int64(sb.S_block_start+(blockIdx*sb.S_block_size)))
                if err != nil {
                    continue
                }
                for _, content := range folderBlock.B_content {
                    name := strings.TrimSpace(strings.Trim(string(content.B_name[:]), "\x00 "))
                    if content.B_inodo == inodeIndex && name != "" && name != "." && name != ".." {
                        foundName = name
                        break
                    }
                }
                if foundName != "" {
                    break
                }
            }
            // Si encontró nombre y es .txt, buscar TODAS las coincidencias en la lista global
            if foundName != "" && strings.HasSuffix(foundName, ".txt") {
                var contenido string
                for _, blockIndex := range inode.I_block {
                    if blockIndex == -1 {
                        break
                    }
                    block := &FileBlock{}
                    err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
                    if err != nil {
                        continue
                    }
                    bloqueContenido := strings.TrimRight(string(block.B_content[:]), "\x00")
                    contenido += bloqueContenido
                }
                // Buscar todas las coincidencias de path en la lista global
                for _, validPath := range validPaths {
                    // Extraer el nombre del archivo del path
                    parts := strings.Split(validPath, "/")
                    if len(parts) == 0 {
                        continue
                    }
                    validName := parts[len(parts)-1]
                    if validName == foundName {
                        TxtFilesExtracted = append(TxtFilesExtracted, TxtFile{
                            Path:      validPath,
                            ID:        partitionID,
                            Contenido: contenido,
                            Size:      inode.I_size,
                        })
                    }
                }
            }
        }
    }
    return nil
}

func (sb *SuperBlock) GetTxtFiles(path string, partitionID string) ([]TxtFile, error) {
    // Solo retorna la lista global ya extraída
    return TxtFilesExtracted, nil
}

// Busca el índice de inodo dado un path absoluto
func (sb *SuperBlock) FindInodeByPath(diskPath string, absPath string) (int32, error) {
    if absPath == "" || absPath == "/" {
        return 0, nil // raíz
    }
    parts := strings.Split(strings.Trim(absPath, "/"), "/")
    inodeIndex := int32(0) // raíz
    for _, part := range parts {
        inode := &Inode{}
        err := inode.Deserialize(diskPath, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
        if err != nil {
            return -1, err
        }
        found := false
        for _, blockIndex := range inode.I_block {
            if blockIndex == -1 {
                continue
            }
            folderBlock := &FolderBlock{}
            err := folderBlock.Deserialize(diskPath, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
            if err != nil {
                continue
            }
            for _, content := range folderBlock.B_content {
                name := strings.Trim(string(content.B_name[:]), "\x00 ")
                // Cambia esta línea:
                // if name == part && content.B_inodo != -1 {
                // Por esta:
                if len(part) > 12 {
                    part = part[:12]
                }
                if len(name) > 12 {
                    name = name[:12]
                }
                if name == part && content.B_inodo != -1 {
                    inodeIndex = content.B_inodo
                    found = true
                    break
                }
            }
            if found {
                break
            }
        }
        if !found {
            return -1, fmt.Errorf("no se encontró el inodo para %s", absPath)
        }
    }
    return inodeIndex, nil
}

// Eliminar un archivo .txt de la lista
func RemoveTxtFileFromExtracted(path string) {
    newList := []TxtFile{}
    for _, file := range TxtFilesExtracted {
        if file.Path != path {
            newList = append(newList, file)
        }
    }
    TxtFilesExtracted = newList
}

// Renombrar un archivo .txt en la lista
// Renombrar un archivo .txt en la lista usando solo el nombre del archivo
func RenameTxtFileInExtractedByName(oldName, newName string) {
    for i, file := range TxtFilesExtracted {
        // Extrae el nombre del archivo del path
        parts := strings.Split(file.Path, "/")
        if len(parts) == 0 {
            continue
        }
        currentName := parts[len(parts)-1]
        if currentName == oldName {
            // Cambia solo el nombre, conserva el path
            parts[len(parts)-1] = newName
            TxtFilesExtracted[i].Path = strings.Join(parts, "/")
            break // Si solo quieres cambiar el primero que coincida
        }
    }
}