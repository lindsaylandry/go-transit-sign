package feed

import (
	"time"

	"github.com/lindsaylandry/go-mta-train-sign/src/busstops"
	"github.com/lindsaylandry/go-mta-train-sign/src/decoder"
)

type BusFeed struct {
	BusStop busstops.BusStop
	Key     string

	Feed gtfs.FeedMessage
}

func NewBusFeed(busstop busstops.BusStop, accessKey, url string) (*BusFeed, error) {
	b := BusFeed{}

	b.Key = accessKey
	b.BusStop = busstop
	feed, err := decoder.Decode(accessKey, url)
	t.Feed = feed

	return &t, err
}

func (b *BusFeed) GetArrivals() []Arrival {
	arrivals := []Arrival{}

	return arrivals
}
