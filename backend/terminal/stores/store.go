package stores

import (
	structures "terminal/structures"
	"errors"
	"fmt"
	"strings"
)
// Carnet de estudiante
const Carnet string = "13" // 202300813

// Declaración de variables globales
var (
	MountedPartitions map[string]string = make(map[string]string)
)
//imprimri para el mounted



// ParseMounted verifica el comando y, si es válido, imprime las particiones montadas


// PrintMountedPartitions devuelve el contenido de MountedPartitions como una cadena y un posible error
func PrintMountedPartitions() (string, error) {
    if len(MountedPartitions) == 0 {
		return "", errors.New("no hay particiones montadas")
    }

    var result strings.Builder
    result.WriteString("MOUNTED: Particiones montadas:\n")
    for id:= range MountedPartitions {
        result.WriteString(fmt.Sprintf("ID: %s\n", id))
    }

    return result.String(), nil
}

//---------------------------
// GetMountedPartition obtiene la partición montada con el id especificado
func GetMountedPartition(id string) (*structures.PARTITION, string, error) {
	// Obtener el path de la partición montada
	path := MountedPartitions[id]
	if path == "" {
		return nil, "", errors.New("la partición no está montada")
	}

	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(path)
	if err != nil {
		return nil, "", err
	}

	// Buscar la partición con el id especificado
	partition, err := mbr.GetPartitionByID(id)
	if partition == nil {
		return nil, "", err
	}

	return partition, path, nil
}

// GetMountedMBR obtiene el MBR de la partición montada con el id especificado
func GetMountedPartitionRep(id string) (*structures.MBR, *structures.SuperBlock, string, error) {
	// Obtener el path de la partición montada
    path, exists := MountedPartitions[id]
    if !exists {
        return nil, nil, "", fmt.Errorf("error: la partición con ID '%s' no está montada. Verifique que el ID sea correcto", id)
    }

	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	//EL DESERIALIZE UTILIZA EL DEL MBR PARA TOMARLO EN CUENTA
	err := mbr.DeserializeMBR(path)
	if err != nil {
		return nil, nil, "", err
	}

	// Buscar la partición con el id especificado
	partition, err := mbr.GetPartitionByID(id)
	if partition == nil {
		return nil, nil, "", err
	}

	// Crear una instancia de SuperBlock
	var sb structures.SuperBlock

	// Deserializar la estructura SuperBlock desde un archivo binario
	err = sb.Deserialize(path, int64(partition.Part_start))
	if err != nil {
		return nil, nil, "", err
	}

	return &mbr, &sb, path, nil
}

// GetMountedPartitionSuperblock obtiene el SuperBlock de la partición montada con el id especificado
func GetMountedPartitionSuperblock(id string) (*structures.SuperBlock, *structures.PARTITION, string, error) {
	// Obtener el path de la partición montada
	path := MountedPartitions[id]
	if path == "" {
		return nil, nil, "", errors.New("la partición no está montada")
	}

	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(path)
	if err != nil {
		return nil, nil, "", err
	}

	// Buscar la partición con el id especificado
	partition, err := mbr.GetPartitionByID(id)
	if partition == nil {
		return nil, nil, "", err
	}

	// Crear una instancia de SuperBlock
	var sb structures.SuperBlock

	// Deserializar la estructura SuperBlock desde un archivo binario
	err = sb.Deserialize(path, int64(partition.Part_start))
	if err != nil {
		return nil, nil, "", err
	}

	return &sb, partition, path, nil
}
func GetActivePartitionID() (string, error) {
    if Auth.IsAuthenticated() {
        return Auth.GetPartitionID(), nil
    }
    return "", errors.New("no se ha iniciado sesión en ninguna partición")
}