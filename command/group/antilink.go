package group

import (
	"strings"

	libgroup "wabot/library/group"
	"wabot/plugin"
)

func init() {
	add(plugin.Command{
		Name: "antilink", Desc: "Hapus link grup dari non-admin", Usage: "on|off", Access: plugin.GroupAdmin,
		Run: func(c *plugin.Ctx) error {
			on := strings.EqualFold(c.Arg, "on")
			libgroup.SetAntilink(c.Chat(), on)
			if on {
				return c.Reply("🔗 Antilink ON")
			}
			return c.Reply("🔗 Antilink OFF")
		},
	})
}
