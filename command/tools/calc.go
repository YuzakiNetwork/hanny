package tools

import (
	libtools "wabot/library/tools"
	"wabot/plugin"
)

func init() {
	add(plugin.Command{
		Name: "calc", Aliases: []string{"hitung"}, Desc: "Kalkulator", Usage: "<ekspresi>",
		Run: func(c *plugin.Ctx) error {
			v, err := libtools.Calc(c.Arg)
			if err != nil {
				return err
			}
			return c.Reply("🧮 " + c.Arg + " = *" + v + "*")
		},
	})
}
