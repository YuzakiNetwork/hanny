// Package group: command kategori "group" (semua khusus admin grup).
package group

import (
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"

	libgroup "wabot/library/group"
	"wabot/plugin"
)

var add = plugin.Category("group")

// member membuat Run untuk aksi anggota (kick/add/promote/demote).
func member(action whatsmeow.ParticipantChange, verb string) func(*plugin.Ctx) error {
	return func(c *plugin.Ctx) error {
		targets, err := libgroup.Targets(c, c.Arg)
		if err != nil {
			return err
		}
		res, err := c.Client.UpdateGroupParticipants(c.Context, c.Chat(), targets, action)
		if err != nil {
			return err
		}
		ok := 0
		for _, r := range res {
			if r.Error == 0 {
				ok++
			}
		}
		return c.Reply(fmt.Sprintf("✅ %d dari %d orang %s", ok, len(targets), verb))
	}
}

// mention mengirim teks dengan mention ke semua anggota grup.
func mention(c *plugin.Ctx, text string, list bool) error {
	info, err := c.Client.GetGroupInfo(c.Context, c.Chat())
	if err != nil {
		return err
	}
	jids := make([]string, 0, len(info.Participants))
	for _, p := range info.Participants {
		jids = append(jids, p.JID.String())
		if list {
			text += "\n@" + p.JID.User
		}
	}
	_, err = c.Send(&waE2E.Message{ExtendedTextMessage: &waE2E.ExtendedTextMessage{
		Text:        proto.String(text),
		ContextInfo: &waE2E.ContextInfo{MentionedJID: jids},
	}})
	return err
}
