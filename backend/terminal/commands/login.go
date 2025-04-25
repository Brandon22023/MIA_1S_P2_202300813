package commands

import (
	stores "terminal/stores"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// LOGIN estructura que representa el comando login con sus parámetros
type LOGIN struct {
	User string // Usuario
	Pass string // Contraseña
	ID   string // ID del disco
}

/*
	login -user=root -pass=123 -id=062A3E2D
*/

func ParseLogin(tokens []string) (string, error) {
	cmd := &LOGIN{} // Crea una nueva instancia de LOGIN

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mkfs
	re := regexp.MustCompile(`-user=[^\s]+|-pass=[^\s]+|-id=[^\s]+`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Verificar que todos los tokens fueron reconocidos por la expresión regular
	if len(matches) != len(tokens) {
		// Identificar el parámetro inválido
		for _, token := range tokens {
			if !re.MatchString(token) {
				return "", fmt.Errorf("parámetro inválido: %s", token)
			}
		}
	}

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
		case "-user":
			if value == "" {
				return "", errors.New("el usuario no puede estar vacío")
			}
			cmd.User = value
		case "-pass":
			if value == "" {
				return "", errors.New("la contraseña no puede estar vacía")
			}
			cmd.Pass = value
		case "-id":
			// Verifica que el id no esté vacío
			if value == "" {
				return "", errors.New("el id no puede estar vacío")
			}
			cmd.ID = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -id haya sido proporcionado
	if cmd.ID == "" {
		return "", errors.New("faltan parámetros requeridos: -id")
	}

	// Si no se proporcionó el tipo, se establece por defecto a "full"
	if cmd.User == "" {
		return "", errors.New("faltan parámetros requeridos: -user")
	}

	// Si no se proporcionó el tipo, se establece por defecto a "full"
	if cmd.Pass == "" {
		return "", errors.New("faltan parámetros requeridos: -pass")
	}

	// Aquí se puede agregar la lógica para ejecutar el comando mkfs con los parámetros proporcionados
	err := CommandLogin(cmd)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("LOGIN: Iniciando sesion\n"+
		"-> Usuario: %s\n"+
		"-> Contraseña: %s\n"+
		"-> ID: %s",
		cmd.User, cmd.Pass, cmd.ID), nil
	
}

func CommandLogin(login *LOGIN) error {
	// Verificar si ya hay una sesión activa
    if stores.Auth.IsAuthenticated() {
        return fmt.Errorf("ya hay una sesión iniciada con el usuario: %s", stores.Auth.Username)
    }
	// Obtener la partición montada
	partitionSuperblock, _, partitionPath, err := stores.GetMountedPartitionSuperblock(login.ID)
	if err != nil {
		return fmt.Errorf("error al obtener la partición montada: %w", err)
	}

	// Obtener el bloque de usuarios
	usersBlock, err := partitionSuperblock.GetUsersBlock(partitionPath)
	if err != nil {
		return fmt.Errorf("error al obtener el bloque de usuarios: %w", err)
	}

	fmt.Println(usersBlock)

	// Convertir el contenido del bloque a string y separar por líneas
	content := strings.Trim(string(usersBlock.B_content[:]), "\x00")
	lines := strings.Split(content, "\n")

	fmt.Println(content)

	// Variables para almacenar la información del usuario
	var foundUser bool
	var userPassword string

	// Buscar el usuario en las líneas
	for _, line := range lines {
		fmt.Println("Línea del bloque de usuarios:", line)
		// Dividir la línea en campos
		fields := strings.Split(line, ",")
		// Limpiar espacios en blanco de cada campo
		for i := range fields {
			fields[i] = strings.TrimSpace(fields[i])
		}

		// Verificar si es una línea de usuario (tipo U)
		if len(fields) == 5 && fields[1] == "U" {
			// Comparar el nombre de usuario (campo 3)
			if strings.EqualFold(fields[3], login.User) {
				foundUser = true
				userPassword = fields[4]
				break
			}
		}
	}

	// Verificar si se encontró el usuario
	if !foundUser {
		
		return fmt.Errorf("el usuario %s no existe", login.User)
	}

	// Verificar la contraseña
	if !strings.EqualFold(userPassword, login.Pass) {
		return fmt.Errorf("la contraseña no coincide")
	}
	fmt.Println("ID proporcionado:", login.ID)

	// If validation succeeds, set the auth state
	stores.Auth.Login(login.User, login.Pass, login.ID)

	return nil
}

func CommandLogout() (string, error) {
    // Verificar si ya hay una sesión activa
    if !stores.Auth.IsAuthenticated() {
        return "", fmt.Errorf("no hay ninguna sesión activa")
    }

    // Cerrar la sesión
    stores.Auth.Logout()
    // Mensaje formateado
	return fmt.Sprintf("LOGOUT: Sesión cerrada exitosamente\n"+
		"-> Mensaje: La sesión ha sido cerrada correctamente"), nil
}