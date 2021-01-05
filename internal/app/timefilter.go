package app

import "time"

func NewTimeFilter(seconds uint8) TimeFilter {
	return TimeFilter{
		lastPass:         make(map[string]time.Time),
		secondsThreshold: seconds,
	}
}

type TimeFilter struct {
	lastPass         map[string]time.Time
	secondsThreshold uint8
}

func (t TimeFilter) Filter(data SensorData) bool {
	now := time.Now()

	if lastReadTime, ok := t.lastPass[data.SensorID()]; ok {
		if now.Sub(lastReadTime).Seconds() > float64(t.secondsThreshold) {
			t.lastPass[data.SensorID()] = now
			return true
		}
	} else {
		t.lastPass[data.SensorID()] = now
		return true
	}

	return false
}
