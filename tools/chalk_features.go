package tools

import (
	"context"
	"os"
	"path/filepath"
	"reflect"

	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
)

type ChalkFeaturesParams struct {
	ProjectRepository string `json:"project_repository" mcp:"required,description=Path to the root of the Chalk project on disk to fetch features for. Should contain a chalk.yml or chalk.yaml file."`
}

type ChalkFeaturesTool struct {
	executor CommandExecutor
}

func NewChalkFeaturesTool(executor CommandExecutor) *ChalkFeaturesTool {
	if executor == nil {
		executor = &DefaultCommandExecutor{}
	}
	return &ChalkFeaturesTool{executor: executor}
}

func (t *ChalkFeaturesTool) Name() string {
	return "chalk_features"
}

func (t *ChalkFeaturesTool) Description() string {
	return "Get the list of features from a chalk project"
}

func (t *ChalkFeaturesTool) ParamsType() reflect.Type {
	return reflect.TypeOf(ChalkFeaturesParams{})
}

func (t *ChalkFeaturesTool) Execute(ctx context.Context, args any) (*mcp.CallToolResult, error) {
	params, ok := args.(*ChalkFeaturesParams)
	if !ok {
		return nil, errors.New("invalid parameter type")
	}

	tmpDir, err := os.MkdirTemp("", "chalk-features-*")
	if err != nil {
		return nil, errors.Wrapf(err, "creating temporary directory")
	}
	defer os.RemoveAll(tmpDir)
	featuresFile := filepath.Join(tmpDir, "features.json")

	cmdArgs := []string{"features", "--out", featuresFile}

	output, err := t.executor.Execute(
		ctx,
		params.ProjectRepository,
		cmdArgs...,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "running chalk command; output: %s", output)
	}

	return mcp.NewToolResultText(string(output)), nil
}
