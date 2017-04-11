package ui

import (
	"github.com/Sirupsen/logrus"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jonas747/giskord/discord"
	"github.com/jonas747/giskord/engine"
)

type ContextKey int

var (
	discordPlugin *discord.Plugin
	logger        *logrus.Entry
)

const (
	ContextKeyDiscordSession ContextKey = iota
)

type Plugin struct {
	logger *logrus.Entry
	UI     *UI
}

func RegisterPlugin() {
	p := &Plugin{}
	engine.RegisterPlugin(p)
}

func (p *Plugin) Name() string {
	return "core.ui"
}

func (p *Plugin) Run(l *logrus.Entry) {
	logger = l
	discordPlugin = engine.FindPlugin("core.discord").(*discord.Plugin)
	p.UI = &UI{}
	p.UI.Setup()

	go func() {
		gtk.Main()
	}()
}
