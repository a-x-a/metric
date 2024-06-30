// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package api

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
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
        "/list": {
            "get": {
                "description": "Возвращает список полученных от сервиса метрик ввиде обычного текста.",
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "list"
                ],
                "summary": "List",
                "operationId": "list",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Performs a health check by pinging the service.",
                "tags": [
                    "ping"
                ],
                "summary": "Ping",
                "operationId": "ping",
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": ""
                    }
                }
            }
        },
        "/update": {
            "post": {
                "description": "Обновляет текущее значение метрики с указанным имененм и типом.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "update"
                ],
                "summary": "UpdateJSON",
                "operationId": "updateJSON",
                "parameters": [
                    {
                        "description": "Параметры метрики: имя, тип, значение",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/adapter.RequestMetric"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "404": {
                        "description": ""
                    },
                    "500": {
                        "description": ""
                    }
                }
            }
        },
        "/update/{kind}/{name}/{value}": {
            "post": {
                "description": "Обновляет текущее значение метрики с указанным имененм и типом.",
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "update"
                ],
                "summary": "Update",
                "operationId": "update",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Тип метрики",
                        "name": "kind",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Имяметрики",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Значение метрики",
                        "name": "value",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": ""
                    }
                }
            }
        },
        "/updates": {
            "post": {
                "description": "Обновляет текущие значения метрик из набора.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "updatebatch"
                ],
                "summary": "UpdateBatch",
                "operationId": "updatebatch",
                "parameters": [
                    {
                        "description": "Набор метрик",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/adapter.RequestMetric"
                            }
                        }
                    }
                ],
                "responses": {
                    "400": {
                        "description": ""
                    },
                    "500": {
                        "description": ""
                    }
                }
            }
        },
        "/value": {
            "get": {
                "description": "Возвращает текущее значение метрики в формате JSON с указанным имененм и типом.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "value"
                ],
                "summary": "GetJSON",
                "operationId": "getJSON",
                "parameters": [
                    {
                        "description": "Параметры метрики: имя, тип",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/adapter.RequestMetric"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": ""
                    },
                    "404": {
                        "description": ""
                    },
                    "500": {
                        "description": ""
                    }
                }
            }
        },
        "/value/{kind}/{name}": {
            "get": {
                "description": "Возвращает текущее значение метрики с указанным имененм и типом.",
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "value"
                ],
                "summary": "Get",
                "operationId": "get",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Тип метрики",
                        "name": "kind",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Имя метрики",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "404": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "adapter.RequestMetric": {
            "type": "object",
            "properties": {
                "delta": {
                    "description": "значение метрики в случае передачи counter",
                    "type": "integer"
                },
                "id": {
                    "description": "имя метрики",
                    "type": "string"
                },
                "type": {
                    "description": "параметр, принимающий значение gauge или counter",
                    "type": "string"
                },
                "value": {
                    "description": "значение метрики в случае передачи gauge",
                    "type": "number"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
