package downloader

import (
	libdl "wabot/library/downloader"
	"wabot/plugin"
)

func init() {
	add(plugin.Command{
		Name: "mp3", Aliases: []string{"audio"}, Desc: "Download audio", Usage: "<url>",
		Run: func(c *plugin.Ctx) error { return libdl.Deliver(c, c.Arg, libdl.Audio) },
	})
}
