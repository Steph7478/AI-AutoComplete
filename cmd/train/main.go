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
	fmt.Println("1. Train new model from corpus")
	fmt.Println("2. Test trained model (sentence completion)")
	fmt.Println("3. Load grammar CSV only")
	fmt.Print("\nChoose (1/2/3): ")

	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	switch option {
	case "1":
		trainModel(reader)
	case "2":
		testModel()
	case "3":
		loadGrammarOnly()
	default:
		fmt.Println("Invalid option")
	}
}

func trainModel(reader *bufio.Reader) {
	fmt.Println("\n=== TRAINING FROM CORPUS ===\n")

	corpusPath := "data/corpus.txt"
	if _, err := os.Stat(corpusPath); os.IsNotExist(err) {
		fmt.Println("No corpus.txt found. Generating from CSV...")
		generateCorpusFromCSV()
	}

	dataBytes, err := os.ReadFile(corpusPath)
	if err != nil {
		log.Fatal("Error reading corpus:", err)
	}

	text := string(dataBytes)
	fmt.Printf("✅ Corpus loaded: %d characters\n", len(text))

	tok := tokenizer.NewTokenizer()
	tok.BuildFromText(text)
	fmt.Printf("✅ Vocabulary: %d unique words\n", tok.VocabSize)

	tokens := tok.Encode(text)
	fmt.Printf("✅ Total tokens: %d\n", len(tokens))

	embeddingDim := 32
	myModel := model.NewNGramModel(tok.VocabSize, embeddingDim)
	fmt.Printf("✅ Model created\n\n")

	fmt.Print("How many epochs? (default 500): ")
	epochsStr, _ := reader.ReadString('\n')
	epochs := 500
	fmt.Sscanf(epochsStr, "%d", &epochs)

	t := trainer.NewNGramTrainer(myModel, 0.01)
	t.Train(tokens, epochs)

	storage.SaveModel(myModel, tok)
	fmt.Println("\n✅ Model saved to internal/data/")

	predictor := model.NewPredictor(myModel, tok)
	cli := ui.NewCLI(predictor)
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

	predictor := model.NewPredictor(myModel, tok)
	cli := ui.NewCLI(predictor)
	cli.Run()
}

func loadGrammarOnly() {
	fmt.Println("\n=== LOADING GRAMMAR CSV ===\n")

	db := data.NewDatabase()
	csvPaths := []string{
		"data/grammar.csv",
		"./data/grammar.csv",
	}

	loaded := false
	for _, path := range csvPaths {
		if err := db.LoadCSV(path); err == nil {
			fmt.Printf("✅ Loaded grammar from: %s\n", path)
			loaded = true
			break
		}
	}

	if !loaded {
		fmt.Println("❌ Could not find grammar.csv")
		return
	}

	grammarCLI(db)
}

func grammarCLI(db *data.Database) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n=== GRAMMAR LOOKUP ===")
	fmt.Println("Commands:")
	fmt.Println("  explain <word>   - Explain a word")
	fmt.Println("  analyze <phrase> - Analyze a phrase")
	fmt.Println("  type <word>      - Get word type")
	fmt.Println("  exit             - Quit\n")

	for {
		fmt.Print("> ")
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

		parts := strings.SplitN(input, " ", 2)
		cmd := parts[0]
		arg := ""
		if len(parts) > 1 {
			arg = parts[1]
		}

		switch cmd {
		case "explain":
			if arg == "" {
				fmt.Println("Usage: explain <word>")
				continue
			}
			fmt.Printf("\n%s\n\n", db.Explain(arg))
		case "analyze":
			if arg == "" {
				fmt.Println("Usage: analyze <phrase>")
				continue
			}
			fmt.Printf("\n%s\n\n", db.AnalyzePhrase(arg))
		case "type":
			if arg == "" {
				fmt.Println("Usage: type <word>")
				continue
			}
			fmt.Printf("\n%s\n\n", db.GetWordType(arg))
		default:
			fmt.Printf("\n%s\n\n", db.Explain(input))
		}
	}
}

func generateCorpusFromCSV() {
	db := data.NewDatabase()
	if err := db.LoadCSV("data/grammar.csv"); err != nil {
		log.Fatal("Error loading CSV:", err)
	}

	file, err := os.Create("data/corpus.txt")
	if err != nil {
		log.Fatal("Error creating corpus:", err)
	}
	defer file.Close()

	for _, entry := range db.GetAllEntries() {
		line := fmt.Sprintf("%s is a %s. %s. Example: %s\n",
			entry.Word, entry.Type, entry.Explanation, entry.Example)
		file.WriteString(line)
	}

	fmt.Printf("Generated corpus.txt with %d entries\n", len(db.GetAllEntries()))
}
