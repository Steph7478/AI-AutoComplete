package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type TranslationDB struct {
	entries    map[string]string
	firstWords map[string]string
}

func NewTranslationDB() *TranslationDB {
	return &TranslationDB{
		entries:    make(map[string]string),
		firstWords: make(map[string]string),
	}
}

func (db *TranslationDB) LoadCSV(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true

	lineNum := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Skipping line %d: %v\n", lineNum, err)
			continue
		}
		lineNum++

		if len(record) < 2 {
			continue
		}

		var japanese, english string

		if len(record) >= 4 && isNumeric(record[0]) && isNumeric(record[2]) {
			japanese = cleanText(record[1])
			english = cleanText(record[3])
		} else if len(record) >= 2 {
			japanese = cleanText(record[0])
			english = cleanText(record[1])
		}

		if japanese != "" && english != "" {
			db.entries[japanese] = english

			words := strings.Fields(japanese)
			if len(words) > 0 {
				firstWord := words[0]
				if _, exists := db.firstWords[firstWord]; !exists {
					db.firstWords[firstWord] = english
				}
			}
		}
	}

	fmt.Printf("Loaded %d translations from %s\n", len(db.entries), path)
	return nil
}

func cleanText(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, `"`)
	s = strings.Trim(s, `'`)
	s = strings.Join(strings.Fields(s), " ")
	return s
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func (db *TranslationDB) GetTranslation(japanese string) string {
	if eng, ok := db.entries[japanese]; ok {
		return eng
	}

	cleaned := strings.TrimRight(japanese, "。！？")
	if eng, ok := db.entries[cleaned]; ok {
		return eng
	}

	words := strings.Fields(japanese)
	if len(words) > 0 {
		if eng, ok := db.firstWords[words[0]]; ok {
			return eng
		}
	}

	return ""
}

func (db *TranslationDB) GetAllEntries() map[string]string {
	return db.entries
}
