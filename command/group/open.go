package group

import "wabot/plugin"

func init() {
	add(plugin.Command{
		Name: "open", Desc: "Semua member boleh chat", Access: plugin.GroupAdmin, BotAdmin: true,
		Run: func(c *plugin.Ctx) error {
			if err := c.Client.SetGroupAnnounce(c.Context, c.Chat(), false); err != nil {
				return err
			}
			return c.React("🔓")
		},
	})
}
