package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type Entry struct {
	Word        string
	Type        string
	Category    string
	Explanation string
	Example     string
	Translation string
}

type Database struct {
	entries map[string]Entry
}

func NewDatabase() *Database {
	return &Database{
		entries: make(map[string]Entry),
	}
}

func (db *Database) LoadCSV(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) < 2 {
		return fmt.Errorf("csv file is empty")
	}

	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) < 6 {
			continue
		}

		entry := Entry{
			Word:        strings.ToLower(strings.TrimSpace(row[0])),
			Type:        strings.TrimSpace(row[1]),
			Category:    strings.TrimSpace(row[2]),
			Explanation: strings.TrimSpace(row[3]),
			Example:     strings.TrimSpace(row[4]),
			Translation: strings.TrimSpace(row[5]),
		}

		if entry.Word != "" {
			db.entries[entry.Word] = entry
		}
	}
	return nil
}

func (db *Database) Lookup(word string) (Entry, bool) {
	word = strings.ToLower(strings.TrimSpace(word))
	entry, ok := db.entries[word]
	return entry, ok
}

func (db *Database) Explain(word string) string {
	word = strings.ToLower(strings.TrimSpace(word))
	entry, ok := db.entries[word]
	if !ok {
		return fmt.Sprintf("'%s' not found in grammar database", word)
	}

	return fmt.Sprintf(`📖 %s
   Type: %s
   Category: %s
   Meaning: %s
   Example: %s
   Translation: %s`,
		entry.Word, entry.Type, entry.Category,
		entry.Explanation, entry.Example, entry.Translation)
}

func (db *Database) AnalyzePhrase(phrase string) string {
	words := strings.Fields(strings.ToLower(phrase))
	var result strings.Builder
	result.WriteString(fmt.Sprintf("🔍 Analyzing: \"%s\"\n\n", phrase))

	for _, word := range words {
		if entry, ok := db.entries[word]; ok {
			result.WriteString(fmt.Sprintf("  • %s → %s (%s)\n", word, entry.Type, entry.Category))
		} else {
			result.WriteString(fmt.Sprintf("  • %s → unknown\n", word))
		}
	}
	return result.String()
}

func (db *Database) GetWordType(word string) string {
	word = strings.ToLower(strings.TrimSpace(word))
	if entry, ok := db.entries[word]; ok {
		return fmt.Sprintf("%s is a %s (%s)", word, entry.Type, entry.Category)
	}
	return fmt.Sprintf("'%s' not found in database", word)
}

func (db *Database) GetAllEntries() []Entry {
	all := make([]Entry, 0, len(db.entries))
	for _, entry := range db.entries {
		all = append(all, entry)
	}
	return all
}

func (db *Database) Search(query string) []Entry {
	query = strings.ToLower(query)
	var results []Entry
	for word, entry := range db.entries {
		if strings.Contains(word, query) {
			results = append(results, entry)
		}
	}
	return results
}
