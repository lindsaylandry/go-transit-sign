package signdata

import (
	"fmt"
	"sort"

	"github.com/lindsaylandry/go-transit-sign/src/feed"
	"github.com/lindsaylandry/go-transit-sign/src/signdata/writer"
)

func PrintArrivalsToStdout(arrivals []feed.Arrival, name string) {
	fmt.Println(name)

	if len(arrivals) == 0 {
		fmt.Println("No trains arriving at this station today")
		return
	}

	sort.Slice(arrivals, func(i, j int) bool { return arrivals[i].Secs < arrivals[j].Secs })
	for _, a := range arrivals {
		if a.Secs < 15 {
			fmt.Printf("%s now\n", a.Label)
		} else {
			fmt.Printf("%s %d min\n", a.Label, a.Secs/60)
		}
	}
	fmt.Println()
}

func PrintArrivals(arrivals []feed.Arrival, name string) error {
  assembly, err := writer.CreateVisualString(name)
	if err != nil {
		return err
	}
	fmt.Println(len(assembly[0]))
  printAssembly(assembly)

	if len(arrivals) == 0 {
    fmt.Println("No trains arriving at this station today")
    return nil
  }

  sort.Slice(arrivals, func(i, j int) bool { return arrivals[i].Secs < arrivals[j].Secs })
	var str string
  for _, a := range arrivals {
    if a.Secs < 15 {
      str = "now"
    } else {
      str = fmt.Sprintf("%d min", a.Secs/60)
    }
		assembly, err = writer.CreateVisualNextArrival(a.Label, str, 64)
  	printAssembly(assembly)
  }

	return nil
}

func printAssembly(assembly [][]uint8) {
	for _, a := range assembly {
    for _, b := range a {
      if b == 0 {
        fmt.Printf(" ")
      } else {
        fmt.Printf("8")
      }
    }
    fmt.Printf("\n")
  }
}
