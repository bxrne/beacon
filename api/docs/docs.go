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
        "/health": {
            "get": {
                "description": "Get the health status of the server",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.healthResponse"
                        }
                    }
                }
            }
        },
        "/metric": {
            "get": {
                "description": "Get metrics for a device",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "metrics"
                ],
                "summary": "Get metrics",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Device ID",
                        "name": "X-DeviceID",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/metrics.DeviceMetrics"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Submit metrics for a device",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "metrics"
                ],
                "summary": "Submit metrics",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Device ID",
                        "name": "X-DeviceID",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Metrics data",
                        "name": "metrics",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/metrics.DeviceMetrics"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/metrics.DeviceMetrics"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "metrics.DeviceMetrics": {
            "type": "object",
            "properties": {
                "metrics": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/metrics.Metric"
                    }
                }
            }
        },
        "metrics.Metric": {
            "type": "object",
            "properties": {
                "type": {
                    "description": "References metric_types.name",
                    "type": "string"
                },
                "unit": {
                    "description": "References units.name",
                    "type": "string"
                },
                "value": {
                    "type": "number"
                }
            }
        },
        "server.errorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "server.healthResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Beacon API",
	Description:      "Collects device and metric data from clients",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
