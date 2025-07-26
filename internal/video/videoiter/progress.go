package videoiter

import (
	"math"
	"time"
)

type Progress struct {
	Done      int
	Max       int
	timestamp time.Time
}

func (p Progress) Percent() float64 {
	return float64(p.Done) / float64(p.Max) * 100
}

func (p Progress) FPS(other Progress) float64 {
	return math.Abs(float64(p.Done-other.Done)) / p.timestamp.Sub(other.timestamp).Seconds()
}
