package pseudo

import (
	"math"
	"time"
)

type Transform func(index int) float32

type Pseudo struct {
	index     int
	delay     time.Duration
	transform Transform
}

func SinTransform(index int) float32 {
	return float32(10 * math.Sin(float64(index)*2*math.Pi/30))
}

func SquareTransform(index int) float32 {
	period := float64(10)
	sin := math.Sin(float64(index) * 2 * math.Pi / period)
	if sin > 0 {
		return 10
	} else {
		return -10
	}
}

func SawtoothTransform(index int) float32 {
	period := 10.0
	percent := float64(index) / period
	return float32(2 * (percent - math.Floor(0.5+percent)))
}

func New(delay time.Duration, transform Transform) *Pseudo {
	return &Pseudo{
		index:     0,
		delay:     delay,
		transform: transform,
	}
}

func (p *Pseudo) ReadSource() (float32, error) {
	time.Sleep(p.delay)
	defer func() {
		p.index++
	}()
	return p.transform(p.index), nil
}
