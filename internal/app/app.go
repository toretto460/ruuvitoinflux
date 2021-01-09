package app

import (
	"errors"
	"fmt"

	"github.com/toretto460/sensor/internal/ruuvi"
)

type Writer interface {
	WriteRecord(string)
}

type SensorData interface {
	SensorID() string
	Temp() float64
	Humidity() float64
	Pressure() uint32
}

type Filter interface {
	Filter(SensorData) bool
}

type App struct {
	writer  Writer
	filters []Filter
}

func NewApp(w Writer, f ...Filter) App {
	return App{w, f}
}

func (a *App) ParseSensorData(data []byte) (SensorData, error) {
	switch {
	case ruuvi.IsRuuvi(data):
		return ruuvi.ParseRuuvi(data)
	}

	return nil, errors.New("Cannot parse sensor data")
}

func (a *App) RecordData(data SensorData) bool {
	pass := true
	for _, f := range a.filters {
		pass = pass && f.Filter(data)
	}

	if pass {
		a.writer.WriteRecord(
			fmt.Sprintf("temperature,sensor=%s,unit=temperature avg=%f", data.SensorID(), data.Temp()),
		)
		a.writer.WriteRecord(
			fmt.Sprintf("humidity,sensor=%s,unit=precentage avg=%f", data.SensorID(), data.Humidity()),
		)
		a.writer.WriteRecord(
			fmt.Sprintf("pressure,sensor=%s,unit=hectopascal avg=%d", data.SensorID(), data.Pressure()),
		)
	}

	return true
}
