package cta

import (
	"strconv"

	"github.com/lindsaylandry/go-transit-sign/src/signdata"
)

type BusFeed struct {
	BusStop  BusStop
	Key      string
	Timezone string

	Feed BusFeedMessage
}

func NewBusFeed(busstop BusStop, accessKey, timezone string) *BusFeed {
	b := BusFeed{}

	b.Key = accessKey
	b.Timezone = timezone
	b.BusStop = busstop

	return &b
}

func (b *BusFeed) GetArrivals() ([]signdata.Arrival, error) {
	arrivals := []signdata.Arrival{}

	feed, err := DecodeBus(b.Key, b.BusStop.StopID, BusFeedURL)
	if err != nil {
		return arrivals, err
	}

	// TODO: read feed errors

	b.Feed = feed

	for _, f := range b.Feed.BusTimeResponse.Prd {
		arr := signdata.Arrival{}
		arr.Label = f.Name

		mins, err := strconv.Atoi(f.PredictedCountdown)
		if err != nil {
			return arrivals, err
		}
		arr.Secs = int64(mins * 60)

		arrivals = append(arrivals, arr)
	}

	return arrivals, nil
}
