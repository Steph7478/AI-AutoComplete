package embedding

import (
	"encoding/json"
	"math/rand"
	"os"
)

type Embedding struct {
	Weights [][]float64
	Dim     int
}

func NewEmbedding(vocabSize, embeddingDim int) *Embedding {
	emb := &Embedding{
		Weights: make([][]float64, vocabSize),
		Dim:     embeddingDim,
	}

	for i := range vocabSize {
		emb.Weights[i] = make([]float64, embeddingDim)
		for j := range embeddingDim {
			emb.Weights[i][j] = (rand.Float64() - 0.5) * 0.1
		}
	}
	return emb
}

func (e *Embedding) Get(token int) []float64 {
	if token >= len(e.Weights) {
		return e.Weights[0]
	}
	return e.Weights[token]
}

func (e *Embedding) Update(token int, gradient []float64, learningRate float64) {
	if token >= len(e.Weights) {
		token = 0
	}
	for i := 0; i < e.Dim; i++ {
		e.Weights[token][i] -= learningRate * gradient[i]
	}
}

func (e *Embedding) UpdateAt(token, dim int, delta float64) {
	if token >= len(e.Weights) {
		token = 0
	}
	if dim < e.Dim {
		e.Weights[token][dim] -= delta
	}
}

func (e *Embedding) SaveToFile(path string) error {
	data, _ := json.Marshal(e.Weights)
	return os.WriteFile(path, data, 0644)
}

func (e *Embedding) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	json.Unmarshal(data, &e.Weights)
	if len(e.Weights) > 0 {
		e.Dim = len(e.Weights[0])
	}
	return nil
}
