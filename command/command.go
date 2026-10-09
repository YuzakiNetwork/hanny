// Package command: pintu masuk semua command. Tiap kategori = satu subfolder (package) yang
// mendaftarkan commandnya lewat init(). Kategori BARU cukup tambah satu baris import di bawah.
// Command baru di kategori yang sudah ada: drop file saja, nggak perlu edit apa-apa.
package command

import (
	_ "wabot/command/downloader"
	_ "wabot/command/group"
	_ "wabot/command/tools"
)
