package main

import (
	"image/png"
	"log"
	"os"

	"github.com/ztstewart/subwayclock/internal/client"
	"github.com/ztstewart/subwayclock/internal/render"
)

// An example stop ID for testing purposes.
const _courtSquareStopID = "719"
const _7NorthDir = "(7) MAIN ST"
const _7SouthDir = "(7) 34 - HUDSON"

func main() {
	nycta, err := client.NewNYCTA(&client.Config{
		APIKey: "",
		FeedID: 51,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	update, err := nycta.GetFeed()
	if err != nil {
		log.Fatalf("err: " + err.Error())
	}

	img := render.SubwayClock(update, _courtSquareStopID, _7NorthDir, _7SouthDir)

	f, err := os.Create("test.png")
	if err != nil {
		log.Fatalf("failed to open file %s", err.Error())
	}

	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		log.Fatal(err.Error())
	}
}
