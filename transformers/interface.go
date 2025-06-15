package transformers

type Transformer interface {
	Compute(data []float32, datum float32) float32
}
