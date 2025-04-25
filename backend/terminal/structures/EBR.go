package structures

import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
    "os"
    "strings"
)

type EBR struct {
    Part_mount [1]byte   // Indica si la partición está montada o no
    Part_fit   [1]byte   // Tipo de ajuste de la partición (B, F, W)
    Part_start int32     // Byte en el que inicia la partición
    Part_size  int32     // Tamaño total de la partición en bytes
    Part_next  int32     // Byte en el que está el próximo EBR (-1 si no hay siguiente)
    Part_name  [16]byte  // Nombre de la partición
}

// SerializeEBR escribe la estructura EBR en un archivo binario
func (ebr *EBR) SerializeEBR(path string, offset int64) error {
    file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return fmt.Errorf("error abriendo el archivo: %v", err)
    }
    defer func() {
        if cerr := file.Close(); cerr != nil {
            fmt.Printf("error cerrando el archivo: %v\n", cerr)
        }
    }()

    // Mover el puntero del archivo al offset especificado
    _, err = file.Seek(offset, 0)
    if err != nil {
        return fmt.Errorf("error moviendo el puntero del archivo: %v", err)
    }

    // Serializar la estructura EBR directamente en el archivo
    err = binary.Write(file, binary.LittleEndian, ebr)
    if err != nil {
        return fmt.Errorf("error serializando el EBR: %v", err)
    }

    return nil
}

// DeserializeEBR lee la estructura EBR desde un archivo binario
func (ebr *EBR) DeserializeEBR(path string, offset int64) error {
    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("error abriendo el archivo: %v", err)
    }
    defer func() {
        if cerr := file.Close(); cerr != nil {
            fmt.Printf("error cerrando el archivo: %v\n", cerr)
        }
    }()

    // Mover el puntero del archivo al offset especificado
    _, err = file.Seek(offset, 0)
    if err != nil {
        return fmt.Errorf("error moviendo el puntero del archivo: %v", err)
    }

    // Leer solo la cantidad de bytes que corresponden al tamaño de la estructura EBR
    buffer := make([]byte, binary.Size(ebr))
    _, err = file.Read(buffer)
    if err != nil {
        return fmt.Errorf("error leyendo el archivo: %v", err)
    }

    // Deserializar los bytes leídos en la estructura EBR
    reader := bytes.NewReader(buffer)
    err = binary.Read(reader, binary.LittleEndian, ebr)
    if err != nil {
        return fmt.Errorf("error deserializando el EBR: %v", err)
    }

    return nil
}

// PrintEBR imprime los valores del EBR
func (ebr *EBR) PrintEBR() {
    fmt.Printf("Part_mount: %c\n", ebr.Part_mount[0])
    fmt.Printf("Part_fit: %c\n", ebr.Part_fit[0])
    fmt.Printf("Part_start: %d\n", ebr.Part_start)
    fmt.Printf("Part_size: %d\n", ebr.Part_size)
    fmt.Printf("Part_next: %d\n", ebr.Part_next)
    fmt.Printf("Part_name: %s\n", strings.Trim(string(ebr.Part_name[:]), "\x00"))
}

// CreatePartition crea una partición lógica con los parámetros proporcionados
func (ebr *EBR) CreatePartition(start int32, size int32, fit string, name string) error {
    if len(fit) != 1 || (fit[0] != 'B' && fit[0] != 'F' && fit[0] != 'W') {
        return errors.New("tipo de ajuste inválido, debe ser 'B', 'F' o 'W'")
    }
    if len(name) > 16 {
        return errors.New("el nombre de la partición no puede exceder los 16 caracteres")
    }

    ebr.Part_mount[0] = 'N' // Inicialmente no montada
    ebr.Part_fit[0] = fit[0]
    ebr.Part_start = start
    ebr.Part_size = size
    ebr.Part_next = -1 // No hay siguiente EBR inicialmente
    copy(ebr.Part_name[:], name)

    return nil
}