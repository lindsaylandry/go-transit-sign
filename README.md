# go-transit-sign
GoLang Project for downloading transit data and broadcasting to LED Matrix signs.

This project currently broadcasts to Adafruit 64x32 LED Matrix for CTA buses and trains, and NYC MTA trains.

## Dependencies

This app depends on the following library:
* https://github.com/tfK1410/go-rpi-rgb-led-matrix

This library works only in Linux-based environments. If running on Mac or Windows,
consider compiling this app in Docker or a Linux VM.

Future commits to this app will include dockerfiles and docker-compose files for local development.

## How to Use

### Building the App

To build the app run:
```bash
go build
```

Running this the first time will result in errors. To fix these, follow these instructions from
[tfK1410/go-rpi-rgb-led-matrix](https://github.com/tfK1410/go-rpi-rgb-led-matrix?tab=readme-ov-file#installation)

### Config File

This app requires a config file to run, located in `configs/config.yaml`
Use the example (config.example.yaml)[configs/config.example.yaml] to build your config for your needs.

### Manual

```
Run transit sign

Usage:
  transit-sign [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  cta         Run CTA data
  help        Help about any command
  nyc-mta     Run NYC MTA data
  testMatrix  Test LED Matrix

Flags:
  -d, --direction string   direction (trains only) (default "N")
  -h, --help               help for transit-sign
  -l, --led                output to led matrix

Use "transit-sign [command] --help" for more information about a command.
```

## How to Get API Keys

### NYC MTA

TODO

### CTA

Create an account and request an API key at https://www.ctabustracker.com/home

## References

### LED Matrix

[Adafruit LED Matrix Pinouts](https://learn.adafruit.com/adafruit-rgb-matrix-bonnet-for-raspberry-pi/pinouts)

### Transit Data

#### GTFS Protocol
https://pkg.go.dev/github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs

#### NYC MTA
[Stations CSV](http://web.mta.info/developers/data/nyct/subway/Stations.csv)

#### Chicago Transit Authority (CTA)
[Train Tracker API Application](https://www.transitchicago.com/developers/traintrackerapply/)
