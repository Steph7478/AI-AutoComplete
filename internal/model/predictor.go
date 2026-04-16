package model

import (
	"AI-AutoComplete/internal/tokenizer"
	"strings"
)

type Predictor struct {
	model     *NGramModel
	tokenizer *tokenizer.Tokenizer
}

func NewPredictor(m *NGramModel, tok *tokenizer.Tokenizer) *Predictor {
	return &Predictor{
		model:     m,
		tokenizer: tok,
	}
}

var naturalEndings = map[string]bool{
	"desu": true, "masu": true, "desu.": true, "masu.": true,
	"gozaimasu": true, "deshita": true, "mashita": true,
	"ka": true, "ka.": true, "ne": true, "yo": true,
}

func (p *Predictor) CompleteSentence(start string, maxSteps int) string {
	words := strings.Fields(start)
	if len(words) == 0 {
		return start
	}

	var result strings.Builder
	result.WriteString(start)

	currentTokens := p.tokenizer.Encode(start)

	if len(currentTokens) == 1 && len(words) == 1 {
		currentTokens = append(currentTokens, currentTokens[0])
	}

	for step := 0; step < maxSteps; step++ {
		if len(currentTokens) < 2 {
			break
		}

		token1 := currentTokens[len(currentTokens)-2]
		token2 := currentTokens[len(currentTokens)-1]
		probs := p.model.Forward(token1, token2)

		bestIdx := 0
		bestProb := 0.0
		for i, prob := range probs {
			if prob > bestProb {
				bestProb = prob
				bestIdx = i
			}
		}

		if bestProb < 0.05 {
			break
		}

		nextWord := p.tokenizer.Decode([]int{bestIdx})

		if nextWord == "<UNK>" {
			break
		}

		result.WriteString(" " + nextWord)
		currentTokens = append(currentTokens, bestIdx)

		if naturalEndings[nextWord] {
			lastWord := p.tokenizer.Decode([]int{currentTokens[len(currentTokens)-1]})
			if naturalEndings[lastWord] {
				break
			}
		}
	}

	return result.String()
}
