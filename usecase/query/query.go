package query

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"
)

type UsecaseQuery interface {
	QueryToOpenAI(query string, outputFormat string, responseLanguage string, filePaths []string) (string, error)
}

type OpenAIClient interface {
	Ask(ctx context.Context, query string, outputFormat string, responseLanguage string, textFiles []InputTextFile) (string, error)
}

type InputTextFile struct {
	FilePath  string
	Content   string
	Extension string
}

type usecaseQuery struct {
	openAIClient OpenAIClient
}

func NewUsecaseQuery(openAIClient OpenAIClient) UsecaseQuery {
	return &usecaseQuery{
		openAIClient: openAIClient,
	}
}

func (uq *usecaseQuery) QueryToOpenAI(
	query string,
	outputFormat string,
	responseLanguage string,
	filePaths []string) (string, error) {
	textFiles := make([]InputTextFile, 0, len(filePaths))
	if len(filePaths) > 0 {
		for _, fp := range filePaths {
			if _, err := os.Stat(fp); os.IsNotExist(err) {
				return "", fmt.Errorf("file does not exist: %s", fp)
			}
			isText, textContent, err := isTextFile(fp)
			if err != nil {
				return "", fmt.Errorf("check if file is text file: %w", err)
			}
			if !isText {
				return "", fmt.Errorf("file is not a text file: %s", fp)
			}
			textFile := InputTextFile{
				FilePath:  fp,
				Content:   textContent,
				Extension: filepath.Ext(fp),
			}
			textFiles = append(textFiles, textFile)
		}
	}
	res, err := uq.openAIClient.Ask(context.Background(), query, outputFormat, responseLanguage, textFiles)
	if err != nil {
		return "", fmt.Errorf("ask to openai: %w", err)
	}
	return res, nil
}

func isTextFile(filePath string) (bool, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var contentBytes []byte
	for scanner.Scan() {
		contentBytes = scanner.Bytes()
		if !utf8.Valid(contentBytes) {
			return false, "", nil
		}
	}
	if err := scanner.Err(); err != nil {
		return true, string(contentBytes), fmt.Errorf("scanning file: %w", err)
	}
	return true, string(contentBytes), nil
}
