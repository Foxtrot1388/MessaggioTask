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
        "/message/create": {
            "post": {
                "description": "Create a message",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "message"
                ],
                "summary": "Create",
                "operationId": "create",
                "parameters": [
                    {
                        "description": "message fo create models",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.MessageToInsert"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "if all message have been create",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.OutputMessage"
                            }
                        }
                    },
                    "400": {
                        "description": "validation error",
                        "schema": {
                            "$ref": "#/definitions/server.response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.response"
                        }
                    }
                }
            }
        },
        "/message/statistics": {
            "get": {
                "description": "Get statistics a message",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "message"
                ],
                "summary": "statistics",
                "operationId": "statistics",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Some date at (format 2006-01-02)",
                        "name": "dateAt",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Some date to (format 2006-01-02)",
                        "name": "dateTo",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.StatMessage"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.MessageToInsert": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "entity.OutputMessage": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "entity.StatMessage": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "day": {
                    "type": "string"
                }
            }
        },
        "server.response": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "89.169.167.99:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Message API",
	Description:      "API Server for messages",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
