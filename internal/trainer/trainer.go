package trainer

import (
	"AI-AutoComplete/internal/model"
	"fmt"
	"math"
	"math/rand"
)

type NGramTrainer struct {
	Model        *model.NGramModel
	LearningRate float64
}

func NewNGramTrainer(m *model.NGramModel, learningRate float64) *NGramTrainer {
	return &NGramTrainer{
		Model:        m,
		LearningRate: learningRate,
	}
}

func (t *NGramTrainer) TrainStep(token1, token2, token3 int) float64 {
	predictions := t.Model.Forward(token1, token2)
	loss := -math.Log(predictions[token3] + 1e-8)

	gradient := make([]float64, len(predictions))
	for i := range predictions {
		if i == token3 {
			gradient[i] = predictions[i] - 1.0
		} else {
			gradient[i] = predictions[i]
		}
	}

	emb1 := t.Model.GetEmbedding().Get(token1)
	emb2 := t.Model.GetEmbedding().Get(token2)

	combined := make([]float64, len(emb1)+len(emb2))
	copy(combined, emb1)
	copy(combined[len(emb1):], emb2)

	linear := t.Model.GetLinear()
	for i := range combined {
		for j := range gradient {
			linear.Weights[i][j] -= t.LearningRate * gradient[j] * combined[i]
		}
	}

	for j := range gradient {
		linear.Biases[j] -= t.LearningRate * gradient[j]
	}

	embGrad1 := make([]float64, len(emb1))
	embGrad2 := make([]float64, len(emb2))

	for i := range emb1 {
		for j := 0; j < len(gradient); j++ {
			embGrad1[i] += gradient[j] * linear.Weights[i][j]
		}
	}

	for i := range emb2 {
		for j := range gradient {
			embGrad2[i] += gradient[j] * linear.Weights[len(emb1)+i][j]
		}
	}

	for i := range embGrad1 {
		t.Model.GetEmbedding().UpdateAt(token1, i, t.LearningRate*embGrad1[i])
	}
	for i := range embGrad2 {
		t.Model.GetEmbedding().UpdateAt(token2, i, t.LearningRate*embGrad2[i])
	}

	return loss
}

func (t *NGramTrainer) Train(tokens []int, epochs int) {
	fmt.Printf("Training N-Gram model with %d tokens...\n", len(tokens))

	for epoch := range epochs {
		totalLoss := 0.0
		count := 0

		for i := 0; i < len(tokens)-2; i++ {
			loss := t.TrainStep(tokens[i], tokens[i+1], tokens[i+2])
			totalLoss += loss
			count++
		}

		if epoch%50 == 0 {
			avgLoss := totalLoss / float64(count)
			fmt.Printf("Epoch %d, Loss: %.4f\n", epoch, avgLoss)
		}
	}
}

func init() {
	rand.Seed(42)
}
