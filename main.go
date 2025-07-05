package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/tfk1410/go-rpi-rgb-led-matrix"

	"github.com/lindsaylandry/go-transit-sign/src/config"
	"github.com/lindsaylandry/go-transit-sign/src/cta"
	"github.com/lindsaylandry/go-transit-sign/src/nycmta"
	"github.com/lindsaylandry/go-transit-sign/src/signdata"
)

var direction string
var led bool
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

	rootCmd.PersistentFlags().StringVarP(&direction, "direction", "d", "N", "direction (trains only)")
	rootCmd.PersistentFlags().BoolVarP(&led, "led", "l", false, "output to led matrix")

	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	conf = *config

	if conf.Emulate {
		os.Setenv("MATRIX_EMULATOR", "1")
	}

	err = rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func CTA() error {
	stps, err := cta.GetBusStops(conf.CTA.Bus.StopIDs)
	if err != nil {
		return err
	}

	stns := []cta.Station{}
	for _, stn := range conf.CTA.Train.StopIDs {
		stn, err := cta.GetStation(stn)
		if err != nil {
			return err
		}
		stns = append(stns, stn)
	}

	bfs := []*cta.BusFeed{}
	for _, s := range stps {
		bf := cta.NewBusFeed(s, conf.CTA.Bus.APIKey)
		bfs = append(bfs, bf)
	}

	tfs := []*cta.TrainFeed{}
	for _, s := range stns {
		tf := cta.NewTrainFeed(s, conf.CTA.Train.APIKey)
		tfs = append(tfs, tf)
	}

	sd, err := signdata.NewSignData()
	if err != nil {
		return err
	}
	sd.Canvas = rgbmatrix.NewCanvas(sd.Matrix)
	defer sd.Canvas.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() error {
		for {
			for _, f := range bfs {
				arrivals, err := f.GetArrivals()
				if err != nil {
					return err
				}

				if led {
					// Print all arrivals
					err = sd.PrintArrivals(arrivals, f.BusStop.Name, f.BusStop.Direction)
					if err != nil {
						return err
					}
				} else {
					signdata.PrintArrivalsToStdout(arrivals, f.BusStop.Name, f.BusStop.Direction)
				}

				time.Sleep(5 * time.Second)
			}

			for _, f := range tfs {
				arrivals, err := f.GetArrivals()
				if err != nil {
					return err
				}

				if led {
					// Print all arrivals
					err = sd.PrintArrivals(arrivals, f.Station.StopName, f.Station.DirectionID)
					if err != nil {
						return err
					}
				} else {
					signdata.PrintArrivalsToStdout(arrivals, f.Station.StopName, f.Station.DirectionID)
				}

				time.Sleep(5 * time.Second)
			}
		}
	}()
	s := <-sigChan
	fmt.Printf("Received signal in main: %s. Shutting down...\n", s)
	return nil
}

func NYCMTA() error {
	//timezone := "America/New_York"
	station, err := nycmta.GetStation(conf.NYCMTA.Train.StopIDs[0])
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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() error {
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
	}()
	s := <-sigChan
	fmt.Printf("Received signal in main: %s. Shutting down...\n", s)
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
