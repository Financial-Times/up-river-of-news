package main

import (
	_ "net/http/pprof"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("up-river-of-news", "Listens for native publication events and posts them to a slack channel")
	//port := app.IntOpt("port", 8080, "Port to listen on")
	//env := app.StringOpt("env", "local", "environment this app is running in")

	app.Action = func() {
		readMessages()
	}

	log.SetLevel(log.InfoLevel)
	log.Infof("Application started with args %s", os.Args)
	app.Run(os.Args)
}

func readMessages() {
	log.Infof("About to read messages")
}
