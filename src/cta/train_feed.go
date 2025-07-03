package cta

import (
	"strconv"

	"github.com/lindsaylandry/go-transit-sign/src/signdata"
)

type TrainFeed struct {
	Station  Station
	Key      string
	Timezone string
}

func NewTrainFeed(station Station, accessKey, timezone string) *TrainFeed {
	b := TrainFeed{}

	b.Key = accessKey
	b.Timezone = timezone
	b.Station = station

	return &b
}

func (t *TrainFeed) GetArrivals() ([]signdata.Arrival, error) {
	arrivals := []signdata.Arrival{}
	feed, err := DecodeTrain(t.Key, t.Station.StopID, TrainFeedURL)
	if err != nil {
		return arrivals, err
	}

	for _, f := range feed.TrainTimeResponse.Eta {
		arr := signdata.Arrival{}
		arr.Label = f.Name

		mins := 0
		if f.PredictedCountdown != "DUE" {
			mins, err = strconv.Atoi(f.PredictedCountdown)
			if err != nil {
				return arrivals, err
			}
			arr.Secs = int64(mins * 60)
		}

		arrivals = append(arrivals, arr)
	}

	return arrivals, nil
}
