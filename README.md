# Desafio Go — Cotação do Dólar

## Sobre o projeto

Projeto desenvolvido como entrega de um desafio da pós-graduação, com o objetivo de aplicar conceitos de **HTTP, Contextos, Banco de Dados e Manipulação de Arquivos na linguagem Go**.

A aplicação é composta por dois programas independentes que se comunicam por HTTP para consultar, persistir e disponibilizar a cotação atual do dólar.

### Componentes

- **`./server/main.go`:** disponibiliza um servidor HTTP na porta `8080`, consulta a API externa de câmbio e persiste os dados da cotação em um banco SQLite.
- **`./client/main.go`:** consome o endpoint local do servidor, extrai o valor atual do dólar e salva o resultado no arquivo `cotacao.txt`.

## Tecnologias utilizadas

- Go
- HTTP e JSON
- Context API (`context`)
- SQLite
- Manipulação de arquivos

**API externa:** [AwesomeAPI — Cotação USD/BRL](https://economia.awesomeapi.com.br/json/last/USD-BRL)

## Requisitos

Antes de executar o projeto, é necessário ter o Go instalado e disponível no terminal.

Para verificar a instalação:

```bash
go version
```

## Como executar

### 1. Clonar o repositório

```bash
git clone https://github.com/andreeduardo/posFullCycleDesafioClientServerApi.git
cd posFullCycleDesafioClientServerApi
```

### 2. Executar o servidor

Em um terminal, execute:

```bash
go run server/main.go
```

O servidor HTTP será iniciado na porta `8080`.

Endpoint disponível:

```text
GET http://localhost:8080/cotacao
```

Ao receber uma requisição, o servidor consulta a API externa de câmbio, registra a cotação no banco SQLite e retorna os dados em formato JSON.

Mantenha esse terminal aberto enquanto utiliza o cliente.

### 3. Executar o cliente

Em outro terminal, na raiz do projeto, execute:

```bash
go run client/main.go
```

O cliente fará uma requisição HTTP ao servidor local, obterá o valor do campo `bid` da cotação e gravará o resultado no arquivo `cotacao.txt`.

O conteúdo do arquivo deverá seguir o formato:

```text
Dólar: 5.1234
```

O valor acima é apenas ilustrativo. A cotação registrada dependerá do retorno da API no momento da execução.

### 4. Consultar o endpoint manualmente

Com o servidor em execução, também é possível consultar a cotação diretamente:

```bash
curl http://localhost:8080/cotacao
```

A resposta deverá ser um JSON contendo os dados retornados pelo servidor para a cotação USD/BRL.

## Timeouts e tratamento de erros

A aplicação utiliza contextos com limites de tempo para evitar operações demoradas.

| Operação | Timeout máximo |
|---|---:|
| Consulta à API externa pelo servidor | 200 ms |
| Persistência da cotação no SQLite | 10 ms |
| Requisição do cliente ao servidor | 300 ms |

Quando uma operação excede seu prazo, o erro correspondente deve ser registrado no console do serviço responsável.

## Persistência dos dados

O servidor utiliza SQLite para registrar as cotações recebidas da API externa.

A persistência é realizada com um contexto limitado a `10 ms`, conforme solicitado no enunciado.

O banco é utilizado pelo servidor e não precisa estar instalado como um serviço independente, desde que o projeto utilize um driver SQLite compatível com a configuração adotada.

## Arquivos gerados

- **Banco SQLite:** armazena o histórico das cotações consultadas pelo servidor.
- **`cotacao.txt`:** armazena o valor atual do dólar obtido pelo cliente.

Os arquivos gerados durante a execução não precisam ser versionados no GitHub.

## Encerramento

Para interromper o servidor ou o cliente em execução no terminal, pressione `Ctrl+C`.

## Autor

André Eduardo da Silva