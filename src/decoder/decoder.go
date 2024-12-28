package decoder

import (
	"encoding/json"
	"errors"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
)

func DecodeNYCMTA(k, url string) (gtfs.FeedMessage, error) {
	feed := gtfs.FeedMessage{}
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("x-api-key", k)
	resp, err := client.Do(req)
	if err != nil {
		return feed, err
	}
	defer resp.Body.Close()

	// read response code
	// TODO: make more robust
	if resp.StatusCode >= 400 {
		return feed, errors.New(http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return feed, err
	}

	err = proto.Unmarshal(body, &feed)
	return feed, err
}

func DecodeCTA(k, stopID, url string) (CTABusFeedMessage, error) {
	bf := CTABusFeedMessage{}
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	q := req.URL.Query()
	q.Add("key", k)
	q.Add("format", "json")
	q.Add("stpid", stopID)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return bf, err
	}
	defer resp.Body.Close()

	// read response code
	// TODO: make more robust
	if resp.StatusCode >= 400 {
		return bf, errors.New(http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return bf, err
	}

	err = json.Unmarshal(body, &bf)

	return bf, nil
}
