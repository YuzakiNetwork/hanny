package group

import (
	"go.mau.fi/whatsmeow"

	"wabot/plugin"
)

func init() {
	add(plugin.Command{
		Name: "promote", Desc: "Jadiin admin", Usage: "@user",
		Access: plugin.GroupAdmin, BotAdmin: true,
		Run: member(whatsmeow.ParticipantChangePromote, "dipromote"),
	})
}
