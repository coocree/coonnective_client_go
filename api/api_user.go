package api

import (
	"time"
)

// UserModel Classe User representa um usuário da API e gerencia a autenticação e autorização desse usuário.
type UserModel struct {
	// Atributos da classe
	id          string
	apiKey      string
	email       string
	authUri     string
	tokenUri    string
	info        string
	auth        string
	services    string
	name        string
	role        string
	firstLetter string
	clientId    string
	storage     *StorageModel
}

var instance *UserModel

func User(config interface{}) *UserModel {
	if instance == nil {
		instance = &UserModel{}
	}
	instance.init(config)
	return instance
}

func (au *UserModel) Storage() *StorageModel {
	storage := Storage()
	au.storage = storage
	return storage
}

func (au *UserModel) init(value interface{}) {
	security := Security()
	Storage()
	au.clientId = security.EncodeSha1("GUEST")

	v, ok := value.(map[string]interface{})
	if !ok {
		return
	}

	if name, ok := v["name"].(string); ok && name != "" {
		au.name = name
		au.firstLetter = name[:1]
	}
	if id, ok := v["_id"].(string); ok && id != "" {
		au.id = id
	}
	if role, ok := v["role"].(string); ok && role != "" {
		au.role = role
	}
	if clientId, ok := v["clientId"].(string); ok && clientId != "" {
		au.clientId = clientId
	}
}

func (au *UserModel) ClientId() string {
	return au.clientId
}
func (au *UserModel) SetClientId(id string) {
	au.clientId = id
}

// Authorize Método para validar o token de autorização e obter o token de acesso e autenticação.
func (au *UserModel) Authorize() {
	accessToken := au.GetAccessToken()
	if accessToken == "" {
		apiError := ApiError([]string{"Nenhum token de 'autorização' foi definido na inicialização da aplicação."},
			"User", "Authorize", "AUTHORIZE_TOKEN_NOT_FOUND", time.Now(), accessToken)

		E(apiError.ToString(), apiError.Code, nil)
		panic(apiError.Code)
	}

	apiTokenAuthorizeUser := TokenAuthorizeUser(accessToken)

	if !apiTokenAuthorizeUser.IsValid() {
		panic("Token de autorização da usuário 'user' inválido.")
	}
	au.handleAuthorization(&apiTokenAuthorizeUser)
}

// Função para lidar com a autorização e obter o token de acesso.
func (au *UserModel) handleAuthorization(apiTokenAuthorizeUser *TokenAuthorizeUserModel) {
	isAuthorizeValid := apiTokenAuthorizeUser.Authorize()
	if !isAuthorizeValid {
		panic("Token de autorização da usuário 'user' inválido.")
	}
	au.handleAccessToken(apiTokenAuthorizeUser)
}

// Função para lidar com o token de acesso e validá-lo.
func (au *UserModel) handleAccessToken(apiTokenAuthorizeUser *TokenAuthorizeUserModel) {
	apiTokenAccessUser := TokenAccessUser(*apiTokenAuthorizeUser)
	isTokenValid := apiTokenAccessUser.Token()
	if !isTokenValid {
		panic("Token de autorização da usuário 'user' inválido.")
	}
}
func (au *UserModel) resume() map[string]interface{} {
	return map[string]interface{}{
		"_id":  au.id,
		"name": au.name,
	}
}

// AccessToken Método para salvar o token de acesso do usuário.
func (au *UserModel) AccessToken(value string) {
	au.storage.Add("accessToken", value)
}

// GetAccessToken Método para obter o token de acesso do usuário.
func (au *UserModel) GetAccessToken() string {
	return au.storage.Read("accessToken")
}

// AuthToken Método para salvar o token de autenticação do usuário.
func (au *UserModel) AuthToken(value string) {
	au.storage.Add("authToken", value)
}

// GetAuthToken Método para obter o token de autenticação do usuário.
func (au *UserModel) GetAuthToken() string {
	return au.storage.Read("authToken")
}

// ServerUri Método para salvar o URI do servidor.
func (au *UserModel) ServerUri(value string) {
	au.storage.Add("serverUri", value)
}

// GetServerUri Método para obter o URI do servidor.
func (au *UserModel) GetServerUri() string {
	return au.storage.Read("serverUri")
}

func (au *UserModel) RefreshToken(value string) {
	au.storage.Add("refreshToken", value)
}
func (au *UserModel) GetRefreshToken() string {
	return au.storage.Read("refreshToken")
}
