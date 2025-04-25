package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os/exec"
	"time"
	"strings"
	//strconv"
	"path/filepath"
	"html" 
	"os"
	utils "terminal/utils"
)

type SuperBlock struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_inodes_count int32
	S_free_blocks_count int32
	S_mtime             float32
	S_umtime            float32
	S_mnt_count         int32
	S_magic             int32
	S_inode_size        int32
	S_block_size        int32
	S_first_ino         int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
	// Total: 68 bytes
}

// Serialize escribe la estructura SuperBlock en un archivo binario en la posición especificada
func (sb *SuperBlock) Serialize(path string, offset int64) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo a la posición especificada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Serializar la estructura SuperBlock directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, sb)
	if err != nil {
		return err
	}

	return nil
}

// Deserialize lee la estructura SuperBlock desde un archivo binario en la posición especificada
func (sb *SuperBlock) Deserialize(path string, offset int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo a la posición especificada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Obtener el tamaño de la estructura SuperBlock
	sbSize := binary.Size(sb)
	if sbSize <= 0 {
		return fmt.Errorf("invalid SuperBlock size: %d", sbSize)
	}

	// Leer solo la cantidad de bytes que corresponden al tamaño de la estructura SuperBlock
	buffer := make([]byte, sbSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserializar los bytes leídos en la estructura SuperBlock
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, sb)
	if err != nil {
		return err
	}

	return nil
}



// PrintSuperBlock imprime los valores de la estructura SuperBlock
func (sb *SuperBlock) Print() {
	// Convertir el tiempo de montaje a una fecha
	mountTime := time.Unix(int64(sb.S_mtime), 0)
	// Convertir el tiempo de desmontaje a una fecha
	unmountTime := time.Unix(int64(sb.S_umtime), 0)

	fmt.Printf("Filesystem Type: %d\n", sb.S_filesystem_type)
	fmt.Printf("Inodes Count: %d\n", sb.S_inodes_count)
	fmt.Printf("Blocks Count: %d\n", sb.S_blocks_count)
	fmt.Printf("Free Inodes Count: %d\n", sb.S_free_inodes_count)
	fmt.Printf("Free Blocks Count: %d\n", sb.S_free_blocks_count)
	fmt.Printf("Mount Time: %s\n", mountTime.Format(time.RFC3339))
	fmt.Printf("Unmount Time: %s\n", unmountTime.Format(time.RFC3339))
	fmt.Printf("Mount Count: %d\n", sb.S_mnt_count)
	fmt.Printf("Magic: %d\n", sb.S_magic)
	fmt.Printf("Inode Size: %d\n", sb.S_inode_size)
	fmt.Printf("Block Size: %d\n", sb.S_block_size)
	fmt.Printf("First Inode: %d\n", sb.S_first_ino)
	fmt.Printf("First Block: %d\n", sb.S_first_blo)
	fmt.Printf("Bitmap Inode Start: %d\n", sb.S_bm_inode_start)
	fmt.Printf("Bitmap Block Start: %d\n", sb.S_bm_block_start)
	fmt.Printf("Inode Start: %d\n", sb.S_inode_start)
	fmt.Printf("Block Start: %d\n", sb.S_block_start)
}

// Imprimir inodos
func (sb *SuperBlock) PrintInodes(path string) error {
	// Imprimir inodos
	fmt.Println("\nInodos\n----------------")
	// Iterar sobre cada inodo
	for i := int32(0); i < sb.S_inodes_count; i++ {
		inode := &Inode{}
		// Deserializar el inodo
		err := inode.Deserialize(path, int64(sb.S_inode_start+(i*sb.S_inode_size)))
		if err != nil {
			return err
		}
		// Imprimir el inodo
		fmt.Printf("\nInodo %d:\n", i)
		inode.Print()
	}

	return nil
}

// Impriir bloques
func (sb *SuperBlock) PrintBlocks(path string) error {
	// Imprimir bloques
	fmt.Println("\nBloques\n----------------")
	// Iterar sobre cada inodo
	for i := int32(0); i < sb.S_inodes_count; i++ {
		inode := &Inode{}
		// Deserializar el inodo
		err := inode.Deserialize(path, int64(sb.S_inode_start+(i*sb.S_inode_size)))
		if err != nil {
			return err
		}
		// Iterar sobre cada bloque del inodo (apuntadores)
		for _, blockIndex := range inode.I_block {
			// Si el bloque no existe, salir
			if blockIndex == -1 {
				break
			}
			// Si el inodo es de tipo carpeta
			if inode.I_type[0] == '0' {
				block := &FolderBlock{}
				// Deserializar el bloque
				err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
				if err != nil {
					return err
				}
				// Imprimir el bloque
				fmt.Printf("\nBloque %d:\n", blockIndex)
				block.Print()
				continue

				// Si el inodo es de tipo archivo
			} else if inode.I_type[0] == '1' {
				block := &FileBlock{}
				// Deserializar el bloque
				err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
				if err != nil {
					return err
				}
				// Imprimir el bloque
				fmt.Printf("\nBloque %d:\n", blockIndex)
				block.Print()
				continue
			}

		}
	}

	return nil
}

func (sb *SuperBlock) GenerateBlocksDot(path string, outputPath string) error {
    dotContent := `digraph G {
        node [shape=plaintext, fontname="Times"]
        edge [color="#4682B4", arrowhead=vee]
        rankdir=LR;
    `

    var connections []string
    var lastBlockIndex int32 = -1

    // Iterar sobre cada inodo
    for i := int32(0); i < sb.S_inodes_count; i++ {
        inode := &Inode{}
        err := inode.Deserialize(path, int64(sb.S_inode_start+(i*sb.S_inode_size)))
        if err != nil {
            return err
        }

        for _, blockIndex := range inode.I_block {
            if blockIndex == -1 {
                break
            }

            var blockContent string
            if inode.I_type[0] == '0' { // Bloque de carpeta
                block := &FolderBlock{}
                err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
                if err != nil {
                    return err
                }

                blockContent = "<table border='0' cellborder='1' cellspacing='0' cellpadding='4'>"
                blockContent += fmt.Sprintf("<tr><td bgcolor='#5F9EA0' style='color:white;' colspan='2'><b>Bloque Carpeta %d</b></td></tr>", blockIndex)
                blockContent += "<tr><td><b>name</b></td><td><b>inodo</b></td></tr>"
                for _, content := range block.B_content {
                    name := strings.TrimRight(string(content.B_name[:]), "\x00")
                    blockContent += fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>", name, content.B_inodo)
                }
                blockContent += "</table>"
            } else if inode.I_type[0] == '1' { // Bloque de archivo
                block := &FileBlock{}
                err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
                if err != nil {
                    return err
                }
                content := strings.TrimRight(string(block.B_content[:]), "\x00")
                blockContent = "<table border='0' cellborder='1' cellspacing='0' cellpadding='4'>"
                blockContent += fmt.Sprintf("<tr><td bgcolor='#5F9EA0' style='color:white;'><b>Bloque Archivo %d</b></td></tr>", blockIndex)
                blockContent += fmt.Sprintf("<tr><td>%s</td></tr>", html.EscapeString(content))
                blockContent += "</table>"
            }

            dotContent += fmt.Sprintf("block%d [label=<%s>];", blockIndex, blockContent)
            
            // Conectar bloques secuencialmente
            if lastBlockIndex != -1 {
                connections = append(connections, fmt.Sprintf("block%d -> block%d;", lastBlockIndex, blockIndex))
            }
            lastBlockIndex = blockIndex
        }
    }

    dotContent += strings.Join(connections, "\n")
    dotContent += "}"

    // Extraer la ruta del directorio y el nombre base del archivo
    dir := filepath.Dir(outputPath)                                      // Carpeta donde se guardará
    fileBase := strings.TrimSuffix(filepath.Base(outputPath), ".png")    // Nombre sin extensión `.png`
    dotFilePath := filepath.Join(dir, fileBase+".dot")                   // Ruta para el archivo `.dot`
    pngFilePath := filepath.Join(dir, fileBase+".png")                   // Ruta final del `.png`

    // Crear el archivo `.dot`
    dotFile, err := os.Create(dotFilePath)
    if err != nil {
        return err
    }
    defer dotFile.Close()

    _, err = dotFile.WriteString(dotContent)
    if err != nil {
        return err
    }

    // Ejecutar Graphviz para convertir el `.dot` en `.png`
    cmd := exec.Command("dot", "-Tpng", dotFilePath, "-o", pngFilePath)
    err = cmd.Run()
    if err != nil {
        return err
    }

    fmt.Println("Diagrama de bloques generado:", pngFilePath)
    return nil
}





// Get users.txt block
func (sb *SuperBlock) GetUsersBlock(path string) (*FileBlock, error) {
	// Ir al inodo 1
	inode := &Inode{}

	// Deserializar el inodo
	err := inode.Deserialize(path, int64(sb.S_inode_start+(1*sb.S_inode_size))) // 1 porque es el inodo 1
	if err != nil {
		return nil, err
	}

	// Iterar sobre cada bloque del inodo (apuntadores)
	for _, blockIndex := range inode.I_block {
		// Si el bloque no existe, salir
		if blockIndex == -1 {
			break
		}
		// Si el inodo es de tipo archivo
		if inode.I_type[0] == '1' {
			block := &FileBlock{}
			// Deserializar el bloque
			err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
			if err != nil {
				return nil, err
			}
			// Deben ir guardando todo el contenido de los bloques en una variable

			// Retornar el bloque por temas explicativos
			return block, nil
		}
	}
	return nil, fmt.Errorf("users.txt block not found")
}
// CreateFolder crea una carpeta en el sistema de archivos
func (sb *SuperBlock) CreateFolder(path string, parentsDir []string, destDir string) error {

	// Validar el sistema de archivos
	if sb.S_filesystem_type == 3 {
		// Si parentsDir está vacío, solo trabajar con el primer inodo que sería el raíz "/"
		if len(parentsDir) == 0 {
			return sb.createFolderInInodeExt3(path, 0, parentsDir, destDir)
		}

		// Iterar sobre cada inodo ya que se necesita buscar el inodo padre
		for i := int32(0); i < sb.S_inodes_count; i++ {
			err := sb.createFolderInInodeExt3(path, i, parentsDir, destDir)
			if err != nil {
				return err
			}
		}
	} else {
		// Si parentsDir está vacío, solo trabajar con el primer inodo que sería el raíz "/"
		if len(parentsDir) == 0 {
			return sb.createFolderInInodeExt2(path, 0, parentsDir, destDir)
		}

		// Iterar sobre cada inodo ya que se necesita buscar el inodo padre
		for i := int32(0); i < sb.S_inodes_count; i++ {
			err := sb.createFolderInInodeExt2(path, i, parentsDir, destDir)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
// CreateFile crea un archivo en el sistema de archivos
func (sb *SuperBlock) CreateFile(path string, parentsDir []string, destFile string, size int, cont []string) error {

	// Si parentsDir está vacío, solo trabajar con el primer inodo que sería el raíz "/"
	if len(parentsDir) == 0 {
		return sb.createFileInInode(path, 0, parentsDir, destFile, size, cont)
	}

	// Iterar sobre cada inodo ya que se necesita buscar el inodo padre
	for i := int32(0); i < sb.S_inodes_count; i++ {
		err := sb.createFileInInode(path, i, parentsDir, destFile, size, cont)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sb *SuperBlock) FolderExists(partitionPath string, folderPath string) (bool, error) {
    // Dividir el path en directorios padres y el directorio destino
    parentDirs, destDir := utils.GetParentDirectories(folderPath)

    // Buscar el inodo correspondiente al directorio destino
    inode, err := sb.FindInode(partitionPath, parentDirs, destDir)
    if err != nil {
        return false, nil // Si no se encuentra, asumimos que no existe
    }

    // Verificar si el inodo encontrado es un directorio
    if inode != nil && inode.I_type[0] == '0' { // '0' indica que es un directorio
        return true, nil
    }

    return false, nil
}
func (sb *SuperBlock) FindInode(partitionPath string, parentDirs []string, destDir string) (*Inode, error) {
    // Comenzar desde el inodo raíz
    currentInode := &Inode{}
    err := currentInode.Deserialize(partitionPath, int64(sb.S_inode_start))
    if err != nil {
        return nil, fmt.Errorf("error al deserializar el inodo raíz: %w", err)
    }

    // Recorrer los directorios padres
    for _, dir := range parentDirs {
        found := false
        for _, blockIndex := range currentInode.I_block {
            if blockIndex == -1 {
                break
            }

            // Leer el bloque de carpeta
            folderBlock := &FolderBlock{}
            err := folderBlock.Deserialize(partitionPath, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
            if err != nil {
                return nil, fmt.Errorf("error al deserializar el bloque de carpeta: %w", err)
            }

            // Buscar el directorio en el bloque
            for _, content := range folderBlock.B_content {
                name := strings.TrimRight(string(content.B_name[:]), "\x00")
                if name == dir {
                    // Cargar el inodo correspondiente
                    currentInode = &Inode{}
                    err := currentInode.Deserialize(partitionPath, int64(sb.S_inode_start+(content.B_inodo*sb.S_inode_size)))
                    if err != nil {
                        return nil, fmt.Errorf("error al deserializar el inodo: %w", err)
                    }
                    found = true
                    break
                }
            }

            if found {
                break
            }
        }

        if !found {
            return nil, fmt.Errorf("directorio '%s' no encontrado", dir)
        }
    }

    // Verificar si el directorio destino existe
    for _, blockIndex := range currentInode.I_block {
        if blockIndex == -1 {
            break
        }

        // Leer el bloque de carpeta
        folderBlock := &FolderBlock{}
        err := folderBlock.Deserialize(partitionPath, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
        if err != nil {
            return nil, fmt.Errorf("error al deserializar el bloque de carpeta: %w", err)
        }

        // Buscar el directorio destino en el bloque
        for _, content := range folderBlock.B_content {
            name := strings.TrimRight(string(content.B_name[:]), "\x00")
            if name == destDir {
                // Cargar el inodo correspondiente
                destInode := &Inode{}
                err := destInode.Deserialize(partitionPath, int64(sb.S_inode_start+(content.B_inodo*sb.S_inode_size)))
                if err != nil {
                    return nil, fmt.Errorf("error al deserializar el inodo: %w", err)
                }
                return destInode, nil
            }
        }
    }

    return nil, fmt.Errorf("directorio destino '%s' no encontrado", destDir)
}

func (sb *SuperBlock) GenerateTreeDot(path string, outputPath string) error {
    dotContent := `digraph EXT2_Tree {
        node [shape=record, style=filled, fontname="Times"];
    `

    // Comenzar desde el inodo raíz
    rootInode := &Inode{}
    err := rootInode.Deserialize(path, int64(sb.S_inode_start))
    if err != nil {
        return fmt.Errorf("error al deserializar el inodo raíz: %w", err)
    }

    // Mapa para registrar los inodos visitados
    visited := make(map[int32]bool)

    // Generar el contenido del árbol recursivamente
    var generateTree func(inode *Inode, inodeIndex int32) string
    generateTree = func(inode *Inode, inodeIndex int32) string {
        if visited[inodeIndex] {
            return ""
        }
        visited[inodeIndex] = true

        // Convertir tiempos a string
        atime := time.Unix(int64(inode.I_atime), 0).Format("02/01/2006 15:04")
        ctime := time.Unix(int64(inode.I_ctime), 0).Format("02/01/2006 15:04")
        mtime := time.Unix(int64(inode.I_mtime), 0).Format("02/01/2006 15:04")

        // Información del inodo
        content := fmt.Sprintf(`"INODO %d" [label="{INODO %d | UID: %d | GID: %d | Size: %d | Atime: %s | Ctime: %s | Mtime: %s | Tipo: %c | Perm: %s`,
            inodeIndex, inodeIndex, inode.I_uid, inode.I_gid, inode.I_size, atime, ctime, mtime, rune(inode.I_type[0]), string(inode.I_perm[:]))

        // Agregar los valores de I_block en filas individuales
        for i, blockIndex := range inode.I_block {
            content += fmt.Sprintf(" | l_block %d: %d", i, blockIndex)
        }
        content += `}"];`

        // Procesar bloques
        for _, blockIndex := range inode.I_block {
            if blockIndex == -1 {
                break
            }

            if inode.I_type[0] == '0' { // Bloque de carpeta
                folderBlock := &FolderBlock{}
                err := folderBlock.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
                if err != nil {
                    continue
                }

                blockContent := fmt.Sprintf(`"BLOQUE %d" [label="{BLOQUE %d`, blockIndex, blockIndex)
                for _, contentEntry := range folderBlock.B_content {
                    name := strings.TrimRight(string(contentEntry.B_name[:]), "\x00")
                    if name != "" {
                        blockContent += fmt.Sprintf(" | %s -> Inodo %d", name, contentEntry.B_inodo)
                    }
                }
                blockContent += `}"];`
                dotContent += blockContent

                // Conectar inodo con bloque
                dotContent += fmt.Sprintf(`"INODO %d" -> "BLOQUE %d";`, inodeIndex, blockIndex)
            } else if inode.I_type[0] == '1' { // Bloque de archivo
                fileBlock := &FileBlock{}
                err := fileBlock.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
                if err != nil {
                    continue
                }

                content := strings.TrimRight(string(fileBlock.B_content[:]), "\x00")
                blockContent := fmt.Sprintf(`"BLOQUE %d" [label="{BLOQUE %d | Contenido: %s}"];`, blockIndex, blockIndex, html.EscapeString(content))
                dotContent += blockContent

                // Conectar inodo con bloque
                dotContent += fmt.Sprintf(`"INODO %d" -> "BLOQUE %d";`, inodeIndex, blockIndex)
            }
        }

        // Procesar hijos recursivamente
        for _, blockIndex := range inode.I_block {
            if blockIndex == -1 {
                break
            }

            folderBlock := &FolderBlock{}
            err := folderBlock.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
            if err != nil {
                continue
            }

            for _, contentEntry := range folderBlock.B_content {
                name := strings.TrimRight(string(contentEntry.B_name[:]), "\x00")
                if name == "" {
                    continue
                }

                childInode := &Inode{}
                err := childInode.Deserialize(path, int64(sb.S_inode_start+(contentEntry.B_inodo*sb.S_inode_size)))
                if err != nil {
                    continue
                }

                dotContent += fmt.Sprintf(`"INODO %d" -> "INODO %d" [label="%s"];`, inodeIndex, contentEntry.B_inodo, name)
                dotContent += generateTree(childInode, contentEntry.B_inodo)
            }
        }

        return content
    }

    dotContent += generateTree(rootInode, 0)
    dotContent += "}"

    // Crear el archivo `.dot`
    dir := filepath.Dir(outputPath)
    fileBase := strings.TrimSuffix(filepath.Base(outputPath), ".png")
    dotFilePath := filepath.Join(dir, fileBase+".dot")
    pngFilePath := filepath.Join(dir, fileBase+".png")

    err = os.MkdirAll(dir, os.ModePerm)
    if err != nil {
        return fmt.Errorf("error al crear las carpetas padre: %w", err)
    }

    dotFile, err := os.Create(dotFilePath)
    if err != nil {
        return fmt.Errorf("error al crear el archivo .dot: %w", err)
    }
    defer dotFile.Close()

    _, err = dotFile.WriteString(dotContent)
    if err != nil {
        return fmt.Errorf("error al escribir en el archivo .dot: %w", err)
    }

    // Ejecutar Graphviz para convertir el `.dot` en `.png`
    cmd := exec.Command("dot", "-Tpng", dotFilePath, "-o", pngFilePath)
    err = cmd.Run()
    if err != nil {
        return fmt.Errorf("error al ejecutar Graphviz: %w", err)
    }

    fmt.Println("Reporte de árbol generado:", pngFilePath)
    return nil
}

