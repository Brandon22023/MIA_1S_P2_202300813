package reports

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "terminal/structures"
    "terminal/utils"
)

// ReportFile genera un reporte basado en el contenido de un archivo y lo guarda como una imagen
func ReportFile(superblock *structures.SuperBlock, diskPath string, path string, pathFileLs string, name string) error {
    // Validar que el archivo especificado en path_file_ls exista
    if _, err := os.Stat(pathFileLs); os.IsNotExist(err) {
        return fmt.Errorf("error: el archivo especificado en path_file_ls '%s' no existe", pathFileLs)
    }

    // Leer el contenido del archivo especificado en path_file_ls
    fileData, err := os.ReadFile(pathFileLs)
    if err != nil {
        return fmt.Errorf("error al leer el archivo '%s': %w", pathFileLs, err)
    }
    content := string(fileData)
    fmt.Printf("Contenido del archivo '%s':\n%s\n", pathFileLs, content)

    // Crear las carpetas padre si no existen
    err = utils.CreateParentDirs(path)
    if err != nil {
        return fmt.Errorf("error al crear las carpetas padre para el path '%s': %w", path, err)
    }

    // Obtener el nombre base del archivo sin la extensión
    dotFileName, outputImage := utils.GetFileNames(path)

    // Obtener el nombre del archivo desde path_file_ls
    fileName := filepath.Base(pathFileLs)

    // Iniciar el contenido DOT
    dotContent := fmt.Sprintf(`digraph NotepadWindow {
    graph [bgcolor=white];
    node [shape=box, style=filled, fillcolor=lightgray];
    
    subgraph cluster_window {
        label="%s: Bloc de notas";
        color=blue;
        style=filled;
        fillcolor=white;
`, fileName)

    // Agregar el contenido del archivo como nodos en el reporte
    lines := strings.Split(content, "\n")
    for i, line := range lines {
        dotContent += fmt.Sprintf(`        content%d [label="%s", shape=box, fillcolor=white];`, i, line)
        if i > 0 {
            dotContent += fmt.Sprintf(" content%d -> content%d;\n", i-1, i)
        } else {
            dotContent += "\n"
        }
    }

    // Cerrar el subgrafo y el contenido DOT
    dotContent += `
    }
}`

    // Crear el archivo DOT
    dotFile, err := os.Create(dotFileName)
    if err != nil {
        return fmt.Errorf("error al crear el archivo DOT '%s': %w", dotFileName, err)
    }
    defer dotFile.Close()

    _, err = dotFile.WriteString(dotContent)
    if err != nil {
        return fmt.Errorf("error al escribir en el archivo DOT '%s': %w", dotFileName, err)
    }

    // Generar la imagen con Graphviz
    cmd := exec.Command("dot", "-Tjpg", dotFileName, "-o", outputImage)
    err = cmd.Run()
    if err != nil {
        return fmt.Errorf("error al generar la imagen con Graphviz: %w", err)
    }

    fmt.Printf("Reporte generado exitosamente en: %s\n", outputImage)
    return nil
}