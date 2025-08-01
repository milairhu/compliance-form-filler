package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GenerateRequest struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Stream  bool   `json:"stream"`
	Context []int  `json:"context"`
}

type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Context  []int  `json:"context,omitempty"`
}

func SendPromptToLLM(url, prompt string, context []int) (string, []int, error) {
	reqBody := GenerateRequest{
		Model:   "mistral",
		Prompt:  prompt,
		Stream:  false,
		Context: context,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", []int{}, fmt.Errorf("failed to marshal prompt: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", []int{}, fmt.Errorf("failed to send request to LLM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", []int{}, fmt.Errorf("LLM responded with status %d: %s", resp.StatusCode, string(body))
	}

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", []int{}, fmt.Errorf("failed to decode LLM response: %w", err)
	}

	if !result.Done {
		return "", []int{}, fmt.Errorf("LLM response generation not done: %s", result.Response)
	}

	return result.Response, result.Context, nil
}
