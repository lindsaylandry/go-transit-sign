package cta

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
)

type Station struct {
	StopID          int    `csv:"STOP_ID"`
	DirectionID     string `csv:"DIRECTION_ID"`
	StopName        string `csv:"STOP_NAME"`
	StationName     string `csv:"STATION_NAME"`
	StationDescName string `csv:"STATION_DESCRIPTIVE_NAME"`
	MapID           int    `csv:"MAP_ID"`
	ADA             bool   `csv:"ADA"`
	Red             bool   `csv:"RED"`
	Blue            bool   `csv:"BLUE"`
	Green           bool   `csv:"G"`
	Brown           bool   `csv:"BRN"`
	Purple          bool   `csv:"P"`
	Pexp            bool   `csv:"Pexp"`
	Yellow          bool   `csv:"Y"`
	Pink            bool   `csv:"Pnk"`
	Orange          bool   `csv:"O"`
	Location        string `csv:"Location"`
}

func GetStation(stopID int) (Station, error) {
	station := Station{}
	stations, err := readStations("data/cta-rail-stations.csv")
	if err != nil {
		return station, err
	}

	// Find station, return error if not found
	for _, s := range stations {
		if s.StopID == stopID {
			return s, nil
		}
	}

	return station, fmt.Errorf("Could not find station %d", stopID)
}

func readStations(filepath string) ([]Station, error) {
	stations := []Station{}
	f, err := os.Open(filepath)
	if err != nil {
		return stations, err
	}
	defer f.Close()

	err = gocsv.UnmarshalFile(f, &stations)

	return stations, err
}
