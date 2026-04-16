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

var stopTokens = map[string]bool{
	".": true, "!": true, "?": true, "\n": true,
	"。": true, "！": true, "？": true,
}

func (p *Predictor) CompleteSentence(start string, maxSteps int) string {
	words := strings.Fields(start)
	if len(words) == 0 {
		return start
	}

	result := start
	currentTokens := p.tokenizer.Encode(start)

	for step := 0; step < maxSteps; step++ {
		if len(currentTokens) < 1 {
			break
		}

		lastToken := currentTokens[len(currentTokens)-1]
		probs := p.model.Forward(lastToken, lastToken)

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

		result += " " + nextWord
		currentTokens = append(currentTokens, bestIdx)

		if p.isStopToken(nextWord) {
			break
		}
	}

	return result
}

func (p *Predictor) PredictNext(word string) string {
	tokens := p.tokenizer.Encode(word)
	if len(tokens) == 0 {
		return "<UNK>"
	}

	lastToken := tokens[len(tokens)-1]
	probs := p.model.Forward(lastToken, lastToken)

	bestIdx := 0
	bestProb := 0.0
	for i, prob := range probs {
		if prob > bestProb {
			bestProb = prob
			bestIdx = i
		}
	}

	return p.tokenizer.Decode([]int{bestIdx})
}

func (p *Predictor) isStopToken(word string) bool {
	return stopTokens[word]
}
