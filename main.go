package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
	"google.golang.org/genai"
)

// Define schemas/types for wrapping our tools to satisfy functiontool.New expectations.

type CreateSandboxInput struct {
	Language string `json:"language" jsonschema:"The language runtime for the sandbox (e.g. 'typescript' or 'python')."`
}

type CreateSandboxOutput struct {
	SandboxID string `json:"sandbox_id" jsonschema:"The unique ID of the created sandbox."`
}

type CloneRepoInput struct {
	SandboxID string `json:"sandbox_id" jsonschema:"The ID of the sandbox to clone the repo into."`
	RepoURL   string `json:"repo_url" jsonschema:"The Git clone URL of the repository."`
	Path      string `json:"path" jsonschema:"The path inside the sandbox to clone the repository to."`
}

type CloneRepoOutput struct {
	Result string `json:"result" jsonschema:"Output status of the git clone."`
}

type RunSandboxCommandInput struct {
	SandboxID string `json:"sandbox_id" jsonschema:"The ID of the sandbox where the command should run."`
	Command   string `json:"command" jsonschema:"The shell command to execute."`
}

type RunSandboxCommandOutput struct {
	Output string `json:"output" jsonschema:"The command output stdout and stderr."`
}

type WriteSandboxFileInput struct {
	SandboxID string `json:"sandbox_id" jsonschema:"The ID of the sandbox."`
	Path      string `json:"path" jsonschema:"The target file path inside the sandbox."`
	Content   string `json:"content" jsonschema:"The text content to write."`
}

type WriteSandboxFileOutput struct {
	Result string `json:"result" jsonschema:"Success confirmation message."`
}

type ReadSandboxFileInput struct {
	SandboxID string `json:"sandbox_id" jsonschema:"The ID of the sandbox."`
	Path      string `json:"path" jsonschema:"The file path inside the sandbox to read."`
}

type ReadSandboxFileOutput struct {
	Content string `json:"content" jsonschema:"The read text content of the file."`
}

type DeleteSandboxInput struct {
	SandboxID string `json:"sandbox_id" jsonschema:"The ID of the sandbox to delete."`
}

type DeleteSandboxOutput struct {
	Result string `json:"result" jsonschema:"Success confirmation message."`
}

type ScrapeWebPageInput struct {
	URL string `json:"url" jsonschema:"The URL of the webpage to extract content from."`
}

type ScrapeWebPageOutput struct {
	Content string `json:"content" jsonschema:"The webpage content."`
}

type SearchWebInput struct {
	Query string `json:"query" jsonschema:"The search query string."`
}

type SearchWebOutput struct {
	Result string `json:"result" jsonschema:"Search query results."`
}

const AgentInstruction = `You are a Red Team Agent Demo orchestrator.
Your goal is to perform a security and vulnerability analysis on a target repository.

You must follow this exact step-by-step workflow:
1. Create a Daytona sandbox environment using 'CreateSandbox' with language 'typescript'.
2. Clone the target repository (https://github.com/OpenCut-app/OpenCut.git) into the sandbox using 'CloneRepo'. Use '/home/daytona/workspace' as the path.
3. Run security and vulnerability tests inside the sandbox using 'RunSandboxCommand'. For this demo, run 'npm install' then a command like 'npm audit' or a mock security scan.
4. Compile a markdown report summarizing the security findings. Suggest changes/remediations in this markdown rather than pushing any commits or PRs.
5. Write the compiled markdown report inside the sandbox to a file named 'vulnerability_report.md' using 'WriteSandboxFile'.
6. Read the report back using 'ReadSandboxFile' to verify it.
7. Print the full markdown report in your final response.
8. Clean up the sandbox using 'DeleteSandbox'.

Ensure all analysis and commands are executed solely inside the created sandbox.`

func main() {
	// Load environment variables
	_ = godotenv.Load()

	doublewordAPIKey := os.Getenv("DOUBLEWORD_API_KEY")
	doublewordModelName := os.Getenv("DOUBLEWORD_MODEL")
	if doublewordAPIKey == "" || doublewordModelName == "" {
		log.Fatal("DOUBLEWORD_API_KEY and DOUBLEWORD_MODEL environment variables are required")
	}

	oxylabsUsername := os.Getenv("OXYLABS_USERNAME")
	oxylabsPassword := os.Getenv("OXYLABS_PASSWORD")
	if oxylabsUsername == "" || oxylabsPassword == "" {
		log.Fatal("OXYLABS_USERNAME and OXYLABS_PASSWORD environment variables are required")
	}

	// 1. Initialize custom Doubleword model implementation
	doublewordModel := NewDoublewordModel(doublewordModelName, doublewordAPIKey)

	// 2. Initialize Daytona tools wrapper
	daytonaTools, err := NewDaytonaToolSet()
	if err != nil {
		log.Fatalf("Failed to initialize Daytona tools: %v", err)
	}

	// 3. Initialize Oxylabs tools wrapper
	oxylabsTools := NewOxylabsToolSet(oxylabsUsername, oxylabsPassword)

	// 4. Wrap handlers into ADK Tool structs using functiontool.New
	createSandboxTool, err := functiontool.New(
		functiontool.Config{
			Name:        "CreateSandbox",
			Description: "Creates an isolated Daytona sandbox. Language can be 'typescript' or 'python'. Returns sandbox_id.",
		},
		func(ctx agent.Context, input CreateSandboxInput) (CreateSandboxOutput, error) {
			id, err := daytonaTools.CreateSandbox(ctx, input.Language)
			return CreateSandboxOutput{SandboxID: id}, err
		},
	)
	if err != nil {
		log.Fatalf("Failed to build CreateSandbox tool: %v", err)
	}

	cloneRepoTool, err := functiontool.New(
		functiontool.Config{
			Name:        "CloneRepo",
			Description: "Clones a Git repository into the sandbox at the specified path. Returns status message.",
		},
		func(ctx agent.Context, input CloneRepoInput) (CloneRepoOutput, error) {
			res, err := daytonaTools.CloneRepo(ctx, input.SandboxID, input.RepoURL, input.Path)
			return CloneRepoOutput{Result: res}, err
		},
	)
	if err != nil {
		log.Fatalf("Failed to build CloneRepo tool: %v", err)
	}

	runCommandTool, err := functiontool.New(
		functiontool.Config{
			Name:        "RunSandboxCommand",
			Description: "Executes a shell command inside the sandbox. Returns command output stdout and stderr.",
		},
		func(ctx agent.Context, input RunSandboxCommandInput) (RunSandboxCommandOutput, error) {
			out, err := daytonaTools.RunSandboxCommand(ctx, input.SandboxID, input.Command)
			return RunSandboxCommandOutput{Output: out}, err
		},
	)
	if err != nil {
		log.Fatalf("Failed to build RunSandboxCommand tool: %v", err)
	}

	writeFileTool, err := functiontool.New(
		functiontool.Config{
			Name:        "WriteSandboxFile",
			Description: "Writes text content to a file path in the sandbox. Returns success message.",
		},
		func(ctx agent.Context, input WriteSandboxFileInput) (WriteSandboxFileOutput, error) {
			res, err := daytonaTools.WriteSandboxFile(ctx, input.SandboxID, input.Path, input.Content)
			return WriteSandboxFileOutput{Result: res}, err
		},
	)
	if err != nil {
		log.Fatalf("Failed to build WriteSandboxFile tool: %v", err)
	}

	readFileTool, err := functiontool.New(
		functiontool.Config{
			Name:        "ReadSandboxFile",
			Description: "Reads text content from a file path in the sandbox. Returns file content.",
		},
		func(ctx agent.Context, input ReadSandboxFileInput) (ReadSandboxFileOutput, error) {
			content, err := daytonaTools.ReadSandboxFile(ctx, input.SandboxID, input.Path)
			return ReadSandboxFileOutput{Content: content}, err
		},
	)
	if err != nil {
		log.Fatalf("Failed to build ReadSandboxFile tool: %v", err)
	}

	deleteSandboxTool, err := functiontool.New(
		functiontool.Config{
			Name:        "DeleteSandbox",
			Description: "Deletes the sandbox with the given sandbox_id and cleans up resources.",
		},
		func(ctx agent.Context, input DeleteSandboxInput) (DeleteSandboxOutput, error) {
			res, err := daytonaTools.DeleteSandbox(ctx, input.SandboxID)
			return DeleteSandboxOutput{Result: res}, err
		},
	)
	if err != nil {
		log.Fatalf("Failed to build DeleteSandbox tool: %v", err)
	}

	scrapeWebPageTool, err := functiontool.New(
		functiontool.Config{
			Name:        "ScrapeWebPage",
			Description: "Scrapes a web page content using Oxylabs Universal Scraper. Returns the raw extracted HTML/JSON text.",
		},
		func(ctx agent.Context, input ScrapeWebPageInput) (ScrapeWebPageOutput, error) {
			content, err := oxylabsTools.ScrapeWebPage(ctx, input.URL)
			return ScrapeWebPageOutput{Content: content}, err
		},
	)
	if err != nil {
		log.Fatalf("Failed to build ScrapeWebPage tool: %v", err)
	}

	searchWebTool, err := functiontool.New(
		functiontool.Config{
			Name:        "SearchWeb",
			Description: "Performs a Google search using the Oxylabs Search API. Returns search results.",
		},
		func(ctx agent.Context, input SearchWebInput) (SearchWebOutput, error) {
			res, err := oxylabsTools.SearchWeb(ctx, input.Query)
			return SearchWebOutput{Result: res}, err
		},
	)
	if err != nil {
		log.Fatalf("Failed to build SearchWeb tool: %v", err)
	}

	// Group tools list
	tools := []tool.Tool{
		createSandboxTool,
		cloneRepoTool,
		runCommandTool,
		writeFileTool,
		readFileTool,
		deleteSandboxTool,
		scrapeWebPageTool,
		searchWebTool,
	}

	// 5. Initialize the LLMAgent
	agentInstance, err := llmagent.New(llmagent.Config{
		Name:        "red_team_agent",
		Description: "An agent that performs security audits inside Daytona sandboxes.",
		Instruction: AgentInstruction,
		Model:       doublewordModel,
		Tools:       tools,
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// 6. Create ADK Runner
	r, err := runner.New(runner.Config{
		AppName:           "red_team_agent_demo",
		Agent:             agentInstance,
		SessionService:    session.InMemoryService(),
		AutoCreateSession: true,
	})
	if err != nil {
		log.Fatalf("Failed to create runner: %v", err)
	}

	ctx := context.Background()
	promptContent := genai.NewContentFromText(
		"Perform a security and vulnerability analysis on the repository https://github.com/OpenCut-app/OpenCut.git using a Daytona sandbox.",
		genai.RoleUser,
	)

	fmt.Println("Starting Red Team Agent Demo...")
	fmt.Println("Doubleword Model:", doublewordModelName)
	fmt.Println("Triggering analysis request workflow...")
	fmt.Println("--------------------------------------------------")

	// 7. Execute the Runner
	events := r.Run(ctx, "user-1", "session-1", promptContent, agent.RunConfig{
		StreamingMode: agent.StreamingModeNone,
	})

	for ev, err := range events {
		if err != nil {
			log.Fatalf("Error during agent execution: %v", err)
		}

		if ev.Content != nil {
			for _, p := range ev.Content.Parts {
				if p.Text != "" {
					fmt.Print(p.Text)
				}
			}
		}
	}

	fmt.Println("\n--------------------------------------------------")
	fmt.Println("Agent run completed successfully.")
}
