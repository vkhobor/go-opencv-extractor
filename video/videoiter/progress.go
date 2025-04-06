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

func (p Progress) MergeWith(other ...Progress) Progress {
	for _, o := range other {
		p.Done += o.Done
		p.Max += o.Max
		if o.timestamp.After(p.timestamp) {
			p.timestamp = o.timestamp
		}

	}
	return p
}
