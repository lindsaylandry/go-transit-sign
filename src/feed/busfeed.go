package feed

import (
	"time"

	"github.com/lindsaylandry/go-transit-sign/src/busstops"
	"github.com/lindsaylandry/go-transit-sign/src/decoder"
)

type BusFeed struct {
	BusStop busstops.CTABusStop
	Key     string
	Timezone string

	Feed decoder.CTABusFeedMessage
}

func NewBusFeed(busstop busstops.CTABusStop, accessKey, timezone string) (*BusFeed, error) {
	b := BusFeed{}

	b.Key = accessKey
	b.Timezone = timezone
	b.BusStop = busstop
	feed, err := decoder.DecodeCTA(accessKey, busstop.StopID, decoder.CTABusFeedURL)
	b.Feed = feed

	return &b, err
}

func (b *BusFeed) GetArrivals() ([]Arrival, error) {
	arrivals := []Arrival{}
	loc, err := time.LoadLocation(b.Timezone)
	if err != nil {
		return arrivals, err
	}

	for _, f := range b.Feed.BusTimeResponse.Prd {
		arr := Arrival{}
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
