package api

import (
	"fmt"
	"time"
)

// ErrorModel A struct ApiError é usada para representar e lidar com erros da API.
type ErrorModel struct {
	// Data e hora em que o erro ocorreu, no formato ISO8601.
	CreatedAt string `json:"createdAt"`

	// Lista de mensagens de erro.
	Messages []string `json:"messages"`

	// Nome do módulo onde ocorreu o erro.
	// Mutation: módulo de mutação.
	// Query: módulo de consulta.
	Module string `json:"module"`

	// Caminho do campo ou recurso onde ocorreu o erro.
	// AclAuthorize: endpoint de autorização de acesso.
	Path string `json:"path"`

	// Tipo do erro, usado para classificar e identificar a natureza do erro.
	// GRAPHQL_VALIDATION_FAILED: Erro de validação GraphQL.
	Code string `json:"code"`

	// Variável adicional fornecida, geralmente usada para fornecer informações adicionais sobre o erro.
	Variables interface{} `json:"variables"`
}

// ApiError Construtor de inicialização de campo.
func ApiError(messages []string, module string, path string, code string, createdAt time.Time, variables interface{}) ErrorModel {
	dateFormat := "2006-01-02 15:04:05"
	_createdAt := createdAt.Format(dateFormat)
	return ErrorModel{CreatedAt: _createdAt, Messages: messages, Module: module, Path: path, Code: code, Variables: variables}
}

func (e ErrorModel) HasError() bool {
	return len(e.Messages) > 0
}

// ToString Retorna uma representação de string da instância ApiError.
func (e ErrorModel) ToString() string {
	return fmt.Sprintf("module:%s, code:%s, path:%s, messages:%v, variables:%v, createdAt:%s", e.Module, e.Code, e.Path, e.Messages, e.Variables, e.CreatedAt)
}
