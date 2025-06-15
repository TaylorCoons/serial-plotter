package passthrough

type Passthrough struct{}

func New() *Passthrough {
	return &Passthrough{}
}

func (p *Passthrough) Compute(data []float32, datum float32) float32 {
	return datum
}
