package api

import (
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

type TokenAccessUserModel struct {
	apiUserAuthorize TokenAuthorizeUserModel
	clientId         string
	tokenType        string
	expiresIn        string
	accessToken      string
	refreshToken     string
	isValid          bool
	user             *UserModel
}

func TokenAccessUser(apiUserAuthorize TokenAuthorizeUserModel) *TokenAccessUserModel {
	return &TokenAccessUserModel{
		apiUserAuthorize: apiUserAuthorize,
		user:             User(nil),
	}
}

type AclTokenResponse struct {
	Data struct {
		AclToken struct {
			Result struct {
				AccessToken  string `bson:"accessToken,omitempty" json:"accessToken"`
				ExpiresIn    string `bson:"expiresIn,omitempty" json:"expiresIn"`
				RefreshToken string `bson:"refreshToken,omitempty" json:"refreshToken"`
				TokenType    string `bson:"tokenType,omitempty" json:"tokenType"`
			} `bson:"result,omitempty" json:"result"`
			Error       ErrorModel
			ElapsedTime string
			Success     bool
		}
	}
}

func (a *TokenAccessUserModel) Token() bool {
	apiConnection, _ := Connection(a.apiUserAuthorize.nonceToken, a.apiUserAuthorize.apiUri)

	const params = `
        mutation aclToken($input: TokenInput!) {
            aclToken(input: $input) {
                result {
                    accessToken
                    expiresIn
                    refreshToken
                    tokenType
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
			"clientId":     a.user.ClientId(),
			"clientSecret": a.apiUserAuthorize.clientSecret,
			"code":         a.apiUserAuthorize.code,
			"codeVerifier": a.apiUserAuthorize.codeVerifier64,
			"grantType":    "authorization_code",
		},
	}

	apiResponse := apiConnection.Mutation(params, variables)

	var response AclTokenResponse
	apiResponse.Endpoint(&response)

	if response.Data.AclToken.Success {
		tokenResult := response.Data.AclToken.Result
		a.tokenType = tokenResult.TokenType
		a.expiresIn = tokenResult.ExpiresIn
		a.accessToken = tokenResult.AccessToken
		a.refreshToken = tokenResult.RefreshToken
		a.isValid = true
		a.user.AccessToken(a.accessToken)
		a.user.RefreshToken(a.refreshToken)

		jwToken, errJwt := jwt.Parse([]byte(a.accessToken))
		if errJwt != nil {
			apiError := ApiError([]string{"Token de acesso de usuário 'user' inválido.", errJwt.Error()},
				"TokenAccessUser", "Token", "ACCESS_TOKEN_FAIL", time.Now(), variables)
			E(apiError.ToString(), apiError.Code, nil)
			panic(apiError.Code)
		}

		iss := jwToken.Issuer()
		if iss != "" {
			a.user.ServerUri(iss)
		} else {
			return false
		}
	}

	return a.isValid
}
