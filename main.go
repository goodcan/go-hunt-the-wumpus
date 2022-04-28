/*
 * @Time     : 2022/4/5 22:13
 * @Author   : cancan
 * @File     : main.go
 * @Function :
 */

package main

import (
	"flag"
	"go-hunt-the-wumpus/client"
	"go-hunt-the-wumpus/server"
	"go-hunt-the-wumpus/single"
	"go-hunt-the-wumpus/storage"
	"go-hunt-the-wumpus/ui"
)

func main() {
	var roleName string
	var mode string
	var serverAddr string
	var language string
	var storageType string
	var inputType string
	var mapLength int
	var mapWidth int
	flag.StringVar(&roleName, "roleName", "tester", "role name")
	flag.StringVar(&mode, "mode", "single", "chose mode: single, server, client")
	flag.StringVar(&serverAddr, "serverAddr", "127.0.0.1:4399", "server addr")
	flag.StringVar(&language, "language", "cn", "language")
	flag.StringVar(&storageType, "storageType", "text", "storage type")
	flag.StringVar(&inputType, "inputType", "terminal", "input type")
	flag.IntVar(&mapLength, "mapLength", 7, "map length")
	flag.IntVar(&mapWidth, "mapWidth", 7, "map width")
	flag.Parse()

	switch mode {
	case "single":
		g := single.Game{
			UI:      ui.Create(language, inputType),
			Storage: storage.Create(storageType),
		}

		g.Init(roleName, mapLength, mapWidth)
		g.Start()
	case "server":
		s := server.Server{
			Storage: storage.Create(storageType),
		}

		s.Init(serverAddr, mapLength, mapWidth)
		s.Start()
	case "client":
		c := client.Client{
			UI: ui.Create(language, inputType),
		}

		c.Init(serverAddr, roleName)
		c.Start()
	}
}
