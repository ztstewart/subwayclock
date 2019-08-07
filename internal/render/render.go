package render

import (
	"image"
	"image/color"
	"strconv"
	"time"

	"github.com/ztstewart/subwayclock/internal/models"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Screen dimensions. Change for larger screen sizes.
const (
	_maxX = 212
	_maxY = 104
)

// Margins. X coordinates; higher numbers move them to the right of the image.
const (
	_leftMargin  = 50
	_rightMargin = 180
)

// Position of the top row label and the bottom row label.
// These are Y coordinates; increasing the number moves them towards the
// bottom of the image.
const (
	_topRow    = 30
	_bottomRow = 70
)

const (
	_gapBetweenUpdates        = 30 // Space between each predicted arrrival.
	_gapBetweenStopAndUpdates = 20 // Space between the top and bottom rows and arrivals.
)

var (
	_black = color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}
	_red = color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	}
	_white = color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
)

// SubwayClock renders an image suitable for display as a subway clock.
// The output will be an image, 212 x 104, that has two sections: one for the
// stop in the "southbound" direction on the top section and one for the
// "northbound" direction in the bottom row.
//
// Each section will be labelled with the input text labels. Below the text,
// this function will render the minutes until the next few trains will
// arrive, rounded up to the nearest minute.
//
// If a train is about to arrive, it may render 0. Depending on the last
// update time, it is possible that some trains may have already departed.
// In that case, those trains will be skipped.
func SubwayClock(update models.FeedUpdate, stopID string, northLabel string, southLabel string) *image.RGBA {
	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: _maxX, Y: _maxY},
	})

	addLabel(img, _leftMargin, _topRow, southLabel, _red)
	addLabel(img, _leftMargin, _bottomRow, northLabel, _red)

	status, ok := update.StationStatus[stopID]
	if !ok { // No stop information, don't render anything else.
		return img
	}

	now := time.Now()

	if updates, ok := status.StopIDToUpdates[stopID+"S"]; ok {
		renderArrivals(updates, img, _leftMargin+_gapBetweenUpdates, _topRow+_gapBetweenStopAndUpdates, now)
	}

	if updates, ok := status.StopIDToUpdates[stopID+"N"]; ok {
		renderArrivals(updates, img, _leftMargin+_gapBetweenUpdates, _bottomRow+_gapBetweenStopAndUpdates, now)
	}

	if len(update.Alerts) != 0 {
		addLabel(img, 0, 13, "SERVICE ALERT - CHECK MTA.INFO", _red)
	}

	return img
}

func renderArrivals(updates []models.StationUpdate, img *image.RGBA, x int, y int, now time.Time) {
	for _, update := range updates {
		minutes := update.Arrival.Sub(now).Round(time.Minute) / time.Minute
		if minutes < 0 {
			continue
		}

		addLabel(img, x, y, strconv.Itoa(int(minutes)), _white)
		x += _gapBetweenUpdates
		if x > _rightMargin {
			break
		}
	}
}

// addLabel adds a label to an image.
func addLabel(img *image.RGBA, x, y int, label string, col color.RGBA) {
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
