package group

import (
	"go.mau.fi/whatsmeow"

	"wabot/plugin"
)

func init() {
	add(plugin.Command{
		Name: "demote", Desc: "Cabut admin", Usage: "@user",
		Access: plugin.GroupAdmin, BotAdmin: true,
		Run: member(whatsmeow.ParticipantChangeDemote, "didemote"),
	})
}
