package internal

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thoas/go-funk"
)

type client struct {
	LastSeen time.Time
}

var clients = make(map[*websocket.Conn]client)
var register = make(chan *websocket.Conn)
var broadcast = make(chan string)
var unregister = make(chan *websocket.Conn)
var outs = []WsMsg{}
var errs = []WsMsg{}

func runHub() {
	for {
		select {
		case connection := <-register:
			clients[connection] = client{LastSeen: time.Now()}
			log.Debug(
				"connection registered",
				"addr",
				connection.RemoteAddr().String(),
				"clients",
				len(clients),
			)

		case message := <-broadcast:

			for connection := range clients {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Error("write error:", err)

					unregister <- connection
					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()

				} else {
					log.Debug("message sent", "addr", connection.RemoteAddr().String(), "msg", message)
				}
			}

		case connection := <-unregister:

			delete(clients, connection)
			log.Debug(
				"connection unregistered",
				"addr",
				connection.RemoteAddr().String(),
				"clients",
				len(clients),
			)

		}
	}
}

func cleanClients() {
	for {
		time.Sleep(1 * time.Second)
		for connection, client := range clients {
			if time.Since(client.LastSeen) > 2*time.Second {
				unregister <- connection

			}
		}
	}
}

func handlePing(c *websocket.Conn) {
	for {
		time.Sleep(1 * time.Second)
		_, _, err := c.ReadMessage()
		if err == nil {
			go func() {
				for connection, client := range clients {

					if connection.RemoteAddr().String() == c.RemoteAddr().String() {
						c := client
						c.LastSeen = time.Now()
						clients[connection] = c

						break
					}
				}
			}()

		}
	}

}

func handleArray(arr []WsMsg, m WsMsg) []WsMsg {
	if m.Id != "" {
		if funk.Find(
			arr,
			func(arrm WsMsg) bool { return arrm.Id == m.Id },
		) == nil {
			arr = append(arr, m)
		} else {
			for i, arrm := range arr {
				if arrm.Id == m.Id {
					arr[i] = m
					break
				}
			}
		}
	}

	return arr
}

func doBroadcast(outs []WsMsg, errs []WsMsg) {
	o := funk.Filter(outs, func(o WsMsg) bool { return o.Out != "" && o.Out != "{}" })
	e := funk.Filter(errs, func(o WsMsg) bool { return o.Err != "" && o.Err != "{}" })
	arr := append(o.([]WsMsg), e.([]WsMsg)...)
	if j, err := json.Marshal(arr); err == nil {
		broadcast <- string(j)

	} else {
		log.Error("socket: marshal", "err", err)
	}

}

func handleChannels(
	outputChan *chan ChanMsg,
	errorChan *chan ChanMsg,

) {
	for {
		select {
		case o := <-*outputChan:
			m := WsMsg{Id: o.Id, Out: o.Msg, Err: "", Time: o.Time}
			outs = handleArray(outs, m)
			doBroadcast(outs, errs)
			break
		case e := <-*errorChan:
			m := WsMsg{Id: e.Id, Out: "", Err: e.Msg, Time: e.Time}
			log.Warn(m)
			errs = handleArray(errs, m)
			doBroadcast(outs, errs)

			break
		}

	}
}

func RunServer(
	scheduler *Scheduler,
	restic *Restic,
	settings *Settings,

	outputChan *chan ChanMsg,
	errorChan *chan ChanMsg,
	version string,
	build string,
) {

	server := fiber.New()
	server.Use(cors.New())
	server.Static("/", "./public")

	cfg := websocket.Config{
		RecoverHandler: func(conn *websocket.Conn) {
			if err := recover(); err != nil {
				conn.WriteJSON(fiber.Map{"customError": "error occurred"})
			}
		},
	}

	api := server.Group("/api")

	api.Use("/ws", func(c *fiber.Ctx) error {

		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	go runHub()
	go cleanClients()
	go handleChannels(outputChan, errorChan)

	api.Get("/ws", websocket.New(func(c *websocket.Conn) {

		defer func() {
			unregister <- c
			c.Close()
		}()

		register <- c

		handlePing(c)

	}, cfg))

	api.Get("/path/autocomplete", func(c *fiber.Ctx) error {
		paths := []string{}
		path := c.Query("path")
		if files, err := os.ReadDir(path); err == nil {
			for _, f := range files {
				if f.IsDir() {
					paths = append(paths, f.Name())
				}
			}
		} else {
			log.Error("reading path", "path", path, "err", err.Error())
		}
		return c.JSON(paths)
	})

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

	api.Get("/version", func(c *fiber.Ctx) error {
		log.Debug(version, build)
		return c.JSON(fiber.Map{"version": version, "build": build})
	})

	api.Get("/logs", func(c *fiber.Ctx) error {
		logs, erros := GetLogFiles()
		return c.JSON(fiber.Map{"logs": logs, "errors": erros})
	})

	api.Get("/logs/:file", func(c *fiber.Ctx) error {
		log, err := GetLogFileContent(c.Params("file"))
		if err != nil {
			c.SendStatus(500)
			return c.SendString(err.Error())
		}
		return c.SendString(string(log))
	})

	api.Post("/check", func(c *fiber.Ctx) error {
		var r Repository
		if err := c.BodyParser(&r); err != nil {
			c.SendStatus(500)
			return c.SendString(err.Error())
		}

		if r.PasswordFile != "" {
			_, err := os.Stat(r.PasswordFile)
			if os.IsNotExist(err) {
				c.SendStatus(500)
				return c.SendString(err.Error())
			}

			data, err := os.ReadFile(r.PasswordFile)
			if err != nil {
				c.SendStatus(500)
				return c.SendString(err.Error())
			}
			if len(data) == 0 {
				c.SendStatus(500)
				return c.SendString("Password file is empty")
			}

		}

		if _, err := restic.Exec(r, []string{"cat", "config"}, []string{}); err != nil {
			if strings.Contains(err.Error(), "key does not exist") ||
				strings.Contains(err.Error(), "/config: no such file") {
				return c.SendString("OK_REPO_EMPTY")
			}
			c.SendStatus(500)
			return c.SendString(err.Error())
		} else {
			return c.SendString("OK_REPO_EXISTING")
		}

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

	config := api.Group("/config")
	backups := api.Group("/backups")
	config.Get("/", func(c *fiber.Ctx) error {
		settings.Refresh()
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
					*settings.Config.GetRepositoryById(id),
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
				*settings.Config.GetRepositoryById(c.Params("id")),
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
				*settings.Config.GetRepositoryById(c.Params("id")),
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
					*settings.Config.GetRepositoryById(c.Params("id")),
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

	server.Listen("0.0.0.0:11278")
}
