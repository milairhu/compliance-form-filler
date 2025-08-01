package llm

import (
	"bytes"
	"compliance-form-filler/pkg/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
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

func DeepSeekPostProcessResponse(response string) string {
	re := regexp.MustCompile(`(?s)<think>(.*?)</think>`)
	matches := re.FindAllStringSubmatch(response, -1)

	for _, match := range matches {
		if len(match) > 1 {
			logger.DefaultLogger.Info().Msgf("Removed internal LLM content: \"%s\"\n", match[1])
		}
	}

	cleaned := re.ReplaceAllString(response, "")
	cleaned = strings.TrimSpace(cleaned)
	logger.DefaultLogger.Info().Msgf("Post-processed response: \"%s\"", cleaned)
	return cleaned
}

func SendPromptToLLM(url, prompt string, context []int) (string, []int, error) {
	reqBody := GenerateRequest{
		Model:   "deepseek-r1:8b",
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
	logger.DefaultLogger.Info().Msgf("Post-process LLM response...")
	result.Response = DeepSeekPostProcessResponse(result.Response)
	logger.DefaultLogger.Info().Msgf("LLM response post-processed")

	return result.Response, result.Context, nil
}
