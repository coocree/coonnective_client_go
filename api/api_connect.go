package api

import (
	"log"
	"time"
)

type ConnectModel struct {
	authToken            string
	TokenAccessClient    TokenAccessClientModel
	TokenAuthorizeClient TokenAuthorizeClientModel
	apiUser              *UserModel
}

var apiInstance *ConnectModel

// Connect Construtor da fábrica ApiConnect.
func Connect(authToken *string) (*ConnectModel, *ErrorModel) {
	// Garante que só exista uma instância da API
	if apiInstance != nil {
		return apiInstance, nil
	}

	if authToken == nil {
		apiError := ApiError([]string{"AuthToken não pode estar vazio."},
			"apiConnect", "connect", "QUERY_ERROR", time.Now(), nil)
		return nil, &apiError
	}

	apiInstance = &ConnectModel{
		authToken:            *authToken,
		TokenAccessClient:    TokenAccessClientModel{},
		TokenAuthorizeClient: TokenAuthorizeClientModel{},
		apiUser:              User(nil),
	}

	// Inicializa a conexão com a API.
	apiError := apiInstance.init(authToken)
	if apiError != nil {
		return nil, apiError
	}

	return apiInstance, nil
}

// Inicializa a conexão com a API.
func (a *ConnectModel) init(authToken *string) *ErrorModel {
	a.apiUser = User(nil)

	// Inicializa TokenAuthorizeClient.
	a.TokenAuthorizeClient = TokenAuthorizeClient(authToken)

	// Verifica se TokenAuthorizeClient é válido.
	if !a.TokenAuthorizeClient.IsValid() {
		apiError := ApiError([]string{"Token de autorização da aplicação 'client' inválido"},
			"apiConnect", "init", "TOKEN_INVALID", time.Now(), nil)
		return &apiError
	}

	// Verifica se TokenAuthorizeClient tem autorização.
	isAclAuthorize := a.TokenAuthorizeClient.AclAuthorize()
	if isAclAuthorize != nil {
		apiError := ApiError([]string{"Token de autorização da aplicação 'client' não tem autorização"},
			"apiConnect", "init", "TOKEN_INVALID", time.Now(), nil)
		return &apiError
	}

	// Obtém AccessToken de TokenAccessClient.
	a.TokenAccessClient = TokenAccessClient(a.TokenAuthorizeClient)
	accessToken := a.TokenAccessClient.Token()

	// Verifica se AccessToken de TokenAccessClient é válido.
	if accessToken == "" {
		apiError := ApiError([]string{"Token de autorização da aplicação 'client' inválido, não foi possível obter AccessToken"},
			"apiConnect", "init", "TOKEN_INVALID", time.Now(), nil)
		return &apiError
	}

	// Inicializa TokenAuthorizeUser.
	apiTokenAuthorizeUser := TokenAuthorizeUser(accessToken)
	// Verifica se TokenAuthorizeUser é válido.
	if !apiTokenAuthorizeUser.IsValid() {
		apiError := ApiError([]string{"Token de autorização do usuário 'user' inválido"},
			"apiConnect", "init", "TOKEN_INVALID", time.Now(), nil)
		return &apiError
	}

	// Verifica se TokenAuthorizeUser tem autorização.
	if !apiTokenAuthorizeUser.Authorize() {
		apiError := ApiError([]string{"Token de autorização do usuário 'user' não tem autorização"},
			"apiConnect", "init", "TOKEN_INVALID", time.Now(), nil)
		return &apiError
	}

	// Obtém AccessToken de TokenAccessUser.
	apiTokenAccessUser := TokenAccessUser(apiTokenAuthorizeUser)

	// Verifica se AccessToken de TokenAccessUser é válido.
	if !apiTokenAccessUser.Token() {
		apiError := ApiError([]string{"Token de autorização da aplicação 'client' inválido"},
			"apiConnect", "init", "TOKEN_INVALID", time.Now(), nil)
		return &apiError
	}

	log.Println("API conectada com sucesso.")
	return nil
}

// Query Realiza uma consulta GraphQL.
func (a *ConnectModel) Query(params string, variables map[string]interface{}) ResponseModel {
	accessToken := a.apiUser.GetAccessToken()
	serverUri := a.apiUser.GetServerUri()

	if accessToken != "" && serverUri != "" {
		apiConnection, _ := Connection(accessToken, serverUri)
		return apiConnection.Query(params, variables)
	}

	apiError := ApiError([]string{"Falha na conexão com servidor API, accessToken ou serverUri inválidos"},
		"apiConnect", "query", "QUERY_ERROR", time.Now(), variables)

	return ResponseModel{
		Success: false,
		Errors:  []ErrorModel{apiError},
	}
}

// Mutation Realiza uma mutação GraphQL.
func (a *ConnectModel) Mutation(params string, variables map[string]interface{}) ResponseModel {
	accessToken := a.apiUser.GetAccessToken()
	serverUri := a.apiUser.GetServerUri()

	if accessToken != "" && serverUri != "" {
		apiConnection, _ := Connection(accessToken, serverUri)
		return apiConnection.Mutation(params, variables)
	}

	apiError := ApiError([]string{"Falha na conexão com servidor API, accessToken ou serverUri inválidos"},
		"apiConnect", "mutation", "MUTATION_ERROR", time.Now(), variables)

	return ResponseModel{
		Success: false,
		Errors:  []ErrorModel{apiError},
	}
}
