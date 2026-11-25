package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	//"github.com/ollama/ollama"
)

type OllamaClient struct {
	client *http.Client
	url    string
}

func NewOllamaClient() *OllamaClient {
	return &OllamaClient{
		client: &http.Client{},
		url:    "http://localhost:11434/api/chat",
	}
}

func (o *OllamaClient) GetProperties() *Properties {
	botProperties := GetBotProperties()
	return botProperties
}

func (o *OllamaClient) GetClient() interface{} {
	return o.client
}

func (o *OllamaClient) GetHistory() []Message {
	return GetHistory(o)
}

func (o *OllamaClient) AddToChatHistory(msg Message) {
	AddToChatHistory(o, msg)
}

func (o *OllamaClient) Ask(query string) (string, error) {
	o.AddToChatHistory(Message{Role: "User", Content: query})
	req := Request{
		Model:    o.GetProperties().Model,
		Stream:   false,
		Messages: o.GetHistory(),
	}
	js, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest(http.MethodPost, o.url, bytes.NewReader(js))
	if err != nil {
		return "", err
	}

	resp, err := o.client.Do(httpReq)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var ollamaResp Response
	err = json.NewDecoder(resp.Body).Decode(&ollamaResp)

	if err != nil {
		return "", err
	}
	o.AddToChatHistory(ollamaResp.Message)
	return ollamaResp.Message.Content, nil
}

func (o *OllamaClient) AskStream(query string) (io.ReadCloser, error) {
	o.AddToChatHistory(Message{Role: "User", Content: query})
	req := Request{
		Model:    o.GetProperties().Model,
		Stream:   true,
		Messages: []Message{{Role: "User", Content: query}},
	}

	js, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, o.url, bytes.NewReader(js))
	if err != nil {
		return nil, err
	}

	httpResp, err := o.client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	// Create a pipe for streaming
	pr, pw := io.Pipe()

	go func() {
		defer httpResp.Body.Close()
		defer pw.Close()

		decoder := json.NewDecoder(httpResp.Body)
		for {
			var ollamaResp Response
			if err := decoder.Decode(&ollamaResp); err == io.EOF {
				break
			} else if err != nil {
				err := pw.CloseWithError(fmt.Errorf("error reading stream: %v", err))
				if err != nil {
					return
				}
				return
			}

			// Write only the message content to the pipe
			if ollamaResp.Message.Content != "" {
				_, _ = pw.Write([]byte(TrimResponse(ollamaResp.Message.Content) + "\n"))
			}

			// Stop if the stream is marked as done
			if ollamaResp.Done {
				break
			}
		}
	}()

	return pr, nil
}

func (o *OllamaClient) RunBot() {
	RunBot(o)
}

func getOllamaModelsPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Unable to determine home directory")
	}
	return filepath.Join(homeDir, ".ollama", "models", "manifests", "registry.ollama.ai", "library")
}

// findOllamaModels searches for locally downloaded Ollama models.
// It returns an array of strings sorted alphabetically.
func findOllamaModels(modelsPath string) ([]string, error) {
	var models []string

	files, err := os.ReadDir(modelsPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			models = append(models, file.Name())
		}
	}

	sort.Strings(models)
	return models, nil
}

func ModelsHandler(w http.ResponseWriter, r *http.Request) {
	modelsPath := getOllamaModelsPath()
	fmt.Println("Models path:", modelsPath)
	models, err := findOllamaModels(modelsPath)
	if err != nil {
		http.Error(w, "Failed to retrieve models", http.StatusInternalServerError)
		return
	}

	// Prepare the response structure with models as an array
	response := struct {
		Models []string `json:"models"`
	}{
		Models: models, // Assuming 'models' is a slice of strings
	}

	// Set response header and return JSON-encoded data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set the status code to 200 OK
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
