package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "modernc.org/sqlite" // driver SQLite pure-Go (tanpa CGO)

	_ "wabot/command" // daftarin semua command
	"wabot/config"
	"wabot/handler"
)

func main() {
	cfg := config.Get() // memuat .env
	ctx := context.Background()

	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=busy_timeout(10000)", cfg.DBPath)
	container, err := sqlstore.New(ctx, "sqlite", dsn, waLog.Stdout("DB", cfg.LogLevel, true))
	if err != nil {
		log.Fatalf("buka database: %v", err)
	}
	device, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Fatalf("ambil device: %v", err)
	}

	client := whatsmeow.NewClient(device, waLog.Stdout("Client", cfg.LogLevel, true))
	client.AddEventHandler(handler.New(client))

	if err := login(ctx, client, cfg.PairPhone); err != nil {
		log.Fatalf("login: %v", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	client.Disconnect()
}

func login(ctx context.Context, client *whatsmeow.Client, pairPhone string) error {
	if client.Store.ID != nil { // sudah pernah login
		return client.Connect()
	}

	if pairPhone != "" { // kode pairing, tanpa scan QR
		if err := client.Connect(); err != nil {
			return err
		}
		time.Sleep(2 * time.Second) // tunggu koneksi siap (lihat doc PairPhone)
		code, err := client.PairPhone(ctx, pairPhone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
		if err != nil {
			return err
		}
		fmt.Printf("\nKode pairing: %s\nWhatsApp > Perangkat Tertaut > Tautkan dengan nomor telepon\n\n", code)
		return nil
	}

	qrChan, err := client.GetQRChannel(ctx) // harus dipanggil sebelum Connect
	if err != nil {
		return err
	}
	if err := client.Connect(); err != nil {
		return err
	}
	go func() {
		for item := range qrChan {
			if item.Event == whatsmeow.QRChannelEventCode {
				fmt.Println()
				qrterminal.GenerateHalfBlock(item.Code, qrterminal.L, os.Stdout)
				fmt.Println("Scan: WhatsApp > Perangkat Tertaut")
			} else {
				log.Printf("QR: %s", item.Event)
			}
		}
	}()
	return nil
}
