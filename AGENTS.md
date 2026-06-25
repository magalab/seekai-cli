# Repository Guidelines

## Project Structure & Module Organization

`seekai` is a Go CLI. The entry point is `main.go`, with Kong command definitions in `cmd/`. Database access and SQL wrappers live in `client/`. Configuration loading is in `config/`, terminal and JSON/table output helpers are in `display/`, and provider defaults are in `provider/`. Release and install helpers are `scripts/build-release.sh`, `install.sh`, and `install.ps1`. GitHub Actions workflows are under `.github/workflows/`. Tests should sit next to the package they cover, for example `provider/models_test.go`.

## Build, Test, and Development Commands

- `make build`: builds the local `seekai` binary with version metadata.
- `make test`: runs `go test ./...`.
- `make fmt`: formats all Go files with `gofmt`.
- `make fmt-check`: fails if any Go file is not formatted.
- `make all`: runs format check, tests, and build.
- `VERSION=v0.1.0 scripts/build-release.sh`: builds a release-style binary into `dist/`.

For one-off local runs, use `go run . --help` or `go run . model list`.

## Coding Style & Naming Conventions

Use standard Go formatting via `gofmt`; do not manually align code. Follow `.editorconfig`: tabs for Go files, two-space indentation for Markdown/YAML/JSON/TOML. Keep package names short and lower-case. Command structs should stay in `cmd/`, while database behavior belongs in `client/`. Prefer explicit SQL constants or small helper methods over duplicating query strings across commands.

## Testing Guidelines

Use Go’s standard `testing` package. Name tests `TestXxx` and keep table-driven tests close to the behavior they protect. Add tests for deterministic logic such as provider URL mappings, config parsing, output normalization, and command validation. Live Seek DB behavior may be documented as manual verification when it cannot run in CI.

## Commit & Pull Request Guidelines

This repository currently has no established commit history. Use short, imperative commit messages such as `add release workflow` or `fix endpoint list query`. Pull requests should include a clear summary, verification commands, and any database/manual test notes. Link related issues when available. Update README files when commands, flags, install paths, release assets, or provider defaults change.

## Security & Configuration Tips

Do not commit real passwords, API keys, or local profile files. Prefer environment expansion in config, for example `password = "${SEEKAI_PASSWORD}"`. GitHub release automation requires repository secrets such as `HOMEBREW_TAP_TOKEN`; document required permissions when changing workflows.
