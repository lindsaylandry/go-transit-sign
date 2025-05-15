package feed

import (
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"time"

	"github.com/lindsaylandry/go-transit-sign/src/decoder"
	"github.com/lindsaylandry/go-transit-sign/src/nycmta"
)

type TrainFeed struct {
	Station   nycmta.Station
	Key       string
	Direction string

	Feed *gtfs.FeedMessage
}

type Arrival struct {
	Label string
	Secs  int64
}

func NewTrainFeed(station nycmta.Station, accessKey, direction, url string) (*TrainFeed, error) {
	t := TrainFeed{}

	t.Key = accessKey
	t.Direction = direction

	t.Station = station

	feed, err := decoder.DecodeNYCMTA(accessKey, url)
	t.Feed = feed

	return &t, err
}

func (t *TrainFeed) GetArrivals() []Arrival {
	stopID := t.Station.GTFSStopID + t.Direction
	now := time.Now()
	arrivals := []Arrival{}
	for _, entity := range t.Feed.Entity {
		trip := entity.GetTripUpdate()
		if trip != nil {
			stopTimes := trip.StopTimeUpdate
			for _, s := range stopTimes {
				if *s.StopId == stopID {
					route := ""
					vehicle := trip.Trip
					if vehicle != nil {
						route = *vehicle.RouteId
					}
					delay := int32(0)
					if s.Arrival.Delay != nil {
						delay = *s.Arrival.Delay
					}

					secs := *s.Arrival.Time + int64(delay) - now.Unix()

					a := Arrival{}
					a.Label = route
					a.Secs = secs

					arrivals = append(arrivals, a)
				}
			}
		}
	}

	return arrivals
}
