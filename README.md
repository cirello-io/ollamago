# ollamago

A Go SDK for [Ollama](https://github.com/ollama/ollama), providing a simple interface to interact with Ollama's API.

## Installation

```sh
go get cirello.io/ollamago
```

## Usage

```go
package main

import (
	"fmt"
	"cirello.io/ollamago"
)

func main() {
	c := ollamago.Client{}
	resp, err := c.GenerateChat(ollamago.ChatRequest{
		Model: "llama3.2",
		Messages: []ollamago.ChatMessage{{
			Role:    "user",
			Content: "say hello world",
		}},
		Stream: true,
	})
	if err != nil {
		panic(err)
	}
	for r := range resp {
		if r.Done {
			break
		}
		fmt.Print(r.Message.Content)
	}
}
```

## Features

Mostly implements the contents of https://github.com/ollama/ollama/blob/main/docs/api.md:

- Generate a completion
- Generate a chat completion
- Create a Model
- List Local Models
- Show Model Information
- Copy a Model
- Delete a Model
- Pull a Model
- Push a Model
- Generate Embeddings
- List Running Models
- Version
