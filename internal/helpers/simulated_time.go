package helpers

import "time"

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
