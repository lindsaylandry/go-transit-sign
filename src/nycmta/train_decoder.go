package nycmta

import (
	"strings"
)

type TrainDecoder struct {
  URL    string
  Trains []string
}

func GetAllMtaTrainDecoders() *[]TrainDecoder {
	f := []TrainDecoder{
		{URL: "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace", Trains: []string{"A", "C", "E"}},
		{URL: "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-g", Trains: []string{"G"}},
		{URL: "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-bdfm", Trains: []string{"B", "D", "F", "M"}},
		{URL: "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-jz", Trains: []string{"J", "Z"}},
		{URL: "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-nqrw", Trains: []string{"N", "Q", "R", "W"}},
		{URL: "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-l", Trains: []string{"L"}},
		{URL: "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs", Trains: []string{"1", "2", "3", "4", "5", "6", "7"}},
		{URL: "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-si", Trains: []string{"SI"}},
	}

	return &f
}

func GetMtaTrainDecoders(trains string) *[]TrainDecoder {
	f := GetAllMtaTrainDecoders()
	fd := []TrainDecoder{}

	// parse space-separated list of trains
	trns := strings.Split(trains, " ")

	for _, u := range *f {
		for _, t := range u.Trains {
			for _, tt := range trns {
				if tt == t {
					found := false
					for _, dd := range fd {
						if dd.URL == u.URL {
							found = true
						}
					}
					if !found {
						fd = append(fd, u)
						break
					}
				}
			}
		}
	}

	return &fd
}
