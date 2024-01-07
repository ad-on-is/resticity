package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func RunServer(
	scheduler *Scheduler,
	restic *Restic,
	settings *Settings,
	errb *bytes.Buffer,
	outb *bytes.Buffer,
) {
	server := fiber.New()
	server.Use(cors.New())

	api := server.Group("/api")

	api.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	api.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index

		for {
			jobs := scheduler.GetRunningJobs()
			data := make(map[string]any)
			data["jobs"] = jobs
			data["out"] = string(outb.Bytes())
			data["err"] = string(errb.Bytes())
			if d, err := json.Marshal(data); err == nil {
				if err = c.WriteMessage(websocket.TextMessage, d); err != nil {
					log.Println("Error writing to socket:", err)
				}
			} else {
				log.Println("Error marshalling data:", err)
			}
			time.Sleep(1 * time.Second)
		}

	}))

	config := api.Group("/config")
	backups := api.Group("/backups")
	config.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(settings.Config)
	})
	config.Post("/", func(c *fiber.Ctx) error {

		s := new(Config)
		if err := c.BodyParser(s); err != nil {
			c.SendStatus(500)
			return c.SendString(err.Error())
		}
		settings.Save(*s)
		scheduler.RescheduleBackups()
		return c.SendString("OK")
	})

	api.Get("/snapshots/:id", func(c *fiber.Ctx) error {
		s := restic.Snapshots(*settings.GetRepositoryById(c.Params("id")))
		fmt.Println(c.Params("id"))
		fmt.Println(s)
		return c.JSON(s)

	})

	api.Get("/schedules/run/:id", func(c *fiber.Ctx) error {
		scheduler.RunJobByName(c.Params("id"))

		return c.SendString("Running schedule in the background")
	})

	backups.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	server.Listen(":11278")
}
