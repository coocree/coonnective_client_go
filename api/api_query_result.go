package api

import (
	"regexp"
	"time"
)

// A struct ApiQueryResult é usada para representar o resultado de uma query da API.
type QueryResultModel struct {
	Errors    []string   // Lista de erros retornados pela API, se houver.
	Timestamp *time.Time // Timestamp do resultado da query.
	Code      string     // Código do resultado da query.
}

// Construtor de inicialização de campo.
func QueryResult(errorMessage string) QueryResultModel {
	errors := extractErrors(errorMessage)
	timestamp := extractTimestamp(errorMessage)
	code := extractCode(errorMessage)
	return QueryResultModel{Errors: errors, Timestamp: timestamp, Code: code}
}

// Extrai os erros da mensagem de erro fornecida.
func extractErrors(errorMessage string) []string {
	re1 := regexp.MustCompile(`errors: \[(.*?)\]`)
	errorsMatches := re1.FindStringSubmatch(errorMessage)
	if len(errorsMatches) == 0 {
		return nil
	}

	errorString := errorsMatches[1]
	re2 := regexp.MustCompile(`message: (.*?),`)
	matches := re2.FindAllStringSubmatch(errorString, -1)
	errors := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) >= 2 {
			errors = append(errors, match[1])
		}
	}
	return errors
}

// Extrai o timestamp da mensagem de erro fornecida.
func extractTimestamp(errorMessage string) *time.Time {
	re := regexp.MustCompile(`timestamp: (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3})`)
	match := re.FindStringSubmatch(errorMessage)
	if len(match) == 0 {
		return nil
	}
	timestampString := match[1]
	timestamp, _ := time.Parse("2006-01-02 15:04:05.000", timestampString)
	return &timestamp
}

// Extrai o código da mensagem de erro fornecida.
func extractCode(errorMessage string) string {
	re := regexp.MustCompile(`code: (\w+)`)
	match := re.FindStringSubmatch(errorMessage)
	if len(match) == 0 {
		return ""
	}
	return match[1]
}
