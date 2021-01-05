package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"

	"github.com/toretto460/sensor/internal/app"
)

func setup(ctx context.Context) context.Context {
	d, err := dev.DefaultDevice()
	if err != nil {
		panic(err)
	}
	ble.SetDefaultDevice(d)

	return ble.WithSigHandler(context.WithCancel(ctx))
}

func handleSignals() chan bool {
	var quit chan bool
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		quit <- true
	}()

	return quit
}

func influxClient() (api.WriteAPI, func()) {
	token := os.Getenv("INFLUX_TOKEN")
	bucket := os.Getenv("INFLUX_BUCKET")
	org := os.Getenv("INFLUX_ORG")
	host := os.Getenv("INFLUX_HOST")

	client := influxdb2.NewClient(host, token)
	writeAPI := client.WriteAPI(org, bucket)

	go func() {
		for {
			time.Sleep(time.Minute)
			writeAPI.Flush()
		}
	}()

	return writeAPI, func() { client.Close() }
}

func main() {
	ctx := setup(context.Background())
	quit := handleSignals()
	writeAPI, closeAPI := influxClient()
	application := app.NewApp(writeAPI, app.NewTimeFilter(60))

	ble.Scan(ctx, true, func(a ble.Advertisement) {
		data, err := application.ParseSensorData(a.ManufacturerData())
		if err == nil {
			application.RecordData(data)
		}
	}, nil)

	<-quit
	closeAPI()
}
