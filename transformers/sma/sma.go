package sma

type Sma struct {
	k int
}

func New(k int) *Sma {
	return &Sma{
		k: k,
	}
}

func (s *Sma) Compute(data []float32, datum float32) float32 {
	if len(data) < 1 {
		return datum
	}
	k := s.k
	if len(data) < s.k {
		k = len(data) + 1
	}
	smaPrev := data[len(data)-1]
	return smaPrev + 1.0/float32(k)*(datum-data[len(data)-k+1])
}
