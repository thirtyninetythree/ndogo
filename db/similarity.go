package similarity

import (
	"fmt"
	"math"
)

type DistanceMetric int

const (
	Euclidean DistanceMetric = iota
	Cosine
	DotProduct
)

func (d DistanceMetric) String() string {
	switch d {
	case Euclidean:
		return "euclidean"
	case Cosine:
		return "cosine"
	case DotProduct:
		return "dot"
	default:
		return fmt.Sprintf("unknown distance metric %d", int(d))
	}
}

func getCacheAttr(metric DistanceMetric, vec []float32) float32 {
	switch metric {
	case DotProduct, Euclidean:
		return 0.0
	case Cosine:
		var sum float32
		for _, x := range vec {
			sum += x * x
		}
		return float32(math.Sqrt(float64(sum)))
	default:
		return 0.0
	}
}

func getDistanceFunc(metric DistanceMetric) func([]float32, []float32, float32) float32 {
	switch metric {
	case Euclidean:
		return euclideanDistance
	case Cosine, DotProduct:
		return dotProduct
	default:
		return func([]float32, []float32, float32) float32 { return 0.0 }
	}
}

func euclideanDistance(a, b []float32, aSumSquares float32) float32 {
	var crossTerms, bSumSquares float32
	for i, x := range a {
		y := b[i]
		crossTerms += x * y
		bSumSquares += y * y
	}
	return float32(math.Sqrt(float64(2.0*(-crossTerms) + aSumSquares + bSumSquares)))
}

func dotProduct(a, b []float32, _ float32) float32 {
	var sum float32
	for i, x := range a {
		sum += x * b[i]
	}
	return sum
}

func normalize(vec []float32) []float32 {
	var sum float32
	for _, x := range vec {
		sum += x * x
	}
	magnitude := float32(math.Sqrt(float64(sum)))
	if magnitude > math.SmallestNonzeroFloat32 {
		result := make([]float32, len(vec))
		for i, x := range vec {
			result[i] = x / magnitude
		}
		return result
	}
	return vec
}

type ScoreIndex struct {
	Score float32
	Index int
}

func (si ScoreIndex) Less(other ScoreIndex) bool {
	return other.Score < si.Score
}

func (si ScoreIndex) Greater(other ScoreIndex) bool {
	return other.Score > si.Score
}
