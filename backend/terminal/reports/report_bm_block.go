package reports

import (
	structures "terminal/structures"
	utils "terminal/utils"
	"fmt"
	"os"
	"strings"
)

// ReportBMInode genera un reporte del bitmap de blockes y lo guarda en la ruta especificada
func ReportBMIblock(superblock *structures.SuperBlock, diskPath string, path string) error {
	// Crear las carpetas padre si no existen
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Abrir el archivo de disco
	file, err := os.Open(diskPath)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo de disco: %v", err)
	}
	defer file.Close()

	// Calcular el número total de block
	totalBlocks := superblock.S_blocks_count + superblock.S_free_blocks_count

	// Obtener el contenido del bitmap de block
	var bitmapContent strings.Builder

	for i := int32(0); i < totalBlocks; i++ {
		// Establecer el puntero
		_, err := file.Seek(int64(superblock.S_bm_block_start+i), 0)
		if err != nil {
			return fmt.Errorf("error al establecer el puntero en el archivo: %v", err)
		}

		// Leer un byte (carácter '0' o '1')
		char := make([]byte, 1)
		_, err = file.Read(char)
		if err != nil {
			return fmt.Errorf("error al leer el byte del archivo: %v", err)
		}

		// Verificar si el carácter leído es una 'X' y reemplazarlo por '1'
        if char[0] == 'X' {
            fmt.Printf("Advertencia: carácter 'X' encontrado en la posición %d. Reemplazando por '1'.\n", superblock.S_bm_inode_start+i)
            char[0] = '1'
        }

		// Agregar el carácter al contenido del bitmap
		bitmapContent.WriteByte(char[0])

		// Agregar un carácter de nueva línea cada 20 caracteres (20 blocks)
		if (i+1)%20 == 0 {
			bitmapContent.WriteString("\n")
		}
	}

	// Crear el archivo TXT
	txtFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error al crear el archivo TXT: %v", err)
	}
	defer txtFile.Close()

	// Escribir el contenido del bitmap en el archivo TXT
	_, err = txtFile.WriteString(bitmapContent.String())
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo TXT: %v", err)
	}

	fmt.Println("Archivo del bitmap de blockes generado:", path)
	return nil
}
