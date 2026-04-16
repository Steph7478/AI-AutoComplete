package model

import (
	"encoding/json"
	"math/rand"
	"os"
)

type LinearLayer struct {
	Weights [][]float64
	Biases  []float64
}

func NewLinearLayer(inputDim, outputDim int) *LinearLayer {
	layer := &LinearLayer{
		Weights: make([][]float64, inputDim),
		Biases:  make([]float64, outputDim),
	}

	for i := range inputDim {
		layer.Weights[i] = make([]float64, outputDim)
		for j := 0; j < outputDim; j++ {
			layer.Weights[i][j] = (rand.Float64() - 0.5) * 0.1
		}
	}
	return layer
}

func (l *LinearLayer) Forward(input []float64) []float64 {
	output := make([]float64, len(l.Biases))
	for j := 0; j < len(l.Biases); j++ {
		sum := l.Biases[j]
		for i := range input {
			sum += input[i] * l.Weights[i][j]
		}
		output[j] = sum
	}
	return output
}

func (l *LinearLayer) SaveToFile(weightsPath, biasesPath string) error {
	wData, _ := json.Marshal(l.Weights)
	bData, _ := json.Marshal(l.Biases)
	os.WriteFile(weightsPath, wData, 0644)
	os.WriteFile(biasesPath, bData, 0644)
	return nil
}

func (l *LinearLayer) LoadFromFile(weightsPath, biasesPath string) error {
	wData, err := os.ReadFile(weightsPath)
	if err != nil {
		return err
	}
	bData, err := os.ReadFile(biasesPath)
	if err != nil {
		return err
	}
	json.Unmarshal(wData, &l.Weights)
	json.Unmarshal(bData, &l.Biases)
	return nil
}
