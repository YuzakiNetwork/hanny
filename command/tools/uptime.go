package tools

import (
	libtools "wabot/library/tools"
	"wabot/plugin"
)

func init() {
	add(plugin.Command{Name: "uptime", Desc: "Lama bot nyala", Run: func(c *plugin.Ctx) error {
		return c.Reply("⏱️ " + libtools.Uptime())
	}})
}
