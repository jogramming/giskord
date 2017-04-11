package ui

import (
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	// "github.com/gotk3/gotk3/pango"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/dutil/dstate"
	"log"
	"time"
)

type MessageView struct {
	UI *UI

	ScrollWindow *gtk.ScrolledWindow
	View         *gtk.TextView
	ChatInput    *gtk.TextView

	Messages map[string]*DisplayedMessage

	CurrentChannel string

	Formatter MessageFormatter

	fetchingMembers   map[string]bool
	requestingHistory map[string]bool
}

type DisplayedMessage struct {
	// Label *gtk.Label
	MS *dstate.MessageState
}

func (mv *MessageView) setup() {
	mv.Messages = make(map[string]*DisplayedMessage)
	mv.fetchingMembers = make(map[string]bool)
	mv.requestingHistory = make(map[string]bool)

	scrollWindow, _ := gtk.ScrolledWindowNew(nil, nil)
	scrollWindow.SetVExpand(true)
	scrollWindow.SetHExpand(true)

	mainViewContainer, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	mainViewContainer.SetVExpand(true)
	mainViewContainer.SetHExpand(true)

	view, _ := gtk.TextViewNew()
	view.SetEditable(false)
	view.SetVExpand(true)
	mv.View = view

	// messagesContainer, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)

	mainViewContainer.Add(scrollWindow)
	scrollWindow.Add(view)
	// scrollWindow.Add(messagesContainer)

	mv.UI.MainGrid.Add(mainViewContainer)

	mv.ScrollWindow = scrollWindow
	// mv.Container = messagesContainer

	sep, _ := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	mainViewContainer.Add(sep)

	entry, _ := gtk.TextViewNew()
	entry.SetHExpand(true)
	entry.SetBorderWidth(5)
	mainViewContainer.Add(entry)
	mv.ChatInput = entry

	mv.Formatter = &StdMessageFormatter{MV: mv}
}

func (mv *MessageView) AddChatMessage(cs *dstate.ChannelState, message *discordgo.Message) {
	if cs == nil {
		cs = discordPlugin.State.Channel(true, message.ChannelID)
	}

	content := mv.Formatter.FormatMessage(cs, message)

	buffer, _ := mv.View.GetBuffer()
	buffer.InsertMarkup(buffer.GetEndIter(), content)
	// l, _ := gtk.LabelNew("")
	// l.SetMarkup(content)
	// l.SetLineWrap(true)
	// l.SetLineWrapMode(pango.WRAP_WORD_CHAR)
	// l.SetHAlign(gtk.ALIGN_START)
	// l.SetJustify(gtk.JUSTIFY_LEFT)
	// l.SetSelectable(true)
	// l.SetTooltipText(fmt.Sprintf("MessageID: %s\nAuthor: %s#%s (%s)", message.ID, message.Author.Username, message.Author.Discriminator, message.Author.ID))
	// l.Show()

	// mv.Container.PackStart(l, false, false, 0)

	dm := &DisplayedMessage{
		// Label: l,
		MS: &dstate.MessageState{
			Message: message,
		},
	}
	dm.MS.ParseTimes()

	mv.Messages[message.ID] = dm

	// We have to wait for it to resize...
	time.AfterFunc(time.Millisecond*100, func() {
		glib.IdleAdd(func() {
			adj := mv.ScrollWindow.GetVAdjustment()
			adj.SetValue(adj.GetUpper())
		})
	})
}

func (mv *MessageView) UpdateMessage(cs *dstate.ChannelState, message *discordgo.Message) {
	_, ok := mv.Messages[message.ID]
	if !ok {
		return
	}

	mv.RecreateMessages()

	// dm.MS.Update(message)

	// if cs == nil {
	// 	cs = discordPlugin.State.Channel(true, message.ChannelID)
	// }

	// formatted := mv.Formatter.FormatMessage(cs, message)
	// dm.Label.SetMarkup(formatted)
}

// Recreates the messages, replacing all current ones with new labels
func (mv *MessageView) RecreateMessages() {
	mv.Clear()

	cs := discordPlugin.State.Channel(true, mv.CurrentChannel)
	if cs == nil {
		return
	}

	cs.Owner.RLock()
	defer cs.Owner.RUnlock()

	currentBuf, _ := mv.View.GetBuffer()
	currentBuf.Delete(currentBuf.GetStartIter(), currentBuf.GetEndIter())

	for _, v := range cs.Messages {
		mv.AddChatMessage(cs, v.Message)
	}
}

func (mv *MessageView) Clear() {
	currentBuf, _ := mv.View.GetBuffer()
	currentBuf.Delete(currentBuf.GetStartIter(), currentBuf.GetEndIter())

	// mv.Container.GetChildren().Foreach(func(item interface{}) {
	// 	cast, ok := item.(*gtk.Widget)
	// 	if !ok {
	// 		return
	// 	}
	// 	mv.Container.Remove(cast)
	// 	// mv.messagesContainer.Remove(cast)
	// })

	mv.Messages = make(map[string]*DisplayedMessage)
}

func (mv *MessageView) FindUserColor(lockGS bool, gs *dstate.GuildState, channelID string, userID string) (color int) {
	if gs == nil {
		return 0
	}

	if lockGS {
		gs.RLock()
		defer gs.RUnlock()
	}

	ms := gs.Member(false, userID)
	if ms == nil || (ms.Member == nil && ms.Presence == nil) {
		if _, ok := mv.fetchingMembers[gs.ID()+":"+userID]; !ok {
			mv.fetchingMembers[gs.ID()+":"+userID] = true
			go mv.FetchMember(gs.ID(), userID)
		}
		return
	}

	var roleIds []string
	if ms.Member != nil {
		roleIds = ms.Member.Roles
	} else {
		if ms.Presence == nil {
			return
		}
		if ms.Presence.Roles != nil {
			roleIds = ms.Presence.Roles
		}
	}

	highest := -1
	for _, v := range roleIds {
		role := gs.Role(false, v)
		if role != nil && role.Color != 0 && role.Position > highest {
			color = role.Color
			highest = role.Position
		}
	}

	return
}

func (mv *MessageView) FetchMember(guildID, userID string) {
	member, err := discordPlugin.Session.GuildMember(guildID, userID)
	defer func() {
		glib.IdleAdd(func() {
			mv.UI.Lock()
			defer mv.UI.Unlock()

			delete(mv.fetchingMembers, guildID+":"+userID)

			if member == nil {
				return
			}

			gs := discordPlugin.State.Guild(true, guildID)
			gs.MemberAddUpdate(true, member)

			cs := gs.Channel(true, mv.CurrentChannel)

			for _, dm := range mv.Messages {
				if dm.MS.Message.Author.ID == userID {
					mv.UpdateMessage(cs, dm.MS.Message)
				}
			}

		})
	}()

	if err != nil {
		log.Println("Failed fetching guildmember:", err)
	}
}

func (mv *MessageView) SelectChannel(id string) {
	mv.CurrentChannel = id
	mv.RecreateMessages()

	cs := discordPlugin.State.Channel(true, id)
	cs.Owner.RLock()
	defer cs.Owner.RUnlock()
	if _, ok := mv.requestingHistory[id]; !ok && len(cs.Messages) < 100 {
		before := ""
		if len(cs.Messages) > 0 {
			before = cs.Messages[0].Message.ID
		}
		go mv.requestHistory(id, before)
	}
}

func (mv *MessageView) requestHistory(cID, before string) {
	history, err := discordPlugin.Session.ChannelMessages(cID, 100, before, "", "")
	if err != nil {
		log.Println("Failed retrieving history")
		return
	}

	cs := discordPlugin.State.Channel(true, cID)
	if cs == nil {
		return
	}

	cs.Owner.Lock()
	defer cs.Owner.Unlock()
	for i := len(history) - 1; i >= 0; i-- {
		cs.MessageAddUpdate(false, history[i], 100, 0)

	}

	glib.IdleAdd(func() {
		mv.UI.Lock()
		defer mv.UI.Unlock()

		// prevent from requesting more history if we have all
		if len(history) >= 100 {
			delete(mv.requestingHistory, cID)
		}

		if mv.CurrentChannel == cID {
			mv.RecreateMessages()
		}
	})
}

type MessageFormatter interface {
	FormatMessage(cs *dstate.ChannelState, msg *discordgo.Message) string
}

type StdMessageFormatter struct {
	MV *MessageView
}

func (f *StdMessageFormatter) FormatMessage(cs *dstate.ChannelState, msg *discordgo.Message) string {
	authorColor := 0
	if !cs.IsPrivate() {
		authorColor = f.MV.FindUserColor(true, cs.Guild, msg.ChannelID, msg.Author.ID)
	}

	authorEscaped := EscapeMarkup(msg.Author.Username)
	messageEscaped := EscapeMarkup(msg.ContentWithMentionsReplaced())

	parsedTime, _ := msg.Timestamp.Parse()
	tStr := ""

	localNow := time.Now().Local()
	localParsed := parsedTime.Local()
	if localNow.Day() == localParsed.Day() {
		tStr = parsedTime.Local().Format(time.Kitchen)
	} else if localNow.Year() == localParsed.Year() {
		tStr = parsedTime.Local().Format(time.Stamp)
	} else {
		tStr = parsedTime.Local().Format(time.RFC822)
	}

	content := fmt.Sprintf("<span foreground=\"#333333\">[%s]:</span> <span foreground=\"#%06X\"><b><u>%s:</u></b></span> %s\n", tStr, authorColor, authorEscaped, messageEscaped)
	return content
}

func EscapeMarkup(in string) (str string) {
	for _, v := range in {
		switch v {
		case '&':
			str += "&amp;"
		case '<':
			str += "&lt;"
		case '>':
			str += "&gt;"
		case '\'':
			str += "&apos;"
		case '"':
			str += "&quot;"
		default:
			if (0x1 <= v && v <= 0x8) ||
				(0xb <= v && v <= 0xc) ||
				(0xe <= v && v <= 0x1f) ||
				(0x7f <= v && v <= 0x84) ||
				(0x86 <= v && v <= 0x9f) {
				str += fmt.Sprintf("&#x%x;", v)
			} else {
				str += string(v)
			}
		}
	}

	return
}
