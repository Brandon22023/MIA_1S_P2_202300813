package global

// Variable global para almacenar el ID de la partición activa
var ActivePartitionID string

// Lista global para almacenar los paths válidos
var ValidPaths []string


var ValidFilePaths_mkfile []string


// Devuelve una copia de la lista de paths válidos de mkfile
func GetValidFilePathsMkfile() []string {
    return ValidFilePaths_mkfile
}

// Asigna un nuevo valor a la lista de paths válidos de mkfile
func SetValidFilePathsMkfile(paths []string) {
    ValidFilePaths_mkfile = paths
}