package signdata

import (
	"fmt"
	"image/color"
	"sort"
	"strings"

	"github.com/tfk1410/go-rpi-rgb-led-matrix"

	"github.com/lindsaylandry/go-transit-sign/src/signdata/writer"
)

type SignData struct {
	Visual [32][64]color.RGBA
	Matrix rgbmatrix.Matrix
	Canvas *rgbmatrix.Canvas
}

func NewSignData() (*SignData, error) {
	sd := SignData{}

	config := &rgbmatrix.DefaultConfig
	config.Rows = len(sd.Visual)
	config.Cols = len(sd.Visual[0])
	config.Parallel = 1
	config.ChainLength = 1
	config.Brightness = 50
	config.HardwareMapping = "adafruit-hat"
	config.ShowRefreshRate = false
	config.InverseColors = false
	config.DisableHardwarePulsing = false

	m, err := rgbmatrix.NewRGBLedMatrix(config)
	if err != nil {
		return &sd, err
	}

	sd.Matrix = m

	return &sd, nil
}

func PrintArrivalsToStdout(arrivals []Arrival, name, direction string) {
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

func (sd *SignData) PrintArrivals(arrivals []Arrival, name, direction string) error {
	// Reset canvas to black
	for i, c := range sd.Visual {
		for j, _ := range c {
			sd.Visual[i][j] = color.RGBA{0, 0, 0, 255}
		}
	}

	assembly, err := writer.CreateVisualString(name)
	if err != nil {
		return err
	}
	sd.addTitle(assembly)

	if len(arrivals) == 0 {
		fmt.Println("None")
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

	return sd.WriteToMatrix()
}

func (sd *SignData) WriteToMatrix() error {
	bounds := sd.Canvas.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			sd.Canvas.Set(x, y, sd.Visual[y][x])
		}
	}
	return sd.Canvas.Render()
}

func (sd *SignData) WriteTestMatrix() error {
	bounds := sd.Canvas.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			sd.Canvas.Set(x, y, color.RGBA{255, 0, 0, 255})
			err := sd.Canvas.Render()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (sd *SignData) addTitle(title [][]uint8) {
	for i, a := range title {
		for j, b := range a {
			// Truncate for now
			if len(sd.Visual[0]) > j {
				if b > 0 {
					sd.Visual[i][j] = color.RGBA{0, 0, 255, 255}
				} else {
					sd.Visual[i][j] = color.RGBA{0, 0, 0, 255}
				}
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
			if len(sd.Visual[0]) > j && len(sd.Visual) > i+start+1 {
				if b > 0 {
					sd.Visual[i+start][j] = color.RGBA{255, 255, 255, 255}
				} else {
					sd.Visual[i+start][j] = color.RGBA{0, 0, 0, 255}
				}
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
			if len(sd.Visual[0]) > j && len(sd.Visual) > i+start+1 {
				if b > 0 {
					sd.Visual[i+start][j] = color.RGBA{255, 0, 0, 255}
				} else {
					sd.Visual[i+start][j] = color.RGBA{0, 0, 0, 255}
				}
			}
		}
	}
}

func getDirection(direction string) string {
	var dir string
	switch strings.ToUpper(direction) {
	case "N", "NB", "NORTH":
		dir = "Northbound"
	case "S", "SB", "SOUTH":
		dir = "Southbound"
	case "W", "WB", "WEST":
		dir = "Westbound"
	case "E", "EB", "EAST":
		dir = "Eastbound"
	case "NW", "NWB", "NORTHWEST":
		dir = "NW-bound"
	case "SW", "SWB", "SOUTHWEST":
		dir = "SW-bound"
	case "NE", "NEB", "NORTHEAST":
		dir = "NE-bound"
	case "SE", "SEB", "SOUTHEAST":
		dir = "SE-bound"
	default:
		dir = direction
	}

	return dir
}
