package adapters

import (
	ant "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/mark3labs/mcp-go/mcp"
)

func ToAnthropicTools(mcpTools []mcp.Tool) []ant.ToolUnionParam {
	out := make([]ant.ToolUnionParam, 0, len(mcpTools))
	for _, t := range mcpTools {
		out = append(out, ant.ToolUnionParam{
			OfTool: &ant.ToolParam{
				Name:        t.Name,
				Description: param.NewOpt(t.Description),
				InputSchema: ant.ToolInputSchemaParam{
					Required:   t.InputSchema.Required,
					Properties: t.InputSchema.Properties,
				},
			},
		})
	}
	return out
}
