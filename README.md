# go-microservices-architecture

Este repositório é um **monorepo** contendo uma arquitetura de microsserviços desenvolvida em **Go**, utilizando uma abordagem **orientada a eventos (EDA)** com **RabbitMQ** como event bus.

O projeto simula um sistema de cadastro de usuários desacoplado, onde diferentes serviços se comunicam por meio de eventos, permitindo escalabilidade, flexibilidade e independência entre os componentes.

## Serviços incluídos

- **auth-service**: Responsável por cadastrar usuários e publicar eventos como `UsuarioCriado`.
- **notification-service**: Escuta eventos e envia notificações (ex: e-mails de boas-vindas).
- **user-worker**: Processa eventos de forma assíncrona para tarefas como auditoria.

## Tecnologias

- Go (Golang)
- RabbitMQ
- Docker + Docker Compose
- Gin
- Gorm (ORM)
- GoMock
- JWT
- MySQL (banco de dados)

## Objetivo

Demonstrar uma arquitetura simples, porém funcional, de microsserviços orientados a eventos com Go. Ideal para estudos, protótipos e como base para projetos maiores.

## Estrutura de pastas

### `auth-service`

```bash
auth-service/
├── __test__           # Testes e2e (End-to-End)
├── cmd/               # Ponto de entrada da aplicação (main.go)
├── docs/              # Documentação específica
├── internal/          # Lógica interna da aplicação
│   ├── config/        # Carregamento e definição de configurações
│   ├── domain/        # Definições de entidades
│   ├── handler/       # Handlers de HTTP
│   ├── mq/            # Implementação do EventBus com RabbitMQ(publicação de eventos)
│   ├── repository/    # Interações com banco de dados
│   ├── router/        # Definição de rotas
│   └── schema/        # Esquemas
├── pkg/               # Pacotes reutilizáveis (exportáveis)
│   ├── error/         # Tratamento de erros personalizados
│   ├── middleware/    # Middlewares de autenticação, error, etc.
│   ├── response/      # Modelos de resposta padronizados
│   └── security/      # Funções de segurança (ex: JWT, hash)
│   └── validate.go    # Lógica auxiliar de validações
├── .env.dev           
├── .env.test          
├── Dorckerfile          
├── go.mod             
└── go.sum             
```

### `notification-service`

```bash
notification-service/
├── cmd/               # Ponto de entrada da aplicação (main.go)
├── internal/          # Lógica interna da aplicação
│   ├── config/        # Carregamento e definição de configurações
│   ├── consumer/      # Configuração e conexão com o RabbitMQ (consumo de eventos)
│   ├── domain/        # Definições de entidades
│   ├── handler/       # Funções que processam eventos recebidos (ex: envio de e-mail)
│   └── messages/      # Responsável pelo envio de e-mail
├── public/            # Lógica interna da aplicação
│   └── validate.go    # Template do e-mail de boas-vindas
├── .env             
├── Dockerfile             
├── go.mod             
└── go.sum   
```

### `user-worker`

```bash
user-worker/
├── __test__           # Testes
│   ├── mocks/         # Mocks gerados para interfaces (usando GoMock)
│   ├── service/       # Testes unitários para serviços (ex: AuditService)
├── cmd/               # Ponto de entrada da aplicação (main.go)
├── internal/          # Lógica interna da aplicação
│   ├── config/        # Carregamento e definição de configurações
│   ├── consumer/      # Consumo de eventos RabbitMQ para auditoria de usuários
│   ├── domain/        # Definições de entidades
│   ├── handler/       # Handlers de HTTP
│   ├── repository/    # Interações com banco de dados
│   ├── router/        # Definição de rotas
│   └── schema/        # Esquemas
│   └── service/       # Regras de negócio
├── .env           
├── Dorckerfile          
├── go.mod             
└── go.sum    
```

## Conceitos abordados

- Event-Driven Architecture (EDA) com RabbitMQ

- Publicação e consumo de eventos com tópicos- 

- Desacoplamento entre serviços com comunicação assíncrona

- Auditoria e rastreabilidade de eventos

- Workers desacoplados para escalabilidade horizontal

- Testes unitários com GoMock

## Fluxo de Cadastro de Usuário (Visão da Arquitetura)

1 - `auth-service` **(Producer)**
- Persiste o novo usuário no banco de dados.

  - Emite dois eventos via RabbitMQ (publicação no exchange user.events):

    - `user.created`: Evento principal informando que o usuário foi criado.
    - `user.audit`: Evento de auditoria com os dados do novo usuário (antes e depois, se aplicável).

2 - `notification-service` **(Consumer)**
- Escuta eventos do tipo user.created com binding para a fila user.welcome.email.

  - Executa o envio de um e-mail de boas-vindas para o endereço informado.

3 - `user-worker` **(Consumer)**
- Escuta eventos do tipo user.audit com binding para a fila user.audit.

  - Salva as informações de auditoria no banco de dados (usuário afetado, tipo de evento, timestamp, mudanças nos campos).

## Benefícios da Arquitetura Event-Driven
- Desacoplamento: auth-service não depende diretamente de notification-service ou user-worker.

- Escalabilidade: serviços podem ser escalados de forma independente.

- Flexibilidade: novos serviços podem ser conectados a eventos existentes sem alterar os produtores.

- Auditabilidade: mudanças em entidades importantes são registradas de forma transparente.

## Como rodar o projeto com Docker

Este projeto utiliza `Docker` e `Docker Compose` para orquestrar todos os microsserviços e dependências. Siga os passos abaixo para rodar localmente.

### 1. Criar o arquivo `.env`

No mesmo diretório do `docker-compose.yml`, crie um arquivo chamado `.env` com o seguinte conteúdo:

```env
MYSQL_ROOT_PASSWORD=
MYSQL_DATABASE=
MYSQL_PORT=

DB_USER=
DB_PASSWORD=
DB_NAME=
DB_NAME_TEST=
DB_HOST=mysql
DB_PORT=3306

JWT_SECRET=

AMQP_USER=
AMQP_PASSWORD=
AMQP_URL=rabbitmq

ROOT_EMAIL=
ROOT_EMAIL_PASSWORD=
SMTP_HOST=
SMTP_PORT=
```
### 2. Subir os containers
Execute o seguinte comando:
```bash
docker-compose up --build
```

