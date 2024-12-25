package trainfeed

import (
	"fmt"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"sort"
	"time"

	"github.com/lindsaylandry/go-mta-train-sign/src/decoder"
	"github.com/lindsaylandry/go-mta-train-sign/src/stations"
)

type TrainFeed struct {
	Station   stations.MtaStation
	Key       string
	Direction string

	Feed gtfs.FeedMessage
}

type Arrival struct {
	Train string
	Secs  int64
}

func NewTrainFeed(station stations.MtaStation, accessKey, direction, url string) (*SignData, error) {
	t := SignData{}

	t.Key = accessKey
	t.Direction = direction

	t.Station = station

	feed, err := decoder.Decode(accessKey, url)
	t.Feed = feed

	return &t, err
}

func (t *TrainFeed) GetArrivals() []Arrival {
	stopID := t.Station.GTFSStopID + t.Direction
	now := time.Now()
	arrivals := []Arrival{}
	for _, entity := range t.Feed.Entity {
		trip := entity.GetTripUpdate()
		if trip != nil {
			stopTimes := trip.StopTimeUpdate
			for _, s := range stopTimes {
				if *s.StopId == stopID {
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

					a := Arrival{}
					a.Train = route
					a.Secs = secs

					arrivals = append(arrivals, a)
				}
			}
		}
	}

	return arrivals
}

func PrintArrivals(arrivals []Arrival, name string) {
	fmt.Printf("STATION %s\n", name)

	if len(arrivals) == 0 {
		fmt.Println("No trains arriving at this station today")
		return
	}

	sort.Slice(arrivals, func(i, j int) bool { return arrivals[i].Secs < arrivals[j].Secs })
	for _, a := range arrivals {
		if a.Secs < 15 {
			fmt.Printf("%s now\n", a.Train)
		} else {
			fmt.Printf("%s %d mins\n", a.Train, a.Secs/60)
		}
	}
	fmt.Println()
}
