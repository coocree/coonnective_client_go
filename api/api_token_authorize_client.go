package api

import (
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

type TokenAuthorizeClientModel struct {
	apiUri         string
	token          string
	clientId       string
	clientSecret   string
	code           string
	codeChallenge  string
	codeVerifier   string
	codeVerifier64 string
	nonceToken     string
	state          string
	context        string
	subject        string
	jwToken        jwt.Token
	apiSecurity    SecurityModel
}

func TokenAuthorizeClient(token *string) TokenAuthorizeClientModel {
	if token == nil {
		apiError := ApiError([]string{"Nenhum token de 'autorização' foi definido na inicialização da aplicação."},
			"TokenAuthorizeUser", "TokenAuthorizeUser", "AUTHORIZE_TOKEN_NOT_FOUND", time.Now(), nil)

		E(apiError.ToString(), apiError.Code, nil)
		panic(apiError.Code)
	}

	_token := *token

	jwToken, errJwt := jwt.Parse([]byte(_token))
	if errJwt != nil {
		apiError := ApiError([]string{"Token de autorização de usuário 'user' inválido.", errJwt.Error()},
			"TokenAuthorizeUser", "TokenAuthorizeUser", "AUTHORIZE_TOKEN_FAIL", time.Now(), _token)

		E(apiError.ToString(), apiError.Code, nil)
		panic(apiError.Code)
	}

	apiSecurity := Security()
	codeVerifier := apiSecurity.RandomBytes(32)
	codeVerifier64 := apiSecurity.Base64URLEncode(codeVerifier)
	codeChallenge := apiSecurity.Base64URLEncode(apiSecurity.EncodeSha256(codeVerifier64))
	state := apiSecurity.RandomBytes(16)

	return TokenAuthorizeClientModel{
		apiUri:         "",
		token:          _token,
		clientId:       "",
		clientSecret:   "",
		code:           "",
		codeChallenge:  codeChallenge,
		codeVerifier:   codeVerifier,
		codeVerifier64: codeVerifier64,
		nonceToken:     "",
		state:          state,
		context:        "",
		subject:        "",
		jwToken:        jwToken,
		apiSecurity:    apiSecurity,
	}
}

func (t *TokenAuthorizeClientModel) IsValid() bool {
	isValid := false
	jwToken := t.jwToken
	if jwToken.Issuer() != "" && jwToken.Subject() == "auth_token" {
		t.apiUri = jwToken.Issuer()
		t.apiUri = "http://localhost:4600/coonective"
		t.clientId = jwToken.PrivateClaims()["cid"].(string)
		t.context = jwToken.PrivateClaims()["ctx"].(string)
		t.subject = jwToken.Subject()
		t.clientSecret = t.apiSecurity.EncodeSha256(t.clientId + t.codeVerifier64)

		//fmt.Println("apiUri: ", t.apiUri, "clientId: ", t.clientId, "clientSecret: ", t.clientSecret, "codeVerifier: ", t.codeVerifier, "codeVerifier64: ", t.codeVerifier64, "codeChallenge: ", t.codeChallenge, "state: ", t.state, "context: ", t.context, "subject: ", t.subject)

		apiStorage := User(nil).Storage()
		if apiStorage != nil && t.clientId == "" {
			apiStorage.Add("clientId", t.clientId)
		}
		isValid = true
	}
	return isValid
}

func (t *TokenAuthorizeClientModel) Ready() bool {
	apiStorage := User(nil).Storage()
	clientId := apiStorage.Read("clientId")
	if clientId == "" {
		return true
	}
	return false
}

type AclAuthorizeResponse struct {
	Data struct {
		AclAuthorize struct {
			Result struct {
				Code       string `bson:"code,omitempty" json:"code"`
				NonceToken string `bson:"nonceToken,omitempty" json:"nonceToken"`
				State      string `bson:"state,omitempty" json:"state"`
			} `bson:"result,omitempty" json:"result"`
			Error       ErrorModel
			ElapsedTime string
			Success     bool
		}
	}
}

func (t *TokenAuthorizeClientModel) AclAuthorize() *ErrorModel {
	apiConnection, _ := Connection(t.token, t.apiUri)
	const params = `
mutation aclAuthorize($input: AuthorizeInput!) {
  aclAuthorize(input: $input) {
    result {
      code
      nonceToken
      state
    }
    success
    error {
      code
      createdAt
      messages
      module
      path
      variables
    }
    elapsedTime
  }
}
`
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"clientId":      t.clientId,
			"codeChallenge": t.codeChallenge,
			"responseType":  "code",
			"scope":         "mut:authorize",
			"state":         t.state,
		},
	}

	apiResponse := apiConnection.Mutation(params, variables)

	var response AclAuthorizeResponse
	apiResponse.Endpoint(&response)

	authorizeResult := response.Data.AclAuthorize.Result
	if authorizeResult.State == t.state {
		t.code = authorizeResult.Code
		t.nonceToken = authorizeResult.NonceToken
		return nil
	}
	return nil
}
