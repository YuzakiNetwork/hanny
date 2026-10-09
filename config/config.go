// Package config: satu-satunya tempat baca environment (.env).
package config

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	BotName       string
	Prefix        string
	DBPath        string
	LogLevel      string
	PairPhone     string
	IgnoreOlder   time.Duration
	TZOffsetHours int
	DLMaxMB       int64
	DLTimeout     time.Duration
	DLConcurrency int
}

var (
	once sync.Once
	cfg  *Config
)

// Get memuat .env sekali lalu mengembalikan config yang sama.
func Get() *Config {
	once.Do(func() {
		_ = godotenv.Load() // .env opsional; env asli tetap dipakai
		cfg = &Config{
			BotName:       str("BOT_NAME", "WaBot"),
			Prefix:        str("PREFIX", "."),
			DBPath:        str("DB_PATH", "whatsapp.db"),
			LogLevel:      strings.ToUpper(str("LOG_LEVEL", "WARN")),
			PairPhone:     str("PAIR_PHONE", ""),
			IgnoreOlder:   time.Duration(num("IGNORE_OLDER_THAN_SECS", 120)) * time.Second,
			TZOffsetHours: int(num("TZ_OFFSET_HOURS", 7)),
			DLMaxMB:       num("DL_MAX_MB", 64),
			DLTimeout:     time.Duration(num("DL_TIMEOUT_SECS", 300)) * time.Second,
			DLConcurrency: int(max(num("DL_CONCURRENCY", 2), 1)),
		}
	})
	return cfg
}

func str(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func num(key string, def int64) int64 {
	if v, err := strconv.ParseInt(strings.TrimSpace(os.Getenv(key)), 10, 64); err == nil {
		return v
	}
	return def
}
