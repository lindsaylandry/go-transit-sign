package signdata

import (
	"fmt"
	"sort"

	"github.com/lindsaylandry/go-transit-sign/src/feed"
)

func PrintArrivals(arrivals []feed.Arrival, name string) {
	fmt.Printf("STATION %s\n", name)

	if len(arrivals) == 0 {
		fmt.Println("No trains arriving at this station today")
		return
	}

	sort.Slice(arrivals, func(i, j int) bool { return arrivals[i].Secs < arrivals[j].Secs })
	for _, a := range arrivals {
		if a.Secs < 15 {
			fmt.Printf("%s now\n", a.Train)
		} else {
			fmt.Printf("%s %d mins\n", a.Train, a.Secs/60)
		}
	}
	fmt.Println()
}
