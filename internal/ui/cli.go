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
	fmt.Println("Type 'exit' to quit")
	fmt.Println("Type 'max N' to change max steps (default 20)\n")

	maxSteps := 20

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

		if strings.HasPrefix(input, "max ") {
			fmt.Sscanf(input, "max %d", &maxSteps)
			fmt.Printf("→ Max steps set to %d\n\n", maxSteps)
			continue
		}

		result := c.predictor.CompleteSentence(input, maxSteps)
		fmt.Printf("→ %s\n\n", result)
	}

	fmt.Println("\nSayonara!")
}
