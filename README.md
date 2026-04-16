# AI-AutoComplete

A lightweight N-Gram language model implementation in Go for Japanese sentence completion and grammar analysis.

## Features

- Train a custom N-Gram model from scratch
- Complete Japanese sentences based on learned patterns
- Grammar lookup from CSV database
- Lightweight - runs on weak PCs (50-100MB RAM)
- No external dependencies or large model downloads

## Execution Flow

STEP 1: Initial Command
go run cmd/train/main.go

STEP 2: Main Menu
1. Train new model from corpus
2. Test trained model (sentence completion)
3. Load grammar CSV only

STEP 3: Training Mode
- Load corpus from data/corpus.txt
- Tokenizer converts text to numbers
- Model creates embeddings (32-dim vectors)
- Trainer runs backpropagation for N epochs
- Model saved to internal/data/

STEP 4: Test Mode
- Type beginning of a sentence
- Model predicts next words
- Continues until stop token or max steps

STEP 5: Grammar Mode
- explain word - Shows grammar explanation
- analyze phrase - Analyzes each word
- type word - Shows word type

## Files and Responsibilities

| File | Responsibility |
|------|----------------|
| cmd/train/main.go | Orchestrates the whole process |
| internal/data/loader.go | Loads CSV grammar database |
| internal/tokenizer/tokenizer.go | Converts text to numbers |
| internal/model/ngram.go | Neural network model |
| internal/model/layer.go | Linear layer |
| internal/model/activation.go | Softmax function |
| internal/model/predictor.go | Sentence completion logic |
| internal/trainer/trainer.go | Training with backpropagation |
| internal/storage/storage.go | Save/load trained model |
| internal/ui/cli.go | User interface |
| data/corpus.txt | Training sentences |
| data/grammar.csv | Grammar reference data |

## Resource Usage on Weak PC

- RAM: 50-100 MB
- CPU: 10-30%
- Disk: 5-10 MB
- Training time: 30-60 seconds
- Inference time: Less than 1ms

## Commands to Run

mkdir -p data
echo "watashi wa nihonjin desu" > data/corpus.txt
go run cmd/train/main.go

## Requirements

- Go 1.21 or higher
- No GPU required
- Works on any PC (Windows, Linux, Mac)
