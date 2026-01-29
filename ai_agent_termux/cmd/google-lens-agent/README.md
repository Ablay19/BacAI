# Google Lens Agent Microservice

**Status:** ⚠️ Needs Configuration Update

This is a standalone HTTP microservice for the Google Lens features. It is not required for the main `ai_agent` CLI tool.

## Current Issue

The agent was designed with a custom configuration structure that differs from the main `config.Config`. It requires:
- `ServerPort`, `MaxConcurrentTasks`, `OllamaBaseURL`, `OllamaModel`, `NtfyServer`, `GoogleCredentialsPath`, `SerpAPIKey`

These fields are not part of the main configuration and need to be added or the agent needs to be refactored.

## Building

To exclude this from the main build, use:
```bash
go build .
```

To build the standalone agent (after fixing configuration):
```bash
go build ./cmd/google-lens-agent
```

## TODO

- [ ] Add missing configuration fields to main config OR
- [ ] Create separate configuration file for this microservice OR
- [ ] Refactor to use existing config fields from the main Config struct
