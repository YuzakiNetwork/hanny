# wabot (Go + whatsmeow)

Bot WhatsApp utility: tools, downloader (yt-dlp), manajemen grup.

## Jalanin
```bash
# butuh Go >= 1.26 (syarat whatsmeow), yt-dlp & ffmpeg di PATH
cp .env.example .env
go mod tidy          # sekali, narik dependency terbaru
go run .             # scan QR, atau isi PAIR_PHONE di .env buat kode pairing
go build -o wabot .  # binary
```

## Struktur
```
main.go              setup client + login
config/              baca .env (satu-satunya tempat)
plugin/              tipe Command, Ctx (Reply/React/Edit/Send), registry
handler/             pesan masuk → prefix → cari command → cek akses → run
command/
  command.go         daftar kategori (blank import)
  menu.go            .menu (otomatis dari registry)
  _template.go       contoh command baru
  tools/ downloader/ group/     1 file = 1 command
library/             logic reusable (tools, downloader, group, media)
```

## Nambah command
1. Salin `command/_template.go` ke `command/<kategori>/<nama>.go`, ganti `package` sesuai folder.
2. Selesai. Otomatis muncul di `.menu`. Akses (`Access`, `BotAdmin`) dicek handler, nggak perlu cek manual.

Kategori baru: bikin folder + file `<kategori>.go` berisi `var add = plugin.Category("<kategori>")`,
lalu tambah 1 baris import di `command/command.go`.
