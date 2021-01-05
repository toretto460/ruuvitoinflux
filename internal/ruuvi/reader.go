package ruuvi

import (
	"encoding/hex"

	"github.com/peterhellberg/ruuvitag"
)

type RuuviData struct {
	ID          string
	Humidity    float64 // 0.0 % to 100 % in 0.5 % increments.
	Temperature float64 // -127.99 °C to +127.99 °C in 0.01 °C increments.
	Pressure    uint32  // 50000 Pa to 115536 Pa in 1 Pa increments.
	Battery     uint16  // 0 mV to 65536 mV in 1 mV increments, practically 1800 ... 3600 mV.
}

func (r RuuviData) SensorID() string {
	return r.ID
}

func (r RuuviData) Temp() float64 {
	return r.Temperature
}

func IsRuuvi(data []byte) bool {
	return ruuvitag.IsRAWv2(data) || ruuvitag.IsRAWv1(data)
}

// ParseRuuvi parses the data []byte into a RuuviData struct
func ParseRuuvi(data []byte) (RuuviData, error) {
	var err error

	switch {
	case ruuvitag.IsRAWv2(data):
		if raw, err := ruuvitag.ParseRAWv2(data); err == nil {
			return RuuviData{
				ID:          hex.EncodeToString(raw.MAC[:]),
				Humidity:    raw.Humidity,
				Temperature: raw.Temperature,
				Pressure:    raw.Pressure,
				Battery:     raw.Battery,
			}, nil
		}
	case ruuvitag.IsRAWv1(data):
		if raw, err := ruuvitag.ParseRAWv1(data); err == nil {
			return RuuviData{
				ID:          "Unkown",
				Humidity:    raw.Humidity,
				Temperature: raw.Temperature,
				Pressure:    raw.Pressure,
				Battery:     raw.Battery,
			}, nil
		}
	}

	return RuuviData{}, err
}
