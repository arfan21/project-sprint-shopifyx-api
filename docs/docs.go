// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
            "name": "API Support",
            "url": "http://www.synapsis.id"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/product": {
            "post": {
                "description": "Create product",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "Create product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "With the bearer started",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Payload product create request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_internal_model.ProductCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_pkg_pkgutil.HTTPResponse"
                        }
                    },
                    "400": {
                        "description": "Error validation field",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_pkg_pkgutil.HTTPResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_pkg_pkgutil.ErrValidationResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_pkg_pkgutil.HTTPResponse"
                        }
                    }
                }
            }
        },
        "/v1/user/login": {
            "post": {
                "description": "Login user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "Payload user Login Request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_internal_model.UserLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_pkg_pkgutil.HTTPResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_internal_model.UserLoginResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Error validation field",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_pkg_pkgutil.HTTPResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_pkg_pkgutil.ErrValidationResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_pkg_pkgutil.HTTPResponse"
                        }
                    }
                }
            }
        },
        "/v1/user/register": {
            "post": {
                "description": "Register user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Register user",
                "parameters": [
                    {
                        "description": "Payload user Register Request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_internal_model.UserRegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_pkg_pkgutil.HTTPResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_internal_model.UserLoginResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Error validation field",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_pkg_pkgutil.HTTPResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_pkg_pkgutil.ErrValidationResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_project-sprint-shopifyx-api_pkg_pkgutil.HTTPResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_arfan21_project-sprint-shopifyx-api_internal_model.ProductCreateRequest": {
            "type": "object",
            "required": [
                "condition",
                "imageUrl",
                "isPurchaseable",
                "name",
                "price",
                "stock"
            ],
            "properties": {
                "condition": {
                    "type": "string",
                    "enum": [
                        "new",
                        "second"
                    ]
                },
                "imageUrl": {
                    "type": "string"
                },
                "isPurchaseable": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string",
                    "maxLength": 60,
                    "minLength": 5
                },
                "price": {
                    "type": "number"
                },
                "stock": {
                    "type": "integer"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "github_com_arfan21_project-sprint-shopifyx-api_internal_model.UserLoginRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "maxLength": 15,
                    "minLength": 5
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "github_com_arfan21_project-sprint-shopifyx-api_internal_model.UserLoginResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "github_com_arfan21_project-sprint-shopifyx-api_internal_model.UserRegisterRequest": {
            "type": "object",
            "required": [
                "name",
                "password",
                "username"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 5
                },
                "password": {
                    "type": "string",
                    "maxLength": 15,
                    "minLength": 5
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "github_com_arfan21_project-sprint-shopifyx-api_pkg_pkgutil.ErrValidationResponse": {
            "type": "object",
            "properties": {
                "field": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_arfan21_project-sprint-shopifyx-api_pkg_pkgutil.HTTPResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string",
                    "example": "Success"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "project-sprint-shopifyx-api",
	Description:      "This is a sample server cell for project-sprint-shopifyx-api.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}