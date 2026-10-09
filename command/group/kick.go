package group

import (
	"go.mau.fi/whatsmeow"

	"wabot/plugin"
)

func init() {
	add(plugin.Command{
		Name: "kick", Desc: "Keluarin member", Usage: "@user|nomor",
		Access: plugin.GroupAdmin, BotAdmin: true,
		Run: member(whatsmeow.ParticipantChangeRemove, "di-kick"),
	})
}
