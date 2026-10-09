// Package group: helper grup (cek admin, resolusi target, antilink).
package group

import (
	"errors"
	"strings"
	"sync"

	"go.mau.fi/whatsmeow/types"

	"wabot/plugin"
)

// ponytail: in-memory, hilang pas restart; pindah ke tabel sqlite kalau perlu persist.
var (
	alMu    sync.Mutex
	alChats = map[string]bool{}
)

func SetAntilink(chat types.JID, on bool) {
	alMu.Lock()
	defer alMu.Unlock()
	if on {
		alChats[chat.String()] = true
	} else {
		delete(alChats, chat.String())
	}
}

func antilinkOn(chat types.JID) bool {
	alMu.Lock()
	defer alMu.Unlock()
	return alChats[chat.String()]
}

// isAdmin: salah satu dari `who` adalah admin grup (grup bisa pakai LID atau nomor HP).
func isAdmin(c *plugin.Ctx, who ...types.JID) (bool, error) {
	info, err := c.Client.GetGroupInfo(c.Context, c.Chat())
	if err != nil {
		return false, err
	}
	for _, p := range info.Participants {
		if !p.IsAdmin && !p.IsSuperAdmin {
			continue
		}
		for _, w := range who {
			if w.User != "" && (w.User == p.JID.User || w.User == p.PhoneNumber.User || w.User == p.LID.User) {
				return true, nil
			}
		}
	}
	return false, nil
}

func SenderIsAdmin(c *plugin.Ctx) (bool, error) {
	return isAdmin(c, c.Sender(), c.Event.Info.SenderAlt)
}

func BotIsAdmin(c *plugin.Ctx) (bool, error) {
	st := c.Client.Store
	var pn types.JID
	if st.ID != nil {
		pn = *st.ID
	}
	return isAdmin(c, pn, st.LID)
}

// Targets: @mention, atau reply ke pesan orang, atau nomor di argumen. Error kalau kosong.
func Targets(c *plugin.Ctx, arg string) ([]types.JID, error) {
	var out []types.JID
	ci := c.Event.Message.GetExtendedTextMessage().GetContextInfo()
	for _, s := range ci.GetMentionedJID() {
		if j, err := types.ParseJID(s); err == nil {
			out = append(out, j.ToNonAD())
		}
	}
	if len(out) == 0 && ci.GetParticipant() != "" {
		if j, err := types.ParseJID(ci.GetParticipant()); err == nil {
			out = append(out, j.ToNonAD())
		}
	}
	if len(out) == 0 {
		for _, w := range strings.Fields(arg) {
			n := strings.Map(func(r rune) rune {
				if r >= '0' && r <= '9' {
					return r
				}
				return -1
			}, w)
			if len(n) >= 8 {
				out = append(out, types.NewJID(n, types.DefaultUserServer))
			}
		}
	}
	if len(out) == 0 {
		return nil, errors.New("Tag orangnya (@user), reply pesannya, atau tulis nomornya")
	}
	return out, nil
}

// Antilink = hook: hapus pesan berisi link grup WA dari non-admin (kalau antilink aktif & bot admin). Best-effort.
func Antilink(c *plugin.Ctx, text string) {
	if !strings.Contains(text, "chat.whatsapp.com/") || !antilinkOn(c.Chat()) {
		return
	}
	if ok, err := SenderIsAdmin(c); err != nil || ok {
		return
	}
	if ok, err := BotIsAdmin(c); err != nil || !ok {
		return
	}
	_, _ = c.Send(c.Client.BuildRevoke(c.Chat(), c.Sender(), c.Event.Info.ID))
}
