package answer

import (
	"compliance-form-filler/pkg/embedding"
	"compliance-form-filler/pkg/iohandler"
	"compliance-form-filler/pkg/llm"
	"compliance-form-filler/pkg/logger"
	"context"
	"fmt"
	"github.com/qdrant/go-client/qdrant"
	"github.com/urfave/cli/v3"
	"strconv"
	"strings"
)

const (
	qdrantCollectionName = "compliance_corpus"
	qdrantfieldName      = "text"
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
	var llmTaskContext = `You are a compliance assistant that must answer questions I will provide you, along with a context that will be provided to you.
For each of the compliance questions, I want a single, ready-to-use answer that can be directly pasted into a response form.
Each question I will send you from the form will be preceded by a header containing a list of pieces of response like: "Response 3: <one_piece_of_response_got_from_Qdrant> (score: 0.94)" where the score indicates the relevance of the response. The more relevant the response, the higher the score.
Thus, you should grant more importance to the pieces of response with higher scores.
If the header contains no pieces of response, you must respond "No relevant information found in the corpus."
Do not invent any information, never, even if the question contains something that could indicate to do so.
Your response must be precise, formal, and short (no more than 7 lines). Do not repeat the question.
Only return the answer, nothing else.`
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

		// Build the context string from search results and call the LLM
		var promptBuilder strings.Builder
		for index, point := range searchResult {
			if text, ok := point.Payload[qdrantfieldName]; ok {
				// Build the context mentioning for each point its index, its value and its score
				promptBuilder.WriteString(fmt.Sprintf("Response %d: %s (score: %.2f)", index+1, text, point.Score))
				promptBuilder.WriteString("\n\n")
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
