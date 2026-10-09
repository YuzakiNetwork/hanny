package tools

import (
	libtools "wabot/library/tools"
	"wabot/plugin"
)

func init() {
	add(plugin.Command{Name: "time", Desc: "Jam sekarang", Run: func(c *plugin.Ctx) error {
		return c.Reply("🕒 " + libtools.Clock())
	}})
}
