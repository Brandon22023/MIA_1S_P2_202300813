package reports

import (
	structures "terminal/structures"
	utils "terminal/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ReportMBR genera un reporte del MBR y lo guarda en la ruta especificada
func ReportMBR(mbr *structures.MBR, path string) error {
	// Crear las carpetas padre si no existen
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := utils.GetFileNames(path)

	// Definir el contenido DOT con una tabla
	dotContent := fmt.Sprintf(`digraph G {
		node [shape=plaintext]
		tabla [label=<
			<table border="0" cellborder="1" cellspacing="0" cellpadding="4" style="rounded; font-family:Arial; font-size:12px;">
				<!-- Encabezado principal -->
				<tr>
					<td colspan="2" bgcolor="#4A024A" style="color:white; font-size:16px; padding:8px; border-top-left-radius:8px; border-top-right-radius:8px;">
						<b>REPORTE DE MBR</b>
					</td>
				</tr>
				
				<!-- Datos del MBR -->
				<tr bgcolor="#EAD3EA">
					<td><b>mbr_tamano</b></td>
					<td>%d</td>
				</tr>
				<tr>
					<td><b>mbr_fecha_creacion</b></td>
					<td>%s</td>
				</tr>
				<tr bgcolor="#EAD3EA">
					<td><b>mbr_disk_signature</b></td>
					<td>%d</td>
				</tr>
				
				<!-- Separador visual -->
				<tr><td colspan="2" height="4" bgcolor="#4A024A"></td></tr>`, 
		mbr.Mbr_size, time.Unix(int64(mbr.Mbr_creation_date), 0), mbr.Mbr_disk_signature)

	// Agregar las particiones a la tabla
	for i, part := range mbr.Mbr_partitions {
		/*
			// Continuar si el tamaño de la partición es -1 (o sea, no está asignada)
			if part.Part_size == -1 {
				continue
			}
		*/

		// Convertir Part_name a string y eliminar los caracteres nulos
		partName := strings.TrimRight(string(part.Part_name[:]), "\x00")
		// Convertir Part_status, Part_type y Part_fit a char
		partStatus := rune(part.Part_status[0])
		partType := rune(part.Part_type[0])
		partFit := rune(part.Part_fit[0])

		// Agregar la partición a la tabla
		dotContent += fmt.Sprintf(`
            <!-- Partición %d -->
            <tr>
                <td colspan="2" bgcolor="#720072" style="color:white; font-size:14px; padding:6px;">
                    <b>PARTICIÓN %d</b>
                </td>
            </tr>
            <tr bgcolor="#F5D0F5">
                <td><b>part_status</b></td>
                <td>%c</td>
            </tr>
            <tr>
                <td><b>part_type</b></td>
                <td>%c</td>
            </tr>
            <tr bgcolor="#F5D0F5">
                <td><b>part_fit</b></td>
                <td>%c</td>
            </tr>
            <tr>
                <td><b>part_start</b></td>
                <td>%d</td>
            </tr>
            <tr bgcolor="#F5D0F5">
                <td><b>part_size</b></td>
                <td>%d</td>
            </tr>
            <tr>
                <td><b>part_name</b></td>
                <td>%s</td>
            </tr>
            
            <!-- Separador visual entre particiones -->
            <tr><td colspan="2" height="4" bgcolor="#4A024A"></td></tr>`,
        i+1, i+1, partStatus, partType, partFit, part.Part_start, part.Part_size, partName)
}

	// Cerrar la tabla y el contenido DOT
	dotContent += `
	<!-- Pie de tabla -->
	<tr>
		<td colspan="2" bgcolor="#4A024A" style="border-bottom-left-radius:8px; border-bottom-right-radius:8px; height:4px;"></td>
	</tr>
	</table>>] }"
`

	// Guardar el contenido DOT en un archivo
	file, err := os.Create(dotFileName)
	if err != nil {
		return fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(dotContent)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}

	// Ejecutar el comando Graphviz para generar la imagen
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el comando Graphviz: %v", err)
	}

	fmt.Println("Imagen de la tabla generada:", outputImage)
	return nil
}

// ReportDiskStructure genera un reporte de la estructura del disco (Primarias y Extendidas)
func ReportDiskStructure(mbr *structures.MBR, path string) error {
    // Crear las carpetas padre si no existen
    err := utils.CreateParentDirs(path)
    if err != nil {
        return err
    }

    // Obtener el nombre base del archivo sin la extensión
    dotFileName, outputImage := utils.GetFileNames(path)

    // Iniciar el contenido DOT
    dotContent := `digraph DiskStructure {
    rankdir=LR;
    node [shape=record, style=filled, fillcolor=lightgray, fontname="Arial", fontsize=12];
    edge [color=black];

    Disk [label="{ MBR |`

    // Variables para calcular el espacio libre
    totalPartitions := 4
    assignedPartitions := 0

    // Iterar sobre las particiones del MBR
    for _, part := range mbr.Mbr_partitions {
        // Si la partición no está asignada, continuar
        if part.Part_size == -1 {
            continue
        }

        // Incrementar el contador de particiones asignadas
        assignedPartitions++

        // Convertir Part_name a string y eliminar los caracteres nulos
        partName := strings.TrimRight(string(part.Part_name[:]), "\x00")
        // Convertir Part_type a char
        partType := rune(part.Part_type[0])

        // Verificar si es una partición Extendida o Primaria
        if partType == 'E' {
            // Partición Extendida
            dotContent += " { Extendida | { EBR } } |"
        } else if partType == 'P' {
            // Partición Primaria
            dotContent += " Primaria (" + partName + ") |"
        }
    }

    // Calcular el espacio libre restante
    freePartitions := totalPartitions - assignedPartitions
    if freePartitions > 0 {
        dotContent += fmt.Sprintf(" Libre (%d%%) }", freePartitions*25)
    } else {
        dotContent += " }"
    }

    // Cerrar el contenido DOT
    dotContent += `", shape=record, style=filled, fillcolor=lightblue, fontcolor=black, penwidth=2];
}`

    // Guardar el contenido DOT en un archivo
    file, err := os.Create(dotFileName)
    if err != nil {
        return fmt.Errorf("error al crear el archivo: %v", err)
    }
    defer file.Close()

    _, err = file.WriteString(dotContent)
    if err != nil {
        return fmt.Errorf("error al escribir en el archivo: %v", err)
    }

    // Ejecutar el comando Graphviz para generar la imagen
    cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
    err = cmd.Run()
    if err != nil {
        return fmt.Errorf("error al ejecutar el comando Graphviz: %v", err)
    }

    fmt.Println("Imagen de la estructura del disco generada:", outputImage)
    return nil
}