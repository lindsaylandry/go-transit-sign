package busstops

type KMLDescription struct {
	Document struct {
		Placemarks []struct {
			Description []byte `xml:"description"`
		} `xml:"Placemark"`
	} `xml:"Document>Folder"`
}

type CTABusStop struct {
	Name      string
	StopID    string
	PositionX float64
	PositionY float64
	Direction string
}
