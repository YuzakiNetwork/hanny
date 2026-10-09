package group

import (
	"go.mau.fi/whatsmeow"

	"wabot/plugin"
)

func init() {
	add(plugin.Command{
		Name: "add", Desc: "Tambah member", Usage: "nomor",
		Access: plugin.GroupAdmin, BotAdmin: true,
		Run: member(whatsmeow.ParticipantChangeAdd, "ditambah"),
	})
}
