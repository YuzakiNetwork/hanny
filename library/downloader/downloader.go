// Package downloader: download lewat yt-dlp (+ ffmpeg), ratusan situs tanpa nulis scraper sendiri.
package downloader

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"wabot/config"
	"wabot/library/media"
	"wabot/plugin"
)

type Kind int

const (
	Video Kind = iota
	Audio
)

var (
	slotsOnce sync.Once
	slots     chan struct{}
)

// Deliver = alur lengkap buat command: react ⏳ → download → kirim → react ✅/❌.
func Deliver(c *plugin.Ctx, url string, kind Kind) error {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("Kirim link valid, contoh: https://...")
	}
	_ = c.React("⏳")
	err := func() error {
		data, err := Fetch(c.Context, url, kind)
		if err != nil {
			return err
		}
		if kind == Audio {
			return media.SendAudio(c, data)
		}
		return media.SendVideo(c, data, "✅ Selesai")
	}()
	if err != nil {
		_ = c.React("❌")
	} else {
		_ = c.React("✅")
	}
	return err
}

// Fetch mengunduh url dan mengembalikan isi file hasilnya.
func Fetch(ctx context.Context, url string, kind Kind) ([]byte, error) {
	cfg := config.Get()
	slotsOnce.Do(func() { slots = make(chan struct{}, cfg.DLConcurrency) })
	select { // batasi download paralel
	case slots <- struct{}{}:
		defer func() { <-slots }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	dir, err := os.MkdirTemp("", "wabot-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir) // selalu bersihin, sukses atau gagal

	ctx, cancel := context.WithTimeout(ctx, cfg.DLTimeout)
	defer cancel()

	args := []string{"--no-playlist", "--no-warnings", "-q", "--max-filesize", fmt.Sprintf("%dM", cfg.DLMaxMB)}
	if kind == Audio {
		args = append(args, "-x", "--audio-format", "mp3")
	} else {
		args = append(args, "-f", "b[ext=mp4]/bv*+ba/b", "-S", "res:720", "--merge-output-format", "mp4")
	}
	// "--" sebelum URL: cegah URL dianggap opsi yt-dlp (argument injection)
	args = append(args, "-o", filepath.Join(dir, "out.%(ext)s"), "--", url)

	if out, err := exec.CommandContext(ctx, "yt-dlp", args...).CombinedOutput(); err != nil {
		switch {
		case errors.Is(err, exec.ErrNotFound):
			return nil, errors.New("yt-dlp belum terpasang di server")
		case errors.Is(ctx.Err(), context.DeadlineExceeded):
			return nil, fmt.Errorf("Timeout (%s)", cfg.DLTimeout)
		}
		log.Printf("yt-dlp gagal: %v\n%s", err, out)
		return nil, fmt.Errorf("Gagal download (link private / > %dMB / tidak didukung)", cfg.DLMaxMB)
	}

	files, err := os.ReadDir(dir)
	if err != nil || len(files) == 0 {
		return nil, errors.New("File hasil tidak ditemukan")
	}
	path := filepath.Join(dir, files[0].Name())
	if st, err := os.Stat(path); err != nil || st.Size() > cfg.DLMaxMB*1024*1024 {
		return nil, fmt.Errorf("File > %dMB", cfg.DLMaxMB)
	}
	return os.ReadFile(path)
}
