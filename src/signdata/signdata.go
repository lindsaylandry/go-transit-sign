package signdata

import (
	"fmt"
	"image/color"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tfk1410/go-rpi-rgb-led-matrix"

	"github.com/lindsaylandry/go-transit-sign/src/signdata/writer"
)

type SignData struct {
	Visual      [32][64]color.RGBA
	Matrix      rgbmatrix.Matrix
	Canvas      *rgbmatrix.Canvas
	MaxArrivals int
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

	// Max number of transit arrival times that fit on screen
	sd.MaxArrivals = len(sd.Visual)/6 - 2

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

	if direction != "" {
		dir := getDirection(direction)
		fmt.Println(dir)
		fmt.Println()
	}

	time.Sleep(5 * time.Second)
}

func (sd *SignData) PrintArrivals(arrivals []Arrival, name, direction string) error {
	// Reset canvas to black
	for i, c := range sd.Visual {
		for j := range c {
			sd.Visual[i][j] = color.RGBA{0, 0, 0, 255}
		}
	}

	if len(arrivals) == 0 {
		arr := Arrival{Label: "None", Secs: -1}
		arrivals = []Arrival{arr}
	}

	sort.Slice(arrivals, func(i, j int) bool { return arrivals[i].Secs < arrivals[j].Secs })
	arrs := 0
	for i, a := range arrivals {
		if arrs < sd.MaxArrivals {
			mins := strconv.FormatInt(a.Secs/60, 10)
			str := fmt.Sprintf("%s min", mins)
			slog.Debug(a.Label, mins, "min")
			if a.Secs == -1 {
				str = "none"
			} else if a.Secs <= 30 {
				str = "now"
			}
			label, color := normalizeTrain(a.Label)
			assembly, timeIndex, err := writer.CreateVisualNextArrival(label, str, 64)
			if err != nil {
				return err
			}
			// TODO: separately add station/bus and arrival time to canvas for more colors
			sd.addArrival(assembly, color, timeIndex, i)
			arrs += 1
		}
	}

	dir := getDirection(direction)

	assembly, err := writer.CreateVisualString(dir)
	if err != nil {
		return err
	}
	sd.addDirection(assembly)

	// Title going last (scroll through title)
	titleAssembly, err := writer.CreateVisualString(name)
	if err != nil {
		return err
	}
	index := -1
	if len(titleAssembly[0]) > len(sd.Visual[0]) {
		index = 0
		// add as many spaces as width of canvas
		for i := 0; i < len(sd.Visual[0]); i += 2 {
			name = name + " "
		}
		titleAssembly, err = writer.CreateVisualString(name)
		if err != nil {
			return err
		}
	}

	// Non-scrolling title
	if index == -1 {
		titleAssembly, err = writer.CreateVisualString(name)
		sd.addTitle(titleAssembly, &index)
		if err != nil {
			return err
		}
		if err := sd.WriteToMatrix(); err != nil {
			return err
		}
	}

	// Scrolling title
	for index >= 0 {
		sd.addTitle(titleAssembly, &index)
		if err := sd.WriteToMatrix(); err != nil {
			return err
		}

		if index == 1 {
			time.Sleep(1 * time.Second)
		}

		time.Sleep(10 * time.Microsecond)

		if index == -1 {
			sd.addTitle(titleAssembly, &index)
			if err := sd.WriteToMatrix(); err != nil {
				return err
			}
		}
	}
	time.Sleep(5 * time.Second)

	return nil
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

func (sd *SignData) addTitle(title [][]uint8, index *int) {
	normIndex := 0
	if *index > 0 {
		normIndex = *index
	}

	for i, a := range title {
		for j := range a {
			if len(sd.Visual[0]) > j {
				if a[j+normIndex] > 0 {
					sd.Visual[i][j] = color.RGBA{0, 255, 255, 255}
				} else {
					sd.Visual[i][j] = color.RGBA{0, 0, 0, 255}
				}
			}
		}
	}
	if *index >= 0 {
		*index += 1
		if *index >= len(title[0])-len(sd.Visual[0]) {
			*index = -1
		}
	}
}

func (sd *SignData) addArrival(arrival [][]uint8, col color.RGBA, timeIndex, index int) {
	// Title takes top 6 pixel rows
	start := 6 * (index + 1)

	for i, a := range arrival {
		for j, b := range a {
			// Truncate for now
			if len(sd.Visual[0]) > j && len(sd.Visual) > i+start+1 {
				if b > 0 {
					if j < timeIndex {
						sd.Visual[i+start][j] = col
					} else {
						sd.Visual[i+start][j] = color.RGBA{255, 255, 255, 255}
					}
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

func normalizeTrain(train string) (string, color.RGBA) {
	var trn string
	var col color.RGBA
	switch strings.ToUpper(train) {
	case "ORG":
		trn = "Orange"
		col = color.RGBA{255, 165, 0, 255}
	case "PNK":
		trn = "Pink"
		col = color.RGBA{255, 209, 220, 255}
	case "GRN":
		trn = "Green"
		col = color.RGBA{0, 255, 0, 255}
	case "BRN":
		trn = "Brown"
		col = color.RGBA{150, 75, 0, 255}
	default:
		trn = train
		col = color.RGBA{255, 255, 255, 255}
	}

	return trn, col
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
