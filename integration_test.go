// Copyright 2023 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package validator

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi-validator/errors"
	"github.com/stretchr/testify/require"
)

// TestValidationErrorIntegration demonstrates the complete error categorization and ValidationSource system
func TestValidationErrorIntegration(t *testing.T) {
	// OpenAPI spec with various validation scenarios
	spec := `
openapi: 3.0.3
info:
  title: Test API
  version: 1.0.0
paths:
  /pets/{petId}:
    get:
      parameters:
        - name: petId
          in: path
          required: true
          schema:
            type: integer
        - name: status
          in: query
          required: true
          schema:
            type: string
            enum: [available, pending, sold]
        - name: api-key
          in: header
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Pet found
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: integer
                  name:
                    type: string
                required: [id, name]
    post:
      parameters:
        - name: petId
          in: path
          required: true
          schema:
            type: integer
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                status:
                  type: string
                  enum: [available, pending, sold]
              required: [name, status]
      responses:
        '201':
          description: Pet created
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: integer
                  name:
                    type: string
                required: [id, name]
`

	// Create validator
	doc, err := libopenapi.NewDocument([]byte(spec))
	require.NoError(t, err)
	
	validator, errs := NewValidator(doc)
	require.Empty(t, errs)

	t.Run("Schema Error - Request Body Validation", func(t *testing.T) {
		// Invalid request body (missing required field)
		body := `{"name": "Fluffy"}`  // missing required "status" field
		req, _ := http.NewRequest("POST", "/pets/123", bytes.NewReader([]byte(body)))
		req.Header.Set("Content-Type", "application/json")
		
		valid, validationErrors := validator.ValidateHttpRequest(req)
		require.False(t, valid)
		require.NotEmpty(t, validationErrors)
		
		// Find the schema validation error
		var schemaError *errors.ValidationError
		for _, ve := range validationErrors {
			if ve.IsSchemaError() {
				schemaError = ve
				break
			}
		}
		require.NotNil(t, schemaError, "Should have a schema validation error")
		
		// Validate error categorization
		require.Equal(t, errors.ErrorCategorySchema, schemaError.ErrorCategory)
		require.True(t, schemaError.IsSchemaError())
		require.False(t, schemaError.IsRetrievalError())
		require.False(t, schemaError.IsStructuralError())
		
		// Validate ValidationSource
		require.NotEmpty(t, schemaError.SchemaValidationErrors)
		for _, sve := range schemaError.SchemaValidationErrors {
			require.Equal(t, errors.ValidationSourceRequestBody, sve.ValidationSource)
		}
	})

	t.Run("Retrieval Error - Missing Required Parameter", func(t *testing.T) {
		// Missing required query parameter
		req, _ := http.NewRequest("GET", "/pets/123", nil)  // missing required "status" query param
		
		valid, validationErrors := validator.ValidateHttpRequest(req)
		require.False(t, valid)
		require.NotEmpty(t, validationErrors)
		
		// Find the retrieval error
		var retrievalError *errors.ValidationError
		for _, ve := range validationErrors {
			if ve.IsRetrievalError() {
				retrievalError = ve
				break
			}
		}
		require.NotNil(t, retrievalError, "Should have a retrieval error")
		
		// Validate error categorization
		require.Equal(t, errors.ErrorCategoryRetrieval, retrievalError.ErrorCategory)
		require.False(t, retrievalError.IsSchemaError())
		require.True(t, retrievalError.IsRetrievalError())
		require.False(t, retrievalError.IsStructuralError())
	})

	t.Run("Schema Error - Parameter Enum Validation", func(t *testing.T) {
		// Invalid enum value for query parameter
		req, _ := http.NewRequest("GET", "/pets/123?status=invalid", nil)
		
		valid, validationErrors := validator.ValidateHttpRequest(req)
		require.False(t, valid)
		require.NotEmpty(t, validationErrors)
		
		// Find the schema validation error for parameter
		var paramSchemaError *errors.ValidationError
		for _, ve := range validationErrors {
			if ve.IsSchemaError() && ve.ValidationType == errors.ValidationTypeQuery {
				paramSchemaError = ve
				break
			}
		}
		require.NotNil(t, paramSchemaError, "Should have a parameter schema validation error")
		
		// Validate error categorization
		require.Equal(t, errors.ErrorCategorySchema, paramSchemaError.ErrorCategory)
		require.True(t, paramSchemaError.IsSchemaError())
		require.Equal(t, errors.ValidationSubTypeEnum, paramSchemaError.ValidationSubType)
		
		// Validate ValidationSource for parameter
		require.NotEmpty(t, paramSchemaError.SchemaValidationErrors)
		for _, sve := range paramSchemaError.SchemaValidationErrors {
			require.Equal(t, errors.ValidationSourceParameter, sve.ValidationSource)
		}
	})

	t.Run("Schema Error - Path Parameter Type Validation", func(t *testing.T) {
		// Invalid type for path parameter
		req, _ := http.NewRequest("GET", "/pets/not-a-number?status=available", nil)
		
		valid, validationErrors := validator.ValidateHttpRequest(req)
		require.False(t, valid)
		require.NotEmpty(t, validationErrors)
		
		// Find the schema validation error for path parameter
		var pathSchemaError *errors.ValidationError
		for _, ve := range validationErrors {
			if ve.IsSchemaError() && ve.ValidationType == errors.ValidationTypePath {
				pathSchemaError = ve
				break
			}
		}
		require.NotNil(t, pathSchemaError, "Should have a path parameter schema validation error")
		
		// Validate error categorization
		require.Equal(t, errors.ErrorCategorySchema, pathSchemaError.ErrorCategory)
		require.True(t, pathSchemaError.IsSchemaError())
		require.Equal(t, errors.ValidationSubTypeType, pathSchemaError.ValidationSubType)
		
		// Validate ValidationSource for parameter
		require.NotEmpty(t, pathSchemaError.SchemaValidationErrors)
		for _, sve := range pathSchemaError.SchemaValidationErrors {
			require.Equal(t, errors.ValidationSourceParameter, sve.ValidationSource)
		}
	})

	t.Run("Retrieval Error - Path Not Found", func(t *testing.T) {
		// Non-existent path
		req, _ := http.NewRequest("GET", "/nonexistent", nil)
		
		valid, validationErrors := validator.ValidateHttpRequest(req)
		require.False(t, valid)
		require.NotEmpty(t, validationErrors)
		
		// Find the path not found error
		var pathError *errors.ValidationError
		for _, ve := range validationErrors {
			if ve.ValidationType == errors.ValidationTypePath && ve.ValidationSubType == errors.ValidationSubTypeMissing {
				pathError = ve
				break
			}
		}
		require.NotNil(t, pathError, "Should have a path not found error")
		
		// Validate error categorization
		require.Equal(t, errors.ErrorCategoryRetrieval, pathError.ErrorCategory)
		require.False(t, pathError.IsSchemaError())
		require.True(t, pathError.IsRetrievalError())
		require.False(t, pathError.IsStructuralError())
	})
}

// TestValidationSourceDetermination tests the automatic ValidationSource determination
func TestValidationSourceDetermination(t *testing.T) {
	testCases := []struct {
		name           string
		validationType string
		expectedSource string
	}{
		{
			name:           "Request validation should map to requestBody",
			validationType: errors.ValidationTypeRequest,
			expectedSource: errors.ValidationSourceRequestBody,
		},
		{
			name:           "Response validation should map to responseBody",
			validationType: errors.ValidationTypeResponse,
			expectedSource: errors.ValidationSourceResponseBody,
		},
		{
			name:           "Parameter validation should map to parameter",
			validationType: errors.ValidationTypeParameter,
			expectedSource: errors.ValidationSourceParameter,
		},
		{
			name:           "Query validation should map to parameter",
			validationType: errors.ValidationTypeQuery,
			expectedSource: errors.ValidationSourceParameter,
		},
		{
			name:           "Header validation should map to parameter",
			validationType: errors.ValidationTypeHeader,
			expectedSource: errors.ValidationSourceParameter,
		},
		{
			name:           "Cookie validation should map to parameter",
			validationType: errors.ValidationTypeCookie,
			expectedSource: errors.ValidationSourceParameter,
		},
		{
			name:           "Schema validation should map to document",
			validationType: errors.ValidationTypeSchema,
			expectedSource: errors.ValidationSourceDocument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ve := &errors.ValidationError{
				ValidationType: tc.validationType,
			}
			
			actualSource := ve.DetermineValidationSource()
			require.Equal(t, tc.expectedSource, actualSource)
		})
	}
}

// TestErrorCategoryClassification tests the automatic error category classification
func TestErrorCategoryClassification(t *testing.T) {
	testCases := []struct {
		name               string
		validationType     string
		validationSubType  string
		hasSchemaErrors    bool
		expectedCategory   string
	}{
		{
			name:             "Schema validation errors should be schema category",
			hasSchemaErrors:  true,
			expectedCategory: errors.ErrorCategorySchema,
		},
		{
			name:               "Path missing should be retrieval category",
			validationType:     errors.ValidationTypePath,
			validationSubType:  errors.ValidationSubTypeMissing,
			expectedCategory:   errors.ErrorCategoryRetrieval,
		},
		{
			name:               "Operation missing should be retrieval category",
			validationType:     errors.ValidationTypePath,
			validationSubType:  errors.ValidationSubTypeMissingOperation,
			expectedCategory:   errors.ErrorCategoryRetrieval,
		},
		{
			name:               "Request content type should be retrieval category",
			validationType:     errors.ValidationTypeRequest,
			validationSubType:  errors.ValidationSubTypeContentType,
			expectedCategory:   errors.ErrorCategoryRetrieval,
		},
		{
			name:               "Request missing should be retrieval category",
			validationType:     errors.ValidationTypeRequest,
			validationSubType:  errors.ValidationSubTypeMissing,
			expectedCategory:   errors.ErrorCategoryRetrieval,
		},
		{
			name:               "Request schema should be structural category",
			validationType:     errors.ValidationTypeRequest,
			validationSubType:  errors.ValidationSubTypeSchema,
			expectedCategory:   errors.ErrorCategoryStructural,
		},
		{
			name:               "Parameter missing should be retrieval category",
			validationType:     errors.ValidationTypeParameter,
			validationSubType:  errors.ValidationSubTypeMissing,
			expectedCategory:   errors.ErrorCategoryRetrieval,
		},
		{
			name:               "Parameter required should be retrieval category",
			validationType:     errors.ValidationTypeParameter,
			validationSubType:  errors.ValidationSubTypeRequired,
			expectedCategory:   errors.ErrorCategoryRetrieval,
		},
		{
			name:               "Parameter format should be structural category",
			validationType:     errors.ValidationTypeParameter,
			validationSubType:  errors.ValidationSubTypeFormat,
			expectedCategory:   errors.ErrorCategoryStructural,
		},
		{
			name:               "Unknown type should default to structural category",
			validationType:     "unknown",
			validationSubType:  "unknown",
			expectedCategory:   errors.ErrorCategoryStructural,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ve := &errors.ValidationError{
				ValidationType:    tc.validationType,
				ValidationSubType: tc.validationSubType,
			}
			
			if tc.hasSchemaErrors {
				ve.SchemaValidationErrors = []*errors.SchemaValidationFailure{
					{
						Reason:   "Test schema error",
						Location: "/test",
					},
				}
			}
			
			ve.SetErrorCategory()
			require.Equal(t, tc.expectedCategory, ve.ErrorCategory)
		})
	}
}
