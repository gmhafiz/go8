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
        "contact": {
            "name": "User Name",
            "url": "https://github.com/gmhafiz/go8",
            "email": "email@example.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/health": {
            "get": {
                "description": "Hits this API to see if API is running in the server",
                "summary": "Checks if API is up",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/health/readiness": {
            "get": {
                "description": "Hits this API to see if both API and Database are running in the server",
                "summary": "Checks if both API and Database are up",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/v1/author": {
            "get": {
                "description": "Lists all authors. By default, it gets first page with 30 items.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Shows all authors",
                "parameters": [
                    {
                        "type": "string",
                        "description": "page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "limit of result",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "result offset",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "search by first_name",
                        "name": "first_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "search by last_name",
                        "name": "last_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "sort by fields name. E.g. first_name,asc",
                        "name": "sort",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/respond.Standard"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Create an author using JSON payload",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create an Author",
                "parameters": [
                    {
                        "description": "Create an author using the following format",
                        "name": "Author",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/author.CreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/author.GetResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/author/{id}": {
            "get": {
                "description": "Get an author by its id.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get an Author",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "author ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/gen.Author"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "Update an author by its model.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update an Author",
                "parameters": [
                    {
                        "description": "Author Request",
                        "name": "Author",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/author.UpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/gen.Author"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete an author by its id.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Delete an Author",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "author ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/book": {
            "get": {
                "description": "Lists all books. By default, it gets first page with 30 items.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Shows all books",
                "parameters": [
                    {
                        "type": "string",
                        "description": "page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "size of result",
                        "name": "size",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "search by title",
                        "name": "title",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "search by description",
                        "name": "description",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/book.Res"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a book using JSON payload",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create a Book",
                "parameters": [
                    {
                        "description": "Create a book using the following format",
                        "name": "Book",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/book.CreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/book.Res"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/book/{bookID}": {
            "get": {
                "description": "Get a book by its id.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get a Book",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "book ID",
                        "name": "bookID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/book.Res"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a book by its model.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update a Book",
                "parameters": [
                    {
                        "description": "Book UpdateRequest",
                        "name": "Book",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/book.UpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/book.Res"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a book by its id.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Delete a Book",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "book ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "author.Book": {
            "type": "object",
            "required": [
                "description",
                "published_date",
                "title"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "published_date": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "author.CreateRequest": {
            "type": "object",
            "required": [
                "first_name",
                "last_name"
            ],
            "properties": {
                "books": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/author.Book"
                    }
                },
                "first_name": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "middle_name": {
                    "type": "string"
                }
            }
        },
        "author.GetResponse": {
            "type": "object",
            "properties": {
                "books": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/book.Schema"
                    }
                },
                "first_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "last_name": {
                    "type": "string"
                },
                "middle_name": {
                    "type": "string"
                }
            }
        },
        "author.UpdateRequest": {
            "type": "object",
            "properties": {
                "first_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "last_name": {
                    "type": "string"
                },
                "middle_name": {
                    "type": "string"
                }
            }
        },
        "book.CreateRequest": {
            "type": "object",
            "required": [
                "description",
                "published_date",
                "title"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "image_url": {
                    "type": "string"
                },
                "published_date": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "book.Res": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "image_url": {
                    "type": "string"
                },
                "published_date": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "book.Schema": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "imageURL": {
                    "type": "string"
                },
                "publishedDate": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "book.UpdateRequest": {
            "type": "object",
            "required": [
                "description",
                "published_date",
                "title"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "image_url": {
                    "type": "string"
                },
                "published_date": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "gen.Author": {
            "type": "object",
            "properties": {
                "edges": {
                    "description": "Edges holds the relations/edges for other nodes in the graph.\nThe values are being populated by the AuthorQuery when eager-loading is set.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/gen.AuthorEdges"
                        }
                    ]
                },
                "first_name": {
                    "description": "FirstName holds the value of the \"first_name\" field.",
                    "type": "string"
                },
                "id": {
                    "description": "ID of the ent.",
                    "type": "integer"
                },
                "last_name": {
                    "description": "LastName holds the value of the \"last_name\" field.",
                    "type": "string"
                },
                "middle_name": {
                    "description": "MiddleName holds the value of the \"middle_name\" field.",
                    "type": "string"
                }
            }
        },
        "gen.AuthorEdges": {
            "type": "object",
            "properties": {
                "books": {
                    "description": "Books holds the value of the books edge.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/gen.Book"
                    }
                }
            }
        },
        "gen.Book": {
            "type": "object",
            "properties": {
                "edges": {
                    "description": "Edges holds the relations/edges for other nodes in the graph.\nThe values are being populated by the BookQuery when eager-loading is set.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/gen.BookEdges"
                        }
                    ]
                },
                "id": {
                    "description": "ID of the ent.",
                    "type": "integer"
                },
                "image_url": {
                    "description": "ImageURL holds the value of the \"image_url\" field.",
                    "type": "string"
                },
                "published_date": {
                    "description": "PublishedDate holds the value of the \"published_date\" field.",
                    "type": "string"
                },
                "title": {
                    "description": "Title holds the value of the \"title\" field.",
                    "type": "string"
                }
            }
        },
        "gen.BookEdges": {
            "type": "object",
            "properties": {
                "authors": {
                    "description": "Authors holds the value of the authors edge.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/gen.Author"
                    }
                }
            }
        },
        "respond.Meta": {
            "type": "object",
            "properties": {
                "size": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "respond.Standard": {
            "type": "object",
            "properties": {
                "data": {},
                "meta": {
                    "$ref": "#/definitions/respond.Meta"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1.0",
	Host:             "localhost:3080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Go8",
	Description:      "Go + Postgres + Chi router + sqlx + ent + Unit Testing starter kit for API development.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
