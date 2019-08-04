package models

import (
	"time"
)

// StationUpdate  records a scheduled (or realtime) arrival and departure for
// a trip at a particular station.
type StationUpdate struct {
	TripID    string
	Arrival   time.Time
	Departure time.Time
}

// A StationStatus records all scheduled arrivals and departures for a given
// stop. Stops on the same line, but in different directions, will be grouped
// into a given StationStatus struct.
type StationStatus struct {
	StopID          string
	StopIDToUpdates map[string][]StationUpdate // multiple directions will be separate keys
}

// An Alert is a type of service disruption, delay, etc.
type Alert struct {
	Effect string
	Header string
}

// FeedUpdate represents the status at a given moment.
type FeedUpdate struct {
	Alerts        []Alert
	StationStatus map[string]StationStatus
}
