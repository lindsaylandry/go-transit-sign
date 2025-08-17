package cta

import (
	"strconv"

	"github.com/lindsaylandry/go-transit-sign/src/signdata"
)

type BusFeed struct {
	BusStop BusStop
	Key     string
}

func NewBusFeed(busstop BusStop, accessKey string) *BusFeed {
	b := BusFeed{}

	b.Key = accessKey
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

	for _, f := range feed.BusTimeResponse.Prd {
		arr := signdata.Arrival{}
		arr.Label = f.Name

		mins := 0
		if f.PredictedCountdown == "DLY" {
			continue
		} else if f.PredictedCountdown != "DUE" {
			mins, err = strconv.Atoi(f.PredictedCountdown)
			if err != nil {
				return arrivals, err
			}
		}
		arr.Secs = int64(mins * 60)

		arrivals = append(arrivals, arr)
	}

	return arrivals, nil
}
