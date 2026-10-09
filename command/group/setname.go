package group

import (
	"errors"

	"wabot/plugin"
)

func init() {
	add(plugin.Command{
		Name: "setname", Desc: "Ganti nama grup", Usage: "<nama>", Access: plugin.GroupAdmin, BotAdmin: true,
		Run: func(c *plugin.Ctx) error {
			if c.Arg == "" {
				return errors.New("Tulis nama barunya")
			}
			if err := c.Client.SetGroupName(c.Context, c.Chat(), c.Arg); err != nil {
				return err
			}
			return c.React("✅")
		},
	})
}
