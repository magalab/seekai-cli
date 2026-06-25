# seekai

`seekai` is a Go CLI for Seek DB AI model management, endpoint management, AI function calls, and raw SQL execution.

## Install

Install with Homebrew on macOS:

```sh
brew tap magalab/tap
brew install seekai
```

Install the latest macOS/Linux release:

```sh
curl -fsSL https://raw.githubusercontent.com/magalab/seekai-cli/main/install.sh | sh
```

Install the latest Windows release from PowerShell:

```powershell
iwr https://raw.githubusercontent.com/magalab/seekai-cli/main/install.ps1 -useb | iex
```

Install a specific version:

```sh
curl -fsSL https://raw.githubusercontent.com/magalab/seekai-cli/main/install.sh | SEEKAI_VERSION=v0.1.0 sh
```

Or install from source:

```sh
go install github.com/magalab/seekai-cli@latest
```

Build a local binary:

```sh
go build -o seekai .
```

## Connection

Every command accepts the global connection flags:

```sh
seekai --host localhost --port 2881 --user root --password '' --database test model list
```

Profiles can be stored at `~/.seekai/config.toml`:

```toml
[default]
host = "localhost"
port = 2881
user = "root"
password = ""
database = "test"

[profiles.production]
host = "192.168.1.100"
port = 2881
user = "admin"
password = "${SEEKAI_PASSWORD}"
database = "test"
```

Use a profile with:

```sh
seekai --profile production model list
```

## Models

```sh
seekai model list
seekai model create ob_embed --type dense_embedding --model-name BAAI/bge-m3
seekai model delete ob_embed
```

`model create` can also create an endpoint when provider data is supplied:

```sh
seekai model create ob_embed \
  --type dense_embedding \
  --model-name BAAI/bge-m3 \
  --provider siliconflow \
  --url https://api.siliconflow.cn/v1/embeddings \
  --access-key "$SILICONFLOW_API_KEY"
```

If required arguments are omitted, `seekai model create` opens an interactive form.

Known provider keys:

| Provider | completion | dense_embedding | rerank |
| --- | --- | --- | --- |
| `aliyun-openAI` | `https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions` | `https://dashscope.aliyuncs.com/compatible-mode/v1/embeddings` | - |
| `aliyun-dashscope` | `https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation` | `https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding` | `https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank` |
| `deepseek` | `https://api.deepseek.com/chat/completions` | - | - |
| `siliconflow` | `https://api.siliconflow.cn/v1/chat/completions` | `https://api.siliconflow.cn/v1/embeddings` | `https://api.siliconflow.cn/v1/rerank` |
| `hunyuan-openAI` | `https://api.hunyuan.cloud.tencent.com/v1/chat/completions` | `https://api.hunyuan.cloud.tencent.com/v1/embeddings` | - |
| `openAI` | `https://api.openai.com/v1/chat/completions` | `https://api.openai.com/v1/embeddings` | - |

## Endpoints

```sh
seekai endpoint list
seekai endpoint create ob_embed_endpoint --model ob_embed --provider siliconflow --url https://api.siliconflow.cn/v1/embeddings --access-key "$SILICONFLOW_API_KEY"
seekai endpoint update ob_embed_endpoint --url https://api.siliconflow.cn/v1/embeddings
seekai endpoint delete ob_embed_endpoint
```

## AI Functions

```sh
seekai ai complete ob_complete "Translate to Chinese: Hello world"
seekai ai complete ob_complete "Long prompt" --pipe
seekai ai complete ob_complete -o json
seekai ai embed ob_embed "Hello world"
cat texts.txt | seekai ai embed ob_embed --stdin -o json
seekai ai rerank ob_rerank "Apple" '["apple","banana","fruit"]'
seekai ai prompt "Summarize: {0}" "Seek DB supports AI functions"
```

Long completion and prompt output opens a terminal pager when stdout is interactive. Piped output stays plain text.

## Raw SQL

```sh
seekai sql 'SELECT * FROM oceanbase.DBA_OB_AI_MODELS'
seekai -o json sql 'SELECT * FROM oceanbase.DBA_OB_AI_MODEL_ENDPOINTS'
```

Quote SQL according to your shell when it contains spaces or special characters.

## Output

Use `-o, --output` with `auto`, `table`, `json`, `yaml`, `toml`, or `text`.

List and rerank commands default to table output. Creation, deletion, completion, embedding, prompt, and raw values default to text unless the command has tabular rows.

Examples:

```sh
seekai -o json model list
seekai -o yaml endpoint list
seekai -o toml ai complete ob_complete "Hello" --pipe
```

TOML output wraps lists under `rows` and scalar values under `value`, because TOML requires a document with keys.

## Release

Pushing a tag named `vX.Y.Z` runs the release workflow and publishes binaries for:

- `seekai_darwin_amd64`
- `seekai_darwin_arm64`
- `seekai_linux_amd64`
- `seekai_linux_arm64`
- `seekai_windows_amd64.exe`

On tag releases, the workflow also updates the macOS Homebrew formula in `magalab/tap`. Configure a repository secret named `HOMEBREW_TAP_TOKEN` with write access to `magalab/tap`.
