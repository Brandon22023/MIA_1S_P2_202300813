package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	analyzer "terminal/analyzer"
	commands "terminal/commands"
	"terminal/global"
	stores "terminal/stores"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Output string `json:"output"`
}




func main() {
	fmt.Println("Servidor Fiber iniciado en http://localhost:3000")
	var paths []string // Lista para almacenar los paths de mkdisk
	
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"time":    time.Now(),
		})
	})

	app.Use(cors.New(cors.Config{}))


	app.Post("/analyze", func(c *fiber.Ctx) error {
		var req CommandRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(CommandResponse{
				Output: "Error: Petición inválida",
			})
		}
		// Imprime el comando recibido
		fmt.Println("Comando recibido:", req.Command)

		commandsList := strings.Split(req.Command, "\n")
		output := ""
		
        // Imprime el comando recibido
    	fmt.Println("Comando recibido:", req.Command)
		for _, cmd := range commandsList {
			if strings.TrimSpace(cmd) == "" {
				continue
			}

			result, err := analyzer.Analyzer(cmd)
			if err != nil {
				output += fmt.Sprintf("Error: %s\n", err.Error())
			} else {
				output += fmt.Sprintf("%s\n", result)


				// Si el comando es mkdir, captura el path
				if strings.HasPrefix(strings.ToLower(cmd), "mkdisk") {
					result, err := commands.ParseMkdisk(strings.Fields(cmd))
					if err == nil {
						// Extraer el path del mensaje devuelto
						lines := strings.Split(result, "\n") // Dividir el mensaje en líneas
						for _, line := range lines {
							if strings.HasPrefix(line, "-> Path:") {
								path := strings.TrimSpace(strings.TrimPrefix(line, "-> Path:"))
								paths = append(paths, path) // Agregar el path a la lista
								break
							}
						}
					}
				}
			}
		}

		if output == "" {
			output = "No se ejecutó ningún comando"
		}
		
		fmt.Println("aqui empezara la salida")
        fmt.Println("---------------------------------")

        // Ruta de la carpeta "info_disk"
		currentDir, err := os.Getwd()
		if err != nil {
			return c.Status(500).JSON(CommandResponse{
				Output: "Error al obtener el directorio actual",
			})
		}
		infoDiskDir := filepath.Join(currentDir, "info_disk")

		// Crear la carpeta si no existe
		if _, err := os.Stat(infoDiskDir); os.IsNotExist(err) {
			err = os.MkdirAll(infoDiskDir, os.ModePerm)
			if err != nil {
				return c.Status(500).JSON(CommandResponse{
					Output: fmt.Sprintf("Error al crear la carpeta info_disk: %s", err.Error()),
				})
			}
		}

		// Exportar información de cada path
		for _, diskPath := range paths {
			err := commands.ExportDiskInfo(diskPath)
			if err != nil {
				output += fmt.Sprintf("\nError al exportar información del disco en %s: %s", diskPath, err.Error())
			} else {
				output += fmt.Sprintf("\nInformación del disco exportada exitosamente para el disco en %s", diskPath)
			}
		}
		fmt.Println("aqui se extrae el txt")

		// Obtener el id de la partición montada activa
		partitionID, err := stores.GetActivePartitionID()
		fmt.Println("ID de la partición activa:", partitionID)
		if err == nil {
			fmt.Println("Obtuve partitionID:", partitionID)
			sb, _, partitionPath, err := stores.GetMountedPartitionSuperblock(partitionID)
			if err == nil {
				fmt.Println("Obtuve superbloque y path:", partitionPath)
				sb.ExtractTxtFiles(partitionPath, partitionID)
				// <-- Reconstruir la lista global

			} else {
				fmt.Println("Error al obtener superbloque:", err)
			}
		} else {
			fmt.Println("Error al obtener partitionID:", err)
		}

		fmt.Println("aqui se termina la extraccion")
        
		fmt.Println("Paths válidos de mkfile guardados:")
		for _, path := range global.GetValidFilePathsMkfile() {
			fmt.Println(path)
		}
		

		return c.JSON(CommandResponse{
			Output: output,
		})
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		var req struct {
			User string `json:"user"`
			Pass string `json:"pass"`
			ID   string `json:"id"`
		}
	
		// Parsear el cuerpo de la solicitud
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Solicitud inválida",
			})
		}
		// Imprimir los datos recibidos en el terminal
		fmt.Printf("Datos recibidos en el backend: %+v\n", req)
	
		// Crear una instancia de LOGIN con los datos recibidos
		login := commands.LOGIN{
			User: req.User,
			Pass: req.Pass,
			ID:   req.ID,
		}
	
		// Ejecutar la lógica de login
		if err := commands.CommandLogin(&login); err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		
	
		// Respuesta exitosa
		return c.JSON(fiber.Map{
			"message": "Inicio de sesión exitoso",
		})
	})

	app.Get("/disks", func(c *fiber.Ctx) error {
		// Ruta de la carpeta "info_disk"
		currentDir, err := os.Getwd()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "No se pudo obtener el directorio actual",
			})
		}
		infoDiskDir := filepath.Join(currentDir, "info_disk")
	
		// Verificar si la carpeta existe
		if _, err := os.Stat(infoDiskDir); os.IsNotExist(err) {
			return c.Status(404).JSON(fiber.Map{
				"error": "La carpeta info_disk no existe",
			})
		}
	
		// Leer los archivos JSON en la carpeta
		files, err := os.ReadDir(infoDiskDir)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "No se pudo leer la carpeta info_disk",
			})
		}
	
		// Crear una lista con la información de los discos
		var disks []map[string]interface{}
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".json") {
				// Leer el contenido del archivo JSON
				filePath := filepath.Join(infoDiskDir, file.Name())
				data, err := os.ReadFile(filePath)
				if err != nil {
					continue // Ignorar archivos que no se puedan leer
				}
	
				// Parsear el JSON
				var diskInfo map[string]interface{}
				if err := json.Unmarshal(data, &diskInfo); err != nil {
					continue // Ignorar archivos con formato inválido
				}
	
				// Calcular el tamaño en MB y bytes
				sizeBytes, ok := diskInfo["size"].(float64)
				if !ok {
					continue // Ignorar si el tamaño no es válido
				}
				sizeMB := sizeBytes / 1000 / 1000 // Convertir a MB usando base 1000
	
				// Manejar particiones
				partitions, ok := diskInfo["partitions"].([]interface{})
				mountedPartitions := "No existen particiones"
				if ok && len(partitions) > 0 {
					count := 0
					for _, partition := range partitions {
						part, ok := partition.(map[string]interface{})
						if ok && part["status"] == "1" {
							count++
						}
					}
					mountedPartitions = fmt.Sprintf("%d", count)
				}
	
				// Agregar la información del disco a la lista
				disks = append(disks, map[string]interface{}{
					"name":              diskInfo["name"],
					"size":              fmt.Sprintf("%.1f MB (%.0f bytes)", sizeMB, sizeBytes),
					"fit":               diskInfo["fit"],
					"mounted_partitions": mountedPartitions,
				})
			}
		}
	
		return c.JSON(fiber.Map{
			"disks": disks,
		})
	})

	app.Get("/partitions/:diskName", func(c *fiber.Ctx) error {
		diskName := c.Params("diskName") // Obtener el nombre del disco desde la URL
		fmt.Println("Disco solicitado:", diskName) // <-- Agrega este log

		// Eliminar la extensión ".mia" si existe
		diskName = strings.TrimSuffix(diskName, ".mia")
	
		// Ruta del archivo JSON del disco
		currentDir, err := os.Getwd()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "No se pudo obtener el directorio actual",
			})
		}
		diskFilePath := filepath.Join(currentDir, "info_disk", diskName+".json")
	
		// Leer el archivo JSON
		data, err := os.ReadFile(diskFilePath)
		if err != nil {
			fmt.Println("Error al leer el archivo JSON:", err) // <-- Agrega este log
			return c.Status(404).JSON(fiber.Map{
				"error": "No se pudo leer el archivo del disco",
			})
		}
	
		// Parsear el JSON
		var diskInfo map[string]interface{}
		if err := json.Unmarshal(data, &diskInfo); err != nil {
			fmt.Println("Error al parsear el JSON:", err) // <-- Agrega este log
			return c.Status(500).JSON(fiber.Map{
				"error": "Error al parsear el archivo JSON",
			})
		}
	
		// Procesar las particiones
		partitions, ok := diskInfo["partitions"].([]interface{})
		if !ok || len(partitions) == 0 {
			fmt.Println("No existen particiones para el disco:", diskName) // <-- Agrega este log
			return c.JSON(fiber.Map{
				"message": "No existen particiones para dicho disco",
			})
		}
	
		var processedPartitions []map[string]interface{}
		for _, partition := range partitions {
			part, ok := partition.(map[string]interface{})
			if !ok {
				continue
			}
	
			// Procesar el nombre
			name, _ := part["name"].(string)
			name = strings.TrimSpace(strings.ReplaceAll(name, "\u0000", ""))
	
			// Procesar el tamaño
			sizeBytes, _ := part["size"].(float64)
			sizeMB := sizeBytes / 1000 / 1000 // Convertir a MB usando base 1000
	
			// Procesar el tipo
			partType, _ := part["type"].(string)
			var typeDescription string
			if partType == "P" {
				typeDescription = "Primaria"
			} else if partType == "E" {
				typeDescription = "Extendida"
			} else {
				typeDescription = "Desconocido"
			}
	
			// Procesar el estado
			status, _ := part["status"].(string)
			var stateDescription string
			if status == "1" {
				stateDescription = "Montada"
			} else {
				stateDescription = "No montada"
			}
	
			// Procesar el ID
			id, _ := part["id"].(string)
			id = strings.TrimSpace(strings.ReplaceAll(id, "\u0000", ""))
			if id == "" || id == "N" {
				id = "No está montada"
			}
	
			// Procesar el fit
			fit, _ := part["fit"].(string)
	
			// Procesar el inicio
			start, _ := part["start"].(float64)
	
			// Agregar la partición procesada
			processedPartitions = append(processedPartitions, map[string]interface{}{
				"name":  name,
				"size":  fmt.Sprintf("%.1f MB (%.0f bytes)", sizeMB, sizeBytes), // Formato similar al de los discos
				"type":  typeDescription,
				"fit":   fit,
				"start": fmt.Sprintf("%d", int(start)),
				"state": stateDescription,
				"id":    id,
			})
		}
	
		fmt.Println("Particiones procesadas:", processedPartitions) // <-- Agrega este log
	
		return c.JSON(fiber.Map{
			"partitions": processedPartitions,
		})
	})

	app.Post("/logout", func(c *fiber.Ctx) error {
		// Llamar al comando de logout
		message, err := commands.CommandLogout()
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	
		// Respuesta exitosa
		return c.JSON(fiber.Map{
			"message": message,
		})
	})

	// Estructura global para almacenar paths válidos y el ID activo

	app.Get("/folders", func(c *fiber.Ctx) error {
		if global.ActivePartitionID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "No hay una partición activa",
			})
		}
	
		var response []map[string]string
		for _, path := range global.ValidPaths {
			// Validar que el path no esté vacío y no contenga segmentos vacíos
			if strings.TrimSpace(path) != "" && !strings.Contains(path, "//") {
				response = append(response, map[string]string{
					"path": path,
					"id":   global.ActivePartitionID,
				})
			}
		}
	
		// Imprimir lo que se está enviando
		fmt.Println("Datos enviados al frontend:", response)
	
		return c.JSON(fiber.Map{
			"carpetas": response,
		})
	})

	app.Get("/txtfiles", func(c *fiber.Ctx) error {
		if global.ActivePartitionID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "No hay una partición activa",
			})
		}
		partitionID := global.ActivePartitionID
		sb, _, partitionPath, err := stores.GetMountedPartitionSuperblock(partitionID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "No se pudo obtener el superbloque",
			})
		}
		fmt.Println("Paths válidos de mkfile guardados:", global.ValidFilePaths_mkfile)
		txtFiles, err := sb.GetTxtFiles(partitionPath, partitionID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "No se pudieron extraer los archivos .txt",
			})
		}
		// Imprimir lo que se está enviando al frontend
		fmt.Println("Archivos .txt enviados al frontend:", txtFiles)
		return c.JSON(fiber.Map{
			"txtfiles": txtFiles,
		})
	})

	

	app.Listen(":3000")
}




