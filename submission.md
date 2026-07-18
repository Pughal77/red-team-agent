# Red Team Agent Demo - Capabilities & Roadmap

This project is a code-first demonstration of an AI agent built on the **Google Agent Development Kit (ADK) v2** in Go. It integrates **Daytona sandboxes** for secure file execution and repository isolation, **Oxylabs** for web search and scraping, and **Doubleword** for intelligent model connection.

## Current Project Capabilities

- **Isolated Code Containment**: Creates fully isolated, ephemeral Daytona sandbox environments dynamically to download, configure, and inspect codebases without risking host environment security.
- **Vulnerability Audits**: Safely executes commands (such as `npm install && npm audit`) in the isolated environment to analyze code configurations, security policies, and package dependency risks.
- **Volume Mirroring & Host Sync**: Mounts and mirrors file changes in real-time. When the agent writes the results inside the sandbox container (at `/home/daytona/workspace/vulnerability_report.md`), the Go toolset intercepts the file write and synchronizes a local copy directly to the host's `./volumes` directory.
- **Doubleword Model Wrapper**: Features an OpenAI-compatible custom `model.LLM` adapter in Go that connects ADK’s internal event loop to Doubleword completions (`https://api.doubleword.ai/v1/chat/completions`).
- **Web Search & Scraping**: Equips the agent with Oxylabs web searching and universal web-page scraping capabilities for external contextual lookups.

---

## The Value of Sandboxes for AI Agents

Using dedicated sandboxes is a crucial architectural prerequisite for deploying autonomous AI agents in DevOps, continuous integration, and security operations. By running agent execution loops inside an isolated container, we establish a secure boundary that prevents destructive side effects (such as recursive file deletion, network loops, or credential leaks) from affecting production systems. This containment allows agents to safely perform deep penetration testing, dependency installations, and shell command debugging in a zero-risk sandbox. Consequently, developers can delegate tedious auditing, maintenance, and diagnostics to AI agents, freeing their time to focus on building features.

---

## Future Roadmap

### 1. Version Control System (VCS) Integrations
- **GitHub, GitLab, & Bitbucket Adapters**: Direct integration with VCS platforms to query repositories dynamically, parse webhook triggers, and fetch pull request diffs directly into sandboxes.
- **Automated PR Auditing**: Run security test suites on every incoming pull request and suggest remediations directly in the VCS.

### 2. VCS Issue Generation
- **Automated Issue Tracking**: When a vulnerability is identified, the agent will automatically open a structured issue in the relevant VCS repository (e.g. GitHub Issues, GitLab Issues).
- **Remediation Context**: Issues will include exact line numbers, code snippets, severity ratings, CWE classifications, and suggested drop-in fix diffs.

### 3. Penetration Testing Encapsulation
- **Advanced System Instructions**: Enhance prompt boundaries and system instructions to support active exploit testing (such as scanning for SQL injection vectors, XSS, and broken access controls).
- **Interactive Security Scanners**: Equip agents with tools to run specialized vulnerability scanners (e.g., OWASP ZAP, semgrep, Trivy) in the sandbox and synthesize the output.
