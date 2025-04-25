package reports

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	"terminal/structures"
	"terminal/utils"
)

// ReportSuperBlock genera un reporte gráfico del SuperBlock y lo guarda en la ruta especificada
func ReportSuperBlock(sb *structures.SuperBlock, path string) error {
	// Crear las carpetas padre si no existen
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := utils.GetFileNames(path)

	// Convertir fechas
	mountTime := time.Unix(int64(sb.S_mtime), 0).Format(time.RFC3339)
	unmountTime := time.Unix(int64(sb.S_umtime), 0).Format(time.RFC3339)

	// Iniciar el contenido DOT
	dotContent := `digraph G {
		node [shape=plaintext, fontname="Times"]
		edge [color="#4682B4", arrowhead=vee]
		superblock [label=<
			<table border="0" cellborder="1" cellspacing="0" cellpadding="4" style="font-family:Times">
				<tr><td colspan="2" bgcolor="#4682B4" style="color:white; font-size:14px; padding:6px;"><b>SUPERBLOCK REPORT</b></td></tr>
				<tr><td><b>Filesystem Type</b></td><td>` + fmt.Sprintf("%d", sb.S_filesystem_type) + `</td></tr>
				<tr><td><b>Inodes Count</b></td><td>` + fmt.Sprintf("%d", sb.S_inodes_count) + `</td></tr>
				<tr><td><b>Blocks Count</b></td><td>` + fmt.Sprintf("%d", sb.S_blocks_count) + `</td></tr>
				<tr><td><b>Free Inodes</b></td><td>` + fmt.Sprintf("%d", sb.S_free_inodes_count) + `</td></tr>
				<tr><td><b>Free Blocks</b></td><td>` + fmt.Sprintf("%d", sb.S_free_blocks_count) + `</td></tr>
				<tr><td><b>Mount Time</b></td><td>` + mountTime + `</td></tr>
				<tr><td><b>Unmount Time</b></td><td>` + unmountTime + `</td></tr>
				<tr><td><b>Mount Count</b></td><td>` + fmt.Sprintf("%d", sb.S_mnt_count) + `</td></tr>
				<tr><td><b>Magic</b></td><td>` + fmt.Sprintf("%d", sb.S_magic) + `</td></tr>
				<tr><td><b>Inode Size</b></td><td>` + fmt.Sprintf("%d", sb.S_inode_size) + `</td></tr>
				<tr><td><b>Block Size</b></td><td>` + fmt.Sprintf("%d", sb.S_block_size) + `</td></tr>
				<tr><td><b>First Inode</b></td><td>` + fmt.Sprintf("%d", sb.S_first_ino) + `</td></tr>
				<tr><td><b>First Block</b></td><td>` + fmt.Sprintf("%d", sb.S_first_blo) + `</td></tr>
				<tr><td><b>Bitmap Inode Start</b></td><td>` + fmt.Sprintf("%d", sb.S_bm_inode_start) + `</td></tr>
				<tr><td><b>Bitmap Block Start</b></td><td>` + fmt.Sprintf("%d", sb.S_bm_block_start) + `</td></tr>
				<tr><td><b>Inode Start</b></td><td>` + fmt.Sprintf("%d", sb.S_inode_start) + `</td></tr>
				<tr><td><b>Block Start</b></td><td>` + fmt.Sprintf("%d", sb.S_block_start) + `</td></tr>
			</table>>];
	}`

	// Crear el archivo DOT
	dotFile, err := os.Create(dotFileName)
	if err != nil {
		return err
	}
	defer dotFile.Close()

	// Escribir el contenido DOT en el archivo
	_, err = dotFile.WriteString(dotContent)
	if err != nil {
		return err
	}

	// Generar la imagen con Graphviz
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
	err = cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("Imagen del SuperBlock generada:", outputImage)
	return nil
}
