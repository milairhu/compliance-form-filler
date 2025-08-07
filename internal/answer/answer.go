package answer

import (
	"compliance-form-filler/pkg/embedding"
	"compliance-form-filler/pkg/iohandler"
	"compliance-form-filler/pkg/llm"
	"compliance-form-filler/pkg/logger"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/qdrant/go-client/qdrant"
	"github.com/urfave/cli/v3"
)

const (
	qdrantCollectionName  = "compliance_corpus"
	qdrantTextFieldName   = "text"
	qdrantSourceFieldName = "source"
)

func Answer(cmd *cli.Command) error {
	// Read the flags and perform the necessary actions
	if cmd == nil {
		return fmt.Errorf("nil command")
	}
	sourceFile := cmd.String("source-file")
	outputFile := cmd.String("output-file")
	qdrantURL := cmd.String("qdrant-url")
	llmURL := cmd.String("llm-url")
	embeddingApiURL := cmd.String("embedding-api-url")

	qdrantHost, qdrantPort, err := ProcessUrl(qdrantURL)
	if err != nil {
		return fmt.Errorf("failed to process Qdrant URL: %w", err)
	}

	// Read the source file and process it
	logger.DefaultLogger.Info().Msgf("Processing source file: %s ...", sourceFile)
	questions, err := iohandler.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}
	logger.DefaultLogger.Info().Msgf("Questions parsed!")

	qdrantClient, err := qdrant.NewClient(&qdrant.Config{
		Host: qdrantHost,
		Port: qdrantPort,
	})
	if err != nil {
		return fmt.Errorf("failed to create Qdrant client: %w", err)
	}

	questionsAnswered := make(map[string]string)
	var llmTaskContext = `You are a compliance assistant. You answer each question **only** using the provided context header (a ranked list of snippets like: "Response 3: <text> (score: 0.94) (source: <title of the source document>").

Rules
1) **Use only the header content.** If the answer cannot be found in the header, reply exactly: **"No information available"**.
2) **Never invent or infer beyond the header.** Do not rely on prior knowledge or assumptions.
3) **Ranking & selection.**
   - Prefer higher score snippets.
   - When snippets conflict, choose the highest-scoring. If still tied, choose the one most specific to the question.
   - If evidence is partial or ambiguous, reply **"No information available"**.
4) **Precision & completeness.**
   - Extract the best answer and compile them into a short, ready-to-use response.
   - If the question implies a "Yes"/"No" question, reply "Yes” or “No” and justify the answer. **Never let an answer be only "Yes" or "No"**. If not clearly supported, reply **"No information available"**.
5) **Output format.**
   - Style: precise, formal, and concise; no preamble.
   - Length: maximum 7 lines.
   - Return **only** the final answer, no restatements of the question, no references to scores or snippets.
6) **Keep in mind the questions are addressed to the company, not to you**. If a questions contains "you", it means the company, not you as an AI assistant.

Process (follow silently)
a) Read the question and header.
b) From the snippets, resolve conflicts (highest score).
c) If a direct answer is present, output it verbatim or lightly edited for grammar; otherwise output **"No information available"**. **Never let an answer be only "Yes" or "No"** .`

	// Prepare and send the context prompt for the LLM
	logger.DefaultLogger.Info().Msgf("Sending context to LLM: %s", llmTaskContext)
	_, llmContext, err := llm.SendPromptToLLM(llmURL, llmTaskContext, []int{})
	if err != nil {
		return fmt.Errorf("failed to send prompt to LLM: %w", err)
	}
	logger.DefaultLogger.Info().Msgf("Searching for answers to %d questions...", len(questions))
	for _, question := range questions {
		// Vectorize the question using the embedding API
		logger.DefaultLogger.Info().Msgf("Embedding question: %s", question)
		vector, err := embedding.EmbedString(question, embeddingApiURL)
		if err != nil {
			logger.DefaultLogger.Error().Msgf("failed to vectorize question: %s", question)
			continue
		}
		logger.DefaultLogger.Info().Msgf("Question vectorized successfully")

		// Search in Qdrant using the vector
		logger.DefaultLogger.Info().Msgf("Searching in Qdrant for question: %s", question)
		var scoreThreshold float32 = 0.4
		searchResult, err := qdrantClient.Query(context.Background(), &qdrant.QueryPoints{
			CollectionName: qdrantCollectionName,
			Query:          qdrant.NewQuery(vector...),
			WithPayload:    qdrant.NewWithPayload(true),
			ScoreThreshold: &scoreThreshold,
		})
		if err != nil {
			logger.DefaultLogger.Error().Msgf("qdrant search failed for question: %s - %s", question, err)
			continue
		}
		logger.DefaultLogger.Info().Msgf("Qdrant search completed")
		if len(searchResult) == 0 {
			logger.DefaultLogger.Warn().Msgf("No results found for this question")
			// Store a default answer if no results found
			questionsAnswered[question] = "No information available"
			continue
		} else {
			// Build the context string from search results and call the LLM
			var promptBuilder strings.Builder
			for index, point := range searchResult {
				if text, ok := point.Payload[qdrantTextFieldName]; ok {
					// Build the context mentioning for each point its index, its value and its score
					if source, ok := point.Payload[qdrantSourceFieldName]; ok {
						promptBuilder.WriteString(fmt.Sprintf("Response %d: %s (score: %.2f) (source: %s)", index+1, text.GetStringValue(), point.Score, source.GetStringValue()))
						promptBuilder.WriteString("\n\n")
					}
				}
			}
			prompt := promptBuilder.String()
			// Prepare the full prompt for the LLM
			prompt = fmt.Sprintf("%s\n\n ===== %s", prompt, question)
			// Call the LLM with the prompt
			logger.DefaultLogger.Info().Msgf("Sending prompt to LLM: %s", prompt)
			answer, _, err := llm.SendPromptToLLM(llmURL, prompt, llmContext)
			if err != nil {
				return fmt.Errorf("failed to send prompt to LLM: %w", err)
			}
			logger.DefaultLogger.Info().Msgf("LLM response received for question: %s", question)
			// Store the answer in the map
			questionsAnswered[question] = answer
		}
	}
	logger.DefaultLogger.Info().Msgf("All questions processed, %d answers generated", len(questionsAnswered))

	// Save the results to the output file
	logger.DefaultLogger.Info().Msgf("Saving answers to output file: %s ...", outputFile)
	err = iohandler.WriteFile(outputFile, questionsAnswered)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}
	logger.DefaultLogger.Info().Msgf("Answers saved to successfully!")

	return nil
}

// ProcessUrl get host and port from the URL
func ProcessUrl(url string) (string, int, error) {
	parts := strings.Split(url, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid URL format, expected 'host:port'")
	}
	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid port number: %w", err)
	}
	return host, port, nil
}
