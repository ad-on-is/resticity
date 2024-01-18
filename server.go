package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thoas/go-funk"
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

type Output struct {
	Id  string `json:"id"`
	Out any    `json:"out"`
}

type MsgJob struct {
	Id       string   `json:"id"`
	Schedule Schedule `json:"schedule"`
	Running  bool     `json:"running"`
	Force    bool     `json:"force"`
}

func RunServer(
	scheduler *Scheduler,
	restic *Restic,
	settings *Settings,
	errb *bytes.Buffer,
	outb *bytes.Buffer,
	outputChan *chan ChanMsg,
) {
	server := fiber.New()
	server.Use(cors.New())
	server.Static("/", "./public")

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

		// outs := []Output{}

		// listen := func() {
		// 	for {
		// 		jobs := scheduler.GetRunningJobs()
		// 		for _, j := range jobs {
		// 			if funk.Find(outs, func(out Output) bool { return out.Id == j.Id }) == nil {
		// 				outs = append(outs, Output{Id: j.Id, Out: map[string]string{}})
		// 			}
		// 			go func(j Job) {
		// 				d := <-j.Chan
		// 				for i, out := range outs {
		// 					if out.Id == j.Id {
		// 						var data any
		// 						if err := json.Unmarshal([]byte(d), &data); err == nil {
		// 							outs[i].Out = data
		// 						} else {
		// 							outs[i].Out = map[string]string{}
		// 						}
		// 						break
		// 					}
		// 				}
		// 				data := make(map[string]any)
		// 				data["jobs"] = funk.Map(jobs, func(j Job) MsgJob {
		// 					return MsgJob{
		// 						Id:       j.Id,
		// 						Schedule: j.Schedule,
		// 						Running:  j.Running,
		// 						Force:    j.Force,
		// 					}
		// 				})
		// 				data["outputs"] = outs
		// 				time.Sleep(100 * time.Millisecond)

		// 				if d, err := json.Marshal(data); err == nil {
		// 					if err = c.WriteMessage(websocket.TextMessage, d); err != nil {
		// 						log.Println("Error writing to socket:", err)
		// 					}
		// 				} else {
		// 					log.Println("Error marshalling data:", err)
		// 				}
		// 				time.Sleep(100 * time.Millisecond)
		// 			}(j)

		// 		}

		// 	}

		// }

		outputs := []ChanMsg{}

		for {
			select {
			case d := <-*outputChan:
				if funk.Find(
					outputs,
					func(out ChanMsg) bool { return out.Id == d.Id },
				) == nil {
					outputs = append(outputs, d)
				} else {
					for i, out := range outputs {
						if out.Id == d.Id {
							outputs[i] = d
							break
						}
					}
				}
				// var data map[string]any
				// data["jobs"] = outputs
				if j, err := json.Marshal(funk.Filter(outputs, func(o ChanMsg) bool { return o.Out != "" && o.Out != "{}" })); err == nil {
					if err = c.WriteMessage(websocket.TextMessage, j); err != nil {
						log.Println("Error writing to socket:", err)
					}
				} else {
					log.Println("Error marshalling data:", err)
				}
			}
		}

		// for {
		// time.Sleep(1 * time.Second)
		// fmt.Println(outputs)

		// data["out"] = string(outb.Bytes())

		// data["outputs"] = funk.Map(filtered, func(o ScheduleOutput) *ScheduleOutput {
		// 	u := scheduler.GetOutputById(o.Id)
		// 	// fmt.Println(*u)
		// 	return u
		// })
		// fmt.Println(data["outputs"])
		// data["err"] = string(errb.Bytes())
		// if d, err := json.Marshal(data); err == nil {
		// 	if err = c.WriteMessage(websocket.TextMessage, d); err != nil {
		// 		log.Println("Error writing to socket:", err)
		// 	}
		// } else {
		// 	log.Println("Error marshalling data:", err)
		// }
		// time.Sleep(300 * time.Millisecond)
		// }

	}))

	api.Get("/schedules/:id/:action", func(c *fiber.Ctx) error {
		switch c.Params("action") {
		case "run":
			scheduler.RunJobById(c.Params("id"))
			break
		case "stop":
			scheduler.StopJobById(c.Params("id"))
			break
		}

		return c.SendString(c.Params("action") + " schedule in the background")
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

			e := exec.Command("/usr/bin/umount", "-l", data.Path)
			e.Output()

			return c.SendString("OK")
		case "snapshots":
			groupBy := c.Query("group_by")
			if groupBy == "" {
				groupBy = "host"
			}
			res, err := restic.Exec(
				*settings.GetRepositoryById(c.Params("id")),
				[]string{act, "--group-by", groupBy},
				[]string{},
			)
			if err != nil {
				c.SendStatus(500)
				return c.SendString(err.Error())
			}
			var data []SnapshotGroup

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
