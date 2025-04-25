package commands

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "terminal/structures"
)

// DiskInfo estructura para representar la información del disco y sus particiones
type DiskInfo struct {
    Name       string           `json:"name"`        // Nombre del disco
    Size       int32            `json:"size"`        // Tamaño del disco
    Signature  int32            `json:"signature"`   // Firma del disco
    Fit        string           `json:"fit"`         // Tipo de ajuste
    Partitions []PartitionInfo  `json:"partitions"`  // Lista de particiones
}

// PartitionInfo estructura para representar la información de una partición
type PartitionInfo struct {
    Status      string `json:"status"`      // Estado de la partición
    Type        string `json:"type"`        // Tipo de partición
    Fit         string `json:"fit"`         // Tipo de ajuste
    Start       int32  `json:"start"`       // Byte de inicio
    Size        int32  `json:"size"`        // Tamaño en bytes
    Name        string `json:"name"`        // Nombre de la partición
    Correlative int32  `json:"correlative"` // Número correlativo
    ID          string `json:"id"`          // ID único
}

// ExportDiskInfo extrae la información del disco y la guarda en un archivo JSON
func ExportDiskInfo(diskPath string) error {
    // Obtener el nombre del disco desde el path
    diskName := filepath.Base(diskPath)

    // Crear la carpeta "info_disk" en la misma ubicación del archivo export.go
    currentDir, err := os.Getwd() // Obtener el directorio actual
    if err != nil {
        return fmt.Errorf("error al obtener el directorio actual: %v", err)
    }
    outDir := filepath.Join(currentDir, "info_disk")

    // Validar si la carpeta existe, si no, crearla
    if _, err := os.Stat(outDir); os.IsNotExist(err) {
        err = os.MkdirAll(outDir, os.ModePerm)
        if err != nil {
            return fmt.Errorf("error al crear la carpeta de destino: %v", err)
        }
    }

    // Leer el MBR del archivo binario
    mbr := &structures.MBR{}
    err = mbr.DeserializeMBR(diskPath)
    if err != nil {
        // Si el disco no existe, devolver un error
        return fmt.Errorf("no se puede exportar el disco %s porque no existe en la ruta %s", diskName, diskPath)
    }

    // Crear la estructura DiskInfo
    diskInfo := DiskInfo{
        Name:      diskName,
        Size:      mbr.Mbr_size,
        Signature: mbr.Mbr_disk_signature,
        Fit:       string(mbr.Mbr_disk_fit[:]),
    }

    // Agregar las particiones al DiskInfo
    for _, partition := range mbr.Mbr_partitions {
        // Ignorar particiones no utilizadas
        if partition.Part_start == -1 {
            continue
        }

        partitionInfo := PartitionInfo{
            Status:      string(partition.Part_status[:]),
            Type:        string(partition.Part_type[:]),
            Fit:         string(partition.Part_fit[:]),
            Start:       partition.Part_start,
            Size:        partition.Part_size,
            Name:        strings.TrimSpace(string(partition.Part_name[:])),
            Correlative: partition.Part_correlative,
            ID:          strings.TrimSpace(string(partition.Part_id[:])),
        }
        diskInfo.Partitions = append(diskInfo.Partitions, partitionInfo)
    }

    // Convertir DiskInfo a JSON
    jsonData, err := json.MarshalIndent(diskInfo, "", "  ")
    if err != nil {
        return fmt.Errorf("error al convertir la información a JSON: %v", err)
    }

    // Guardar el JSON en un archivo dentro de la carpeta "info_disk"
    outputFilePath := filepath.Join(outDir, strings.ReplaceAll(diskName, ".mia", ".json"))
    err = os.WriteFile(outputFilePath, jsonData, 0644)
    if err != nil {
        return fmt.Errorf("error al guardar el archivo JSON: %v", err)
    }

    fmt.Printf("Información del disco exportada exitosamente a %s\n", outputFilePath)
    return nil
}