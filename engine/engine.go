package engine

import (
	"github.com/Sirupsen/logrus"
)

var (
	stopChan = make(chan string)
)

func Run() chan string {
	logrus.Info("Starting up")
	started = true
	for _, v := range Plugins {
		logrus.Info("Starting ", v.Name())
		v.Run(logrus.WithField("p", v.Name()))
	}

	return stopChan
}

func Stop() {
	logrus.Info("Shutting down...")
	for _, v := range Plugins {
		if stopper, ok := v.(PluginStopper); ok {
			logrus.Info("Stopping ", v.Name())
			stopper.Stop()
		}
	}

	stopChan <- "Shut down"
}
