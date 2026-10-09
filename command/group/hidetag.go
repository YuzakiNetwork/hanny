package group

import (
	"errors"

	"wabot/plugin"
)

func init() {
	add(plugin.Command{
		Name: "hidetag", Desc: "Mention semua tanpa daftar nama", Usage: "<teks>", Access: plugin.GroupAdmin,
		Run: func(c *plugin.Ctx) error {
			if c.Arg == "" {
				return errors.New("Tulis teksnya")
			}
			return mention(c, c.Arg, false)
		},
	})
}
