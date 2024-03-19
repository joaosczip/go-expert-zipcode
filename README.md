# Zipcode Temperature

## Iniciando a aplicação

1. Clone o repositório para a sua máquina local usando `git clone`.
2. Navegue até o diretório do projeto.
3. Execute `docker-compose up` para iniciar a aplicação.

## Como usar a aplicação

A aplicação fornece uma API para buscar a temperatura atual de uma localidade baseada no CEP fornecido.

### Endpoint

`GET /temperature?zipcode={zipcode}`

Substitua `{zipcode}` pelo CEP desejado.

A resposta será um objeto JSON com a temperatura atual na localidade do CEP fornecido.

### Testes

Para executar os testes, navegue até o diretório do projeto e execute `go test ./....`

## Configuração

As configurações da aplicação estão localizadas no arquivo `.env`. Você pode usar o arquivo `.env.example` como referência para criar o seu próprio arquivo `.env`, substituindo as variáveis pelos seus valores reais.

## Endereço da aplicação

A aplicação está executando no Google CloudRun e pode ser requisitada através da url https://zipcode-temp-conttfglyq-uc.a.run.app/