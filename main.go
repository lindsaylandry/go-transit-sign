package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tfk1410/go-rpi-rgb-led-matrix"

	"github.com/lindsaylandry/go-transit-sign/src/cta"
	"github.com/lindsaylandry/go-transit-sign/src/config"
	"github.com/lindsaylandry/go-transit-sign/src/signdata"
	"github.com/lindsaylandry/go-transit-sign/src/nycmta"
)

var stop, key, direction string
var train, led bool
var conf config.Config

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
		Use:   "testMatrix",
		Short: "Test LED Matrix",
		RunE: func(cmd *cobra.Command, args []string) error {
			return TestMatrix()
		},
	}

	rootCmd.AddCommand(nycMtaCmd)
	rootCmd.AddCommand(ctaCmd)
	rootCmd.AddCommand(testMatrix)

	rootCmd.PersistentFlags().StringVarP(&stop, "stop", "s", "D30", "stop to parse")
	//rootCmd.PersistentFlags().StringVarP(&key, "key", "k", "foobar", "API access key")
	rootCmd.PersistentFlags().BoolVarP(&train, "train", "t", true, "train or bus (train=true, bus=false)")
	rootCmd.PersistentFlags().StringVarP(&direction, "direction", "d", "N", "direction (trains only)")
	rootCmd.PersistentFlags().BoolVarP(&led, "led", "l", false, "output to led matrix")

	conf, err := config.NewConfig()
  if err != nil {
    panic(err)
  }
	fmt.Println(conf)

	err = rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func CTA() error {
	timezone := "America/Chicago"
	stp, err := cta.GetBusStop(stop)
	if err != nil {
		return err
	}

	// TODO: add cta trains
	bf, err := cta.NewBusFeed(stp, conf.CTA.Bus.APIKey, timezone)
	if err != nil {
		return err
	}

	sd, err := signdata.NewSignData()
	if err != nil {
		return err
	}
	sd.Canvas = rgbmatrix.NewCanvas(sd.Matrix)
	defer sd.Canvas.Close()

	for {
		arrivals, err := bf.GetArrivals()
		if err != nil {
			return err
		}

		if led {
			// Print all arrivals
			err = sd.PrintArrivals(arrivals, stp.Name, stp.Direction)
			if err != nil {
				return err
			}
		} else {
			signdata.PrintArrivalsToStdout(arrivals, stp.Name, stp.Direction)
		}

		time.Sleep(5 * time.Second)
	}
	return nil
}

func NYCMTA() error {
	//timezone := "America/New_York"
	station, err := nycmta.GetStation(stop)
	if err != nil {
		return err
	}

	// Get subway feeds from station trains
	feeds := nycmta.GetMtaTrainDecoders(station.DaytimeRoutes)

	sd, err := signdata.NewSignData()
	if err != nil {
		return err
	}
	sd.Canvas = rgbmatrix.NewCanvas(sd.Matrix)
	defer sd.Canvas.Close()

	for {
		arrivals := []signdata.Arrival{}
		for _, f := range *feeds {
			// TODO: get buses
			t, err := nycmta.NewTrainFeed(station, conf.NYCMTA.APIKey, direction, f.URL)
			if err != nil {
				return err
			}

			arr := t.GetArrivals()
			arrivals = append(arrivals, arr...)
		}

		// Print all arrivals
		if led {
			err = sd.PrintArrivals(arrivals, station.StopName, direction)
			if err != nil {
				return err
			}
		} else {
			signdata.PrintArrivalsToStdout(arrivals, station.StopName, direction)
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

func TestMatrix() error {
	sd, err := signdata.NewSignData()
	if err != nil {
		return err
	}

	sd.Canvas = rgbmatrix.NewCanvas(sd.Matrix)
	defer sd.Canvas.Close()

	return sd.WriteTestMatrix()
}
