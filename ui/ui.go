package ui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/giskord/discord"
	"github.com/jonas747/giskord/engine"
	"github.com/pkg/errors"
	"log"
	"sort"
	"strings"
	"sync"
)

type UI struct {
	sync.RWMutex

	Window   *gtk.Window
	MainGrid *gtk.Grid

	Guilds   *ListWindow
	Channels *ListWindow

	guildRows   map[int]string
	channelRows map[int]string

	// ChatBox   *gtk.Box
	// ChatInput *gtk.TextView

	// variables for the current state of positioning objects, can be used by plugins
	GridMinX, GridMinY int
	GridMaxX, GridMaxY int

	MessageView *MessageView

	CurrentGuild   string
	guildPositions []string
}

type ListWindow struct {
	List   *gtk.ListBox
	Window *gtk.ScrolledWindow
}

func (ui *UI) Setup() error {
	ui.guildRows = make(map[int]string)
	ui.channelRows = make(map[int]string)

	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	err := ui.createMainWindow()
	if err != nil {
		return err
	}

	engine.AddHandler(ui.GTKThreadEventHandler(ui.LockedEventHandler(ui.handleReady)), discord.EventReady)

	engine.AddHandler(ui.GTKThreadEventHandler(ui.LockedEventHandler(ui.handleGuildEvent)), discord.EventGuildCreate,
		discord.EventGuildDelete, discord.EventGuildUpdate, discord.EventChannelCreate, discord.EventChannelDelete,
		discord.EventChannelUpdate, discord.EventMessageCreate)

	return nil
}

func (ui *UI) createMainWindow() error {
	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return errors.Wrap(err, "Unable to create window")
	}
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.SetTitle("Giscord - GTK Discord client")
	win.Connect("destroy", func() {
		gtk.MainQuit()
		go engine.Stop()
	})
	win.SetEvents(int(gdk.KEY_PRESS_MASK))
	win.Connect("key-press-event", ui.handleKeyPress)

	ui.Window = win

	err = ui.createMainWindowContents()
	if err != nil {
		return err
	}

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	return nil
}

func (ui *UI) handleKeyPress(w *gtk.Window, event *gdk.Event) bool {
	log.Println("HANDLING KEYPRESS")

	eventKey := gdk.EventKeyNewFromEvent(event)
	if !ui.MessageView.ChatInput.HasFocus() {

		if eventKey.KeyVal() == 103 && eventKey.State() == 4 {
			ui.refreshGuildsList()
		}

		log.Println("No chat focus")
		return false
	}

	if eventKey.KeyVal() != 65293 || eventKey.State() != 0 {
		log.Println("Keyval", eventKey.KeyVal(), eventKey.State())
		return false
	}

	buffer, _ := ui.MessageView.ChatInput.GetBuffer()
	text, _ := buffer.GetText(buffer.GetStartIter(), buffer.GetEndIter(), true)
	if strings.TrimSpace(text) == "" {
		log.Println("empty buffer")
		return false // Empty buffer
	}
	buffer.SetText("")

	go func(c string) {
		_, err := discordPlugin.Session.ChannelMessageSend(c, text)
		if err != nil {
			log.Println(err)
		}

	}(ui.MessageView.CurrentChannel)

	return true
}

func (ui *UI) createMainWindowContents() error {
	// Create a new grid widget to arrange child widgets
	grid, err := gtk.GridNew()
	if err != nil {
		return errors.Wrap(err, "Unable to create grid")
	}

	grid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	ui.MainGrid = grid
	ui.Window.Add(grid)

	ui.createServersChannelsLists()

	// ui.MainGrid.Attach(entry, 2, 1, 3, 1)
	ui.MessageView = &MessageView{
		UI: ui,
	}
	ui.MessageView.setup()

	return nil
}

func (ui *UI) createServersChannelsLists() {

	serversList := ui.createListBox("Servers", func(lb *gtk.ListBox, lr *gtk.ListBoxRow) {
		log.Println("Selected server")
		index := lr.GetIndex()
		if index < 0 {
			return
		}
		ui.SelectServer(ui.guildRows[index])
	})

	channelList := ui.createListBox("Channels", func(lb *gtk.ListBox, lr *gtk.ListBoxRow) {
		index := lr.GetIndex()
		if index < 0 {
			return
		}
		log.Println("Selected a channel row", index, ui.channelRows[index])
		ui.SelectChannel(ui.channelRows[index])
	})

	ui.MainGrid.Add(serversList.Window)

	sep1, _ := gtk.SeparatorNew(gtk.ORIENTATION_VERTICAL)
	ui.MainGrid.Add(sep1)

	ui.MainGrid.Add(channelList.Window)

	sep2, _ := gtk.SeparatorNew(gtk.ORIENTATION_VERTICAL)
	ui.MainGrid.Add(sep2)

	ui.Guilds = serversList
	ui.Channels = channelList
}

func (ui *UI) SelectServer(id string) {
	ui.CurrentGuild = id
	ui.refreshChannelsList()
}

func (ui *UI) SelectChannel(id string) {
	ui.MessageView.SelectChannel(id)
}

func (u *UI) createListBox(columnName string, callback func(*gtk.ListBox, *gtk.ListBoxRow)) *ListWindow {
	window, _ := gtk.ScrolledWindowNew(nil, nil)
	window.SetVExpand(true)

	window.SetSizeRequest(150, 0)

	listBox, _ := gtk.ListBoxNew()
	listBox.SetSelectionMode(gtk.SELECTION_BROWSE)

	window.Add(listBox)
	listBox.Connect("row-selected", callback)

	return &ListWindow{
		List:   listBox,
		Window: window,
	}
}

// LockedEventHandler will lock ui before calling the event handler
func (ui *UI) LockedEventHandler(inner engine.HandlerFunc) engine.HandlerFunc {
	return func(evt *engine.EventData) {
		ui.Lock()
		defer ui.Unlock()

		inner(evt)
	}
}

// RLockedEventHandler will rlock ui before calling the event handler
func (ui *UI) RLockedEventHandler(inner engine.HandlerFunc) engine.HandlerFunc {
	return func(evt *engine.EventData) {
		ui.RLock()
		defer ui.RUnlock()

		inner(evt)
	}
}

// GTKThreadEventHandler calls the inner event handler on the gtk thread
func (ui *UI) GTKThreadEventHandler(inner engine.HandlerFunc) engine.HandlerFunc {
	return func(evt *engine.EventData) {
		glib.IdleAdd(inner, evt)
	}
}

func (ui *UI) handleReady(evt *engine.EventData) {
	r := evt.Evt.(*discordgo.Ready)
	ui.guildPositions = r.Settings.GuildPositions

	ui.refreshGuildsList()
}

func (ui *UI) handleGuildEvent(evt *engine.EventData) {
	logger.Info("Handling event", evt.Name)
	switch t := evt.Evt.(type) {
	case *discordgo.GuildCreate:
		if len(ui.guildPositions) < 1 {
			ui.addGuild(t.Guild)
		} else {
			ui.refreshGuildsList()
		}
	case *discordgo.GuildDelete:
		ui.refreshGuildsList()
	case *discordgo.GuildUpdate:
		ui.refreshGuildsList()
	case *discordgo.ChannelCreate:
		if t.Channel.GuildID == ui.CurrentGuild {
			ui.refreshChannelsList()
		}
	case *discordgo.ChannelDelete:
		if t.Channel.GuildID == ui.CurrentGuild {
			ui.refreshChannelsList()
		}
	case *discordgo.ChannelUpdate:
		if t.Channel.GuildID == ui.CurrentGuild {
			ui.refreshChannelsList()
		}
	case *discordgo.MessageCreate:
		if t.ChannelID == ui.MessageView.CurrentChannel {
			ui.MessageView.AddChatMessage(nil, t.Message)
		}
	}
}

// addGuild adds a guild to the displayed guild list
func (ui *UI) refreshGuildsList() {

	// Remove all current
	ui.Guilds.List.GetChildren().Foreach(func(item interface{}) {
		cast, ok := item.(*gtk.Widget)
		if !ok {
			return
		}

		ui.Guilds.List.Remove(cast)
	})
	ui.guildRows = make(map[int]string)

	// First check guild positions
	if len(ui.guildPositions) > 1 {
		for _, gID := range ui.guildPositions {
			ui.addGuild(discordPlugin.State.Guild(true, gID).LightCopy(true))
		}
		return
	}

	// If not fallback to raw state
	discordPlugin.State.RLock()
	guilds := make([]*discordgo.Guild, 0)
	for _, v := range discordPlugin.State.Guilds {
		if v.Guild.Unavailable {
			continue
		}
		guilds = append(guilds, v.LightCopy(true))
	}

	sort.Slice(guilds, func(i, j int) bool {
		return guilds[i].Name < guilds[j].Name
	})

	for _, v := range guilds {
		ui.addGuild(v)
	}

	discordPlugin.State.RUnlock()
}

func (ui *UI) addGuild(guild *discordgo.Guild) {

	item, _ := gtk.ListBoxRowNew()
	l, _ := gtk.LabelNew(guild.Name)
	l.SetJustify(gtk.JUSTIFY_LEFT)
	l.SetHAlign(gtk.ALIGN_START)
	item.Add(l)

	item.ShowAll()

	ui.Guilds.List.Add(item)
	ui.guildRows[item.GetIndex()] = guild.ID
}

func (ui *UI) refreshChannelsList() {
	log.Println("Refreshing channels")
	if ui.CurrentGuild == "" {
		return
	}

	currentChannel := ui.MessageView.CurrentChannel

	// Remove all current
	ui.Channels.List.GetChildren().Foreach(func(item interface{}) {
		cast, ok := item.(*gtk.Widget)
		if !ok {
			return
		}

		ui.Channels.List.Remove(cast)
	})
	ui.channelRows = make(map[int]string)

	gs := discordPlugin.State.Guild(true, ui.CurrentGuild)

	gs.RLock()
	channels := make([]*discordgo.Channel, 0, len(gs.Channels))
	for _, v := range gs.Channels {
		if v.Type() != "text" {
			continue
		}
		channels = append(channels, v.Copy(false, false))
	}
	gs.RUnlock()

	sort.Slice(channels, func(i, j int) bool {
		return channels[i].Position < channels[j].Position
	})

	for _, v := range channels {
		ui.addChannel(v, v.ID == currentChannel)
	}
}

// addChannel adds a chanenl to the displayed channel list
// if the channel does not belong to the current guild/dms then this does nothing
func (ui *UI) addChannel(channel *discordgo.Channel, selected bool) {
	if ui.CurrentGuild == "" || ui.CurrentGuild != channel.GuildID {
		return
	}

	item, _ := gtk.ListBoxRowNew()
	l, _ := gtk.LabelNew("#" + channel.Name)
	l.SetJustify(gtk.JUSTIFY_LEFT)
	l.SetHAlign(gtk.ALIGN_START)
	item.Add(l)

	item.ShowAll()

	ui.Channels.List.Add(item)
	ui.channelRows[item.GetIndex()] = channel.ID
	if selected {
		ui.Channels.List.SelectRow(item)
	}
}

// addChannel removes a channel from the displayed channel list
// if the channel does not belong to the current guild/dms then this does nothing
func (ui *UI) removeChannel(channel *discordgo.Channel) {
	if ui.CurrentGuild == "" || ui.CurrentGuild != channel.GuildID {
		return
	}
}
