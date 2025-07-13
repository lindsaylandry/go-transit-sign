package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

	switch config.Log {
	case 0:
		slog.SetLogLoggerLevel(slog.Level(-8))
	case 1:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case 2:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case 3:
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case 4:
		slog.SetLogLoggerLevel(slog.LevelError)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
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

	go func() {
		for {
			for _, f := range bfs {
				arrivals, err := f.GetArrivals()
				if err != nil {
					panic(err)
				}

				if err := printArrivals(sd, arrivals, f.BusStop.Name, f.BusStop.Direction); err != nil {
					panic(err)
				}
			}

			for _, f := range tfs {
				arrivals, err := f.GetArrivals()
				if err != nil {
					panic(err)
				}

				if err := printArrivals(sd, arrivals, f.Station.StopName, f.Station.DirectionID); err != nil {
					panic(err)
				}
			}
		}
	}()
	s := <-sigChan
	fmt.Printf("Received signal in main: %s. Shutting down\n", s)
	return nil
}

func NYCMTA() error {
	stations, err := nycmta.GetStations(conf.NYCMTA.Train.StopIDs)
	if err != nil {
		return err
	}

	// Get subway feeds from station trains
	// TODO: get feeds from all stations
	feeds := nycmta.GetMtaTrainDecoders(stations[0].DaytimeRoutes)

	tfs := []*nycmta.TrainFeed{}
	for _, f := range *feeds {
    tf := nycmta.NewTrainFeed(stations[0], conf.NYCMTA.APIKey, direction, f.URL)
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

	go func() {
		for {
			arrivals := []signdata.Arrival{}
			for _, tf := range tfs {
				arr, err := tf.GetArrivals()
				if err != nil {
					panic(err)
				}
				arrivals = append(arrivals, arr...)
			
				if err := printArrivals(sd, arrivals, tf.Station.StopName, direction); err != nil {
					panic(err)
				}
			}
		}
	}()
	s := <-sigChan
	fmt.Printf("Received signal in main: %s. Shutting down\n", s)
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

func printArrivals(sd *signdata.SignData, arrivals []signdata.Arrival, name, direction string) error {
	if led {
		if err := sd.PrintArrivals(arrivals, name, direction); err != nil {
			return err
		}
	} else {
		signdata.PrintArrivalsToStdout(arrivals, name, direction)
	}
	return nil
}
