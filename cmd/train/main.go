package main

import (
	"AI-AutoComplete/internal/data"
	"AI-AutoComplete/internal/model"
	"AI-AutoComplete/internal/storage"
	"AI-AutoComplete/internal/tokenizer"
	"AI-AutoComplete/internal/trainer"
	"AI-AutoComplete/internal/ui"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== TRAIN N-GRAM MODEL FROM CSV ===\n")
	fmt.Println("1. Train new model from CSV")
	fmt.Println("2. Test trained model (sentence completion)")
	fmt.Println("3. Search CSV only (no model)")
	fmt.Print("\nChoose (1/2/3): ")

	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	switch option {
	case "1":
		trainModel(reader)
	case "2":
		testModel()
	case "3":
		searchCSVOnly()
	default:
		fmt.Println("Invalid option")
	}
}

func trainModel(reader *bufio.Reader) {
	fmt.Println("\n=== TRAINING FROM CSV ===")

	transDB := data.NewTranslationDB()
	csvPath := "data/translations_small.csv"

	if err := transDB.LoadCSV(csvPath); err != nil {
		log.Fatal("Error loading CSV:", err)
	}

	var corpusBuilder strings.Builder
	for japanese := range transDB.GetAllEntries() {
		corpusBuilder.WriteString(japanese + "\n")
	}
	text := corpusBuilder.String()

	fmt.Printf("✅ Generated corpus with %d sentences\n", len(transDB.GetAllEntries()))
	fmt.Printf("✅ Corpus loaded: %d characters\n", len(text))

	tok := tokenizer.NewTokenizer()
	tok.BuildFromText(text)
	fmt.Printf("✅ Vocabulary: %d unique words\n", tok.VocabSize)

	tokens := tok.Encode(text)
	fmt.Printf("✅ Total tokens: %d\n", len(tokens))

	embeddingDim := 32
	myModel := model.NewNGramModel(tok.VocabSize, embeddingDim)
	fmt.Printf("✅ Model created\n\n")

	fmt.Print("How many epochs? (default 50 for quick training): ")
	epochsStr, _ := reader.ReadString('\n')
	epochs := 50
	fmt.Sscanf(epochsStr, "%d", &epochs)

	t := trainer.NewNGramTrainer(myModel, 0.05)
	t.Train(tokens, epochs)

	storage.SaveModel(myModel, tok)
	fmt.Println("\n✅ Model saved to internal/data/")

	predictor := model.NewPredictor(myModel, tok, transDB)
	cli := ui.NewCLI(predictor, transDB)
	cli.Run()
}

func testModel() {
	fmt.Println("\n=== LOADING MODEL ===")

	myModel, tok, err := storage.LoadModel()
	if err != nil {
		fmt.Println("No saved model found. Please train first (option 1).")
		return
	}

	fmt.Printf("✅ Model loaded\n\n")

	transDB := data.NewTranslationDB()
	csvPath := "data/translations_small.csv"
	if err := transDB.LoadCSV(csvPath); err != nil {
		fmt.Printf("⚠️ Could not load translations: %v\n", err)
	}

	predictor := model.NewPredictor(myModel, tok, transDB)
	cli := ui.NewCLI(predictor, transDB)
	cli.Run()
}

func searchCSVOnly() {
	fmt.Println("\n=== CSV SEARCH MODE ===\n")

	transDB := data.NewTranslationDB()
	csvPath := "data/translations_small.csv"

	if err := transDB.LoadCSV(csvPath); err != nil {
		log.Fatal("Error loading CSV:", err)
	}

	fmt.Printf("✅ Loaded %d translations\n", len(transDB.GetAllEntries()))
	fmt.Println("Type a Japanese phrase to get the English translation")
	fmt.Println("Type 'exit' to quit")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("JP: ")
		input, err := reader.ReadString('\n')
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

		translation := transDB.GetTranslation(input)
		if translation != "" {
			fmt.Printf("EN: %s\n", translation)
		} else {
			found := false
			for jp, en := range transDB.GetAllEntries() {
				if strings.Contains(jp, input) {
					fmt.Printf("JP: %s\n", jp)
					fmt.Printf("EN: %s\n", en)
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("EN: (not found in CSV)\n")
			}
		}
		fmt.Println()
	}

	fmt.Println("\nSayonara!")
}
