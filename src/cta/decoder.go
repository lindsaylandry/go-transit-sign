package cta

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type BusFeedMessage struct {
	BusTimeResponse struct {
		Error string `json:"error"`
		Prd   []struct {
			RouteDir           string `json:"rtdir"`
			Name               string `json:"rt"`
			PredictedCountdown string `json:"prdctdn"`
		} `json:"prd"`
	} `json:"bustime-response"`
}

type TrainFeedMessage struct {
	TrainTimeResponse struct {
		Eta []struct {
			RouteDir           string `xml:"rtdir"`
			Name               string `xml:"rt"`
			PredictedCountdown string `xml:"prdctdn"`
		} `xml:"eta"`
	} `xml:"ctatt"`
}

var BusFeedURL = "http://www.ctabustracker.com/bustime/api/v2/getpredictions"
var TrainFeedURL = "http://lapi.transitchicago.com/api/1.0/ttarrivals.aspx"

func DecodeBus(k string, stopID int, url string) (BusFeedMessage, error) {
	bf := BusFeedMessage{}
	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return bf, err
	}

	q := req.URL.Query()
	q.Add("key", k)
	q.Add("format", "json")
	q.Add("stpid", strconv.Itoa(stopID))
	q.Add("top", strconv.Itoa(5))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return bf, err
	}
	defer resp.Body.Close()

	// read response code
	if resp.StatusCode >= 400 {
		return bf, errors.New(http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return bf, err
	}

	err = json.Unmarshal(body, &bf)

	return bf, err
}

func DecodeTrain(k string, stopID int, url string) (TrainFeedMessage, error) {
	tf := TrainFeedMessage{}
	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return tf, err
	}

	q := req.URL.Query()
	q.Add("key", k)
	q.Add("max", "5")
	q.Add("stpid", strconv.Itoa(stopID))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return tf, err
	}
	defer resp.Body.Close()

	// read response code
	// TODO: make more robust
	if resp.StatusCode >= 400 {
		return tf, errors.New(http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf, err
	}

	err = xml.Unmarshal(body, &tf)

	return tf, err
}
