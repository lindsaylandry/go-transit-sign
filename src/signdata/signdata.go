package signdata

import (
	"fmt"
	"image"
	"sort"

	"github.com/lindsaylandry/go-transit-sign/src/feed"
	"github.com/lindsaylandry/go-transit-sign/src/signdata/writer"
)

type SignData struct {
	Visual [32][64]uint8
	Image  image.Image
}

func NewSignData() *SignData {
	sd := SignData{}

	sd.Image = image.NewRGBA(image.Rect(0, 0, 64, 32))

	return &sd
}

func PrintArrivalsToStdout(arrivals []feed.Arrival, name, direction string) {
	fmt.Println(name)

	if len(arrivals) == 0 {
		fmt.Println("No trains arriving at this station today")
		return
	}

	sort.Slice(arrivals, func(i, j int) bool { return arrivals[i].Secs < arrivals[j].Secs })
	for _, a := range arrivals {
		if a.Secs <= 30 {
			fmt.Printf("%s now\n", a.Label)
		} else {
			fmt.Printf("%s %d min\n", a.Label, a.Secs/60)
		}
	}
	dir := getDirection(direction)

	fmt.Println(dir)
	fmt.Println()
}

func (sd *SignData) PrintArrivals(arrivals []feed.Arrival, name, direction string) error {
	assembly, err := writer.CreateVisualString(name)
	if err != nil {
		return err
	}
	sd.addTitle(assembly)

	if len(arrivals) == 0 {
		fmt.Println("No trains arriving at this station today")
		return nil
	}

	sort.Slice(arrivals, func(i, j int) bool { return arrivals[i].Secs < arrivals[j].Secs })
	var str string
	for i, a := range arrivals {
		if a.Secs < 30 {
			str = "now"
		} else {
			str = fmt.Sprintf("%d min", a.Secs/60)
		}
		assembly, err = writer.CreateVisualNextArrival(a.Label, str, 64)
		if err != nil {
			return err
		}
		sd.addArrival(assembly, i)
	}

	dir := getDirection(direction)

	assembly, err = writer.CreateVisualString(dir)
	if err != nil {
		return err
	}
	sd.addDirection(assembly)

	sd.printMatrix()

	return nil
}

func (sd *SignData) addTitle(title [][]uint8) {
	for i, a := range title {
		for j, b := range a {
			// Truncate for now
			if len(sd.Visual[0]) > j {
				sd.Visual[i][j] = b
			}
		}
	}
}

func (sd *SignData) addArrival(arrival [][]uint8, index int) {
	// Title takes top 6 pixel rows
	start := 6 * (index + 1)

	for i, a := range arrival {
		for j, b := range a {
			// Truncate for now
			if len(sd.Visual[0]) > j && len(sd.Visual[0]) > i+start {
				sd.Visual[i+start][j] = b
			}
		}
	}
}

func (sd *SignData) addDirection(direction [][]uint8) {
	// Direction takes bottom 5 pixel rows
	start := len(sd.Visual) - 6

	for i, a := range direction {
		for j, b := range a {
			// Truncate for now
			if len(sd.Visual[0]) > j && len(sd.Visual) > i+start {
				sd.Visual[i+start][j] = b
			}
		}
	}
}

func (sd *SignData) printMatrix() {
	for _, a := range sd.Visual {
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

func getDirection(direction string) string {
	var dir string
	switch direction {
	case "N", "NB", "NORTH", "North":
		dir = "Northbound"
	case "S", "SB", "SOUTH", "South":
		dir = "Southbound"
	case "W", "WB", "WEST", "West":
		dir = "Westbound"
	case "E", "EB", "EAST", "East":
		dir = "Eastbound"
	case "NW", "NWB", "NORTHWEST", "Northwest":
		dir = "Northwestbound"
	case "SW", "SWB", "SOUTHWEST", "Southwest":
		dir = "Southwestbound"
	case "NE", "NEB", "NORTHEAST", "Northeast":
		dir = "Northeastbound"
	case "SE", "SEB", "SOUTHEAST", "Southeast":
		dir = "Southeastbound"
	default:
		dir = direction
	}

	return dir
}
