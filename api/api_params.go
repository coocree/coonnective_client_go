package api

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type ParamsModel struct {
	Module      string
	Path        string
	Errors      []string
	queryString string
}

func Params(queryString string) ParamsModel {
	queryString = strings.TrimSpace(queryString)
	module := extractModule(queryString)
	path := extractPath(queryString)
	return ParamsModel{Module: module, Path: path, queryString: queryString}
}

// IsValid Verifica se a string de consulta é válida
func (m *ParamsModel) IsValid() *ErrorModel {
	queryString := m.queryString

	if !strings.Contains(queryString, "result {") {
		m.Errors = append(m.Errors, `Campo "result" ausente na string de consulta`)
	}

	if !strings.Contains(queryString, "success") {
		m.Errors = append(m.Errors, `Campo "success" ausente na string de consulta`)
	}

	if strings.Contains(queryString, "error {") {
		if !strings.Contains(queryString, "code") {
			m.Errors = append(m.Errors, `Campo "code" ausente no bloco "error" da string de consulta`)
		}
		if !strings.Contains(queryString, "createdAt") {
			m.Errors = append(m.Errors, `Campo "createdAt" ausente no bloco "error" da string de consulta`)
		}
		if !strings.Contains(queryString, "messages") {
			m.Errors = append(m.Errors, `Campo "messages" ausente no bloco "error" da string de consulta`)
		}
		if !strings.Contains(queryString, "module") {
			m.Errors = append(m.Errors, `Campo "module" ausente no bloco "error" da string de consulta`)
		}
		if !strings.Contains(queryString, "path") {
			m.Errors = append(m.Errors, `Campo "path" ausente no bloco "error" da string de consulta`)
		}
		if !strings.Contains(queryString, "variables") {
			m.Errors = append(m.Errors, `Campo "variables" ausente no bloco "error" da string de consulta`)
		}
	} else {
		m.Errors = append(m.Errors, `Bloco "error" ausente na string de consulta`)
	}

	if !strings.Contains(queryString, "elapsedTime") {
		m.Errors = append(m.Errors, `Campo "elapsedTime" ausente na string de consulta`)
	}

	if len(m.Errors) > 0 {
		apiError := ApiError(m.Errors, m.Module, m.Path, "INVALID_API_PARAMS", time.Now(), nil)
		return &apiError
	}

	if strings.ToLower(m.Module) != "query" && strings.ToLower(m.Module) != "mutation" {
		apiError := ApiError([]string{"Não foi encontrado definção de module"}, "apiParams", "isValid", "INVALID_MODULE", time.Now(), nil)
		return &apiError
	}
	if strings.ToLower(m.Path) == "" {
		apiError := ApiError([]string{"Não foi encontrado definção de path"}, "apiParams", "isValid", "INVALID_PATH", time.Now(), nil)
		return &apiError
	}

	return nil
}

func extractModule(queryString string) string {
	fields := strings.Fields(queryString)
	if len(fields) > 0 {
		return capitalize(fields[0])
	}
	return ""
}

func extractPath(queryString string) string {
	re := regexp.MustCompile(`([a-zA-Z_]+\()`)
	match := re.FindStringSubmatch(queryString)
	if len(match) > 1 {
		operation := match[1]
		return operation[:len(operation)-1]
	}
	return ""
}

func capitalize(input string) string {
	if input == "" {
		return input
	}
	return strings.ToUpper(input[:1]) + input[1:]
}

func (p ParamsModel) ToString() string {
	return fmt.Sprintf("Instance of Params(module:%s, path:%s, errors:%v)", p.Module, p.Path, p.Errors)
}
