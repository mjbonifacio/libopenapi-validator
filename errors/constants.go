// Copyright 2023 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package errors

// Error Categories - High-level classification of validation errors
const (
	// ErrorCategorySchema indicates JSON schema validation failures where data doesn't match schema constraints
	ErrorCategorySchema = "schema"
	
	// ErrorCategoryRetrieval indicates missing resources or structural issues before content validation
	ErrorCategoryRetrieval = "retrieval"
	
	// ErrorCategoryStructural indicates format, encoding, or parsing issues with the request/response structure
	ErrorCategoryStructural = "structural"
)

// Validation Sources - Where the validation is being applied
const (
	// ValidationSourceRequestBody indicates validation of HTTP request body content
	ValidationSourceRequestBody = "requestBody"
	
	// ValidationSourceResponseBody indicates validation of HTTP response body content  
	ValidationSourceResponseBody = "responseBody"
	
	// ValidationSourceParameter indicates validation of HTTP parameters (query, path, header, cookie)
	ValidationSourceParameter = "parameter"
	
	// ValidationSourceDocument indicates validation of the OpenAPI document itself
	ValidationSourceDocument = "document"
)

// Parameter Types - Specific parameter locations
const (
	// ParameterTypeQuery indicates URL query parameters (?param=value)
	ParameterTypeQuery = "query"
	
	// ParameterTypePath indicates URL path parameters (/users/{id})
	ParameterTypePath = "path"
	
	// ParameterTypeHeader indicates HTTP header parameters
	ParameterTypeHeader = "header"
	
	// ParameterTypeCookie indicates HTTP cookie parameters
	ParameterTypeCookie = "cookie"
)

// Validation Types - Existing magic strings now as constants
const (
	// ValidationTypePath indicates path-related validation
	ValidationTypePath = "path"
	
	// ValidationTypeRequest indicates request-related validation
	ValidationTypeRequest = "request"
	
	// ValidationTypeResponse indicates response-related validation
	ValidationTypeResponse = "response"
	
	// ValidationTypeParameter indicates parameter-related validation
	ValidationTypeParameter = "parameter"
	
	// ValidationTypeQuery indicates query parameter validation
	ValidationTypeQuery = "query"
	
	// ValidationTypeHeader indicates header validation
	ValidationTypeHeader = "header"
	
	// ValidationTypeCookie indicates cookie validation
	ValidationTypeCookie = "cookie"
	
	// ValidationTypeBody indicates request/response body validation
	ValidationTypeBody = "body"
	
	// ValidationTypeSchema indicates schema validation
	ValidationTypeSchema = "schema"
)

// Validation SubTypes - Existing magic strings now as constants
const (
	// ValidationSubTypeMissing indicates a required resource is missing
	ValidationSubTypeMissing = "missing"
	
	// ValidationSubTypeMissingOperation indicates an HTTP operation is not defined
	ValidationSubTypeMissingOperation = "missingOperation"
	
	// ValidationSubTypeInvalid indicates invalid format or value
	ValidationSubTypeInvalid = "invalid"
	
	// ValidationSubTypeRequired indicates a required field/parameter is missing
	ValidationSubTypeRequired = "required"
	
	// ValidationSubTypeEnum indicates value not in allowed enumeration
	ValidationSubTypeEnum = "enum"
	
	// ValidationSubTypeType indicates incorrect data type
	ValidationSubTypeType = "type"
	
	// ValidationSubTypeFormat indicates incorrect format (e.g., date, email)
	ValidationSubTypeFormat = "format"
	
	// ValidationSubTypeEncoding indicates incorrect parameter encoding
	ValidationSubTypeEncoding = "encoding"
	
	// ValidationSubTypeContentType indicates unsupported or missing content type
	ValidationSubTypeContentType = "contentType"
	
	// ValidationSubTypeSchema indicates schema compilation or processing failure
	ValidationSubTypeSchema = "schema"
)

// HTTP Methods - Standard HTTP methods as constants
const (
	// HTTPMethodGet represents HTTP GET method
	HTTPMethodGet = "GET"
	
	// HTTPMethodPost represents HTTP POST method
	HTTPMethodPost = "POST"
	
	// HTTPMethodPut represents HTTP PUT method
	HTTPMethodPut = "PUT"
	
	// HTTPMethodPatch represents HTTP PATCH method
	HTTPMethodPatch = "PATCH"
	
	// HTTPMethodDelete represents HTTP DELETE method
	HTTPMethodDelete = "DELETE"
	
	// HTTPMethodHead represents HTTP HEAD method
	HTTPMethodHead = "HEAD"
	
	// HTTPMethodOptions represents HTTP OPTIONS method
	HTTPMethodOptions = "OPTIONS"
	
	// HTTPMethodTrace represents HTTP TRACE method
	HTTPMethodTrace = "TRACE"
)
