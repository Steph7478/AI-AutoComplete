package model

import (
	"encoding/json"
	"math/rand"
	"os"
)

type NGramModel struct {
	Embedding *Embedding
	Linear    *LinearLayer
}

func NewNGramModel(vocabSize, embeddingDim int) *NGramModel {
	return &NGramModel{
		Embedding: NewEmbedding(vocabSize, embeddingDim),
		Linear:    NewLinearLayer(embeddingDim*2, vocabSize),
	}
}

func (m *NGramModel) Forward(token1, token2 int) []float64 {
	emb1 := m.Embedding.Get(token1)
	emb2 := m.Embedding.Get(token2)

	combined := make([]float64, len(emb1)+len(emb2))
	copy(combined, emb1)
	copy(combined[len(emb1):], emb2)

	logits := m.Linear.Forward(combined)
	return Softmax(logits)
}

func (m *NGramModel) GetEmbedding() *Embedding {
	return m.Embedding
}

func (m *NGramModel) GetLinear() *LinearLayer {
	return m.Linear
}

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
