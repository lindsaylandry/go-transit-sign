package main

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/lindsaylandry/go-transit-sign/src/busstops"
	"github.com/lindsaylandry/go-transit-sign/src/decoder"
	"github.com/lindsaylandry/go-transit-sign/src/feed"
	"github.com/lindsaylandry/go-transit-sign/src/signdata"
	"github.com/lindsaylandry/go-transit-sign/src/signdata/ledMatrix"
	"github.com/lindsaylandry/go-transit-sign/src/stations"
)

var stop, key, direction string
var cont, train, led bool

func main() {
	rootCmd := &cobra.Command{
		Use:   "transit-sign",
		Short: "Run transit sign",
	}

	nycMtaCmd := &cobra.Command{
		Use:   "nyc-mta",
		Short: "Run NYC MTA data",
		RunE: func(cmd *cobra.Command, args []string) error {
			return NYCMTA()
		},
	}

	ctaCmd := &cobra.Command{
		Use:   "cta",
		Short: "Run CTA data",
		RunE: func(cmd *cobra.Command, args []string) error {
			return CTA()
		},
	}

	testMatrix := &cobra.Command{
		Use: "testMatrix",
		Short: "Test LED Matrix",
		RunE: func(cmd *cobra.Command, args []string) error {
      return TestMatrix()
    }, 
	}

	rootCmd.AddCommand(nycMtaCmd)
	rootCmd.AddCommand(ctaCmd)
	rootCmd.AddCommand(testMatrix)

	rootCmd.PersistentFlags().StringVarP(&stop, "stop", "s", "D30", "stop to parse")
	rootCmd.PersistentFlags().StringVarP(&key, "key", "k", "foobar", "API access key")
	rootCmd.PersistentFlags().BoolVarP(&cont, "continue", "c", true, "continue printing arrivals")
	rootCmd.PersistentFlags().BoolVarP(&train, "train", "t", true, "train or bus (train=true, bus=false)")
	rootCmd.PersistentFlags().StringVarP(&direction, "direction", "d", "N", "direction (trains only)")
	rootCmd.PersistentFlags().BoolVarP(&led, "led", "l", false, "output to led matrix")

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func CTA() error {
	stp, err := busstops.GetBusStop(stop)
	if err != nil {
		return err
	}

	for {
		bf, err := feed.NewBusFeed(stp, key)
		if err != nil {
			return err
		}

		arrivals, err := bf.GetArrivals()
		if err != nil {
			return err
		}

		// Print all arrivals
		if led {
			sd := signdata.NewSignData()
			err = sd.PrintArrivals(arrivals, stp.Name, stp.Direction)
			if err != nil {
				return err
			}
		} else {
			signdata.PrintArrivalsToStdout(arrivals, stp.Name, stp.Direction)
		}

		if !cont {
			break
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

func NYCMTA() error {
	station, err := stations.GetStation(stop)
	if err != nil {
		return err
	}

	// Get subway feeds from station trains
	feeds := decoder.GetMtaTrainDecoders(station.DaytimeRoutes)

	for {
		arrivals := []feed.Arrival{}
		for _, f := range *feeds {
			t, err := feed.NewTrainFeed(station, key, direction, f.URL)
			if err != nil {
				return err
			}

			arr := t.GetArrivals()
			arrivals = append(arrivals, arr...)
		}

		// Print all arrivals
		if led {
			sd := signdata.NewSignData()
			err := sd.PrintArrivals(arrivals, station.StopName, direction)
			if err != nil {
				return err
			}
		} else {
			signdata.PrintArrivalsToStdout(arrivals, station.StopName, direction)
		}

		if !cont {
			break
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

func TestMatrix() error {
	return ledMatrix.WriteToMatrix()
}
