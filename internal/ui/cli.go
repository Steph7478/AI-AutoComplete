package ui

import (
	"AI-AutoComplete/internal/data"
	"AI-AutoComplete/internal/model"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CLI struct {
	predictor    *model.Predictor
	translations *data.TranslationDB
	reader       *bufio.Reader
}

func NewCLI(predictor *model.Predictor, transDB *data.TranslationDB) *CLI {
	return &CLI{
		predictor:    predictor,
		translations: transDB,
		reader:       bufio.NewReader(os.Stdin),
	}
}

func (c *CLI) Run() {
	fmt.Println("\n=== JAPANESE SENTENCE COMPLETER ===")
	fmt.Println("Type the beginning of a Japanese sentence")
	fmt.Println("The model will complete it and show English translation")
	fmt.Println("Type 'exit' to quit")
	fmt.Println("Example: 最近, 私は, すぐに, 愛して, 何して, おはよう")
	fmt.Print("\nMax tokens to generate? (default 10): ")

	maxStr, _ := c.reader.ReadString('\n')
	maxSteps := 10
	if maxStr = strings.TrimSpace(maxStr); maxStr != "" {
		fmt.Sscanf(maxStr, "%d", &maxSteps)
	}
	fmt.Printf("Will generate up to %d tokens\n\n", maxSteps)

	for {
		fmt.Print("You: ")
		input, err := c.reader.ReadString('\n')
		if err != nil {
			break
		}

		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		if input == "" {
			continue
		}

		completed := c.predictor.CompleteSentence(input, maxSteps)
		fmt.Printf("JP: %s\n", completed)

		translation := c.translations.GetTranslation(completed)
		if translation != "" {
			fmt.Printf("EN: %s\n", translation)
		} else {
			found := false
			for jp, en := range c.translations.GetAllEntries() {
				if strings.Contains(jp, completed) || strings.Contains(completed, jp) {
					fmt.Printf("EN: %s\n", en)
					found = true
					break
				}
			}
			if !found {
				for jp, en := range c.translations.GetAllEntries() {
					if strings.Contains(jp, input) {
						fmt.Printf("EN: %s\n", en)
						found = true
						break
					}
				}
			}
			if !found {
				fmt.Printf("EN: (translation not found)\n")
			}
		}
		fmt.Println()
	}

	fmt.Println("\nSayonara!")
}
