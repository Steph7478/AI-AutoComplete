package storage

import (
	"AI-AutoComplete/internal/model"
	"AI-AutoComplete/internal/tokenizer"
	"os"
)

func SaveModel(m *model.NGramModel, tok *tokenizer.Tokenizer) error {
	os.MkdirAll("internal/data", 0755)

	tok.SaveToFile("internal/data/tokenizer.json")
	m.GetEmbedding().SaveToFile("internal/data/embeddings.json")
	m.GetLinear().SaveToFile("internal/data/linear_weights.json", "internal/data/linear_biases.json")
	return nil
}

func LoadModel() (*model.NGramModel, *tokenizer.Tokenizer, error) {
	tok := tokenizer.NewTokenizer()
	if err := tok.LoadFromFile("internal/data/tokenizer.json"); err != nil {
		return nil, nil, err
	}

	emb := model.NewEmbedding(tok.VocabSize, 32)
	if err := emb.LoadFromFile("internal/data/embeddings.json"); err != nil {
		return nil, nil, err
	}

	linear := model.NewLinearLayer(64, tok.VocabSize)
	if err := linear.LoadFromFile("internal/data/linear_weights.json", "internal/data/linear_biases.json"); err != nil {
		return nil, nil, err
	}

	m := &model.NGramModel{
		Embedding: emb,
		Linear:    linear,
	}

	return m, tok, nil
}
