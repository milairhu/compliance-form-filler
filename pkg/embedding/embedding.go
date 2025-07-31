package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type EmbedRequest struct {
	Texts []string `json:"texts"`
}

type EmbedResponse struct {
	Vectors [][]float32 `json:"vectors"`
}

func EmbedString(str string, url string) ([]float32, error) {
	payload, err := json.Marshal(EmbedRequest{Texts: []string{str}})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to reach embedding API: %w", err)
	}
	defer resp.Body.Close()

	var res EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("invalid embedding response: %w", err)
	}

	if len(res.Vectors) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return res.Vectors[0], nil
}
