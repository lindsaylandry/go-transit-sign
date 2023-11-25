package traininfo

import (
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"fmt"
	"time"

	"github.com/lindsaylandry/go-mta-train-sign/src/decoder"
	"github.com/lindsaylandry/go-mta-train-sign/src/stations"
)

type TrainInfo struct {
	Station stations.MtaStation
	Key string
	Direction string

	Feed gtfs.FeedMessage
}

func NewTrainInfo(station stations.MtaStation, accessKey, direction string) (*TrainInfo, error) {
	t := TrainInfo{}

	t.Key = accessKey
	t.Direction = direction

	t.Station = station

	feed, err := decoder.Decode(accessKey)
	t.Feed = feed

	return &t, err
}

func (t *TrainInfo) GetVehicleStopInfo() (error) {
	for _, entity := range t.Feed.Entity {
		vehicle := entity.GetVehicle()
		if vehicle != nil {
			stopId := vehicle.StopId
			msg := fmt.Sprintf("Stop ID: %s", *stopId)
			status := vehicle.CurrentStatus
			if status != nil {
				msg = fmt.Sprintf("%s %s", msg, *status)
			}	else {
				currentStatus := vehicle.CurrentStatus
				if currentStatus == nil {
					msg = fmt.Sprintf("%s IN_TRANSIT_TO", msg)
				} else {
					switch *currentStatus {
					case gtfs.VehiclePosition_INCOMING_AT:
						msg = fmt.Sprintf("%s INCOMING_AT", msg)
					case gtfs.VehiclePosition_STOPPED_AT:
						msg = fmt.Sprintf("%s STOPPED_AT", msg)
					case gtfs.VehiclePosition_IN_TRANSIT_TO:
						msg = fmt.Sprintf("%s IN_TRANSIT_TO", msg)
					// If current_status is missing IN_TRANSIT_TO is assumed.
					default:
						msg = fmt.Sprintf("%s IN_TRANSIT_TO", msg)
					}
				}
			}
			desc := vehicle.Vehicle
			if desc != nil {
				msg = fmt.Sprintf("Label: %s %s", *desc.Label, msg)
			}
			fmt.Printf("%s\n", msg)
		}
	}
	return nil
}

func (t *TrainInfo) GetTripUpdateInfo(stopId string) (error) {
	now := time.Now()
	for _, entity := range t.Feed.Entity {
		trip := entity.GetTripUpdate()
		if trip != nil {
			stopTimes := trip.StopTimeUpdate
			for _, time := range stopTimes {
				if *time.StopId == stopId {
					route := ""
					vehicle := trip.Trip
					if vehicle != nil {
						route = *vehicle.RouteId
					}
					secs := *time.Arrival.Time - now.Unix()
					mins := secs/60
					fmt.Printf("%s %s %d mins\n", route, *time.StopId, mins)
				}
			}
		}
  }

	return nil
}