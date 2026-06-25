package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"

	"github.com/magalab/seekai-cli/client"
	"github.com/magalab/seekai-cli/display"
	"github.com/magalab/seekai-cli/provider"
)

type CLI struct {
	Globals
	Model    ModelCmd      `cmd:"" help:"Manage AI models."`
	Endpoint EndpointCmd   `cmd:"" help:"Manage AI model endpoints."`
	AI       AICmd         `cmd:"" name:"ai" help:"Call Seek DB AI functions."`
	SQL      SQLCmd        `cmd:"" name:"sql" help:"Execute raw SQL."`
	Complete CompletionCmd `cmd:"" name:"completion" help:"Generate shell completion scripts."`
}

type ModelCmd struct {
	List   ModelListCmd   `cmd:"" help:"List AI models."`
	Create ModelCreateCmd `cmd:"" help:"Create an AI model."`
	Delete ModelDeleteCmd `cmd:"" help:"Delete an AI model."`
}

type EndpointCmd struct {
	List   EndpointListCmd   `cmd:"" help:"List AI model endpoints."`
	Create EndpointCreateCmd `cmd:"" help:"Create an AI model endpoint."`
	Update EndpointUpdateCmd `cmd:"" help:"Update an AI model endpoint."`
	Delete EndpointDeleteCmd `cmd:"" help:"Delete an AI model endpoint."`
}

type AICmd struct {
	Complete AICompleteCmd `cmd:"" help:"Call AI_COMPLETE."`
	Embed    AIEmbedCmd    `cmd:"" help:"Call AI_EMBED."`
	Rerank   AIRerankCmd   `cmd:"" help:"Call AI_RERANK."`
	Prompt   AIPromptCmd   `cmd:"" help:"Call AI_PROMPT."`
}

type ModelListCmd struct{}

func (c *ModelListCmd) Run(g *Globals) error {
	db, err := open(g)
	if err != nil {
		return err
	}
	defer db.Close()
	result, err := db.ListModels()
	if err != nil {
		return err
	}
	return display.Rows(result.Rows, display.ResolveFormat(string(g.Output), "table"), result.Columns)
}

type ModelCreateCmd struct {
	Name      string `arg:"" optional:"" help:"Model key/name."`
	Type      string `help:"Model type: completion, dense_embedding, or rerank."`
	ModelName string `name:"model-name" help:"Provider model name."`
	Provider  string `help:"Provider for optional endpoint creation."`
	URL       string `help:"Endpoint URL for optional endpoint creation."`
	AccessKey string `name:"access-key" help:"Endpoint access key."`
}

func (c *ModelCreateCmd) Run(g *Globals) error {
	if err := c.fillInteractive(); err != nil {
		return err
	}
	if err := c.validate(); err != nil {
		return err
	}
	db, err := open(g)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.CreateModel(c.Name, c.Type, c.ModelName); err != nil {
		return err
	}
	endpointCreated := false
	if c.endpointRequested() {
		name := c.Name + "_endpoint"
		if err := db.CreateEndpoint(name, client.EndpointConfig{
			AIModelName: c.Name,
			URL:         c.URL,
			Provider:    c.Provider,
			AccessKey:   c.AccessKey,
		}); err != nil {
			return err
		}
		endpointCreated = true
	}

	return display.Value(map[string]any{
		"status":           "created",
		"model":            c.Name,
		"endpoint_created": endpointCreated,
	}, display.ResolveFormat(string(g.Output), "text"))
}

func (c *ModelCreateCmd) fillInteractive() error {
	if c.Name == "" || c.Type == "" || c.ModelName == "" {
		providers := provider.Names()
		if c.Provider == "" && len(providers) > 0 {
			c.Provider = providers[0]
		}
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Model Name").Value(&c.Name),
				huh.NewSelect[string]().Title("Type").Options(
					huh.NewOption("completion", "completion"),
					huh.NewOption("dense_embedding", "dense_embedding"),
					huh.NewOption("rerank", "rerank"),
				).Value(&c.Type),
				huh.NewInput().Title("Provider Model Name").Value(&c.ModelName),
				huh.NewSelect[string]().Title("Provider").Options(providerOptions()...).Value(&c.Provider),
				huh.NewInput().Title("Endpoint URL").Value(&c.URL),
				huh.NewInput().Title("API Key").EchoMode(huh.EchoModePassword).Value(&c.AccessKey),
			),
		)
		if err := form.Run(); err != nil {
			return err
		}
	}
	if c.Name == "" || c.Type == "" || c.ModelName == "" {
		return fmt.Errorf("model name, type, and --model-name are required")
	}
	if c.URL == "" && c.Provider != "" {
		c.URL = provider.DefaultURL(c.Provider, c.Type)
	}
	return nil
}

func (c ModelCreateCmd) validate() error {
	if c.Name == "" || c.Type == "" || c.ModelName == "" {
		return fmt.Errorf("model name, type, and --model-name are required")
	}
	if !validModelType(c.Type) {
		return fmt.Errorf("invalid model type %q: expected completion, dense_embedding, or rerank", c.Type)
	}
	if c.endpointRequested() && c.URL == "" {
		return fmt.Errorf("endpoint URL is required: provider %q has no default URL for model type %q", c.Provider, c.Type)
	}
	return nil
}

func (c ModelCreateCmd) endpointRequested() bool {
	return c.URL != "" || c.Provider != "" || c.AccessKey != ""
}

func validModelType(value string) bool {
	switch value {
	case "completion", "dense_embedding", "rerank":
		return true
	default:
		return false
	}
}

func providerOptions() []huh.Option[string] {
	names := provider.Names()
	options := make([]huh.Option[string], 0, len(names))
	for _, name := range names {
		options = append(options, huh.NewOption(name, name))
	}
	return options
}

type ModelDeleteCmd struct {
	Name string `arg:"" help:"Model key/name."`
}

func (c *ModelDeleteCmd) Run(g *Globals) error {
	db, err := open(g)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.DeleteModel(c.Name); err != nil {
		return err
	}
	return display.Value(map[string]any{"status": "deleted", "model": c.Name}, display.ResolveFormat(string(g.Output), "text"))
}

type EndpointListCmd struct{}

func (c *EndpointListCmd) Run(g *Globals) error {
	db, err := open(g)
	if err != nil {
		return err
	}
	defer db.Close()
	result, err := db.ListEndpoints()
	if err != nil {
		return err
	}
	return display.Rows(result.Rows, display.ResolveFormat(string(g.Output), "table"), result.Columns)
}

type EndpointCreateCmd struct {
	Name      string `arg:"" help:"Endpoint name."`
	Model     string `help:"AI model name." required:""`
	URL       string `help:"Endpoint URL." required:""`
	Provider  string `help:"Provider name." required:""`
	AccessKey string `name:"access-key" help:"Endpoint access key."`
}

func (c *EndpointCreateCmd) Run(g *Globals) error {
	db, err := open(g)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.CreateEndpoint(c.Name, client.EndpointConfig{
		AIModelName: c.Model,
		URL:         c.URL,
		Provider:    c.Provider,
		AccessKey:   c.AccessKey,
	})
	if err != nil {
		return err
	}
	return display.Value(map[string]any{"status": "created", "endpoint": c.Name}, display.ResolveFormat(string(g.Output), "text"))
}

type EndpointUpdateCmd struct {
	Name      string `arg:"" help:"Endpoint name."`
	URL       string `help:"Endpoint URL."`
	AccessKey string `name:"access-key" help:"Endpoint access key."`
}

func (c *EndpointUpdateCmd) Run(g *Globals) error {
	db, err := open(g)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.UpdateEndpoint(c.Name, client.EndpointUpdate{URL: c.URL, AccessKey: c.AccessKey}); err != nil {
		return err
	}
	return display.Value(map[string]any{"status": "updated", "endpoint": c.Name}, display.ResolveFormat(string(g.Output), "text"))
}

type EndpointDeleteCmd struct {
	Name string `arg:"" help:"Endpoint name."`
}

func (c *EndpointDeleteCmd) Run(g *Globals) error {
	db, err := open(g)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.DeleteEndpoint(c.Name); err != nil {
		return err
	}
	return display.Value(map[string]any{"status": "deleted", "endpoint": c.Name}, display.ResolveFormat(string(g.Output), "text"))
}

type AICompleteCmd struct {
	ModelKey   string `arg:"" help:"Model key."`
	Prompt     string `arg:"" optional:"" help:"Prompt text."`
	Parameters string `help:"JSON parameters."`
	Pipe       bool   `help:"Force plain text output instead of the terminal pager."`
}

func (c *AICompleteCmd) Run(g *Globals) error {
	if c.Prompt == "" {
		if err := huh.NewText().Title("Prompt").Value(&c.Prompt).Run(); err != nil {
			return err
		}
	}
	db, err := open(g)
	if err != nil {
		return err
	}
	defer db.Close()
	result, duration, err := db.Complete(c.ModelKey, c.Prompt, c.Parameters)
	if err != nil {
		return err
	}
	format := display.ResolveFormat(string(g.Output), "text")
	if display.IsStructured(format) {
		return display.Value(map[string]any{"result": result, "model_key": c.ModelKey, "duration_ms": duration.Milliseconds()}, format)
	}
	if format == "table" {
		return display.Rows([]map[string]any{{"result": result}}, "table")
	}
	if c.Pipe {
		return display.Value(result, "text")
	}
	return display.LongText(result, os.Stdout)
}

type AIEmbedCmd struct {
	ModelKey string `arg:"" help:"Model key."`
	Text     string `arg:"" optional:"" help:"Text to embed."`
	Dim      int    `help:"Embedding dimension."`
	Stdin    bool   `help:"Read newline-delimited texts from stdin."`
}

func (c *AIEmbedCmd) Run(g *Globals) error {
	db, err := open(g)
	if err != nil {
		return err
	}
	defer db.Close()
	if c.Stdin {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		lines := nonEmptyLines(string(data))
		embeddings, err := db.EmbedBatch(c.ModelKey, lines, c.Dim)
		if err != nil {
			return err
		}
		return display.Value(map[string]any{"embeddings": embeddings, "model_key": c.ModelKey, "dim": c.Dim}, display.ResolveFormat(string(g.Output), "json"))
	}
	if c.Text == "" {
		return fmt.Errorf("text is required unless --stdin is used")
	}
	embedding, err := db.Embed(c.ModelKey, c.Text, c.Dim)
	if err != nil {
		return err
	}
	format := display.ResolveFormat(string(g.Output), "text")
	if format == "table" {
		return display.Rows([]map[string]any{{"embedding": embedding}}, "table")
	}
	return display.Value(embedding, format)
}

type AIRerankCmd struct {
	ModelKey  string `arg:"" help:"Model key."`
	Query     string `arg:"" help:"Query."`
	Documents string `arg:"" help:"Documents JSON array."`
}

func (c *AIRerankCmd) Run(g *Globals) error {
	db, err := open(g)
	if err != nil {
		return err
	}
	defer db.Close()
	result, err := db.Rerank(c.ModelKey, c.Query, c.Documents)
	if err != nil {
		return err
	}
	format := display.ResolveFormat(string(g.Output), "table")
	if display.IsStructured(format) {
		return display.Value(parseJSONValue(result), format)
	}
	rows := rerankRows(result)
	return display.Rows(rows, "table")
}

func parseJSONValue(raw string) any {
	var parsed any
	if json.Unmarshal([]byte(raw), &parsed) == nil {
		return parsed
	}
	return raw
}

func rerankRows(raw string) []map[string]any {
	var parsed []map[string]any
	if json.Unmarshal([]byte(raw), &parsed) == nil {
		sort.SliceStable(parsed, func(i, j int) bool {
			return scoreValue(parsed[i]) > scoreValue(parsed[j])
		})
		return parsed
	}
	return []map[string]any{{"result": raw}}
}

func scoreValue(row map[string]any) float64 {
	for _, key := range []string{"relevance_score", "score", "relevanceScore"} {
		if value, ok := row[key]; ok {
			switch v := value.(type) {
			case float64:
				return v
			case float32:
				return float64(v)
			case int:
				return float64(v)
			case json.Number:
				score, _ := v.Float64()
				return score
			}
		}
	}
	return 0
}

type AIPromptCmd struct {
	Template string   `arg:"" help:"Prompt template."`
	Args     []string `arg:"" optional:"" help:"Template arguments."`
}

func (c *AIPromptCmd) Run(g *Globals) error {
	db, err := open(g)
	if err != nil {
		return err
	}
	defer db.Close()
	result, duration, err := db.Prompt(c.Template, c.Args)
	if err != nil {
		return err
	}
	format := display.ResolveFormat(string(g.Output), "text")
	if display.IsStructured(format) {
		return display.Value(map[string]any{"result": result, "duration_ms": duration.Milliseconds()}, format)
	}
	return display.LongText(result, os.Stdout)
}

type SQLCmd struct {
	Statement []string `arg:"" help:"SQL statement."`
}

func (c *SQLCmd) Run(g *Globals) error {
	db, err := open(g)
	if err != nil {
		return err
	}
	defer db.Close()
	result, err := db.QueryResult(strings.Join(c.Statement, " "))
	if err != nil {
		return err
	}
	return display.Rows(result.Rows, display.ResolveFormat(string(g.Output), "table"), result.Columns)
}

func open(g *Globals) (*client.Client, error) {
	cfg, err := g.ClientConfig()
	if err != nil {
		return nil, err
	}
	return client.Open(cfg)
}

func nonEmptyLines(value string) []string {
	raw := strings.Split(value, "\n")
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}
