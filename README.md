# AI-AutoComplete

**⚠️ This is an educational project for learning Machine Learning from scratch in Go.**

A lightweight N-Gram language model implementation in Go for Japanese sentence completion and grammar analysis. Built to understand the fundamentals of neural networks, embeddings, and backpropagation.

## 🎯 Educational Purpose

This project was created to learn by doing how a language model works, without using pre-built libraries like TensorFlow or PyTorch. You will see:

- How words are converted into numbers (tokenization)
- How to create embeddings (vectors that represent meaning)
- How a simple neural network works (linear layer + softmax)
- How training with backpropagation adjusts weights
- How an N-Gram model predicts the next word

## ⚠️ Technical Limitations

- N-Gram Model: Does not understand long context, only repeats statistical patterns
- Limited Learning: With small datasets, the model does not generalize well
- For real autocomplete: Consider using Ollama or larger pre-trained models
- Training data: The model only learns what is in your CSV/corpus file

## 🚀 Current Features

### 1. Train Model from CSV
- Load Japanese sentences from CSV file
- Tokenize and convert to numerical vectors
- Train N-Gram model with configurable epochs
- Save trained model weights

### 2. Test Model (Autocomplete)
- Type the beginning of a Japanese sentence
- Model predicts and completes the sentence
- Shows English translation if available in CSV

### 3. CSV Search Mode (No ML)
- Direct lookup in CSV file
- Fast and reliable for exact matches
- Partial matching for similar phrases

## 🔧 How It Works

### Step 1: Tokenization
Text is split into words and each word is mapped to a unique number.

### Step 2: Embedding
Each word becomes a vector of 32 numbers (embeddings) that represent its meaning.

### Step 3: Neural Network
- Takes two word embeddings as input
- Passes through a linear layer
- Applies softmax to get probabilities for the next word

### Step 4: Training (Backpropagation)
- Calculates loss (cross-entropy)
- Computes gradients
- Updates weights and embeddings
- Repeats for N epochs

### Step 5: Prediction
- User types beginning of a sentence
- Model predicts the most likely next word
- Continues until max steps or low probability

## 📊 Training Results Example

With 562 sentences and 50 epochs:
- Initial loss: ~6.36
- Final loss: ~0.19
- Model learns to complete common patterns

## 💻 Commands

go run cmd/train/main.go
- Choose option 1 to train new model
- Choose option 2 to test trained model
- Choose option 3 for CSV search only (fast, no ML)

## 📝 Example Usage

You: 最近
JP: 最近考えることが多過ぎる。
EN: I have too many things on my mind these days.

You: 私は
JP: 私はキャビアを食べた。
EN: I ate caviar.

You: すぐに
JP: すぐに戻ります。
EN: I will be back soon.

## 🔬 Technical Components

- Tokenizer: Converts text to numbers and back
- Embedding: 32-dimensional word vectors
- Linear Layer: Matrix multiplication (32 × vocab_size)
- Softmax: Converts logits to probabilities
- Cross-Entropy Loss: Measures prediction error
- Backpropagation: Updates weights using gradients

## 💾 Resource Usage

- RAM: 50-100 MB
- CPU: 10-30%
- Disk: 5-10 MB
- Training time: 30-60 seconds (562 sentences, 50 epochs)
- Inference time: Less than 1ms

## 📋 Requirements

- Go 1.21 or higher
- No GPU required
- Works on any PC (Windows, Linux, Mac)

## 🎓 What You Learn

- Fundamentals of word embeddings
- How neural networks process text
- Training loops and loss functions
- Gradient descent and backpropagation
- Limitations of simple language models

## 📚 Acknowledgments

This is a simplified implementation for educational purposes. For production use, consider established frameworks like TensorFlow, PyTorch, or Ollama.