package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/lindsaylandry/go-mta-train-sign/src/busstops"
	"github.com/lindsaylandry/go-mta-train-sign/src/decoder"
	"github.com/lindsaylandry/go-mta-train-sign/src/feed"
	"github.com/lindsaylandry/go-mta-train-sign/src/signdata"
	"github.com/lindsaylandry/go-mta-train-sign/src/stations"
)

var stop, key, direction string
var cont, train bool

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

	rootCmd.AddCommand(nycMtaCmd)
	rootCmd.AddCommand(ctaCmd)

	rootCmd.PersistentFlags().StringVarP(&stop, "stop", "s", "D30", "stop to parse")
	rootCmd.PersistentFlags().StringVarP(&key, "key", "k", "foobar", "API access key")
	rootCmd.PersistentFlags().BoolVarP(&cont, "continue", "c", true, "continue printing arrivals")
	rootCmd.PersistentFlags().BoolVarP(&train, "train", "t", true, "train or bus (train=true, bus=false)")
	rootCmd.PersistentFlags().StringVarP(&direction, "direction", "d", "N", "direction (trains only)")

	rootCmd.Execute()
}

func CTA() error {
	stop, err := busstops.GetBusStop(stop)
	if err != nil {
		return err
	}

	fmt.Println(stop)

	return nil
}

func NYCMTA() error {
	station, err := stations.GetStation(stop)
	if err != nil {
		return err
	}

	// Get subway feeds from station trains
	feeds := decoder.GetMtaFeeds(station.DaytimeRoutes)

	for {
		arrivals := []feed.Arrival{}
		for _, f := range *feeds {
			t, err := feed.NewTrainFeed(station, key, direction, f.URL)
			if err != nil {
				return err
			}

			arr := t.GetArrivals()
			for _, a := range arr {
				arrivals = append(arrivals, a)
			}
		}

		// Print all arrivals
		signdata.PrintArrivals(arrivals, station.StopName)

		if !cont {
			break
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}
