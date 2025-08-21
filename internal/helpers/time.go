package helpers

import (
	"errors"
	"time"
)

const PlaybackSpeed = time.Duration(5)

var startSimulatedTime time.Time
var baseSimulatedTime time.Time

func StartSimulatedTime() {
	startSimulatedTime = time.Now().UTC()

	// Opening in UTC time
	baseSimulatedTime = time.Date(
		startSimulatedTime.Year(), startSimulatedTime.Month(), startSimulatedTime.Day(),
		13, 30, 0, 0, startSimulatedTime.Location())
	baseSimulatedTime = baseSimulatedTime.Add(-24 * time.Hour)
}

func GetSimulatedTime(t time.Time) time.Time {
	elapsed := t.Sub(startSimulatedTime) * PlaybackSpeed
	return baseSimulatedTime.Add(elapsed)
}

func TimeframeToDuration(timeframe string) (time.Duration, error) {
	allowed_timeframes := map[string]time.Duration{
		"1m":  1 * time.Minute,
		"5m":  5 * time.Minute,
		"15m": 15 * time.Minute,
		"30m": 30 * time.Minute,
		"1h":  time.Hour,
		"2h":  2 * time.Hour,
		"4h":  4 * time.Hour,
	}

	duration, exists := allowed_timeframes[timeframe]
	if !exists {
		return 0, errors.New("Timeframe not allowed")
	}

	return duration, nil
}

func CandleTime(t time.Time, timeframe time.Duration) time.Time {
	rounded := t.Truncate(timeframe)
	if t.Equal(rounded) {
		return rounded.Add(-timeframe)
	}
	return rounded
}
