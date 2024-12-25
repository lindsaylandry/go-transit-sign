package decoder

import (
	"strings"
)

func GetCTABusDecoder() *BusDecoder {
	f := BusDecoder{
		{URL: "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace"},
	}

	return &bd
}
