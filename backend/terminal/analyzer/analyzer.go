package analyzer

import (
	"errors"
	"fmt"
	"strings"
	commands "terminal/commands"
	"terminal/stores"
)

// Analyzer analiza el comando de entrada y ejecuta la acción correspondiente
func Analyzer(input string) (string, error) {
	// Divide la entrada en tokens usando espacios en blanco como delimitadores
	tokens := strings.Fields(input)

	// Si no se proporcionó ningún comando, devuelve un error
	if len(tokens) == 0 {
		return "", errors.New("no se proporcionó ningún comando")
	}

	// Switch para manejar diferentes comandos
	switch tokens[0] {
	case "mkdir":
		// Llama a la función Mkdir del paquete commands con los argumentos restantes
		return commands.ParseMkdir(tokens[1:])
    case "mkdisk":
		// Llama a la función ParseMkdisk del paquete commands con los argumentos restantes
		return commands.ParseMkdisk(tokens[1:])
	case "rmdisk":
	    return commands.ParseRmdisk(tokens[1:])
    case "fdisk":
		// Llama a la función CommandFdisk del paquete commands con los argumentos restantes
		return commands.ParseFdisk(tokens[1:])
	case "mount":
		// Llama a la función CommandMount del paquete commands con los argumentos restantes
		return commands.ParseMount(tokens[1:])
	case "mkfs":
		// Llama a la función CommandMkfs del paquete commands con los argumentos restantes
		return commands.ParseMkfs(tokens[1:])
	case "rep":
		// Llama a la función CommandRep del paquete commands con los argumentos restantes
		return commands.ParseRep(tokens[1:])
	case "login":
		return commands.ParseLogin(tokens[1:])
	case "mounted":
		return stores.PrintMountedPartitions()
	case "logout":
		return commands.CommandLogout()
	case "mkfile":
		return commands.ParserMkfile(tokens[1:])
	case "unmount":
		return commands.ParseUnmount(tokens[1:])	
	case "remove":
        return commands.ParseRemove(tokens[1:])
	case "rename":
		return commands.ParseRename(tokens[1:])
	case "copy":
        return commands.ParseCopy(tokens[1:])
	case "move":
        return commands.ParseMove(tokens[1:])
	case "find":
        return commands.ParseFind(tokens[1:])
	default:
		// Si el comando no es reconocido, devuelve un error
		return "", fmt.Errorf("comando desconocido: %s", tokens[0])
	}
}
