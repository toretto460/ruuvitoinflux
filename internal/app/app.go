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
		a.writer.WriteRecord(fmt.Sprintf("temperature,sensor=%s,unit=temperature avg=%f", data.SensorID(), data.Temp()))
	}

	return true
}
