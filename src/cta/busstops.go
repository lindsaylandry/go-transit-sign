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

type KMLDescription struct {
  Document struct {
    Placemarks []struct {
      Description []byte `xml:"description"`
    } `xml:"Placemark"`
  } `xml:"Document>Folder"`
}

type BusStop struct {
  Name      string
  StopID    string
  PositionX float64
  PositionY float64
  Direction string
}

func GetBusStop(stopID string) (BusStop, error) {
	stop := BusStop{}
	stops, err := readBusStops("data/cta-bus-stations.kml")
	if err != nil {
		return stop, err
	}

	for _, s := range stops {
		if s.StopID == stopID {
			return s, nil
		}
	}

	return stop, fmt.Errorf("Could not find bus stop %s", stopID)
}

func readBusStops(filepath string) ([]BusStop, error) {
	busStops := []BusStop{}
	descriptions := KMLDescription{}
	data, err := os.ReadFile(filepath)
	if err != nil {
		return busStops, err
	}

	if err = xml.Unmarshal(data, &descriptions); err != nil {
		return busStops, err
	}

	for _, d := range descriptions.Document.Placemarks {
		bus := BusStop{}

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
			case "SYSTEMSTOP":
				bus.StopID = content[i+1]
			case "PUBLIC_NAME":
				bus.Name = content[i+1]
			case "DIR":
				bus.Direction = content[i+1]
			case "POINT_X":
				bus.PositionX, err = strconv.ParseFloat(content[i+1], 64)
			case "POINT_Y":
				bus.PositionY, err = strconv.ParseFloat(content[i+1], 64)
			default:
				continue
			}

			if err != nil {
				return busStops, err
			}
		}

		busStops = append(busStops, bus)
	}

	return busStops, err
}
