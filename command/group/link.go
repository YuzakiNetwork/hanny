package group

import "wabot/plugin"

func init() {
	add(plugin.Command{
		Name: "link", Desc: "Link invite grup", Access: plugin.GroupAdmin, BotAdmin: true,
		Run: func(c *plugin.Ctx) error {
			link, err := c.Client.GetGroupInviteLink(c.Context, c.Chat(), false)
			if err != nil {
				return err
			}
			return c.Reply(link)
		},
	})
}
