// Package handler: alur pesan masuk → parse prefix → cari command → cek akses → jalanin.
package handler

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"

	"wabot/config"
	libgroup "wabot/library/group"
	"wabot/plugin"
)

func New(cli *whatsmeow.Client) whatsmeow.EventHandler {
	return func(evt any) {
		switch e := evt.(type) {
		case *events.Message:
			go handle(cli, e) // satu goroutine per pesan: download lama nggak ngeblok yang lain
		case *events.Connected:
			log.Printf("✅ %s online", config.Get().BotName)
		case *events.LoggedOut:
			log.Println("⚠️ Sesi logout dari HP. Hapus database sesi lalu pairing ulang.")
		}
	}
}

func handle(cli *whatsmeow.Client, e *events.Message) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic di handler: %v", r)
		}
	}()
	cfg := config.Get()
	if e.Info.IsFromMe || time.Since(e.Info.Timestamp) > cfg.IgnoreOlder {
		return // pesan sendiri / backlog lama
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	c := &plugin.Ctx{Context: ctx, Client: cli, Event: e}

	if err := process(c, cfg.Prefix); err != nil {
		_ = c.Reply("⚠️ " + err.Error())
	}
}

func process(c *plugin.Ctx, prefix string) error {
	text := strings.TrimSpace(c.Text())

	// Hook non-command (jalan di setiap pesan grup)
	if c.IsGroup() {
		libgroup.Antilink(c, text)
	}

	rest, ok := strings.CutPrefix(text, prefix)
	if !ok {
		return nil
	}
	fields := strings.Fields(rest)
	if len(fields) == 0 {
		return nil
	}
	cmd := plugin.Find(fields[0])
	if cmd == nil {
		return nil
	}
	c.Arg = strings.TrimSpace(strings.TrimPrefix(rest, fields[0]))

	// Cek akses terpusat
	if cmd.Access != plugin.Everyone && !c.IsGroup() {
		return errors.New("Command ini cuma buat grup")
	}
	if cmd.Access == plugin.GroupAdmin {
		if ok, err := libgroup.SenderIsAdmin(c); err != nil {
			return err
		} else if !ok {
			return errors.New("Khusus admin grup")
		}
	}
	if cmd.BotAdmin {
		if ok, err := libgroup.BotIsAdmin(c); err != nil {
			return err
		} else if !ok {
			return errors.New("Jadiin bot admin dulu")
		}
	}
	return cmd.Run(c)
}
