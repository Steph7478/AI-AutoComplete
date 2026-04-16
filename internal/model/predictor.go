package model

import (
	"AI-AutoComplete/internal/data"
	"AI-AutoComplete/internal/tokenizer"
	"strings"
)

type Predictor struct {
	model        *NGramModel
	tokenizer    *tokenizer.Tokenizer
	translations *data.TranslationDB
}

func NewPredictor(m *NGramModel, tok *tokenizer.Tokenizer, transDB *data.TranslationDB) *Predictor {
	return &Predictor{
		model:        m,
		tokenizer:    tok,
		translations: transDB,
	}
}

func (p *Predictor) CompleteSentence(start string, maxSteps int) string {
	words := strings.Fields(start)
	if len(words) == 0 {
		return start
	}

	if p.translations != nil {
		for jp := range p.translations.GetAllEntries() {
			if strings.HasPrefix(jp, start) {
				return jp
			}
		}
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
	}

	return result.String()
}
