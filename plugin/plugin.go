// Package plugin: tipe Command, konteks pesan (Ctx), dan registry. Baca ini kalau mau bikin command baru.
package plugin

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// Access dicek otomatis oleh handler, command nggak perlu ngecek sendiri.
type Access int

const (
	Everyone   Access = iota // siapa saja, di mana saja
	GroupOnly                // cuma di grup
	GroupAdmin               // cuma di grup + pengirim harus admin
)

type Command struct {
	Name     string
	Aliases  []string
	Desc     string
	Usage    string
	Access   Access
	BotAdmin bool // bot harus admin grup (kick, add, promote, dll)
	Run      func(c *Ctx) error
}

type Entry struct {
	Category string
	Cmd      *Command
}

var (
	mu      sync.RWMutex
	entries []Entry
	index   = map[string]*Command{}
)

// Register mendaftarkan command. Panic saat startup kalau nama/alias bentrok (fail-fast).
func Register(category string, cmd Command) {
	if cmd.Name == "" || cmd.Run == nil {
		panic("plugin: command butuh Name dan Run")
	}
	mu.Lock()
	defer mu.Unlock()
	c := &cmd
	for _, n := range append([]string{cmd.Name}, cmd.Aliases...) {
		n = strings.ToLower(n)
		if _, dup := index[n]; dup {
			panic(fmt.Sprintf("plugin: nama command dobel: %q", n))
		}
		index[n] = c
	}
	entries = append(entries, Entry{category, c})
}

// Category mengembalikan fungsi pendaftar untuk satu kategori: var add = plugin.Category("tools")
func Category(name string) func(Command) {
	return func(cmd Command) { Register(name, cmd) }
}

func Find(name string) *Command {
	mu.RLock()
	defer mu.RUnlock()
	return index[strings.ToLower(name)]
}

// All mengembalikan semua command terurut per kategori lalu nama.
func All() []Entry {
	mu.RLock()
	defer mu.RUnlock()
	out := append([]Entry(nil), entries...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Category != out[j].Category {
			return out[i].Category < out[j].Category
		}
		return out[i].Cmd.Name < out[j].Cmd.Name
	})
	return out
}

// Ctx = semua yang dibutuhin command untuk satu pesan masuk.
type Ctx struct {
	Context context.Context
	Client  *whatsmeow.Client
	Event   *events.Message
	Arg     string // teks setelah nama command
}

func (c *Ctx) Chat() types.JID   { return c.Event.Info.Chat }
func (c *Ctx) Sender() types.JID { return c.Event.Info.Sender }
func (c *Ctx) IsGroup() bool     { return c.Event.Info.IsGroup }
func (c *Ctx) Text() string      { return TextOf(c.Event.Message) }

func (c *Ctx) Send(msg *waE2E.Message) (whatsmeow.SendResponse, error) {
	return c.Client.SendMessage(c.Context, c.Chat(), msg)
}

// ReplyID membalas (quote) pesan pengirim dan mengembalikan ID pesan balasan (buat Edit).
func (c *Ctx) ReplyID(text string) (types.MessageID, error) {
	ci := &waE2E.ContextInfo{
		StanzaID:      proto.String(c.Event.Info.ID),
		QuotedMessage: c.Event.Message,
	}
	if c.IsGroup() {
		ci.Participant = proto.String(c.Sender().String())
	}
	resp, err := c.Send(&waE2E.Message{ExtendedTextMessage: &waE2E.ExtendedTextMessage{
		Text:        proto.String(text),
		ContextInfo: ci,
	}})
	return resp.ID, err
}

func (c *Ctx) Reply(text string) error {
	_, err := c.ReplyID(text)
	return err
}

func (c *Ctx) React(emoji string) error {
	_, err := c.Send(c.Client.BuildReaction(c.Chat(), c.Sender(), c.Event.Info.ID, emoji))
	return err
}

func (c *Ctx) Edit(id types.MessageID, text string) error {
	_, err := c.Send(c.Client.BuildEdit(c.Chat(), id, &waE2E.Message{Conversation: proto.String(text)}))
	return err
}

// TextOf mengambil teks dari berbagai tipe pesan (teks biasa, reply, caption gambar/video).
func TextOf(m *waE2E.Message) string {
	switch {
	case m.GetConversation() != "":
		return m.GetConversation()
	case m.GetExtendedTextMessage().GetText() != "":
		return m.GetExtendedTextMessage().GetText()
	case m.GetImageMessage().GetCaption() != "":
		return m.GetImageMessage().GetCaption()
	case m.GetVideoMessage().GetCaption() != "":
		return m.GetVideoMessage().GetCaption()
	}
	return ""
}
