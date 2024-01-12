package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type BrowseData struct {
	Path string `json:"path"`
}

type MountData struct {
	Path string `json:"path"`
}

type RestoreData struct {
	RootPath string `json:"root_path"`
	FromPath string `json:"from_path"`
	ToPath   string `json:"to_path"`
}

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

		defer c.Close()
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

	api.Get("/schedules/:id/run", func(c *fiber.Ctx) error {
		scheduler.RunJobByName(c.Params("id"))

		return c.SendString("Running schedule in the background")
	})

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
			if _, err := restic.Exec(r, []string{"cat", "config"}, []string{}); err != nil {
				c.SendStatus(500)
				return c.SendString(err.Error())
			} else {
				return c.SendString("OK_REPO_EXISTING")
			}
		}

		return c.SendString("OK_REPO_EMPTY")
	})
	api.Post("/init", func(c *fiber.Ctx) error {
		var r Repository
		if err := c.BodyParser(&r); err != nil {
			c.SendStatus(500)
			return c.SendString(err.Error())
		}
		if _, err := restic.Exec(r, []string{"init"}, []string{}); err != nil {
			c.SendStatus(500)
			return c.SendString(err.Error())
		}
		return c.SendString("OK")
	})

	api.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStream(bytes.NewReader(outb.Bytes()))
	})

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

	repositories := api.Group("/repositories")

	repositories.Post("/:id/:action", func(c *fiber.Ctx) error {
		act := c.Params("action")

		switch act {
		case "mount":
			var data MountData
			if err := c.BodyParser(&data); err != nil {
				c.SendStatus(500)
				return c.SendString(err.Error())
			}
			fmt.Println(data.Path)

			go func(id string) {
				restic.Exec(
					*settings.GetRepositoryById(id),
					[]string{act, data.Path},
					[]string{},
				)
			}(c.Params("id"))

			return c.SendString("OK")
		case "unmount":
			var data MountData
			if err := c.BodyParser(&data); err != nil {
				c.SendStatus(500)
				return c.SendString(err.Error())
			}
			fmt.Println("unmounting")

			e := exec.Command("/usr/bin/umount", "-l", data.Path)
			e.Output()

			return c.SendString("OK")
		case "snapshots":
			res, err := restic.Exec(
				*settings.GetRepositoryById(c.Params("id")),
				[]string{act},
				[]string{},
			)
			if err != nil {
				c.SendStatus(500)
				return c.SendString(err.Error())
			}
			var data []Snapshot

			if err := json.Unmarshal([]byte(res), &data); err != nil {
				c.SendStatus(500)
				return c.SendString(err.Error())
			}
			return c.JSON(data)
		}

		return c.SendString("Unknown action")

	})

	repositories.Post("/:id/snapshots/:snapshot_id/:action", func(c *fiber.Ctx) error {
		switch c.Params("action") {
		case "browse":
			var data BrowseData
			if err := c.BodyParser(&data); err != nil {
				c.SendStatus(500)
				return c.SendString(err.Error())
			}
			res, err := restic.BrowseSnapshot(
				*settings.GetRepositoryById(c.Params("id")),
				c.Params("snapshot_id"),
				data.Path,
			)
			if err != nil {
				c.SendStatus(500)
				return c.SendString(err.Error())

			}
			return c.JSON(res)

		case "restore":
			var data RestoreData
			if err := c.BodyParser(&data); err != nil {
				c.SendStatus(500)
				return c.SendString(err.Error())
			} else {
				if _, err := restic.Exec(
					*settings.GetRepositoryById(c.Params("id")),
					[]string{"restore",
						c.Params("snapshot_id") + ":" + data.RootPath,
						"--target",
						data.ToPath,
						"--include", data.FromPath}, []string{},
				); err != nil {
					c.SendStatus(500)
					return c.SendString(err.Error())
				}
				return c.SendString("OK")
			}
		}

		return c.SendString(c.Params("action"))
	})

	backups.Get("/", func(c *fiber.Ctx) error {

		return c.SendString("Hello, World!")
	})

	server.Listen(":11278")
}
