package logger

import (
	"time"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	LogDebug.Printf("func(%s) took %s", name, elapsed)
}
