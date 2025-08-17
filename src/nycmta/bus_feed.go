package nycmta

import (
	"strconv"
	"time"

	"github.com/lindsaylandry/go-transit-sign/src/signdata"
)

type BusFeed struct {
	BusStop BusStop
	Key     string
	URL     string
}

func NewBusFeed(busstop BusStop, accessKey string) *BusFeed {
	b := BusFeed{}

	b.Key = accessKey
	b.BusStop = busstop
	b.URL = BusFeedURL

	return &b
}

func (b *BusFeed) GetArrivals() ([]signdata.Arrival, error) {
	stopID := b.BusStop.StopID
	now := time.Now()
	arrivals := []signdata.Arrival{}

	feed, err := DecodeNYCMTA(b.Key, b.URL)
	if err != nil {
		return arrivals, err
	}
	for _, entity := range feed.Entity {
		trip := entity.GetTripUpdate()
		if trip != nil {
			stopTimes := trip.StopTimeUpdate
			for _, s := range stopTimes {
				if *s.StopId == strconv.Itoa(stopID) {
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

	return arrivals, nil
}
