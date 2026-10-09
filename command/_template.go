//go:build ignore

// TEMPLATE command baru. Salin ke command/<kategori>/<nama>.go (file berawalan "_" diabaikan Go).
// Ganti `package` sesuai nama folder kategori. Folder baru = kategori baru (+1 baris import di command.go).
package tools

import "wabot/plugin"

func init() {
	add(plugin.Command{
		Name:  "contoh",
		Desc:  "Deskripsi singkat",
		Usage: "<arg>",
		// Aliases:  []string{"ct"},
		// Access:   plugin.GroupAdmin, // Everyone (default) | GroupOnly | GroupAdmin
		// BotAdmin: true,              // bot harus admin grup
		Run: func(c *plugin.Ctx) error {
			return c.Reply("Halo! arg: " + c.Arg)
		},
	})
}
