package commands

import (
    "errors"
    "fmt"
    "os"
    "regexp"
    "strings"
)

// RMDISK estructura que representa el comando rmdisk con sus parámetros
type RMDISK struct {
    path string // Ruta del archivo del disco
}

/*
   rmdisk -path="/home/mis discos/Disco4.mia"
*/

func ParseRmdisk(tokens []string) (string, error) {
    cmd := &RMDISK{} // Crea una nueva instancia de RMDISK

    // Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
    args := strings.Join(tokens, " ")
    // Expresión regular para encontrar los parámetros del comando rmdisk
    re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+`)
    // Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
    matches := re.FindAllString(args, -1)

    // Itera sobre cada coincidencia encontrada
    for _, match := range matches {
        // Divide cada parte en clave y valor usando "=" como delimitador
        kv := strings.SplitN(match, "=", 2)
        if len(kv) != 2 {
            return "", fmt.Errorf("formato de parámetro inválido: %s", match)
        }
        key, value := strings.ToLower(kv[0]), kv[1]

        // Remove quotes from value if present
        if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
            value = strings.Trim(value, "\"")
        }

        // Switch para manejar diferentes parámetros
        switch key {
        case "-path":
            // Verifica que el path no esté vacío
            if value == "" {
                return "", errors.New("el path no puede estar vacío")
            }
            cmd.path = value
        default:
            // Si el parámetro no es reconocido, devuelve un error
            return "", fmt.Errorf("parámetro desconocido: %s", key)
        }
    }

    // Verifica que el parámetro -path haya sido proporcionado
    if cmd.path == "" {
        return "", errors.New("faltan parámetros requeridos: -path")
    }

    // Eliminar el disco con los parámetros proporcionados
    err := commandRmdisk(cmd)
    if err != nil {
        return "", err
    }

    // Devuelve un mensaje de éxito con los detalles del disco creado
	return fmt.Sprintf("RMDISK: Disco eliminado exitosamente\n"+
    "-> Path: %s\n",
    cmd.path), nil
}

func commandRmdisk(rmdisk *RMDISK) error {
    // Verificar si el archivo existe
    if _, err := os.Stat(rmdisk.path); os.IsNotExist(err) {
        return fmt.Errorf("el archivo no existe: %s", rmdisk.path)
    }

    // Eliminar el archivo
    err := os.Remove(rmdisk.path)
    if err != nil {
        return fmt.Errorf("error al eliminar el archivo: %v", err)
    }

    return nil
}