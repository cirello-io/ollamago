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

package ollamago

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type CompletionRequest struct {
	Model   string          `json:"model"`
	Prompt  string          `json:"prompt,omitempty"`
	Options ModelParameters `json:"options,omitempty"`
	Stream  bool            `json:"stream,omitempty"`
}

type CompletionResponse struct {
	Model         string        `json:"model"`
	Response      string        `json:"response"`
	Done          bool          `json:"done"`
	TotalDuration time.Duration `json:"total_duration"`
	Error         error         `json:"error,omitempty"`
}

func (c *Client) baseURL() string {
	if c.BaseURL == "" {
		return "http://localhost:11434"
	}
	return c.BaseURL
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient == nil {
		return http.DefaultClient
	}
	return c.HTTPClient
}

func (c *Client) GenerateCompletion(ctx context.Context, req CompletionRequest) (<-chan CompletionResponse, error) {
	url := c.baseURL() + "/api/generate"
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare CompletionRequest: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("cannot prepare HTTP CompletionRequest: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("cannot execute HTTP CompletionRequest: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, errors.New("failed to generate completion: " + resp.Status)
	}
	out := make(chan CompletionResponse)
	go func() {
		defer resp.Body.Close()
		defer close(out)
		dec := json.NewDecoder(resp.Body)
		for {
			var res CompletionResponse
			err := dec.Decode(&res)
			if errors.Is(err, io.EOF) {
				out <- res
				return
			} else if err != nil {
				res.Error = err
				out <- res
				return
			}
			out <- res
		}
	}()
	return out, nil
}

type EmbedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type EmbedResponse struct {
	Model      string        `json:"model"`
	Embeddings [][]float64   `json:"embeddings"`
	Duration   time.Duration `json:"total_duration"`
}

func (c *Client) GenerateEmbeddings(ctx context.Context, req EmbedRequest) (*EmbedResponse, error) {
	url := c.baseURL() + "/api/embed"
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare EmbedRequest: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("cannot prepare HTTP EmbedRequest: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("cannot execute HTTP EmbedRequest: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to generate embeddings: %s", resp.Status)
	}
	var embedResp EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, fmt.Errorf("cannot decode embed response: %w", err)
	}
	return &embedResp, nil
}

type ChatRequest struct {
	Model    string          `json:"model"`
	Messages []ChatMessage   `json:"messages"`
	Stream   bool            `json:"stream,omitempty"`
	Options  ModelParameters `json:"options,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Model         string        `json:"model"`
	Message       ChatMessage   `json:"message"`
	Done          bool          `json:"done"`
	TotalDuration time.Duration `json:"total_duration"`
	Error         error         `json:"error,omitempty"`
}

func (c *Client) GenerateChat(ctx context.Context, req ChatRequest) (<-chan ChatResponse, error) {
	url := c.baseURL() + "/api/chat"
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare ChatRequest: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("cannot prepare HTTP ChatRequest: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("cannot execute HTTP ChatRequest: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to generate chat: %s", resp.Status)
	}
	out := make(chan ChatResponse)
	go func() {
		defer resp.Body.Close()
		defer close(out)
		dec := json.NewDecoder(resp.Body)
		for {
			var res ChatResponse
			err := dec.Decode(&res)
			if errors.Is(err, io.EOF) {
				out <- res
				return
			} else if err != nil {
				res.Error = err
				out <- res
				return
			}
			out <- res
		}
	}()
	return out, nil
}

type ModelInfo struct {
	Name       string    `json:"name"`
	ModifiedAt time.Time `json:"modified_at"`
	Size       int64     `json:"size"`
}

type ListModelsResponse struct {
	Models []ModelInfo `json:"models"`
}

func (c *Client) ListModels(ctx context.Context) (*ListModelsResponse, error) {
	url := c.baseURL() + "/api/tags"
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare HTTP request: %w", err)
	}
	resp, err := c.httpClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("cannot execute HTTP request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list models: %s", resp.Status)
	}
	var listResp ListModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("cannot decode response: %w", err)
	}
	return &listResp, nil
}

type ShowModelRequest struct {
	Model   string `json:"model"`
	Verbose bool   `json:"verbose,omitempty"`
}

type ShowModelResponse struct {
	Modelfile string `json:"modelfile"`
	Details   struct {
		Format        string   `json:"format"`
		ParameterSize string   `json:"parameter_size"`
		Quantization  string   `json:"quantization_level"`
		Family        string   `json:"family"`
		Families      []string `json:"families"`
	} `json:"details"`
}

func (c *Client) ShowModelInfo(ctx context.Context, req ShowModelRequest) (*ShowModelResponse, error) {
	url := c.baseURL() + "/api/show"
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare ShowModelRequest: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("cannot prepare HTTP ShowModelRequest: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("cannot execute HTTP ShowModelRequest: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to show model info: %s", resp.Status)
	}
	var showResp ShowModelResponse
	if err := json.NewDecoder(resp.Body).Decode(&showResp); err != nil {
		return nil, fmt.Errorf("cannot decode show response: %w", err)
	}
	return &showResp, nil
}

type DeleteModelRequest struct {
	Model string `json:"model"`
}

func (c *Client) DeleteModel(ctx context.Context, req DeleteModelRequest) error {
	url := c.baseURL() + "/api/delete"
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("cannot prepare DeleteModelRequest: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("cannot prepare HTTP DeleteModelRequest: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient().Do(httpReq)
	if err != nil {
		return fmt.Errorf("cannot execute HTTP DeleteModelRequest: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete model: %s", resp.Status)
	}
	return nil
}

func (c *Client) Version(ctx context.Context) (string, error) {
	url := c.baseURL() + "/api/version"
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("cannot prepare HTTP request: %w", err)
	}
	resp, err := c.httpClient().Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("cannot execute HTTP request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get version: %s", resp.Status)
	}
	var versionResp struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
		return "", fmt.Errorf("cannot decode version response: %w", err)
	}
	return versionResp.Version, nil
}

type ModelParameters struct {
	// Mirostat enables Mirostat sampling for controlling perplexity.
	// (0 = disabled, 1 = Mirostat, 2 = Mirostat 2.0)
	Mirostat int `json:"mirostat,omitempty"`

	// MirostatEta influences how quickly the algorithm responds to
	// feedback. Lower learning rate = slower adjustments, higher = more
	// responsive.
	MirostatEta float64 `json:"mirostat_eta,omitempty"`

	// MirostatTau controls balance between coherence and diversity.
	// Lower value results in more focused and coherent text.
	MirostatTau float64 `json:"mirostat_tau,omitempty"`

	// NumCtx sets the size of the context window for next token generation.
	NumCtx int `json:"num_ctx,omitempty"`

	// RepeatLastN sets how far back to look to prevent repetition.
	// (0 = disabled, -1 = num_ctx)
	RepeatLastN int `json:"repeat_last_n,omitempty"`

	// RepeatPenalty sets how strongly to penalize repetitions.
	// Higher value penalizes more strongly, lower is more lenient.
	RepeatPenalty float64 `json:"repeat_penalty,omitempty"`

	// Temperature controls creativity of the model's responses.
	// Higher temperature increases creativity.
	Temperature float64 `json:"temperature,omitempty"`

	// Seed sets the random number seed for generation.
	// Specific seed generates same text for same prompt.
	Seed int `json:"seed,omitempty"`

	// Stop sets the stop sequences to use.
	// Model stops generating when this pattern is encountered.
	Stop string `json:"stop,omitempty"`

	// TfsZ controls tail free sampling to reduce impact of less probable tokens.
	// Higher value reduces impact more, 1.0 disables.
	TfsZ float64 `json:"tfs_z,omitempty"`

	// NumPredict sets maximum number of tokens to predict.
	// -1 for infinite generation.
	NumPredict int `json:"num_predict,omitempty"`

	// TopK reduces probability of nonsense generation.
	// Higher value gives more diverse answers, lower is more conservative.
	TopK int `json:"top_k,omitempty"`

	// TopP works with top-k for diversity control.
	// Higher value leads to more diverse text, lower is more focused.
	TopP float64 `json:"top_p,omitempty"`

	// MinP sets minimum probability for token consideration.
	// Alternative to top_p for balancing quality and variety.
	MinP float64 `json:"min_p,omitempty"`
}
