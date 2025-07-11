// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/usuarios": {
            "post": {
                "description": "Cadastra um novo usuário no sistema",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Usuários"
                ],
                "summary": "Cadastrar um novo usuário",
                "parameters": [
                    {
                        "description": "Dados do usuário",
                        "name": "usuario",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Usuario"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Usuário cadastrado com sucesso",
                        "schema": {
                            "$ref": "#/definitions/schema.UsuarioCreatedEnvelope"
                        }
                    },
                    "400": {
                        "description": "Requisição inválida",
                        "schema": {
                            "$ref": "#/definitions/error.AppError"
                        }
                    },
                    "409": {
                        "description": "Conflito: usuário já existe",
                        "schema": {
                            "$ref": "#/definitions/error.AppError"
                        }
                    },
                    "500": {
                        "description": "Erro interno do servidor",
                        "schema": {
                            "$ref": "#/definitions/error.AppError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Usuario": {
            "type": "object",
            "required": [
                "email",
                "nome",
                "senha"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "nome": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 3
                },
                "senha": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 6
                }
            }
        },
        "error.AppError": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "ex.: USER_NOT_FOUND",
                    "type": "string"
                },
                "message": {
                    "description": "msg “safe” para cliente",
                    "type": "string"
                }
            }
        },
        "schema.UsuarioCreatedEnvelope": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object",
                    "properties": {
                        "id": {
                            "type": "integer"
                        }
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
