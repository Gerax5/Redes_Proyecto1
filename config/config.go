// Package config loads and validates MCP server configuration from a TOML file.
package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// MCPServer defines a single MCP server entry in the configuration.
type MCPServer struct {
	Name    string            `toml:"name"`
	Type    string            `toml:"type"`    // "http" | "stdio"
	URL     string            `toml:"url"`     // si type=http
	Command string            `toml:"command"` // si type=stdio
	Args    []string          `toml:"args"`
	Env     map[string]string `toml:"env"`
}

// MCPBlock groups multiple MCP servers.
type MCPBlock struct {
	Servers []MCPServer `toml:"servers"`
}

// Config is the root configuration struct.
type Config struct {
	MCP MCPBlock `toml:"mcp"`
}

// Load reads and validates a configuration from a TOML file.
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
