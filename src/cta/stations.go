package cta

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type Station struct {
  Name      string
	RailLine  string
  StationID string
  PositionX float64
  PositionY float64
  Direction string
}

func GetStation(stationID string) (Station, error) {
	station := Station{}
	stations, err := readStations("data/cta-rail-stations.kml")
	if err != nil {
		return station, err
	}

	for _, s := range stations {
		if s.StationID == stationID {
			return s, nil
		}
	}

	return station, fmt.Errorf("Could not find station %s", stationID)
}

func readStations(filepath string) ([]Station, error) {
	stations := []Station{}
	descriptions := KMLDescription{}
	data, err := os.ReadFile(filepath)
	if err != nil {
		return stations, err
	}

	if err = xml.Unmarshal(data, &descriptions); err != nil {
		return stations, err
	}

	for _, d := range descriptions.Document.Placemarks {
		rail := Station{}

		z := html.NewTokenizer(bytes.NewReader(d.Description))
		content := []string{}

		for z.Token().Data != "html" {
			tt := z.Next()
			if tt == html.StartTagToken {
				t := z.Token()
				if t.Data == "td" {
					inner := z.Next()
					if inner == html.TextToken {
						text := (string)(z.Text())
						if text != "" {
							t := strings.TrimSpace(text)
							content = append(content, t)
						}
					}
				}
			}
		}

		for i, v := range content {
			switch v {
			case "Station ID":
				rail.StationID = content[i+1]
			case "Station Name":
				rail.Name = content[i+1]
			case "Rail Line":
				rail.RailLine = content[i+1]
			case "DIR":
				rail.Direction = content[i+1]
			case "POINT_X":
				rail.PositionX, err = strconv.ParseFloat(content[i+1], 64)
			case "POINT_Y":
				rail.PositionY, err = strconv.ParseFloat(content[i+1], 64)
			default:
				continue
			}

			if err != nil {
				return stations, err
			}
		}

		stations = append(stations, rail)
	}

	return stations, err
}
