package cta

import (
	"time"

	"github.com/lindsaylandry/go-transit-sign/src/signdata"
)

type BusFeed struct {
	BusStop  BusStop
	Key      string
	Timezone string

	Feed CTABusFeedMessage
}

func NewBusFeed(busstop BusStop, accessKey, timezone string) (*BusFeed, error) {
	b := BusFeed{}

	b.Key = accessKey
	b.Timezone = timezone
	b.BusStop = busstop
	feed, err := DecodeCTA(accessKey, busstop.StopID, CTABusFeedURL)
	b.Feed = feed

	return &b, err
}

func (b *BusFeed) GetArrivals() ([]signdata.Arrival, error) {
	arrivals := []signdata.Arrival{}
	loc, err := time.LoadLocation(b.Timezone)
	if err != nil {
		return arrivals, err
	}

	for _, f := range b.Feed.BusTimeResponse.Prd {
		arr := signdata.Arrival{}
		arr.Label = f.Name

		// find time
		t, err := time.ParseInLocation("20060102 15:04", f.PredictedTime, loc)
		if err != nil {
			return arrivals, err
		}

		now := time.Now()
		secs := t.Unix() - now.Unix()
		arr.Secs = secs

		arrivals = append(arrivals, arr)
	}

	return arrivals, nil
}
