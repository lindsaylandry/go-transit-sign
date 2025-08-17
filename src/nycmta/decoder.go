package nycmta

import (
	"errors"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"strings"
)

type TrainDecoder struct {
	URL    string
	Trains []string
}

var BusFeedURL = "https://gtfsrt.prod.obanyc.com/tripUpdates"

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

func DecodeNYCMTA(k, url string) (*gtfs.FeedMessage, error) {
	feed := gtfs.FeedMessage{}
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &feed, err
	}

	req.Header.Add("x-api-key", k)
	resp, err := client.Do(req)
	if err != nil {
		return &feed, err
	}
	defer resp.Body.Close()

	// read response code
	// TODO: make more robust
	if resp.StatusCode >= 400 {
		return &feed, errors.New(http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &feed, err
	}

	err = proto.Unmarshal(body, &feed)
	return &feed, err
}
