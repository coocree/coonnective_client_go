# coonnective_client_go

O coonnective_client_go é um pacote de conexão GraphQL desenvolvido em GO, que permite a comunicação com um servidor GraphQL para enviar e receber dados. Ele é utilizado para facilitar o processo de comunicação entre o cliente e o servidor, além de garantir que as respostas recebidas sejam tratadas de forma adequada.

```dart
func EventReset() api.ResponseModel {
	variables := map[string]interface{}{}
	graphQL := `
mutation EventReset {
	EventReset(filter: {idEvent: "123"}) {
		result {
			idEvent
			status
		}
		error {
			code
			createdAt
			messages
			module
			path
			variables
		}
		elapsedTime
		success
	}
}
`
	return api.Dao(graphQL, variables)
}
```
No código apresentado, temos a função EventReset() que realiza uma mutação no servidor GraphQL. O objetivo dessa mutação é resetar um evento identificado pelo ID "123". Para isso, a função define uma variável vazia e uma string contendo a query GraphQL a ser executada. A query é definida dentro da variável "graphQL", que contém o código GraphQL necessário para a execução da mutação.

```dart
func main() {

	token := ""
	api.Connect(&token)

	apiResponse := EventReset()
	if !apiResponse.IsValid() {
		apiResponse.ThrowException()
	}
	apiEndpoint := apiResponse.Endpoint("EventReset")
	if !apiEndpoint.IsValid() {
		apiEndpoint.ThrowException()
	}

	fmt.Println(apiEndpoint.Result)
}
```
Após a definição da query, a função chama a função "api.Dao()" que é responsável por executar a query. Essa função recebe a query definida anteriormente e as variáveis definidas no início da função. Ela retorna um objeto do tipo "api.ResponseModel" que contém a resposta recebida do servidor GraphQL.

Em seguida, o código faz uma verificação para garantir que a resposta recebida seja válida. Caso contrário, a função "apiResponse.ThrowException()" é chamada, o que lança uma exceção informando que algo deu errado na execução da query.

Após a verificação da resposta, a função "apiResponse.Endpoint()" é chamada para obter o resultado da mutação. A variável "apiEndpoint" recebe o resultado retornado pela função. Em seguida, outra verificação é realizada para garantir que o resultado obtido seja válido. Caso contrário, a função "apiEndpoint.ThrowException()" é chamada, o que também lança uma exceção informando que algo deu errado na execução da query.

Por fim, o resultado é impresso na tela através do comando "fmt.Println(apiEndpoint.Result)". Esse resultado contém o status do evento após o reset, que foi obtido como resposta da mutação realizada no servidor GraphQL.