package main

import (
	"image"
	"log"
	"time"

	"github.com/ztstewart/subwayclock/internal/client"
	"github.com/ztstewart/subwayclock/internal/inky"
	"github.com/ztstewart/subwayclock/internal/render"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

// An example stop ID for testing purposes.
const _courtSquareStopID = "719"
const _7NorthDir = "(7) MAIN ST"
const _7SouthDir = "(7) 34 - HUDSON"

const _updateFrequency = 1 * time.Minute

func main() {
	nycta, err := client.NewNYCTA(&client.Config{
		APIKey: "",
		FeedID: 51,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	b, err := spireg.Open("SPI0.0")
	if err != nil {
		log.Fatal(err)
	}

	dc := gpioreg.ByName("22")
	reset := gpioreg.ByName("27")
	busy := gpioreg.ByName("17")

	log.Println("initializing InkyPHAT display")

	dev, err := inky.New(b, dc, reset, busy, &inky.Opts{
		Model:       inky.PHAT,
		ModelColor:  inky.Red,
		BorderColor: inky.Black,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("display initialized")

	ticker := time.NewTicker(_updateFrequency)

	updateScreen := func() {
		log.Println("fetching new subway status")

		update, err := nycta.GetFeed()
		if err != nil {
			log.Fatalf("err: " + err.Error())
		}

		img := render.SubwayClock(update, _courtSquareStopID, _7NorthDir, _7SouthDir)

		if err := dev.Draw(img.Bounds(), img, image.ZP); err != nil {
			log.Fatal(err)
		}
	}

	updateScreen()

	for range ticker.C {
		updateScreen()
	}
}
