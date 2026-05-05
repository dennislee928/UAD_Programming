# Dev Container Configuration

This directory contains the configuration for VS Code Dev Containers and GitHub Codespaces.

## Features

- **Base Image**: Go 1.21 official dev container
- **Pre-installed Tools**:
  - Go 1.21+
  - Git
  - GitHub CLI
  - golangci-lint (installed via post-create)
- **VS Code Extensions**:
  - Go language support
  - Makefile tools
  - GitHub Copilot
  - GitLens
  - Markdown support

## Usage

### GitHub Codespaces

1. Go to the repository on GitHub
2. Click "Code" → "Codespaces" → "Create codespace on dev"
3. Wait for the container to build and initialize
4. Run `make build` to compile
5. Run `make test` to run tests

### VS Code Dev Containers

1. Install "Dev Containers" extension in VS Code
2. Open the repository
3. Press F1 and select "Dev Containers: Reopen in Container"
4. Wait for the container to build
5. The workspace will open inside the container

## Post-Create Actions

The container automatically runs:
```bash
make deps  # Download and tidy Go dependencies
make build # Build all binaries
```

## Port Forwarding

- Port 6060: Go documentation server (`godoc -http=:6060`)

## Build Artifacts

The `bin/` directory is mounted as a bind mount for better performance when building binaries.

## Customization

Edit `devcontainer.json` to:
- Add more VS Code extensions
- Change Go version
- Add additional tools
- Modify environment variables
- Configure post-create commands


