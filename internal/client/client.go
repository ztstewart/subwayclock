package client

import (
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/ztstewart/subwayclock/internal/client/transit_realtime"
)

// The MTA considers one physical station to be multiple stop IDs depending
// on the direction a train is travelling. For example, Grand Central on the 7
// line would have two stop IDs: 723N and 723S. 723N would be Grand Central on
// the 7 line in the direction of travel in Queens.
const (
	_northboundSuffx = 'N'
	_soutboundSuffx  = 'S'
)

const _avgNumStopsPerLine = 30

const _baseURL = "http://datamine.mta.info/mta_esi.php?key="

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
	Description string
}

// FeedUpdate represents the status at a given moment.
type FeedUpdate struct {
	Alerts        []Alert
	StationStatus map[string]StationStatus
}

// Config defines how to configure the subway client.
type Config struct {
	APIKey string `yaml:"api_key" json:"api_key"`
	FeedID int    `yaml:"feed_id" json:"feed_id"`
}

// NYCTA is a client for the New York City Transit Authority's realtime feed.
type NYCTA struct {
	cfg *Config
	url string
}

// NewNYCTA creates a new New York City Transit Authority client.
// An error will be returned if the configuration is invalid.
func NewNYCTA(cfg *Config) (*NYCTA, error) {

	url := "http://datamine.mta.info/mta_esi.php?key=" + cfg.APIKey + "&feed_id=" + strconv.Itoa(cfg.FeedID)
	return &NYCTA{
		cfg: cfg,
		url: url,
	}, nil
}

// GetFeed retrieves the current feed information.
// Currently for testing purposes it returns a JSON string.
func (n *NYCTA) GetFeed() (FeedUpdate, error) {
	resp, err := http.Get(n.url)
	if err != nil {
		return FeedUpdate{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return FeedUpdate{}, err
	}

	feedMessage := &transit_realtime.FeedMessage{}
	if err := proto.Unmarshal(body, feedMessage); err != nil {
		return FeedUpdate{}, err
	}

	return n.parseStatus(feedMessage)
}

func (n *NYCTA) parseStatus(feedMessage *transit_realtime.FeedMessage) (FeedUpdate, error) {
	stopToTimestamp := make(map[string][]StationUpdate, _avgNumStopsPerLine)
	var alerts []*transit_realtime.Alert

	for _, e := range feedMessage.Entity {
		if e.Alert != nil {
			alerts = append(alerts, e.Alert)
		}

		if e.TripUpdate == nil || e.TripUpdate.Trip == nil {
			continue
		}

		for _, stu := range e.TripUpdate.StopTimeUpdate {
			if stu.StopId == nil || stu.Arrival == nil || stu.Departure == nil {
				continue
			}

			stopID := *stu.StopId
			stopToTimestamp[stopID] = append(stopToTimestamp[stopID], StationUpdate{
				TripID:    *e.TripUpdate.Trip.TripId,
				Arrival:   time.Unix(*stu.Arrival.Time, 0),
				Departure: time.Unix(*stu.Departure.Time, 0),
			})
		}
	}

	for _, timestamps := range stopToTimestamp {
		sort.Slice(timestamps, func(i, j int) bool {
			return timestamps[i].Arrival.Before(timestamps[j].Arrival)
		})
	}

	update := FeedUpdate{
		StationStatus: make(map[string]StationStatus, len(stopToTimestamp)),
	}
	for k, v := range stopToTimestamp {
		last := k[len(k)-1]

		// Stop ID has a direction suffix. Don't use it as the stop ID.
		if len(k) > 1 && (last == _northboundSuffx || last == _soutboundSuffx) {
			stopID := k[:len(k)-1]
			stop := update.StationStatus[stopID]
			stop.StopID = stopID

			if stop.StopIDToUpdates == nil {
				stop.StopIDToUpdates = make(map[string][]StationUpdate, 2)
			}
			stop.StopIDToUpdates[k] = v

			update.StationStatus[stopID] = stop

			continue
		}

		// No direction suffix.
		stopID := k
		stop := update.StationStatus[stopID]
		stop.StopID = stopID

		if stop.StopIDToUpdates == nil {
			stop.StopIDToUpdates = make(map[string][]StationUpdate, 2)
		}
		stop.StopIDToUpdates[k] = v

		update.StationStatus[stopID] = stop

	}

	return update, nil
}
