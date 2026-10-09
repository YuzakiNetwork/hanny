// Package media: upload + kirim media ke chat asal pesan.
package media

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"

	"wabot/plugin"
)

func SendVideo(c *plugin.Ctx, data []byte, caption string) error {
	up, err := c.Client.Upload(c.Context, data, whatsmeow.MediaVideo)
	if err != nil {
		return err
	}
	_, err = c.Send(&waE2E.Message{VideoMessage: &waE2E.VideoMessage{
		URL:           &up.URL,
		DirectPath:    &up.DirectPath,
		MediaKey:      up.MediaKey,
		Mimetype:      proto.String("video/mp4"),
		FileEncSHA256: up.FileEncSHA256,
		FileSHA256:    up.FileSHA256,
		FileLength:    &up.FileLength,
		Caption:       proto.String(caption),
	}})
	return err
}

func SendAudio(c *plugin.Ctx, data []byte) error {
	up, err := c.Client.Upload(c.Context, data, whatsmeow.MediaAudio)
	if err != nil {
		return err
	}
	_, err = c.Send(&waE2E.Message{AudioMessage: &waE2E.AudioMessage{
		URL:           &up.URL,
		DirectPath:    &up.DirectPath,
		MediaKey:      up.MediaKey,
		Mimetype:      proto.String("audio/mpeg"),
		FileEncSHA256: up.FileEncSHA256,
		FileSHA256:    up.FileSHA256,
		FileLength:    &up.FileLength,
		PTT:           proto.Bool(false),
	}})
	return err
}
