package tools

import (
	"fmt"

	"wabot/plugin"
)

func init() {
	add(plugin.Command{Name: "jid", Desc: "Lihat JID chat & pengirim", Run: func(c *plugin.Ctx) error {
		return c.Reply(fmt.Sprintf("chat: %s\npengirim: %s", c.Chat(), c.Sender()))
	}})
}
