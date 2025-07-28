package iohandler

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// escapeCSVField escapes a field for CSV output by doubling quotes and wrapping the field in quotes, so that if it contains commas or newlines, it will be correctly interpreted by CSV parsers.
func escapeCSVField(field string) string {
	field = strings.ReplaceAll(field, `"`, `""`)
	return `"` + field + `"`
}

// WriteFile generates the CSV file providing responses to the questions
func WriteFile(destPath string, answers map[string]string) error {
	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for question, answer := range answers {
		escapedQuestion := escapeCSVField(question)
		escapedAnswer := escapeCSVField(answer)
		if _, err := writer.WriteString(fmt.Sprintf("%s,%s\n", escapedQuestion, escapedAnswer)); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}
