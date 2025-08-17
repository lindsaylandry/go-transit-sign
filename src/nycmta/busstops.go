package nycmta

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
)

type BusStop struct {
	StopID        int    `csv:"stop_id"`
	StopName      string `csv:"stop_name"`
	Description   string `csv:"desc"`
	StopLat       string `csv:"stop_lat"`
	StopLon       string `csv:"stop_lon"`
	ZoneID        string `csv:"zone_id"`
	StopURL       string `csv:"stop_url"`
	LocationType  int    `csv:"location_type"`
	ParentStation string `csv:"parent_station"`
}

func GetBusStops(stopIDs []int) ([]BusStop, error) {
	busstops := []BusStop{}
	stps, err := readBusStops("data/nyc-busstops.csv")
	if err != nil {
		return busstops, err
	}

	// Find station, return error if not found
	for _, id := range stopIDs {
		found := false
		for _, s := range stps {
			if s.StopID == id {
				busstops = append(busstops, s)
				found = true
				break
			}
		}
		if !found {
			return busstops, fmt.Errorf("Could not find station %d", id)
		}
	}

	return busstops, nil
}

func readBusStops(filepath string) ([]BusStop, error) {
	busstops := []BusStop{}
	f, err := os.Open(filepath)
	if err != nil {
		return busstops, err
	}
	defer f.Close()

	err = gocsv.UnmarshalFile(f, &busstops)

	return busstops, err
}
