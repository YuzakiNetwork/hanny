package command

import (
	"fmt"
	"strings"

	"wabot/config"
	"wabot/plugin"
)

func init() {
	// Kategori "" = command root (nggak masuk daftar kategori di menu)
	plugin.Register("", plugin.Command{
		Name:    "menu",
		Aliases: []string{"help"},
		Desc:    "Daftar command",
		Run:     menu,
	})
}

func icon(cat string) string {
	switch cat {
	case "tools":
		return "🛠️"
	case "downloader":
		return "📥"
	case "group":
		return "👥"
	}
	return "📁"
}

func menu(c *plugin.Ctx) error {
	cfg := config.Get()
	var b strings.Builder
	fmt.Fprintf(&b, "*🤖 %s*\n", cfg.BotName)

	cur := ""
	for _, e := range plugin.All() {
		if e.Category == "" {
			continue
		}
		if e.Category != cur {
			cur = e.Category
			fmt.Fprintf(&b, "\n*%s %s*\n", icon(cur), strings.ToUpper(cur))
		}
		lock := ""
		if e.Cmd.Access == plugin.GroupAdmin {
			lock = " 👮"
		}
		fmt.Fprintf(&b, "%s%s %s — %s%s\n", cfg.Prefix, e.Cmd.Name, e.Cmd.Usage, e.Cmd.Desc, lock)
	}
	return c.Reply(strings.TrimSpace(b.String()))
}
