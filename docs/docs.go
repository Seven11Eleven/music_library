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
        "/music": {
            "post": {
                "description": "Save a new music record",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Music"
                ],
                "summary": "Save new music",
                "parameters": [
                    {
                        "description": "Music input",
                        "name": "music",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.MusicQuery"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Music"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/music/info": {
            "get": {
                "description": "Retrieve a list of music based on provided filters",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Music"
                ],
                "summary": "Get list of music",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Release date in format YYYY-MM-DD",
                        "name": "release_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Music link",
                        "name": "link",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Name of the song",
                        "name": "song_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Name of the group",
                        "name": "group_name",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page number for pagination",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of records per page",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Music"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid date format",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/music/verses": {
            "get": {
                "description": "Retrieve verses of a music track with pagination",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Music"
                ],
                "summary": "Get verses of music",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Music ID",
                        "name": "music_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page number for pagination",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of records per page",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Verse"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/music/{id}": {
            "put": {
                "description": "Update an existing music record",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Music"
                ],
                "summary": "Update music",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Music ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated music object",
                        "name": "music",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Music"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Music"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a music record by ID",
                "tags": [
                    "Music"
                ],
                "summary": "Delete music",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Music ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Music deleted successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Music": {
            "type": "object",
            "properties": {
                "group_name": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "link": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string"
                },
                "song_name": {
                    "type": "string"
                },
                "verses": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Verse"
                    }
                }
            }
        },
        "models.MusicQuery": {
            "type": "object",
            "properties": {
                "group_name": {
                    "type": "string"
                },
                "song_name": {
                    "type": "string"
                }
            }
        },
        "models.Verse": {
            "type": "object",
            "properties": {
                "number": {
                    "type": "integer"
                },
                "text": {
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
