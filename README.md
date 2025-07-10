# go-microservices-architecture

Este repositório é um **monorepo** contendo uma arquitetura de microsserviços desenvolvida em **Go**, utilizando uma abordagem **orientada a eventos (EDA)** com **RabbitMQ** como event bus.

O projeto simula um sistema de cadastro de usuários desacoplado, onde diferentes serviços se comunicam por meio de eventos, permitindo escalabilidade, flexibilidade e independência entre os componentes.

## Serviços incluídos

- **auth-service**: Responsável por cadastrar usuários e publicar eventos como `UsuarioCriado`.
- **notification-service**: Escuta eventos e envia notificações (ex: e-mails de boas-vindas).
- **user-worker**: Processa eventos de forma assíncrona para tarefas como auditoria ou onboarding.

## Tecnologias

- Go (Golang)
- RabbitMQ
- Docker + Docker Compose

## Objetivo

Demonstrar uma arquitetura simples, porém funcional, de microsserviços orientados a eventos com Go. Ideal para estudos, protótipos e como base para projetos maiores.
