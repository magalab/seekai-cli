# Contributing

## Development

Requirements:

- Go 1.22 or newer
- Access to a Seek DB instance for integration testing

Common commands:

```sh
make fmt
make test
make build
```

Before opening a pull request, run:

```sh
make all
```

## Pull Requests

- Keep changes focused on one behavior or maintenance topic.
- Include documentation updates when commands, flags, output, or configuration change.
- Add tests for logic that can be verified without a live database.
- For database behavior, include the exact SQL or manual verification steps in the PR description.

## Release

Create a tag named `vX.Y.Z` and push it. GitHub Actions will build release artifacts for:

- `seekai_darwin_amd64`
- `seekai_darwin_arm64`
- `seekai_linux_amd64`
- `seekai_linux_arm64`
- `seekai_windows_amd64.exe`

Tag releases also update the macOS Homebrew formula in `magalab/homebrew-tap`, which users install as `brew tap magalab/tap`. The repository must have a `HOMEBREW_TAP_TOKEN` secret with write access to that tap repository.
