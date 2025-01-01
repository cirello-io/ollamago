// Copyright 2024 cirello.io/ollamago & U. Cirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ollamago_test

import (
	"context"
	"fmt"

	"cirello.io/ollamago"
)

func Example_Client_GenerateChat() {
	c := ollamago.Client{}
	resp, err := c.GenerateChat(context.Background(), ollamago.ChatRequest{
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
	// Output:
	// Hello World!
}
