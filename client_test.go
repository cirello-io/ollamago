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
	"net/http"
	"net/http/httptest"
	"testing"

	"cirello.io/ollamago"
	"github.com/stretchr/testify/require"
)

func TestGenerateCompletion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/generate", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"model":"test","response":"test response","done":true,"total_duration":1000}`))
	}))
	t.Cleanup(server.Close)
	client := ollamago.Client{BaseURL: server.URL}
	respChan, err := client.GenerateCompletion(ollamago.CompletionRequest{
		Model:  "test",
		Prompt: "test prompt",
		Stream: false,
	})
	require.NoError(t, err)
	resp := <-respChan
	require.Equal(t, "test response", resp.Response)
	require.Equal(t, "test", resp.Model)
}

func TestGenerateEmbeddings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/embed", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"model":"test","embeddings":[[1.0,2.0,3.0]],"total_duration":1000}`))
	}))
	t.Cleanup(server.Close)
	client := ollamago.Client{BaseURL: server.URL}
	resp, err := client.GenerateEmbeddings(ollamago.EmbedRequest{
		Model: "test",
		Input: []string{"test input"},
	})
	require.NoError(t, err)
	require.Equal(t, "test", resp.Model)
	require.Len(t, resp.Embeddings, 1)
	require.Equal(t, []float64{1.0, 2.0, 3.0}, resp.Embeddings[0])
}

func TestGenerateChat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/chat", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"model":"test","message":{"role":"assistant","content":"hello"},"done":true,"total_duration":1000}`))
	}))
	t.Cleanup(server.Close)
	client := ollamago.Client{BaseURL: server.URL}
	respChan, err := client.GenerateChat(ollamago.ChatRequest{
		Model: "test",
		Messages: []ollamago.ChatMessage{{
			Role:    "user",
			Content: "hi",
		}},
		Stream: false,
	})
	require.NoError(t, err)
	resp := <-respChan
	require.Equal(t, "test", resp.Model)
	require.Equal(t, "hello", resp.Message.Content)
	require.Equal(t, "assistant", resp.Message.Role)
}

func TestListModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/tags", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"models":[{"name":"model1","modified_at":"2023-01-01T00:00:00Z","size":1024}]}`))
	}))
	t.Cleanup(server.Close)
	client := ollamago.Client{BaseURL: server.URL}
	resp, err := client.ListModels()
	require.NoError(t, err)
	require.Len(t, resp.Models, 1)
	require.Equal(t, "model1", resp.Models[0].Name)
	require.Equal(t, int64(1024), resp.Models[0].Size)
}

func TestShowModelInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/show", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"modelfile":"test file","details":{"format":"gguf","parameter_size":"7B"}}`))
	}))
	t.Cleanup(server.Close)
	client := ollamago.Client{BaseURL: server.URL}
	resp, err := client.ShowModelInfo(ollamago.ShowModelRequest{
		Model: "test",
	})
	require.NoError(t, err)
	require.Equal(t, "test file", resp.Modelfile)
	require.Equal(t, "gguf", resp.Details.Format)
	require.Equal(t, "7B", resp.Details.ParameterSize)
}

func TestDeleteModel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/delete", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	client := ollamago.Client{BaseURL: server.URL}
	err := client.DeleteModel(ollamago.DeleteModelRequest{
		Model: "test",
	})
	require.NoError(t, err)
}

func TestVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/version", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version":"1.0.0"}`))
	}))
	t.Cleanup(server.Close)
	client := ollamago.Client{BaseURL: server.URL}
	version, err := client.Version()
	require.NoError(t, err)
	require.Equal(t, "1.0.0", version)
}
