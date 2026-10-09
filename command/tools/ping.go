package tools

import (
	"fmt"
	"time"

	"wabot/plugin"
)

func init() {
	add(plugin.Command{Name: "ping", Desc: "Cek latensi bot", Run: func(c *plugin.Ctx) error {
		t := time.Now()
		id, err := c.ReplyID("🏓 ...")
		if err != nil {
			return err
		}
		return c.Edit(id, fmt.Sprintf("🏓 Pong! %d ms", time.Since(t).Milliseconds()))
	}})
}
