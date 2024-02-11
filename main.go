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

	b.Handle(tele.OnMedia, func(c tele.Context) error {
		// return c.Send("Photo received!")
		media := c.Message().Media()
		if media.MediaType() == "sticker" {
			sticker_data := media.(*tele.Sticker)
			sticker_file := sticker_data.File
			b.Download(&sticker_file, "sticker.webp")
			ImageOutline("sticker.webp", "sticker_outlined.png")
			return c.Send(fmt.Sprintf("Height: %d, Width: %d", sticker_data.Height, sticker_data.Width))
		}
		return c.Send("Not a sticker")
	})

	b.Start()
}
