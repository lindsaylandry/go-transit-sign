package decoder

import (
	"errors"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
)

func Decode(k, url string) (gtfs.FeedMessage, error) {
	client := &http.Client{}
	feed := gtfs.FeedMessage{}

	// TODO: get URL from train line

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("x-api-key", k)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return feed, err
	}

	// read response code
	// TODO: make more robust
	if resp.StatusCode >= 400 {
		return feed, errors.New(http.StatusText(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return feed, err
	}

	err = proto.Unmarshal(body, &feed)
	return feed, err
}
