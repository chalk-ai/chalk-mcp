package tools

import (
	"context"
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

func TestParseMCPTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: map[string]string{},
		},
		{
			name:  "single required flag",
			input: "required",
			expected: map[string]string{
				"required": "true",
			},
		},
		{
			name:  "simple description",
			input: "description=Simple description",
			expected: map[string]string{
				"description": "Simple description",
			},
		},
		{
			name:  "required with description",
			input: "required,description=Path to the root directory",
			expected: map[string]string{
				"required":    "true",
				"description": "Path to the root directory",
			},
		},
		{
			name:  "description with complex punctuation",
			input: "required,description=Another complex description with commas, semicolons; and other stuff.",
			expected: map[string]string{
				"required":    "true",
				"description": "Another complex description with commas, semicolons; and other stuff.",
			},
		},
		{
			name:  "whitespace handling",
			input: " required , type=string , description=Spaced out description ",
			expected: map[string]string{
				"required":    "true",
				"type":        "string",
				"description": "Spaced out description ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := parseMCPTags(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateMetadataInterface(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		toolFunc func() Tool
		expected mcp.Tool
	}{
		{
			name: "simple tool with required string",
			toolFunc: func() Tool {
				return &testTool{
					name:        "test_tool",
					description: "A test tool",
					paramsType:  reflect.TypeOf(simpleParams{}),
				}
			},
			expected: mcp.NewTool("test_tool",
				mcp.WithDescription("A test tool"),
				mcp.WithString("field1", mcp.Required(), mcp.Description("A required field")),
			),
		},
		{
			name: "simple tool with required int",
			toolFunc: func() Tool {
				return &testTool{
					name:        "test_tool",
					description: "A test tool",
					paramsType:  reflect.TypeOf(simpleIntParams{}),
				}
			},
			expected: mcp.NewTool("test_tool",
				mcp.WithDescription("A test tool"),
				mcp.WithNumber("field1", mcp.Required(), mcp.Description("A required field")),
			),
		},
		{
			name: "simple tool with required float",
			toolFunc: func() Tool {
				return &testTool{
					name:        "test_tool",
					description: "A test tool",
					paramsType:  reflect.TypeOf(simpleFloatParams{}),
				}
			},
			expected: mcp.NewTool("test_tool",
				mcp.WithDescription("A test tool"),
				mcp.WithNumber("field1", mcp.Required(), mcp.Description("A required field")),
			),
		},
		{
			name: "simple tool with required bool",
			toolFunc: func() Tool {
				return &testTool{
					name:        "test_tool",
					description: "A test tool",
					paramsType:  reflect.TypeOf(simpleBoolParams{}),
				}
			},
			expected: mcp.NewTool("test_tool",
				mcp.WithDescription("A test tool"),
				mcp.WithBoolean("field1", mcp.Required(), mcp.Description("A required field")),
			),
		},
		{
			name: "tool with optional string",
			toolFunc: func() Tool {
				return &testTool{
					name:        "optional_tool",
					description: "Tool with optional field",
					paramsType:  reflect.TypeOf(optionalParams{}),
				}
			},
			expected: mcp.NewTool("optional_tool",
				mcp.WithDescription("Tool with optional field"),
				mcp.WithString("field1", mcp.Description("An optional field")),
			),
		},
		{
			name: "tool with map field",
			toolFunc: func() Tool {
				return &testTool{
					name:        "map_tool",
					description: "Tool with map field",
					paramsType:  reflect.TypeOf(mapParams{}),
				}
			},
			expected: mcp.NewTool("map_tool",
				mcp.WithDescription("Tool with map field"),
				mcp.WithObject("features",
					mcp.Description("Map of features"),
					mcp.AdditionalProperties(map[string]any{"type": "string"}),
				),
			),
		},
		{
			name: "tool with slice field",
			toolFunc: func() Tool {
				return &testTool{
					name:        "slice_tool",
					description: "Tool with slice field",
					paramsType:  reflect.TypeOf(sliceParams{}),
				}
			},
			expected: mcp.NewTool("slice_tool",
				mcp.WithDescription("Tool with slice field"),
				mcp.WithArray("items",
					mcp.Description("List of items"),
					mcp.Items(map[string]any{"type": "string"}),
				),
			),
		},
		{
			name: "tool with mixed fields",
			toolFunc: func() Tool {
				return &testTool{
					name:        "mixed_tool",
					description: "Tool with mixed field types",
					paramsType:  reflect.TypeOf(mixedParams{}),
				}
			},
			expected: mcp.NewTool("mixed_tool",
				mcp.WithDescription("Tool with mixed field types"),
				mcp.WithString("required_field", mcp.Required(), mcp.Description("Required string field")),
				mcp.WithString("optional_field", mcp.Description("Optional string field")),
				mcp.WithObject("features",
					mcp.Description("Map of features"),
					mcp.AdditionalProperties(map[string]any{"type": "string"}),
				),
				mcp.WithArray("items",
					mcp.Description("List of items"),
					mcp.Items(map[string]any{"type": "string"}),
				),
			),
		},
		{
			name: "tool with skipped field",
			toolFunc: func() Tool {
				return &testTool{
					name:        "skipped_tool",
					description: "Tool with skipped field",
					paramsType:  reflect.TypeOf(skipParams{}),
				}
			},
			expected: mcp.NewTool("skipped_tool",
				mcp.WithDescription("Tool with skipped field"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tool := tt.toolFunc()
			result := GenerateMetadata(tool)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		args       map[string]any
		paramsType reflect.Type
		expected   any
		wantErr    bool
	}{
		{
			name: "successful simple params",
			args: map[string]any{
				"field1": "test value",
			},
			paramsType: reflect.TypeOf(simpleParams{}),
			expected: &simpleParams{
				Field1: "test value",
			},
		},
		{
			name:       "missing required field",
			args:       map[string]any{},
			paramsType: reflect.TypeOf(simpleParams{}),
			expected:   nil,
			wantErr:    true,
		},
		{
			name: "coalesced type for simple field",
			args: map[string]any{
				"field1": 123,
			},
			paramsType: reflect.TypeOf(simpleParams{}),
			expected: &simpleParams{
				Field1: "123",
			},
		},
		{
			name: "successful simple int",
			args: map[string]any{
				"field1": 123,
			},
			paramsType: reflect.TypeOf(simpleIntParams{}),
			expected: &simpleIntParams{
				Field1: 123,
			},
		},
		{
			name: "coalesced type for simple int",
			args: map[string]any{
				"field1": "123",
			},
			paramsType: reflect.TypeOf(simpleIntParams{}),
			expected: &simpleIntParams{
				Field1: 123,
			},
		},
		{
			name: "coalesced type for simple failure",
			args: map[string]any{
				"field1": "abc",
			},
			paramsType: reflect.TypeOf(simpleIntParams{}),
			expected:   nil,
			wantErr:    true,
		},
		{
			name: "successful simple float",
			args: map[string]any{
				"field1": 123.456,
			},
			paramsType: reflect.TypeOf(simpleFloatParams{}),
			expected: &simpleFloatParams{
				Field1: 123.456,
			},
		},
		{
			name: "coalesced type for simple float",
			args: map[string]any{
				"field1": "123.456",
			},
			paramsType: reflect.TypeOf(simpleFloatParams{}),
			expected: &simpleFloatParams{
				Field1: 123.456,
			},
		},
		{
			name: "coalesced type for simple float failure",
			args: map[string]any{
				"field1": "abc",
			},
			paramsType: reflect.TypeOf(simpleFloatParams{}),
			expected:   nil,
			wantErr:    true,
		},
		{
			name: "successful simple bool",
			args: map[string]any{
				"field1": true,
			},
			paramsType: reflect.TypeOf(simpleBoolParams{}),
			expected: &simpleBoolParams{
				Field1: true,
			},
		},
		{
			name: "coalesced type for simple bool",
			args: map[string]any{
				"field1": "true",
			},
			paramsType: reflect.TypeOf(simpleBoolParams{}),
			expected: &simpleBoolParams{
				Field1: true,
			},
		},
		{
			name: "coalesced type for simple bool failure",
			args: map[string]any{
				"field1": "abc",
			},
			paramsType: reflect.TypeOf(simpleBoolParams{}),
			expected:   nil,
			wantErr:    true,
		},
		{
			name:       "missing optional field",
			args:       map[string]any{},
			paramsType: reflect.TypeOf(optionalParams{}),
			expected:   &optionalParams{},
		},
		{
			name: "successful map params",
			args: map[string]any{
				"features": map[string]any{
					"feature1": "test value",
				},
			},
			paramsType: reflect.TypeOf(mapParams{}),
			expected: &mapParams{
				Features: map[string]any{
					"feature1": "test value",
				},
			},
		},
		{
			name: "coalesced type for map field",
			args: map[string]any{
				"features": []map[string]any{
					{
						"feature1": 123,
						"feature2": 456,
					},
				},
			},
			paramsType: reflect.TypeOf(mapParams{}),
			expected: &mapParams{
				Features: map[string]any{
					"feature1": "123",
					"feature2": "456",
				},
			},
		},
		{
			name: "successful slice params",
			args: map[string]any{
				"items": []string{"test value"},
			},
			paramsType: reflect.TypeOf(sliceParams{}),
			expected: &sliceParams{
				Items: []string{"test value"},
			},
		},
		{
			name: "coalesced type for slice field",
			args: map[string]any{
				"items": []int{123, 456},
			},
			paramsType: reflect.TypeOf(sliceParams{}),
			expected: &sliceParams{
				Items: []string{"123", "456"},
			},
		},
		{
			name: "successful skipped field",
			args: map[string]any{
				"field1": "test value",
			},
			paramsType: reflect.TypeOf(skipParams{}),
			expected:   &skipParams{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := parseParams(tt.paramsType, tt.args)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// Test helpers

type simpleParams struct {
	Field1 string `json:"field1" mcp:"required,description=A required field"`
}

type simpleIntParams struct {
	Field1 int `json:"field1" mcp:"required,description=A required field"`
}

type simpleFloatParams struct {
	Field1 float64 `json:"field1" mcp:"required,description=A required field"`
}

type simpleBoolParams struct {
	Field1 bool `json:"field1" mcp:"required,description=A required field"`
}

type optionalParams struct {
	Field1 string `json:"field1" mcp:"description=An optional field"`
}

type mapParams struct {
	Features map[string]any `json:"features" mcp:"description=Map of features"`
}

type sliceParams struct {
	Items []string `json:"items" mcp:"description=List of items"`
}

type skipParams struct {
	Field1 string `json:"field1" mcp:"-"`
}

type mixedParams struct {
	RequiredField string         `json:"required_field" mcp:"required,description=Required string field"`
	OptionalField string         `json:"optional_field" mcp:"description=Optional string field"`
	Features      map[string]any `json:"features" mcp:"description=Map of features"`
	Items         []string       `json:"items" mcp:"description=List of items"`
}

type testTool struct {
	name        string
	description string
	paramsType  reflect.Type
}

func (t *testTool) Name() string {
	return t.name
}

func (t *testTool) Description() string {
	return t.description
}

func (t *testTool) ParamsType() reflect.Type {
	return t.paramsType
}

func (t *testTool) Execute(ctx context.Context, args any) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("test result"), nil
}
