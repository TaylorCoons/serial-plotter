package dummy

import (
	"math"
	"time"
)

type Function func(index int) float32

type Dummy struct {
	index    int
	delay    time.Duration
	function Function
}

func ConstantFunction(index int) float32 {
	return 8
}

func SinFunction(index int) float32 {
	return float32(10 * math.Sin(float64(index)*2*math.Pi/30))
}

func XSinXFunction(index int) float32 {
	x := float64(index)
	return float32(0.25 * x * math.Sin(x*math.Pi/15))
}

func Neg100X(index int) float32 {
	return -100 * float32(index)
}

func SquareFunction(index int) float32 {
	period := float64(50)
	sin := math.Sin(float64(index) * 2 * math.Pi / period)
	if sin > 0 {
		return 10
	} else {
		return -10
	}
}

func SawtoothFunction(index int) float32 {
	period := 10.0
	percent := float64(index) / period
	return 10 * float32(2*(percent-math.Floor(0.5+percent)))
}

func New(delay time.Duration, function Function) *Dummy {
	return &Dummy{
		index:    0,
		delay:    delay,
		function: function,
	}
}

func (p *Dummy) ResetIndex() {
	p.index = 0
}

func (p *Dummy) SetFunction(function Function) {
	p.function = function
}

func (p *Dummy) ReadSource() (float32, error) {
	time.Sleep(p.delay)
	defer func() {
		p.index++
	}()
	return p.function(p.index), nil
}
