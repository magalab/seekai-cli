package cmd

import (
	"fmt"
	"os"
)

type CompletionShell string

const (
	CompletionBash CompletionShell = "bash"
	CompletionZsh  CompletionShell = "zsh"
	CompletionFish CompletionShell = "fish"
)

type CompletionCmd struct {
	Shell CompletionShell `arg:"" enum:"bash,zsh,fish" help:"Shell to generate completion for: bash, zsh, or fish."`
}

func (c *CompletionCmd) Run() error {
	switch c.Shell {
	case CompletionBash:
		_, err := fmt.Fprint(os.Stdout, bashCompletion)
		return err
	case CompletionZsh:
		_, err := fmt.Fprint(os.Stdout, zshCompletion)
		return err
	case CompletionFish:
		_, err := fmt.Fprint(os.Stdout, fishCompletion)
		return err
	default:
		return fmt.Errorf("unsupported shell %q: expected bash, zsh, or fish", c.Shell)
	}
}

const bashCompletion = `# bash completion for seekai
_seekai_completion() {
  local cur prev words cword
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"
  words="${COMP_WORDS[*]}"

  case "${prev}" in
    --output|-o)
      COMPREPLY=( $(compgen -W "auto table json yaml toml text" -- "${cur}") )
      return 0
      ;;
    --type)
      COMPREPLY=( $(compgen -W "completion dense_embedding rerank" -- "${cur}") )
      return 0
      ;;
    --provider)
      COMPREPLY=( $(compgen -W "aliyun-openAI aliyun-dashscope deepseek siliconflow hunyuan-openAI openAI" -- "${cur}") )
      return 0
      ;;
    completion)
      COMPREPLY=( $(compgen -W "bash zsh fish" -- "${cur}") )
      return 0
      ;;
  esac

  case "${cur}" in
    -*)
      COMPREPLY=( $(compgen -W "--host --port --user --password --database --profile --output -o --help" -- "${cur}") )
      return 0
      ;;
  esac

  if [[ "${COMP_CWORD}" -eq 1 ]]; then
    COMPREPLY=( $(compgen -W "model endpoint ai sql completion" -- "${cur}") )
    return 0
  fi

  case "${words}" in
    *" model "*)
      COMPREPLY=( $(compgen -W "list create delete" -- "${cur}") )
      ;;
    *" endpoint "*)
      COMPREPLY=( $(compgen -W "list create update delete" -- "${cur}") )
      ;;
    *" ai "*)
      COMPREPLY=( $(compgen -W "complete embed rerank prompt" -- "${cur}") )
      ;;
    *)
      COMPREPLY=( $(compgen -W "model endpoint ai sql completion" -- "${cur}") )
      ;;
  esac
}

complete -F _seekai_completion seekai
`

const zshCompletion = `#compdef seekai

_seekai() {
  local -a commands model_commands endpoint_commands ai_commands global_options output_formats model_types providers shells

  commands=(
    'model:Manage AI models'
    'endpoint:Manage AI model endpoints'
    'ai:Call Seek DB AI functions'
    'sql:Execute raw SQL'
    'completion:Generate shell completion scripts'
  )
  model_commands=('list:List AI models' 'create:Create an AI model' 'delete:Delete an AI model')
  endpoint_commands=('list:List AI model endpoints' 'create:Create an AI model endpoint' 'update:Update an AI model endpoint' 'delete:Delete an AI model endpoint')
  ai_commands=('complete:Call AI_COMPLETE' 'embed:Call AI_EMBED' 'rerank:Call AI_RERANK' 'prompt:Call AI_PROMPT')
  output_formats=(auto table json yaml toml text)
  model_types=(completion dense_embedding rerank)
  providers=(aliyun-openAI aliyun-dashscope deepseek siliconflow hunyuan-openAI openAI)
  shells=(bash zsh fish)

  global_options=(
    '--host[Seek DB host]:host:'
    '--port[Seek DB MySQL protocol port]:port:'
    '--user[Seek DB user]:user:'
    '--password[Seek DB password]:password:'
    '--database[Default database]:database:'
    '--profile[Profile name from ~/.seekai/config.toml]:profile:'
    '(-o --output)'{-o,--output}'[Output format]:format:->output'
    '--help[Show help]'
  )

  _arguments -C \
    $global_options \
    '1:command:->command' \
    '*::arg:->args'

  case "$state" in
    command)
      _describe 'command' commands
      ;;
    output)
      _describe 'format' output_formats
      ;;
    args)
      case "${words[1]}" in
        model)
          if (( CURRENT == 3 )); then
            _describe 'model command' model_commands
          else
            _arguments \
              '--type[Model type]:type:->model_type' \
              '--model-name[Provider model name]:model:' \
              '--provider[Provider]:provider:->provider' \
              '--url[Endpoint URL]:url:' \
              '--access-key[Endpoint access key]:key:'
          fi
          ;;
        endpoint)
          if (( CURRENT == 3 )); then
            _describe 'endpoint command' endpoint_commands
          else
            _arguments \
              '--model[AI model name]:model:' \
              '--provider[Provider]:provider:->provider' \
              '--url[Endpoint URL]:url:' \
              '--access-key[Endpoint access key]:key:'
          fi
          ;;
        ai)
          if (( CURRENT == 3 )); then
            _describe 'AI command' ai_commands
          else
            _arguments \
              '--parameters[JSON parameters]:json:' \
              '--pipe[Force plain text output]' \
              '--dim[Embedding dimension]:dim:' \
              '--stdin[Read newline-delimited texts from stdin]'
          fi
          ;;
        completion)
          _describe 'shell' shells
          ;;
      esac

      case "$state" in
        model_type) _describe 'type' model_types ;;
        provider) _describe 'provider' providers ;;
      esac
      ;;
  esac
}

_seekai "$@"
`

const fishCompletion = `# fish completion for seekai
complete -c seekai -f

complete -c seekai -l host -d 'Seek DB host' -r
complete -c seekai -l port -d 'Seek DB MySQL protocol port' -r
complete -c seekai -l user -d 'Seek DB user' -r
complete -c seekai -l password -d 'Seek DB password' -r
complete -c seekai -l database -d 'Default database' -r
complete -c seekai -l profile -d 'Profile name from ~/.seekai/config.toml' -r
complete -c seekai -s o -l output -d 'Output format' -xa 'auto table json yaml toml text'

complete -c seekai -n '__fish_use_subcommand' -xa 'model' -d 'Manage AI models'
complete -c seekai -n '__fish_use_subcommand' -xa 'endpoint' -d 'Manage AI model endpoints'
complete -c seekai -n '__fish_use_subcommand' -xa 'ai' -d 'Call Seek DB AI functions'
complete -c seekai -n '__fish_use_subcommand' -xa 'sql' -d 'Execute raw SQL'
complete -c seekai -n '__fish_use_subcommand' -xa 'completion' -d 'Generate shell completion scripts'

complete -c seekai -n '__fish_seen_subcommand_from model; and not __fish_seen_subcommand_from list create delete' -xa 'list create delete'
complete -c seekai -n '__fish_seen_subcommand_from endpoint; and not __fish_seen_subcommand_from list create update delete' -xa 'list create update delete'
complete -c seekai -n '__fish_seen_subcommand_from ai; and not __fish_seen_subcommand_from complete embed rerank prompt' -xa 'complete embed rerank prompt'
complete -c seekai -n '__fish_seen_subcommand_from completion' -xa 'bash zsh fish'

complete -c seekai -n '__fish_seen_subcommand_from create' -l type -d 'Model type' -xa 'completion dense_embedding rerank'
complete -c seekai -n '__fish_seen_subcommand_from create' -l model-name -d 'Provider model name' -r
complete -c seekai -n '__fish_seen_subcommand_from create update' -l url -d 'Endpoint URL' -r
complete -c seekai -n '__fish_seen_subcommand_from create' -l model -d 'AI model name' -r
complete -c seekai -n '__fish_seen_subcommand_from create' -l provider -d 'Provider name' -xa 'aliyun-openAI aliyun-dashscope deepseek siliconflow hunyuan-openAI openAI'
complete -c seekai -n '__fish_seen_subcommand_from create update' -l access-key -d 'Endpoint access key' -r

complete -c seekai -n '__fish_seen_subcommand_from complete' -l parameters -d 'JSON parameters' -r
complete -c seekai -n '__fish_seen_subcommand_from complete' -l pipe -d 'Force plain text output'
complete -c seekai -n '__fish_seen_subcommand_from embed' -l dim -d 'Embedding dimension' -r
complete -c seekai -n '__fish_seen_subcommand_from embed' -l stdin -d 'Read newline-delimited texts from stdin'
`
