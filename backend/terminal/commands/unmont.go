package commands

import (
    "errors"
    "fmt"
    "strings"
    stores "terminal/stores"
    structures "terminal/structures"
	"regexp" 
)

// UNMOUNT estructura para representar el comando unmount
type UNMOUNT struct {
    id string // ID de la partición montada
}

// ParseUnmount procesa el comando unmount y devuelve un mensaje de éxito o error
func ParseUnmount(tokens []string) (string, error) {
    cmd := &UNMOUNT{} // Crea una nueva instancia de UNMOUNT

    // Unir tokens en una sola cadena y luego dividir por espacios
    args := strings.Join(tokens, " ")
    // Expresión regular para encontrar el parámetro -id
    re := regexp.MustCompile(`-id=[^\s]+`)
    matches := re.FindAllString(args, -1)

    // Validar que se haya proporcionado el parámetro -id
    if len(matches) == 0 {
        return "", errors.New("faltan parámetros requeridos: -id")
    }

    // Extraer el valor del parámetro -id
    kv := strings.SplitN(matches[0], "=", 2)
    if len(kv) != 2 || kv[1] == "" {
        return "", errors.New("formato de parámetro inválido: -id")
    }
    cmd.id = kv[1]

    // Llamar a la función para desmontar la partición
    err := commandUnmount(cmd)
    if err != nil {
        return "", err
    }

    // Mensaje de éxito
    return fmt.Sprintf("UNMOUNT: Partición desmontada exitosamente\n-> ID: %s", cmd.id), nil
}

func commandUnmount(unmount *UNMOUNT) error {
    // Verificar si el ID existe en las particiones montadas
    path, exists := stores.MountedPartitions[unmount.id]
    if !exists {
        return fmt.Errorf("error: la partición con ID '%s' no está montada", unmount.id)
    }

    // Crear una instancia de MBR
    var mbr structures.MBR

    // Deserializar la estructura MBR desde el archivo binario
    err := mbr.DeserializeMBR(path)
    if err != nil {
        return fmt.Errorf("error deserializando el MBR desde el archivo '%s': %v", path, err)
    }

    // Buscar la partición con el ID especificado
    partition, err := mbr.GetPartitionByID(unmount.id)
    if err != nil {
        return fmt.Errorf("error obteniendo la partición con ID '%s': %v", unmount.id, err)
    }
    if partition == nil {
        return fmt.Errorf("error: la partición con ID '%s' no existe en el disco", unmount.id)
    }

    // Cambiar el estado de la partición a desmontada
    partition.Part_status[0] = '0' // Estado desmontado
    partition.Part_correlative = 0 // Correlativo inicial
    // Limpiar el ID de la partición y establecerlo como "N"
	for i := range partition.Part_id {
		partition.Part_id[i] = 0 // Limpiar todos los bytes
	}
	partition.Part_id[0] = 'N' // Asignar explícitamente "N"

    // Serializar el MBR actualizado en el archivo binario
    err = mbr.SerializeMBR(path)
    if err != nil {
        return fmt.Errorf("error serializando el MBR en el archivo '%s': %v", path, err)
    }

    // Eliminar la partición del mapa de particiones montadas
    delete(stores.MountedPartitions, unmount.id)

    return nil
}