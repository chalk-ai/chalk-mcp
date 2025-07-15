package tools

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/chalk-ai/chalk-mcp/utils"
	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/mock"
)

type Tool interface {
	// Name returns the tool name
	Name() string

	// Description returns the tool description
	Description() string

	// Execute runs the tool with parsed arguments
	Execute(ctx context.Context, args any) (*mcp.CallToolResult, error)

	// ParamsType returns the type of the params struct
	ParamsType() reflect.Type
}

type CommandExecutor interface {
	Execute(ctx context.Context, workDir string, args ...string) ([]byte, error)
}

// DefaultCommandExecutor calls chalk command with the given args
type DefaultCommandExecutor struct{}

func (e *DefaultCommandExecutor) Execute(ctx context.Context, workDir string, args ...string) ([]byte, error) {
	cmd, err := utils.GetChalkCommand(ctx, workDir, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get chalk command")
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errors.New("chalk command timed out")
		}
		return nil, errors.Wrapf(err, "failed to execute chalk command; output: %s", string(output))
	}
	return output, nil
}

// MockCommandExecutor for testing
type MockCommandExecutor struct {
	mock.Mock
}

func (m *MockCommandExecutor) Execute(ctx context.Context, workDir string, args ...string) ([]byte, error) {
	argList := m.Called(ctx, workDir, args)
	return argList.Get(0).([]byte), argList.Error(1)
}

type fieldInfo struct {
	shouldSkip bool
	field      reflect.StructField
	fieldName  string
	tags       map[string]string
}

// Cache for parsed field info by reflect.Type
var (
	fieldInfoCache = make(map[reflect.Type][]fieldInfo)
	fieldInfoMutex = sync.Mutex{}
)

func parseFieldInfo(field reflect.StructField) fieldInfo {
	mcpTag := field.Tag.Get("mcp")
	if mcpTag == "-" { // Skip this field
		return fieldInfo{
			shouldSkip: true,
		}
	}

	jsonTag := field.Tag.Get("json")
	fieldName := strings.Split(jsonTag, ",")[0]
	if fieldName == "" {
		fieldName = field.Name
	}

	tags := parseMCPTags(mcpTag)

	return fieldInfo{
		field:     field,
		fieldName: fieldName,
		tags:      tags,
	}
}

// getFieldInfos returns cached field info for a type or computes and caches it
func getFieldInfos(paramsType reflect.Type) []fieldInfo {
	fieldInfoMutex.Lock()
	defer fieldInfoMutex.Unlock()

	if cached, exists := fieldInfoCache[paramsType]; exists {
		return cached
	}

	var fieldInfos []fieldInfo
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)
		if field.PkgPath != "" { // Skip unexported fields
			continue
		}
		fieldInfos = append(fieldInfos, parseFieldInfo(field))
	}

	fieldInfoCache[paramsType] = fieldInfos
	return fieldInfos
}

func GenerateMetadata(tool Tool) mcp.Tool {
	paramsType := tool.ParamsType()
	toolOptions := []mcp.ToolOption{mcp.WithDescription(tool.Description())}

	fieldInfos := getFieldInfos(paramsType)
	for _, info := range fieldInfos {
		if info.shouldSkip {
			continue
		}

		switch info.field.Type.Kind() {
		case reflect.String:
			opts := []mcp.PropertyOption{mcp.Description(info.tags["description"])}
			if info.tags["required"] == "true" {
				opts = append(opts, mcp.Required())
			}
			toolOptions = append(toolOptions, mcp.WithString(info.fieldName, opts...))

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64:
			opts := []mcp.PropertyOption{mcp.Description(info.tags["description"])}
			if info.tags["required"] == "true" {
				opts = append(opts, mcp.Required())
			}
			toolOptions = append(toolOptions, mcp.WithNumber(info.fieldName, opts...))

		case reflect.Bool:
			opts := []mcp.PropertyOption{mcp.Description(info.tags["description"])}
			if info.tags["required"] == "true" {
				opts = append(opts, mcp.Required())
			}
			toolOptions = append(toolOptions, mcp.WithBoolean(info.fieldName, opts...))

		case reflect.Map:
			if info.field.Type.Key().Kind() == reflect.String {
				toolOptions = append(toolOptions, mcp.WithObject(info.fieldName,
					mcp.Description(info.tags["description"]),
					mcp.AdditionalProperties(map[string]any{"type": "string"}),
				))
			}

		case reflect.Slice:
			toolOptions = append(toolOptions, mcp.WithArray(info.fieldName,
				mcp.Description(info.tags["description"]),
				mcp.Items(map[string]any{"type": "string"}),
			))
		}
	}

	return mcp.NewTool(tool.Name(), toolOptions...)
}

// Assumption: `description` is the last tag since it can contain commas
func parseMCPTags(tag string) map[string]string {
	tags := make(map[string]string)
	if tag == "" {
		return tags
	}

	beforeDesc := tag
	descIndex := strings.Index(tag, "description=")
	if descIndex != -1 {
		beforeDesc = strings.TrimSuffix(tag[:descIndex], ",")
		desc := tag[descIndex+len("description="):]
		tags["description"] = desc
	}

	parts := strings.Split(beforeDesc, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "=") {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 {
				tags[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		} else {
			// Simple flag like "required"
			tags[part] = "true"
		}
	}

	return tags
}

func parseParams(paramsType reflect.Type, args map[string]any) (any, error) {
	params := reflect.New(paramsType).Interface()
	paramsValue := reflect.ValueOf(params).Elem()

	fieldInfos := getFieldInfos(paramsType)
	for _, info := range fieldInfos {
		if info.shouldSkip {
			continue
		}

		value, exists := args[info.fieldName]
		if !exists {
			if info.tags["required"] == "true" {
				return nil, errors.Newf("%s is required", info.fieldName)
			}
			continue
		}

		fieldValue := paramsValue.FieldByName(info.field.Name)
		if err := setFieldValue(fieldValue, value); err != nil {
			return nil, errors.Wrapf(err, "error setting %s", info.fieldName)
		}
	}

	return params, nil
}

func toStringValue(v any) string {
	return fmt.Sprintf("%v", v)
}

func populateMap(newMap reflect.Value, source any) {
	switch src := source.(type) {
	case map[string]any:
		for k, v := range src {
			newMap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(toStringValue(v)))
		}
	case []map[string]any:
		for _, m := range src {
			populateMap(newMap, m)
		}
	case []any:
		for _, item := range src {
			if m, ok := item.(map[string]any); ok {
				populateMap(newMap, m)
			}
		}
	default:
		strValue := toStringValue(source)
		newMap.SetMapIndex(reflect.ValueOf(strValue), reflect.ValueOf(strValue))
	}
}

func populateSlice(newSlice reflect.Value, source any) reflect.Value {
	switch src := source.(type) {
	case []string:
		for _, item := range src {
			newSlice = reflect.Append(newSlice, reflect.ValueOf(item))
		}
	default:
		valueType := reflect.TypeOf(source)
		if valueType.Kind() == reflect.Slice {
			sourceSlice := reflect.ValueOf(source)
			for i := 0; i < sourceSlice.Len(); i++ {
				item := sourceSlice.Index(i).Interface()
				newSlice = reflect.Append(newSlice, reflect.ValueOf(toStringValue(item)))
			}
		} else {
			newSlice = reflect.Append(newSlice, reflect.ValueOf(toStringValue(source)))
		}
	}
	return newSlice
}

func setFieldValue(field reflect.Value, value any) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(toStringValue(value))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intVal, ok := value.(int); ok {
			field.SetInt(int64(intVal))
		} else if floatVal, ok := value.(float64); ok {
			field.SetInt(int64(floatVal))
		} else if strVal, ok := value.(string); ok {
			if parsedInt, err := strconv.ParseInt(strVal, 10, 64); err == nil {
				field.SetInt(parsedInt)
			} else {
				return errors.Newf("cannot convert string %q to int", strVal)
			}
		} else {
			return errors.Newf("cannot convert %T to int", value)
		}

	case reflect.Float32, reflect.Float64:
		if floatVal, ok := value.(float64); ok {
			field.SetFloat(floatVal)
		} else if intVal, ok := value.(int); ok {
			field.SetFloat(float64(intVal))
		} else if strVal, ok := value.(string); ok {
			if parsedFloat, err := strconv.ParseFloat(strVal, 64); err == nil {
				field.SetFloat(parsedFloat)
			} else {
				return errors.Newf("cannot convert string %q to float", strVal)
			}
		} else {
			return errors.Newf("cannot convert %T to float", value)
		}

	case reflect.Bool:
		if boolVal, ok := value.(bool); ok {
			field.SetBool(boolVal)
		} else if strVal, ok := value.(string); ok {
			if parsedBool, err := strconv.ParseBool(strVal); err == nil {
				field.SetBool(parsedBool)
			} else {
				return errors.Newf("cannot convert string %q to bool", strVal)
			}
		} else {
			return errors.Newf("cannot convert %T to bool", value)
		}

	case reflect.Map:
		if field.Type().Key().Kind() == reflect.String {
			newMap := reflect.MakeMap(field.Type())
			populateMap(newMap, value)
			field.Set(newMap)
		}

	case reflect.Slice:
		newSlice := reflect.MakeSlice(field.Type(), 0, 0)
		newSlice = populateSlice(newSlice, value)
		field.Set(newSlice)

	default:
		return errors.Newf("unsupported field type: %v", field.Kind())
	}

	return nil
}

func CreateHandler(tool Tool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, err := parseParams(tool.ParamsType(), request.Params.Arguments)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse params")
		}
		result, err := tool.Execute(ctx, params)
		if err != nil {
			return nil, errors.Wrap(err, "failed to execute tool")
		}
		return result, nil
	}
}
