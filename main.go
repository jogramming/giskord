package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/jonas747/giskord/discord"
	"github.com/jonas747/giskord/engine"
	"github.com/jonas747/giskord/ui"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		logrus.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	discord.RegisterPlugin()
	ui.RegisterPlugin()

	logrus.Println("Running...")
	msg := <-engine.Run()
	logrus.Println("Stopped: ", msg)
}
