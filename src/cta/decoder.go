package cta

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type BusFeedMessage struct {
	BusTimeResponse struct {
		Error []struct {
			StopID  string `json:"stpid"`
			Message string `json:"msg"`
		} `json:"error"`
		Prd []struct {
			RouteDir           string `json:"rtdir"`
			Name               string `json:"rt"`
			PredictedCountdown string `json:"prdctdn"`
		} `json:"prd"`
	} `json:"bustime-response"`
}

type TrainFeedMessage struct {
	TrainTimeResponse struct {
		Timestamp string `json:"tmst"`
		Eta       []struct {
			Destination string `json:"destNm"`
			Name        string `json:"rt"`
			ArrivalTime string `json:"arrT"`
		} `json:"eta"`
		Error string `json:"errNm"`
	} `json:"ctatt"`
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

	slog.Debug(string(body))

	err = json.Unmarshal(body, &bf)

	jsonData, err := json.Marshal(bf)
	slog.Debug(string(jsonData))

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
	q.Add("outputType", "JSON")
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

	slog.Debug(string(body))

	err = json.Unmarshal(body, &tf)

	jsonData, err := json.Marshal(tf)
	slog.Debug(string(jsonData))

	return tf, err
}
