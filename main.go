package main

import (
	"context"
	"fmt"

	"github.com/lmorg/mcp-scrape-page/internal"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		internal.Name,
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add tool
	tool := mcp.NewTool(internal.Name,
		mcp.WithDescription(internal.Description),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL to scrape"),
		),
	)

	// Add tool handler
	s.AddTool(tool, scraperHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func scraperHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, err := request.RequireString("url")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	response, err := internal.Scrape(ctx, url)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(response), nil
}
