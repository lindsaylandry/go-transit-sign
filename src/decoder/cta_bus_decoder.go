package decoder

type CTABusFeedMessage struct {
	BusTimeResponse struct {
		Prd []struct {
			Timestamp     string `json:"tmstmp"`
			RouteDir      string `json:"rtdir"`
			Name          string `json:"rt"`
			PredictedTime string `json:"prdtm"`
		} `json:"prd"`
	} `json:"bustime-response"`
}

var CTABusFeedURL = "http://www.ctabustracker.com/bustime/api/v2/getpredictions"
