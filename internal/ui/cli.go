package ui

import (
	"AI-AutoComplete/internal/model"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CLI struct {
	predictor *model.Predictor
	reader    *bufio.Reader
}

func NewCLI(predictor *model.Predictor) *CLI {
	return &CLI{
		predictor: predictor,
		reader:    bufio.NewReader(os.Stdin),
	}
}

func (c *CLI) Run() {
	fmt.Println("\n=== JAPANESE SENTENCE COMPLETER ===")
	fmt.Println("Type the beginning of a sentence, I will complete it")
	fmt.Println("The model stops when uncertain or reaches max tokens")
	fmt.Println("Type 'exit' to quit")
	fmt.Print("Max tokens to generate? (default 15): ")

	maxStr, _ := c.reader.ReadString('\n')
	maxSteps := 15
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

		result := c.predictor.CompleteSentence(input, maxSteps)
		fmt.Printf("→ %s\n\n", result)
	}

	fmt.Println("\nSayonara!")
}
