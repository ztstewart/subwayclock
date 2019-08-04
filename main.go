package main

import (
	"encoding/json"
	"fmt"

	"github.com/ztstewart/subwayclock/internal/client"
)

// An example stop ID for testing purposes.
const _courtSquareStopID = "719"

func main() {
	nycta, _ := client.NewNYCTA(&client.Config{
		APIKey: "",
		FeedID: 51,
	})

	update, err := nycta.GetFeed()
	if err != nil {
		fmt.Println("err: " + err.Error())
		return
	}

	// Print out a JSON formatted example of our update.
	if status, ok := update.StationStatus[_courtSquareStopID]; ok {
		bytes, _ := json.Marshal(status)
		fmt.Println(string(bytes))
	}

}
