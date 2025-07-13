package nycmta

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
)

type Station struct {
	StationID           int     `csv:"Station ID"`
	ComplexID           int     `csv:"Complex ID"`
	GTFSStopID          string  `csv:"GTFS Stop ID"`
	Division            string  `csv:"Division"`
	Line                string  `csv:"Line"`
	StopName            string  `csv:"Stop Name"`
	Borough             string  `csv:"Borough"`
	DaytimeRoutes       string  `csv:"Daytime Routes"`
	Structure           string  `csv:"Structure"`
	GTFSLatitude        float64 `csv:"GTFS Latitude"`
	GTFSLongitude       float64 `csv:"GTFS Longitude"`
	NorthDirectionLabel string  `csv:"North Direction Label"`
	SouthDirectionLabel string  `csv:"South Direction Label"`
	ADA                 int     `csv:"ADA"`
	ADADirectionNotes   string  `csv:"ADA Direction Notes"`
	ADANB               int     `csv:"ADA NB"`
	ADASB               int     `csv:"ADA SB"`
	CapitalOutageNB     string  `csv:"Capital Outage NB"`
	CapitalOutageSB     string  `csv:"Capital Outage SB"`
}

func GetStations(stopIDs []string) ([]Station, error) {
	stations := []Station{}
	stns, err := readStations("data/nyc-subway-stations.csv")
	if err != nil {
		return stations, err
	}

	// Find station, return error if not found
	for _, id := range stopIDs {
		found := false
		for _, s := range stns {
			if s.GTFSStopID == id {
				stations = append(stations, s)
				found = true
				break
			}
		}
		if !found {
			return stations, fmt.Errorf("Could not find station %d", id)
		}
	}

	return stations, nil
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
