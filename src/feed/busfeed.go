package feed

import (
	//	"os"
	"fmt"
	"time"

	"github.com/lindsaylandry/go-transit-sign/src/busstops"
	"github.com/lindsaylandry/go-transit-sign/src/decoder"
)

type BusFeed struct {
	BusStop busstops.CTABusStop
	Key     string

	Feed decoder.CTABusFeedMessage
}

func NewBusFeed(busstop busstops.CTABusStop, accessKey string) (*BusFeed, error) {
	b := BusFeed{}

	b.Key = accessKey
	b.BusStop = busstop
	feed, err := decoder.DecodeCTA(accessKey, busstop.StopID, decoder.CTABusFeedURL)
	b.Feed = feed

	return &b, err
}

func (b *BusFeed) GetArrivals() ([]Arrival, error) {
	arrivals := []Arrival{}

	for _, f := range b.Feed.BusTimeResponse.Prd {
		arr := Arrival{}
		arr.Label = f.Name

		// find time
		t, err := time.ParseInLocation("20060102 15:04", f.PredictedTime, time.Local)
		if err != nil {
			return arrivals, err
		}

		now := time.Now()
		secs := t.Unix() - now.Unix()
		arr.Secs = secs

		fmt.Println(f.PredictedTime)

		arrivals = append(arrivals, arr)
	}

	return arrivals, nil
}
