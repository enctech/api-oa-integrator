// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
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
        "/auth/login": {
            "post": {
                "description": "For user to login into admin",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "user login",
                "parameters": [
                    {
                        "description": "Request Body",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/auth.LoginRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/auth/user": {
            "post": {
                "description": "For admin to create new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "create new user",
                "parameters": [
                    {
                        "description": "Request Body",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/auth.CreateUserRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/auth/user/{id}": {
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "For admin to delete user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "delete user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/config/integrator-config": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Get configurations for all integrators",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "summary": "Get configs for all integrator",
                "responses": {}
            },
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Create configuration required for OA to send data to integrator.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "summary": "Create config for integrator",
                "parameters": [
                    {
                        "description": "Request Body",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/config.IntegratorConfig"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/config/integrator-config/{id}": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Create configuration required for OA to send data to integrator.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "summary": "Create config for integrator",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request Body",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/config.IntegratorConfig"
                        }
                    }
                ],
                "responses": {}
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Create configuration required for OA to send data to integrator.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "summary": "Delete config for integrator",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/config/snb-config": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Get all configuration required for OA to works.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "summary": "Get all config for snb",
                "responses": {}
            },
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Create configuration required for OA to works.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "summary": "Create config for snb",
                "parameters": [
                    {
                        "description": "Request Body",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/config.SnbConfig"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/config/snb-config/{id}": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Get configuration required for OA to works.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "summary": "Get config for snb",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Update configuration required for OA to works.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "summary": "Update config for snb",
                "parameters": [
                    {
                        "description": "Request Body",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/config.SnbConfig"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Delete configuration required for OA to works.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "summary": "Delete config for snb",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/health": {
            "get": {
                "description": "To check overall system health",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "check system health",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Facility",
                        "name": "facility",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Device",
                        "name": "device",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/oa/{vendor}/AuthorizationService3rdParty/version": {
            "put": {
                "description": "get the version and configuration available",
                "consumes": [
                    "application/xml"
                ],
                "produces": [
                    "application/xml"
                ],
                "tags": [
                    "oa"
                ],
                "summary": "check version",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Vendor",
                        "name": "vendor",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request Body",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/oa.VersionRequestWrapper"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/oa/{vendor}/AuthorizationService3rdParty/{facility}/{device}/{jobId}": {
            "post": {
                "description": "Creates new job and sends the required information as URI and \u003cjob\u003e element to 3rd party system.",
                "consumes": [
                    "application/xml"
                ],
                "produces": [
                    "application/xml"
                ],
                "tags": [
                    "oa"
                ],
                "summary": "S\u0026B creates new job",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Facility",
                        "name": "facility",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Device",
                        "name": "device",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Job ID",
                        "name": "jobId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Vendor",
                        "name": "vendor",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request Body",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/oa.JobWrapper"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/oa/{vendor}/AuthorizationService3rdParty/{facility}/{device}/{jobId}/cancel": {
            "put": {
                "description": "This request cancels a running job on the 3rd party side. The job is identified by its resource /facility/device/jobid",
                "consumes": [
                    "application/xml"
                ],
                "produces": [
                    "application/xml"
                ],
                "tags": [
                    "oa"
                ],
                "summary": "Cancels a running job",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Facility",
                        "name": "facility",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Device",
                        "name": "device",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Job ID",
                        "name": "jobId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Vendor",
                        "name": "vendor",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request Body",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/oa.CancelJobWrapper"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/oa/{vendor}/AuthorizationService3rdParty/{facility}/{device}/{jobId}/finalmessage": {
            "put": {
                "description": "This request sends the last message for a job. The job is identified by its resources /facility/device/jobid",
                "consumes": [
                    "application/xml"
                ],
                "produces": [
                    "application/xml"
                ],
                "tags": [
                    "oa"
                ],
                "summary": "Receive Final Message from S\u0026B",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Facility",
                        "name": "facility",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Device",
                        "name": "device",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Job ID",
                        "name": "jobId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Vendor",
                        "name": "vendor",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request Body",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/oa.FinalMessageSBWrapper"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/oa/{vendor}/AuthorizationService3rdParty/{facility}/{device}/{jobId}/medialist": {
            "post": {
                "description": "Creates new media data for an existing job and sends the required information as a \u003cmediaData\u003e element to the 3rd party system.",
                "consumes": [
                    "application/xml"
                ],
                "produces": [
                    "application/xml"
                ],
                "tags": [
                    "oa"
                ],
                "summary": "Creates new media data",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Facility",
                        "name": "facility",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Device",
                        "name": "device",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Job ID",
                        "name": "jobId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Vendor",
                        "name": "vendor",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request Body",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/oa.MediaDataWrapper"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/transactions/logs": {
            "get": {
                "description": "To check overall system health",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transactions"
                ],
                "summary": "get all logs",
                "parameters": [
                    {
                        "type": "string",
                        "format": "dateTime",
                        "description": "Before",
                        "name": "before",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "dateTime",
                        "description": "After",
                        "name": "after",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Message",
                        "name": "message",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Fields",
                        "name": "fields",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "PerPage",
                        "name": "perPage",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "auth.CreateUserRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "permission": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "auth.LoginRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "config.IntegratorConfig": {
            "type": "object",
            "properties": {
                "clientId": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "insecureSkipVerify": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "plazaIdMap": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "providerId": {
                    "type": "integer"
                },
                "serviceProviderId": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "config.SnbConfig": {
            "type": "object",
            "properties": {
                "devices": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "endpoint": {
                    "type": "string"
                },
                "facilities": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "oa.BusinessTransaction": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "oa.CancelJobWrapper": {
            "type": "object",
            "properties": {
                "cancel": {
                    "type": "object",
                    "properties": {
                        "reason": {
                            "type": "object",
                            "properties": {
                                "cancelCode": {
                                    "type": "string"
                                },
                                "reasonText": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "oa.Configuration": {
            "type": "object",
            "properties": {
                "supportedFunctions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "oa.Customer": {
            "type": "object",
            "properties": {
                "customerGroup": {
                    "type": "string"
                },
                "customerId": {
                    "type": "string"
                },
                "customerName": {
                    "type": "string"
                }
            }
        },
        "oa.CustomerInformation": {
            "type": "object",
            "properties": {
                "customer": {
                    "$ref": "#/definitions/oa.Customer"
                }
            }
        },
        "oa.FinalMessageSB": {
            "type": "object",
            "properties": {
                "finalState": {
                    "type": "string"
                },
                "paymentMedia": {
                    "type": "string"
                }
            }
        },
        "oa.FinalMessageSBWrapper": {
            "type": "object",
            "properties": {
                "finalMessageSB": {
                    "$ref": "#/definitions/oa.FinalMessageSB"
                }
            }
        },
        "oa.Job": {
            "type": "object",
            "properties": {
                "businessTransaction": {
                    "$ref": "#/definitions/oa.BusinessTransaction"
                },
                "customerInformation": {
                    "$ref": "#/definitions/oa.CustomerInformation"
                },
                "jobId": {
                    "type": "object",
                    "properties": {
                        "id": {
                            "type": "string"
                        }
                    }
                },
                "jobType": {
                    "type": "string"
                },
                "mediaDataList": {
                    "type": "object",
                    "properties": {
                        "identifier": {
                            "type": "object",
                            "properties": {
                                "name": {
                                    "type": "string"
                                }
                            }
                        },
                        "mediaType": {
                            "type": "string"
                        }
                    }
                },
                "paymentData": {
                    "$ref": "#/definitions/oa.PaymentData"
                },
                "providerInformation": {
                    "type": "object",
                    "properties": {
                        "provider": {
                            "type": "object",
                            "properties": {
                                "providerId": {
                                    "type": "string"
                                },
                                "providerName": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                },
                "timeAndPlace": {
                    "type": "object",
                    "properties": {
                        "computer": {
                            "type": "object",
                            "properties": {
                                "computerNumber": {
                                    "type": "string"
                                }
                            }
                        },
                        "device": {
                            "type": "object",
                            "properties": {
                                "deviceNumber": {
                                    "type": "string"
                                },
                                "deviceType": {
                                    "type": "string"
                                }
                            }
                        },
                        "facility": {
                            "type": "object",
                            "properties": {
                                "facilityNumber": {
                                    "type": "string"
                                }
                            }
                        },
                        "operator": {
                            "type": "object",
                            "properties": {
                                "operatorNumber": {
                                    "type": "string"
                                }
                            }
                        },
                        "transactionTimeStamp": {
                            "type": "object",
                            "properties": {
                                "timeStamp": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "oa.JobWrapper": {
            "type": "object",
            "properties": {
                "job": {
                    "$ref": "#/definitions/oa.Job"
                }
            }
        },
        "oa.MediaDataWrapper": {
            "type": "object",
            "properties": {
                "mediaData": {
                    "type": "object",
                    "properties": {
                        "hashValue": {
                            "type": "object",
                            "properties": {
                                "value": {
                                    "type": "string"
                                }
                            }
                        },
                        "mediaType": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "oa.OriginalAmount": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "string"
                },
                "vatRate": {
                    "type": "string"
                }
            }
        },
        "oa.PaymentData": {
            "type": "object",
            "properties": {
                "originalAmount": {
                    "$ref": "#/definitions/oa.OriginalAmount"
                },
                "remainingAmount": {
                    "type": "object",
                    "properties": {
                        "amount": {
                            "type": "string"
                        },
                        "text": {
                            "type": "string"
                        },
                        "vatRate": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "oa.VersionRequestWrapper": {
            "type": "object",
            "properties": {
                "version": {
                    "type": "object",
                    "properties": {
                        "configuration": {
                            "$ref": "#/definitions/oa.Configuration"
                        },
                        "entervoVersion": {
                            "type": "string"
                        },
                        "sbAuthorizationVersion": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "description": "Type \"Bearer\" followed by a space and JWT token.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Swagger OA Integrator API",
	Description:      "This is a server OA integrator.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
