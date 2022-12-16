// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Web-Service Support",
            "email": "aggregator_lets_go@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/last_news": {
            "get": {
                "description": "Get array of Articles for the last 7 days",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "news"
                ],
                "summary": "Retrieves news from last 7 days",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.DBArticle"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.DBArticle": {
            "description": "Article's information with autogenerated id",
            "type": "object",
            "properties": {
                "author": {
                    "description": "Authors of an article",
                    "type": "string"
                },
                "created": {
                    "description": "Date and time of publication of an article, format ISO 8601, example \"2022-12-22T20:57:12Z\"",
                    "type": "string"
                },
                "description": {
                    "description": "Descripton of an article",
                    "type": "string"
                },
                "id": {
                    "description": "Autogenerated id, format depends of dbms type"
                },
                "title": {
                    "description": "Title of an article",
                    "type": "string"
                },
                "url": {
                    "description": "Permanent link to an article",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.1",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Web-Service Swagger",
	Description:      "Swagger Web-Service for Let's Go Aggregator",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}