package group

import "wabot/plugin"

func init() {
	add(plugin.Command{
		Name: "setdesc", Desc: "Ganti deskripsi grup", Usage: "<teks>", Access: plugin.GroupAdmin, BotAdmin: true,
		Run: func(c *plugin.Ctx) error {
			if err := c.Client.SetGroupTopic(c.Context, c.Chat(), "", "", c.Arg); err != nil {
				return err
			}
			return c.React("✅")
		},
	})
}
