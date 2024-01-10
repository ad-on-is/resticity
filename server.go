package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
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

	api.Get("/schedules/:id/run", func(c *fiber.Ctx) error {
		scheduler.RunJobByName(c.Params("id"))

		return c.SendString("Running schedule in the background")
	})

	repositories := api.Group("/repositories")
	api.Post("/check", func(c *fiber.Ctx) error {
		var r Repository
		if err := c.BodyParser(&r); err != nil {
			c.SendStatus(500)
			return c.SendString(err.Error())
		}

		files, err := os.ReadDir(r.Path)
		if err != nil {
			c.SendStatus(500)
			return c.SendString(err.Error())
		}
		if len(files) > 0 {
			fmt.Println("containing files")
			if err := restic.Verify(r); err != nil {
				c.SendStatus(500)
				return c.SendString(err.Error())
			} else {
				return c.SendString("REPO_OK_EXISTING")
			}
		}

		return c.SendString("REPO_OK_EMPTY")
	})
	api.Post("/init", func(c *fiber.Ctx) error {
		var r Repository
		if err := c.BodyParser(&r); err != nil {
			c.SendStatus(500)
			return c.SendString(err.Error())
		}
		if err := restic.Initialize(r); err != nil {
			c.SendStatus(500)
			return c.SendString(err.Error())
		}
		return c.SendString("OK")
	})
	repositories.Get("/:id/snapshots", func(c *fiber.Ctx) error {
		s := restic.ListSnapshots(*settings.GetRepositoryById(c.Params("id")))
		return c.JSON(s)

	})

	type MountData struct {
		Path string `json:"path"`
	}

	repositories.Post("/:id/snapshots/:snapshot_id/:action", func(c *fiber.Ctx) error {
		switch c.Params("action") {
		case "browse":
			var data MountData
			if err := c.BodyParser(&data); err != nil {
				c.SendStatus(500)
				return c.SendString(err.Error())
			} else {
				return c.JSON(restic.BrowseSnapshot(
					*settings.GetRepositoryById(c.Params("id")),
					c.Params("snapshot_id"),
					data.Path,
				))
			}
		}

		return c.SendString(c.Params("action"))
	})

	backups.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	server.Listen(":11278")
}
