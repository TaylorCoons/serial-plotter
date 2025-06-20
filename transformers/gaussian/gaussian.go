package gaussian

import (
	"math/rand"
	"time"

	"gonum.org/v1/gonum/stat/distuv"
)

type Guassian struct {
	normalDistribution distuv.Normal
}

func New(mean float32, standardDeviation float32) *Guassian {
	normalDistribution := distuv.Normal{
		Mu:    float64(mean),
		Sigma: float64(standardDeviation),
		Src:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return &Guassian{
		normalDistribution: normalDistribution,
	}
}

func (g *Guassian) Compute(data []float32, datum float32) float32 {
	return datum + float32(g.normalDistribution.Rand())
}
