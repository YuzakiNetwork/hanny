package group

import "wabot/plugin"

func init() {
	add(plugin.Command{
		Name: "tagall", Desc: "Mention semua member", Usage: "[teks]", Access: plugin.GroupAdmin,
		Run: func(c *plugin.Ctx) error { return mention(c, "📢 "+c.Arg, true) },
	})
}
