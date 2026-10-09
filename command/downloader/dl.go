package downloader

import (
	libdl "wabot/library/downloader"
	"wabot/plugin"
)

func init() {
	add(plugin.Command{
		Name: "dl", Aliases: []string{"video"}, Desc: "Download video (YT/TikTok/IG/FB/X, dll)", Usage: "<url>",
		Run: func(c *plugin.Ctx) error { return libdl.Deliver(c, c.Arg, libdl.Video) },
	})
}
