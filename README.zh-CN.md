# seekai

`seekai` 是一个 Go CLI 工具，用于 Seek DB AI 模型管理、端点管理、AI 函数调用和原始 SQL 执行。

## 安装

在 macOS 上通过 Homebrew 安装：

```sh
brew tap magalab/tap
brew install seekai
```

安装最新 macOS/Linux 版本：

```sh
curl -fsSL https://raw.githubusercontent.com/magalab/seekai-cli/main/install.sh | sh
```

在 Windows PowerShell 中安装最新版本：

```powershell
iwr https://raw.githubusercontent.com/magalab/seekai-cli/main/install.ps1 -useb | iex
```

安装指定版本：

```sh
curl -fsSL https://raw.githubusercontent.com/magalab/seekai-cli/main/install.sh | SEEKAI_VERSION=v0.1.0 sh
```

也可以从源码安装：

```sh
go install github.com/magalab/seekai-cli@latest
```

构建本地二进制文件：

```sh
go build -o seekai .
```

## 连接

每条命令都支持全局连接参数：

```sh
seekai --host localhost --port 2881 --user root --password '' --database test model list
```

配置文件可存放在 `~/.seekai/config.toml`：

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

使用配置文件的连接方式：

```sh
seekai --profile production model list
```

## 模型

```sh
seekai model list
seekai model create ob_embed --type dense_embedding --model-name BAAI/bge-m3
seekai model delete ob_embed
```

`model create` 在提供供应商数据时也可以创建端点：

```sh
seekai model create ob_embed \
  --type dense_embedding \
  --model-name BAAI/bge-m3 \
  --provider siliconflow \
  --url https://api.siliconflow.cn/v1/embeddings \
  --access-key "$SILICONFLOW_API_KEY"
```

如果省略必需参数，`seekai model create` 会打开交互式表单。

已知 provider key：

| Provider | completion | dense_embedding | rerank |
| --- | --- | --- | --- |
| `aliyun-openAI` | `https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions` | `https://dashscope.aliyuncs.com/compatible-mode/v1/embeddings` | - |
| `aliyun-dashscope` | `https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation` | `https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding` | `https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank` |
| `deepseek` | `https://api.deepseek.com/chat/completions` | - | - |
| `siliconflow` | `https://api.siliconflow.cn/v1/chat/completions` | `https://api.siliconflow.cn/v1/embeddings` | `https://api.siliconflow.cn/v1/rerank` |
| `hunyuan-openAI` | `https://api.hunyuan.cloud.tencent.com/v1/chat/completions` | `https://api.hunyuan.cloud.tencent.com/v1/embeddings` | - |
| `openAI` | `https://api.openai.com/v1/chat/completions` | `https://api.openai.com/v1/embeddings` | - |

## 端点

```sh
seekai endpoint list
seekai endpoint create ob_embed_endpoint --model ob_embed --provider siliconflow --url https://api.siliconflow.cn/v1/embeddings --access-key "$SILICONFLOW_API_KEY"
seekai endpoint update ob_embed_endpoint --url https://api.siliconflow.cn/v1/embeddings
seekai endpoint delete ob_embed_endpoint
```

## AI 函数

```sh
seekai ai complete ob_complete "Translate to Chinese: Hello world"
seekai ai complete ob_complete "长文本提示" --pipe
seekai ai complete ob_complete -o json
seekai ai embed ob_embed "Hello world"
cat texts.txt | seekai ai embed ob_embed --stdin -o json
seekai ai rerank ob_rerank "Apple" '["apple","banana","fruit"]'
seekai ai prompt "Summarize: {0}" "Seek DB supports AI functions"
```

当标准输出为交互式终端时，长文本的补全和提示输出会打开终端分页器。管道输出保持纯文本。

## 原始 SQL

```sh
seekai sql 'SELECT * FROM oceanbase.DBA_OB_AI_MODELS'
seekai -o json sql 'SELECT * FROM oceanbase.DBA_OB_AI_MODEL_ENDPOINTS'
```

当 SQL 包含空格或特殊字符时，请根据你的 shell 规则进行引号包裹。

## 输出格式

使用 `-o, --output` 参数，可选值：`auto`、`table`、`json`、`yaml`、`toml` 或 `text`。

列表和重排序命令默认使用表格输出。创建、删除、补全、嵌入、提示和原始值命令默认使用文本输出，除非命令包含表格形式的数据行。

示例：

```sh
seekai -o json model list
seekai -o yaml endpoint list
seekai -o toml ai complete ob_complete "Hello" --pipe
```

TOML 输出会把列表包装到 `rows` 字段下，把标量值包装到 `value` 字段下，因为 TOML 文档需要具名字段。

## Shell 补全

使用以下命令生成 shell completion 脚本：

```sh
seekai completion bash
seekai completion zsh
seekai completion fish
```

安装示例：

```sh
seekai completion bash > /usr/local/etc/bash_completion.d/seekai
seekai completion zsh > "${fpath[1]}/_seekai"
seekai completion fish > ~/.config/fish/completions/seekai.fish
```

## 发布

推送 `vX.Y.Z` 格式的 tag 会触发 release workflow，并发布以下平台的二进制文件：

- `seekai_darwin_amd64`
- `seekai_darwin_arm64`
- `seekai_linux_amd64`
- `seekai_linux_arm64`
- `seekai_windows_amd64.exe`

tag release 时，workflow 还会更新 `magalab/homebrew-tap` 仓库中的 macOS Homebrew formula，用户侧仍通过 `brew tap magalab/tap` 使用。需要在本仓库配置名为 `HOMEBREW_TAP_TOKEN` 的 secret，并授予它写入 `magalab/homebrew-tap` 的权限。
