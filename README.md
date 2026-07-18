# Red Team Agent Demo with Google ADK in Golang

This project is a code-first demonstration of an AI Agent built using the **Google Agent Development Kit (ADK) v2** in Go. It integrates with **Daytona** sandboxes for secure, isolated code execution, **Oxylabs** for web scraping and search capabilities, and **Doubleword** as the LLM model provider.

## Workflow Overview

1. **Doubleword LLM Integration**: Uses a custom `model.LLM` interface implementation (`DoublewordModel`) in `doubleword_model.go` to communicate with the Doubleword API endpoint (`https://api.doubleword.ai/v1/chat/completions`) using OpenAI-compatible payload schemas.
2. **Daytona Sandboxes**: Spins up an isolated sandbox workspace where the repository is cloned. Any commands or security scripts are run safely inside this sandbox.
3. **Oxylabs Scraping & Searching**: The agent is equipped with web-page scraping (Universal API) and search tools (Search API) via Oxylabs.
4. **Vulnerability Report**: Code execution results are compiled into a markdown report inside the sandbox rather than making commits or pull requests.

## Prerequisites

- **Go 1.25+** (Go toolchain automatically manages this in `go.mod`).
- Daytona Account & API Key.
- Oxylabs Account & Username/Password.
- Doubleword Account & API Key.

## Setup & Configuration

1. Copy the `.env` template and set your API keys:
   ```bash
   cp .env.example .env # or copy and edit the created .env
   ```

2. Fill in the values inside `.env`:
   ```bash
   DOUBLEWORD_API_KEY="your_doubleword_api_key"
   DOUBLEWORD_MODEL="deepseek-ai/DeepSeek-V4-Flash"
   DAYTONA_API_KEY="your_daytona_api_key"
   DAYTONA_SERVER_URL="https://app.daytona.io/api"
   OXYLABS_USERNAME="your_oxylabs_username"
   OXYLABS_PASSWORD="your_oxylabs_password"
   ```

## Running the Demo

To run the agent orchestration:

```bash
go run .
```

This will trigger the agent to:
- Spin up a TypeScript Daytona sandbox.
- Clone the target codebase (`https://github.com/OpenCut-app/OpenCut.git`).
- Perform the specified analysis steps inside the sandbox.
- Write the report `vulnerability_report.md` inside the sandbox.
- Read it back to print on stdout, then clean up by deleting the sandbox.
