package nycmta

import (
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"

	"github.com/lindsaylandry/go-transit-sign/src/signdata"
)

type TrainFeed struct {
	Station   Station
	Key       string
	Direction string

	Feed *gtfs.FeedMessage
}

func NewTrainFeed(station Station, accessKey, direction, url string) (*TrainFeed, error) {
	t := TrainFeed{}

	t.Key = accessKey
	t.Direction = direction

	t.Station = station

	feed, err := DecodeNYCMTA(accessKey, url)
	t.Feed = feed

	return &t, err
}

func (t *TrainFeed) GetArrivals() []signdata.Arrival {
	stopID := t.Station.GTFSStopID + t.Direction
	now := time.Now()
	arrivals := []signdata.Arrival{}
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

					a := signdata.Arrival{}
					a.Label = route
					a.Secs = secs

					arrivals = append(arrivals, a)
				}
			}
		}
	}

	return arrivals
}
