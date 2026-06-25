package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/magalab/seekai-cli/client"
	"github.com/magalab/seekai-cli/config"
)

type OutputFormat string

const (
	OutputAuto  OutputFormat = "auto"
	OutputTable OutputFormat = "table"
	OutputJSON  OutputFormat = "json"
	OutputYAML  OutputFormat = "yaml"
	OutputTOML  OutputFormat = "toml"
	OutputText  OutputFormat = "text"
)

type Globals struct {
	Host     string       `help:"Seek DB host." default:"localhost"`
	Port     int          `help:"Seek DB MySQL protocol port." default:"2881"`
	User     string       `help:"Seek DB user." default:"root"`
	Password string       `help:"Seek DB password."`
	Database string       `help:"Default database." default:"test"`
	Profile  string       `help:"Profile name from ~/.seekai/config.toml."`
	Output   OutputFormat `short:"o" enum:"auto,table,json,yaml,toml,text" help:"Output format: auto, table, json, yaml, toml, or text." default:"auto"`
}

func (g Globals) ClientConfig() (client.Config, error) {
	cfg := client.Config{
		Host:     g.Host,
		Port:     g.Port,
		User:     g.User,
		Password: g.Password,
		Database: g.Database,
	}

	file, err := config.LoadDefault()
	if err != nil {
		return cfg, err
	}
	if file == nil {
		return cfg, nil
	}

	profileName := g.Profile
	profile := file.Default
	if profileName != "" {
		selected, ok := file.Profiles[profileName]
		if !ok {
			return cfg, fmt.Errorf("profile %q not found in %s", profileName, config.DefaultPath())
		}
		profile = selected
	}

	applyProfile(&cfg, profile)
	applyExplicitFlags(&cfg, g)
	return cfg, nil
}

func applyProfile(cfg *client.Config, p config.Profile) {
	if p.Host != "" {
		cfg.Host = p.Host
	}
	if p.Port != 0 {
		cfg.Port = p.Port
	}
	if p.User != "" {
		cfg.User = p.User
	}
	if p.Password != "" {
		cfg.Password = expandEnv(p.Password)
	}
	if p.Database != "" {
		cfg.Database = p.Database
	}
}

func applyExplicitFlags(cfg *client.Config, g Globals) {
	args := os.Args[1:]
	set := map[string]bool{}
	for i, arg := range args {
		if strings.HasPrefix(arg, "--") {
			key := strings.TrimPrefix(strings.SplitN(arg, "=", 2)[0], "--")
			set[key] = true
			continue
		}
		if arg == "-o" && i+1 < len(args) {
			set["output"] = true
		}
	}
	if set["host"] {
		cfg.Host = g.Host
	}
	if set["port"] {
		cfg.Port = g.Port
	}
	if set["user"] {
		cfg.User = g.User
	}
	if set["password"] {
		cfg.Password = g.Password
	}
	if set["database"] {
		cfg.Database = g.Database
	}
}

func expandEnv(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		return os.Getenv(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}"))
	}
	return os.ExpandEnv(value)
}

func PrintError(err error) {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "connect: connection refused"),
		strings.Contains(msg, "i/o timeout"),
		strings.Contains(msg, "no such host"):
		fmt.Fprintf(os.Stderr, "connection failed: %s\ncheck --host, --port, user permissions, and that Seek DB is reachable\n", msg)
	case strings.Contains(strings.ToLower(msg), "access ai model"),
		strings.Contains(strings.ToLower(msg), "permission"):
		fmt.Fprintf(os.Stderr, "permission denied: %s\nthis operation may require ACCESS AI MODEL privileges\n", msg)
	case strings.Contains(strings.ToLower(msg), "syntax"):
		fmt.Fprintf(os.Stderr, "sql error: %s\ntry `seekai sql` to debug the statement directly\n", msg)
	default:
		fmt.Fprintln(os.Stderr, msg)
	}
}
