package discord

//go:generate go run gen/events_gen.go -o events.go

import (
	"github.com/Sirupsen/logrus"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/dutil/dstate"
	"github.com/jonas747/giskord/engine"
	"os"
)

type ContextKey int

const (
	ContextKeyDiscordSession ContextKey = iota
)

type Plugin struct {
	logger *logrus.Entry

	State        *dstate.State
	Session      *discordgo.Session
	StateHandler *engine.HandlerFunc
}

func RegisterPlugin() {
	p := &Plugin{}
	engine.RegisterPlugin(p)
}

func (p *Plugin) Name() string {
	return "core.discord"
}

func (p *Plugin) Run(logger *logrus.Entry) {
	p.logger = logger

	state := dstate.NewState()
	state.MaxChannelMessages = 100
	p.State = state

	session, err := discordgo.New(os.Getenv("DG_TOKEN"))
	if err != nil {
		logger.WithError(err).Error("Failed to initialise discord session")
	}

	// The discord event converter
	session.AddHandler(HandleEvent)
	session.LogLevel = discordgo.LogInformational
	session.StateEnabled = false
	err = session.Open()
	if err != nil {
		logger.WithError(err).Error("Failed to open gateway connection")
	}
	p.Session = session

	p.StateHandler = engine.AddHandler(p.stateHandler, AllDiscordEvents...)
	engine.AddHandler(p.ready, EventReady)
}

func (p *Plugin) stateHandler(evt *engine.EventData) {
	p.State.HandleEvent(p.Session, evt.Evt)
}

func (p *Plugin) ready(evt *engine.EventData) {
	p.logger.Info("Ready!")
}
