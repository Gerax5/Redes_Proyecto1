package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type MCPServer struct {
	Name    string            `toml:"name"`
	Type    string            `toml:"type"`    // "http" | "stdio"
	URL     string            `toml:"url"`     // si type=http
	Command string            `toml:"command"` // si type=stdio
	Args    []string          `toml:"args"`
	Env     map[string]string `toml:"env"`
}

type MCPBlock struct {
	Servers []MCPServer `toml:"servers"`
}

type Config struct {
	MCP MCPBlock `toml:"mcp"`
}

func Load(path string) (Config, error) {
	var cfg Config
	if _, err := os.Stat(path); err != nil {
		return cfg, fmt.Errorf("no se encontró %s: %w", path, err)
	}
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return cfg, err
	}
	// Validación mínima
	for i, s := range cfg.MCP.Servers {
		switch s.Type {
		case "http":
			if s.URL == "" {
				return cfg, fmt.Errorf("mcp.servers[%d] (%s): falta url para type=http", i, s.Name)
			}
		case "stdio":
			if s.Command == "" {
				return cfg, fmt.Errorf("mcp.servers[%d] (%s): falta command para type=stdio", i, s.Name)
			}
		default:
			return cfg, fmt.Errorf("mcp.servers[%d] (%s): type inválido: %s", i, s.Name, s.Type)
		}
	}
	return cfg, nil
}
