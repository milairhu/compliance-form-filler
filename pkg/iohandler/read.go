package iohandler

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadFile reads a file and returns its content as a slice of strings.
// WARNING: We assume one line is one question to be asked to the LLM.
func ReadFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var questions []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			questions = append(questions, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file: %w", err)
	}

	return questions, nil
}
