# go-transit-sign
GoLang Project for downloading transit data and broadcasting to LED Matrix signs.

Currently will print to STDOUT for MTA train lines.

Broadcasting to LED Matrix and ability to use more than MTA data is a work in progress.

## How to Use

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

Flags:
  -c, --continue           continue printing arrivals (default true)
  -d, --direction string   direction (trains only) (default "N")
  -h, --help               help for transit-sign
  -k, --key string         API access key (default "foobar")
  -s, --stop string        stop to parse (default "D30")
  -t, --train              train or bus (train=true, bus=false) (default true)

Use "transit-sign [command] --help" for more information about a command.
```

### Compile

To build the binary run:
`go build`

## How to Get API Keys

### NYC MTA

TODO

### CTA

Create an account and request an API key at https://www.ctabustracker.com/home

## References

### GTFS Protocol
https://pkg.go.dev/github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs

### NYC MTA
[Stagions CSV](http://web.mta.info/developers/data/nyct/subway/Stations.csv)

### Chicago Transit Authority (CTA)
[Train Tracker API Application](https://www.transitchicago.com/developers/traintrackerapply/)
