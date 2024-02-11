package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

var startTime = time.Now()

func main() {
	godotenv.Load(".env")
	// discards updates older than client start time
	poller := &tele.MiddlewarePoller{
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},

		Filter: func(u *tele.Update) bool {
			return u.Message.Time().After(startTime)
		},
	}
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: poller,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(c tele.Context) error {
		return c.Send(`Hello! I'm bot that can outline stickers. 
		Just send me a sticker and I'll outline it for you, then forward the resulted image to the @Stickers to make it as a sticker.
		<b>As a prerequisite you should already have a stiker pack configured in @Stickers, perform "adsticker" command and make @Stickers wait for the image.</b>`, tele.ModeHTML)
	})

	b.Handle(tele.OnMedia, func(c tele.Context) error {
		media := c.Message().Media()
		if media.MediaType() == "sticker" {
			sticker_data := media.(*tele.Sticker)
			sticker_file := sticker_data.File
			webpFilename := fmt.Sprintf("%s.webp", sticker_file.FileID[0:10])
			pngFilename := fmt.Sprintf("%s.png", sticker_file.FileID[0:10])
			b.Download(&sticker_file, webpFilename)
			outlinedSticker := ImageOutline(webpFilename, pngFilename)
			file := &tele.Document{File: tele.FromDisk(outlinedSticker), MIME: "image/png", FileName: pngFilename}
			c.Send(file)
			os.Remove(webpFilename)
			os.Remove(pngFilename)
			return c.Send("Now just forward this image to the Stickers chat to make it as a sticker")
		}
		return c.Send("Not a sticker")
	})

	b.Start()
}
