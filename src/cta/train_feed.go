package cta

import (
	"errors"
	"time"

	"github.com/lindsaylandry/go-transit-sign/src/signdata"
)

type TrainFeed struct {
	Station Station
	Key     string
}

func NewTrainFeed(station Station, accessKey string) *TrainFeed {
	b := TrainFeed{}

	b.Key = accessKey
	b.Station = station

	return &b
}

func (t *TrainFeed) GetArrivals() ([]signdata.Arrival, error) {
	arrivals := []signdata.Arrival{}
	feed, err := DecodeTrain(t.Key, t.Station.StopID, TrainFeedURL)
	if err != nil {
		return arrivals, err
	}
	if feed.TrainTimeResponse.Error != "" {
		return arrivals, errors.New(feed.TrainTimeResponse.Error)
	}

	tmst, err := time.Parse("2006-01-02T15:04:05", feed.TrainTimeResponse.Timestamp)
	if err != nil {
		return arrivals, err
	}

	for _, f := range feed.TrainTimeResponse.Eta {
		arr := signdata.Arrival{}
		arr.Label = f.Name

		// use time to get minutes
		prdtm, err := time.Parse("2006-01-02T15:04:05", f.ArrivalTime)
		if err != nil {
			return arrivals, err
		}

		arr.Secs = prdtm.Unix() - tmst.Unix()

		arrivals = append(arrivals, arr)
	}

	return arrivals, nil
}
