// GENERATED using events_gen.go

// Custom event handlers that adds a redis connection to the handler
// They will also recover from panics

package discord

import (
	"context"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/giskord/engine"
)

const (
	EventMemberFetched           = "d_MemberFetched"
	EventChannelCreate           = "d_ChannelCreate"
	EventChannelDelete           = "d_ChannelDelete"
	EventChannelPinsUpdate       = "d_ChannelPinsUpdate"
	EventChannelUpdate           = "d_ChannelUpdate"
	EventConnect                 = "d_Connect"
	EventDisconnect              = "d_Disconnect"
	EventGuildBanAdd             = "d_GuildBanAdd"
	EventGuildBanRemove          = "d_GuildBanRemove"
	EventGuildCreate             = "d_GuildCreate"
	EventGuildDelete             = "d_GuildDelete"
	EventGuildEmojisUpdate       = "d_GuildEmojisUpdate"
	EventGuildIntegrationsUpdate = "d_GuildIntegrationsUpdate"
	EventGuildMemberAdd          = "d_GuildMemberAdd"
	EventGuildMemberRemove       = "d_GuildMemberRemove"
	EventGuildMemberUpdate       = "d_GuildMemberUpdate"
	EventGuildMembersChunk       = "d_GuildMembersChunk"
	EventGuildRoleCreate         = "d_GuildRoleCreate"
	EventGuildRoleDelete         = "d_GuildRoleDelete"
	EventGuildRoleUpdate         = "d_GuildRoleUpdate"
	EventGuildUpdate             = "d_GuildUpdate"
	EventMessageAck              = "d_MessageAck"
	EventMessageCreate           = "d_MessageCreate"
	EventMessageDelete           = "d_MessageDelete"
	EventMessageDeleteBulk       = "d_MessageDeleteBulk"
	EventMessageReactionAdd      = "d_MessageReactionAdd"
	EventMessageReactionRemove   = "d_MessageReactionRemove"
	EventMessageUpdate           = "d_MessageUpdate"
	EventPresenceUpdate          = "d_PresenceUpdate"
	EventPresencesReplace        = "d_PresencesReplace"
	EventRateLimit               = "d_RateLimit"
	EventReady                   = "d_Ready"
	EventRelationshipAdd         = "d_RelationshipAdd"
	EventRelationshipRemove      = "d_RelationshipRemove"
	EventResumed                 = "d_Resumed"
	EventTypingStart             = "d_TypingStart"
	EventUserGuildSettingsUpdate = "d_UserGuildSettingsUpdate"
	EventUserSettingsUpdate      = "d_UserSettingsUpdate"
	EventUserUpdate              = "d_UserUpdate"
	EventVoiceServerUpdate       = "d_VoiceServerUpdate"
	EventVoiceStateUpdate        = "d_VoiceStateUpdate"
)

var AllDiscordEvents = []string{
	EventMemberFetched,
	EventChannelCreate,
	EventChannelDelete,
	EventChannelPinsUpdate,
	EventChannelUpdate,
	EventConnect,
	EventDisconnect,
	EventGuildBanAdd,
	EventGuildBanRemove,
	EventGuildCreate,
	EventGuildDelete,
	EventGuildEmojisUpdate,
	EventGuildIntegrationsUpdate,
	EventGuildMemberAdd,
	EventGuildMemberRemove,
	EventGuildMemberUpdate,
	EventGuildMembersChunk,
	EventGuildRoleCreate,
	EventGuildRoleDelete,
	EventGuildRoleUpdate,
	EventGuildUpdate,
	EventMessageAck,
	EventMessageCreate,
	EventMessageDelete,
	EventMessageDeleteBulk,
	EventMessageReactionAdd,
	EventMessageReactionRemove,
	EventMessageUpdate,
	EventPresenceUpdate,
	EventPresencesReplace,
	EventRateLimit,
	EventReady,
	EventRelationshipAdd,
	EventRelationshipRemove,
	EventResumed,
	EventTypingStart,
	EventUserGuildSettingsUpdate,
	EventUserSettingsUpdate,
	EventUserUpdate,
	EventVoiceServerUpdate,
	EventVoiceStateUpdate,
}

func HandleEvent(s *discordgo.Session, evt interface{}) {

	name := ""

	switch evt.(type) {
	case *discordgo.ChannelCreate:
		name = "d_ChannelCreate"
	case *discordgo.ChannelDelete:
		name = "d_ChannelDelete"
	case *discordgo.ChannelPinsUpdate:
		name = "d_ChannelPinsUpdate"
	case *discordgo.ChannelUpdate:
		name = "d_ChannelUpdate"
	case *discordgo.Connect:
		name = "d_Connect"
	case *discordgo.Disconnect:
		name = "d_Disconnect"
	case *discordgo.GuildBanAdd:
		name = "d_GuildBanAdd"
	case *discordgo.GuildBanRemove:
		name = "d_GuildBanRemove"
	case *discordgo.GuildCreate:
		name = "d_GuildCreate"
	case *discordgo.GuildDelete:
		name = "d_GuildDelete"
	case *discordgo.GuildEmojisUpdate:
		name = "d_GuildEmojisUpdate"
	case *discordgo.GuildIntegrationsUpdate:
		name = "d_GuildIntegrationsUpdate"
	case *discordgo.GuildMemberAdd:
		name = "d_GuildMemberAdd"
	case *discordgo.GuildMemberRemove:
		name = "d_GuildMemberRemove"
	case *discordgo.GuildMemberUpdate:
		name = "d_GuildMemberUpdate"
	case *discordgo.GuildMembersChunk:
		name = "d_GuildMembersChunk"
	case *discordgo.GuildRoleCreate:
		name = "d_GuildRoleCreate"
	case *discordgo.GuildRoleDelete:
		name = "d_GuildRoleDelete"
	case *discordgo.GuildRoleUpdate:
		name = "d_GuildRoleUpdate"
	case *discordgo.GuildUpdate:
		name = "d_GuildUpdate"
	case *discordgo.MessageAck:
		name = "d_MessageAck"
	case *discordgo.MessageCreate:
		name = "d_MessageCreate"
	case *discordgo.MessageDelete:
		name = "d_MessageDelete"
	case *discordgo.MessageDeleteBulk:
		name = "d_MessageDeleteBulk"
	case *discordgo.MessageReactionAdd:
		name = "d_MessageReactionAdd"
	case *discordgo.MessageReactionRemove:
		name = "d_MessageReactionRemove"
	case *discordgo.MessageUpdate:
		name = "d_MessageUpdate"
	case *discordgo.PresenceUpdate:
		name = "d_PresenceUpdate"
	case *discordgo.PresencesReplace:
		name = "d_PresencesReplace"
	case *discordgo.RateLimit:
		name = "d_RateLimit"
	case *discordgo.Ready:
		name = "d_Ready"
	case *discordgo.RelationshipAdd:
		name = "d_RelationshipAdd"
	case *discordgo.RelationshipRemove:
		name = "d_RelationshipRemove"
	case *discordgo.Resumed:
		name = "d_Resumed"
	case *discordgo.TypingStart:
		name = "d_TypingStart"
	case *discordgo.UserGuildSettingsUpdate:
		name = "d_UserGuildSettingsUpdate"
	case *discordgo.UserSettingsUpdate:
		name = "d_UserSettingsUpdate"
	case *discordgo.UserUpdate:
		name = "d_UserUpdate"
	case *discordgo.VoiceServerUpdate:
		name = "d_VoiceServerUpdate"
	case *discordgo.VoiceStateUpdate:
		name = "d_VoiceStateUpdate"
	default:
		return
	}

	ctx := context.WithValue(context.Background(), ContextKeyDiscordSession, s)
	data := engine.NewEventData(name, evt, ctx)
	engine.EmitEvent(name, data)
}
