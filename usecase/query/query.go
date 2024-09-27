//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package query

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
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
			if _, err := os.Stat(fp); err != nil {
				if os.IsNotExist(err) {
					return "", fmt.Errorf("filepath %s does not exists: %w", fp, err)
				}
				return "", fmt.Errorf("while stat file: %w", err)
			}
			isText, err := isTextFile(fp)
			if err != nil {
				return "", fmt.Errorf("checking if file %s is text: %w", fp, err)
			}
			if !isText {
				slog.Warn("File is not a text file. Skipping...", "filepath", fp, "error", err)
				continue
			}
			// will get an extension like .c, .go, .txt, ..., and empty string if no extension
			ext := filepath.Ext(fp)
			if len(ext) > 0 && ext[0] == '.' {
				ext = ext[1:]
			}
			textContentBytes, err := os.ReadFile(fp)
			if err != nil {
				return "", fmt.Errorf("read file: %w", err)
			}
			textFile := InputTextFile{
				FilePath:  fp,
				Content:   string(textContentBytes),
				Extension: ext,
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

func isTextFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if !utf8.Valid(scanner.Bytes()) {
			return false, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return true, fmt.Errorf("scanning file: %w", err)
	}
	return true, nil
}
