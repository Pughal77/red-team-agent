package main

import (
	"context"
	"fmt"

	"github.com/daytonaio/daytona/libs/sdk-go/pkg/daytona"
	"github.com/daytonaio/daytona/libs/sdk-go/pkg/types"
)

// DaytonaToolSet wraps Daytona SDK client capabilities.
type DaytonaToolSet struct {
	client *daytona.Client
}

// NewDaytonaToolSet initializes a new DaytonaToolSet.
func NewDaytonaToolSet() (*DaytonaToolSet, error) {
	client, err := daytona.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create daytona client: %w", err)
	}
	return &DaytonaToolSet{client: client}, nil
}

// CreateSandbox creates an isolated Daytona sandbox.
func (d *DaytonaToolSet) CreateSandbox(ctx context.Context, language string) (string, error) {
	langCode := types.CodeLanguage(language)
	if langCode == "" {
		langCode = types.CodeLanguageTypeScript // Default to typescript as requested
	}

	params := types.SnapshotParams{
		SandboxBaseParams: types.SandboxBaseParams{
			Language: langCode,
		},
	}

	sandbox, err := d.client.Create(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to create sandbox: %w", err)
	}

	return sandbox.ID, nil
}

// CloneRepo clones a git repository into the specified sandbox path.
func (d *DaytonaToolSet) CloneRepo(ctx context.Context, sandboxId string, repoUrl string, path string) (string, error) {
	sandbox, err := d.client.Get(ctx, sandboxId)
	if err != nil {
		return "", fmt.Errorf("failed to fetch sandbox %s: %w", sandboxId, err)
	}

	if path == "" {
		path = "/home/daytona/workspace"
	}

	err = sandbox.Git.Clone(ctx, repoUrl, path)
	if err != nil {
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	return fmt.Sprintf("Successfully cloned %s to %s", repoUrl, path), nil
}

// RunSandboxCommand executes a shell command inside the sandbox.
func (d *DaytonaToolSet) RunSandboxCommand(ctx context.Context, sandboxId string, command string) (string, error) {
	sandbox, err := d.client.Get(ctx, sandboxId)
	if err != nil {
		return "", fmt.Errorf("failed to fetch sandbox %s: %w", sandboxId, err)
	}

	resp, err := sandbox.Process.ExecuteCommand(ctx, command)
	if err != nil {
		return "", fmt.Errorf("failed to execute command in sandbox: %w", err)
	}

	return resp.Result, nil
}

// WriteSandboxFile writes text content to a remote path in the sandbox.
func (d *DaytonaToolSet) WriteSandboxFile(ctx context.Context, sandboxId string, path string, content string) (string, error) {
	sandbox, err := d.client.Get(ctx, sandboxId)
	if err != nil {
		return "", fmt.Errorf("failed to fetch sandbox %s: %w", sandboxId, err)
	}

	err = sandbox.FileSystem.UploadFile(ctx, []byte(content), path)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return fmt.Sprintf("Successfully wrote file to %s", path), nil
}

// ReadSandboxFile reads text content from a path in the sandbox.
func (d *DaytonaToolSet) ReadSandboxFile(ctx context.Context, sandboxId string, path string) (string, error) {
	sandbox, err := d.client.Get(ctx, sandboxId)
	if err != nil {
		return "", fmt.Errorf("failed to fetch sandbox %s: %w", sandboxId, err)
	}

	data, err := sandbox.FileSystem.DownloadFile(ctx, path, nil)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}

	return string(data), nil
}

// DeleteSandbox terminates and cleans up the sandbox.
func (d *DaytonaToolSet) DeleteSandbox(ctx context.Context, sandboxId string) (string, error) {
	sandbox, err := d.client.Get(ctx, sandboxId)
	if err != nil {
		return "", fmt.Errorf("failed to fetch sandbox %s: %w", sandboxId, err)
	}
	err = sandbox.Delete(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to delete sandbox %s: %w", sandboxId, err)
	}
	return fmt.Sprintf("Sandbox %s deleted successfully", sandboxId), nil
}
