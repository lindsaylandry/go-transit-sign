package stations

type MtaStation struct {
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
