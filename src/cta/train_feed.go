package cta

import (
	"time"

	"github.com/lindsaylandry/go-transit-sign/src/signdata"
)

type TrainFeed struct {
	Station  Station
	Key      string
	Timezone string

	Feed TrainFeedMessage
}

func NewTrainFeed(station Station, accessKey, timezone string) (*TrainFeed, error) {
	b := TrainFeed{}

	b.Key = accessKey
	b.Timezone = timezone
	b.Station = station
	feed, err := DecodeTrain(accessKey, station.StopID, TrainFeedURL)
	b.Feed = feed

	return &b, err
}

func (b *TrainFeed) GetArrivals() ([]signdata.Arrival, error) {
	arrivals := []signdata.Arrival{}
	loc, err := time.LoadLocation(b.Timezone)
	if err != nil {
		return arrivals, err
	}

	for _, f := range b.Feed.TrainTimeResponse.Eta {
		arr := signdata.Arrival{}
		arr.Label = f.Name

		// find time
		t, err := time.ParseInLocation("20060102 15:04:05", f.PredictedTime, loc)
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
