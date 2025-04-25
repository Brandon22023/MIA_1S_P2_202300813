package main

import (
	analyzer "terminal/analyzer"
	commands "terminal/commands"
	"fmt"
	"strings"

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
	app := fiber.New()

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

		commands := strings.Split(req.Command, "\n")
		output := ""
        // Imprime el comando recibido
    	fmt.Println("Comando recibido:", req.Command)
		for _, cmd := range commands {
			if strings.TrimSpace(cmd) == "" {
				continue
			}

			result, err := analyzer.Analyzer(cmd)
			if err != nil {
				output += fmt.Sprintf("Error: %s\n", err.Error())
			} else {
				output += fmt.Sprintf("%s\n", result)
			}
		}

		if output == "" {
			output = "No se ejecutó ningún comando"
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
	

	app.Listen(":3000")
}

