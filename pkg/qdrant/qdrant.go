package qdrant

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v3"
	"net/http"
	"os"
	"strings"
)

const (
	collectionName = "compliance_corpus"
)

type EmbedRequest struct {
	Text string `json:"text"`
}

type EmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

type SearchRequest struct {
	Vector     []float32 `json:"vector"`
	Top        int       `json:"top"`
	Collection string    `json:"collection"`
}

type SearchResponse struct {
	Result []struct {
		Payload map[string]interface{} `json:"payload"`
	} `json:"result"`
}

type LLMRequest struct {
	Question string `json:"question"`
	Context  string `json:"context"`
}

type LLMResponse struct {
	Answer string `json:"answer"`
}

func fillForm(cmd *cli.Command) error {
	sourceFile := cmd.String("source-file")
	outputFile := cmd.String("output-file")
	qdrantURL := cmd.String("qdrant-url")
	llmURL := cmd.String("llm-url")

	input, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer input.Close()

	reader := csv.NewReader(input)
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read csv: %w", err)
	}

	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer output.Close()
	writer := csv.NewWriter(output)
	defer writer.Flush()

	// Assuming the first column is the question, append the generated answer
	headers := append(rows[0], "Answer")
	_ = writer.Write(headers)

	for _, row := range rows[1:] {
		question := row[0]

		// Embed the question
		embedReq := EmbedRequest{Text: question}
		embedBody, _ := json.Marshal(embedReq)
		embedResp, err := http.Post(llmURL+"/embed", "application/json", bytes.NewReader(embedBody))
		if err != nil {
			return fmt.Errorf("embedding request failed: %w", err)
		}
		var embedResult EmbedResponse
		if err := json.NewDecoder(embedResp.Body).Decode(&embedResult); err != nil {
			return fmt.Errorf("failed to decode embed response: %w", err)
		}
		embedResp.Body.Close()

		// Search Qdrant
		qdrantReq := SearchRequest{
			Vector:     embedResult.Embedding,
			Top:        5,
			Collection: collectionName,
		}
		qdrantBody, _ := json.Marshal(qdrantReq)
		qdrantResp, err := http.Post(qdrantURL+"/search", "application/json", bytes.NewReader(qdrantBody))
		if err != nil {
			return fmt.Errorf("qdrant search failed: %w", err)
		}
		var qdrantResult SearchResponse
		if err := json.NewDecoder(qdrantResp.Body).Decode(&qdrantResult); err != nil {
			return fmt.Errorf("failed to decode qdrant response: %w", err)
		}
		qdrantResp.Body.Close()

		// Collect context from Qdrant results
		var contextParts []string
		for _, item := range qdrantResult.Result {
			if text, ok := item.Payload["text"].(string); ok {
				contextParts = append(contextParts, text)
			}
		}
		context := strings.Join(contextParts, "\n")

		// Ask LLM
		llmReq := LLMRequest{Question: question, Context: context}
		llmBody, _ := json.Marshal(llmReq)
		llmResp, err := http.Post(llmURL+"/generate", "application/json", bytes.NewReader(llmBody))
		if err != nil {
			return fmt.Errorf("llm generation failed: %w", err)
		}
		var llmResult LLMResponse
		if err := json.NewDecoder(llmResp.Body).Decode(&llmResult); err != nil {
			return fmt.Errorf("failed to decode llm response: %w", err)
		}
		llmResp.Body.Close()

		// Write row with answer
		newRow := append(row, llmResult.Answer)
		_ = writer.Write(newRow)
	}

	return nil
}
