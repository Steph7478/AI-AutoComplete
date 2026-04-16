package tokenizer

import (
	"encoding/json"
	"os"
	"strings"
)

type Tokenizer struct {
	WordToIdx map[string]int
	IdxToWord map[int]string
	VocabSize int
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		WordToIdx: make(map[string]int),
		IdxToWord: make(map[int]string),
	}
}

func (t *Tokenizer) BuildFromText(text string) {
	words := strings.Fields(strings.ToLower(text))

	t.WordToIdx["<UNK>"] = 0
	t.IdxToWord[0] = "<UNK>"
	t.VocabSize = 1

	for _, word := range words {
		if _, exists := t.WordToIdx[word]; !exists {
			t.WordToIdx[word] = t.VocabSize
			t.IdxToWord[t.VocabSize] = word
			t.VocabSize++
		}
	}
}

func (t *Tokenizer) Encode(text string) []int {
	words := strings.Fields(strings.ToLower(text))
	tokens := make([]int, len(words))

	for i, word := range words {
		if idx, exists := t.WordToIdx[word]; exists {
			tokens[i] = idx
		} else {
			tokens[i] = 0
		}
	}
	return tokens
}

func (t *Tokenizer) Decode(tokens []int) string {
	words := make([]string, len(tokens))
	for i, token := range tokens {
		if word, exists := t.IdxToWord[token]; exists {
			words[i] = word
		} else {
			words[i] = "<UNK>"
		}
	}
	return strings.Join(words, " ")
}

func (t *Tokenizer) SaveToFile(path string) error {
	data := struct {
		WordToIdx map[string]int `json:"word_to_idx"`
		IdxToWord map[int]string `json:"idx_to_word"`
		VocabSize int            `json:"vocab_size"`
	}{
		WordToIdx: t.WordToIdx,
		IdxToWord: t.IdxToWord,
		VocabSize: t.VocabSize,
	}
	file, _ := json.Marshal(data)
	return os.WriteFile(path, file, 0644)
}

func (t *Tokenizer) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var loaded struct {
		WordToIdx map[string]int `json:"word_to_idx"`
		IdxToWord map[int]string `json:"idx_to_word"`
		VocabSize int            `json:"vocab_size"`
	}
	json.Unmarshal(data, &loaded)
	t.WordToIdx = loaded.WordToIdx
	t.IdxToWord = loaded.IdxToWord
	t.VocabSize = loaded.VocabSize
	return nil
}
