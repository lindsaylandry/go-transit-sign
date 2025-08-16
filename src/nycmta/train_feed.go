package nycmta

import (
	"time"

	"github.com/lindsaylandry/go-transit-sign/src/signdata"
)

type TrainFeed struct {
	Station   Station
	Key       string
	Direction string
	Decoders  *[]TrainDecoder
}

func NewTrainFeed(station Station, accessKey, direction, daytimeRoutes string) *TrainFeed {
	t := TrainFeed{}

	t.Key = accessKey
	t.Direction = direction
	t.Station = station
	t.Decoders = GetMtaTrainDecoders(daytimeRoutes)

	return &t
}

func (t *TrainFeed) GetArrivals() ([]signdata.Arrival, error) {
	stopID := t.Station.GTFSStopID + t.Direction
	now := time.Now()
	arrivals := []signdata.Arrival{}

	for _, d := range *t.Decoders {
		feed, err := DecodeNYCMTA(t.Key, d.URL)
		if err != nil {
			return arrivals, err
		}
		for _, entity := range feed.Entity {
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
	}

	return arrivals, nil
}
