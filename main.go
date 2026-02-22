package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sarthak/pprof-mcp-server/profiler"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer("Pprof Inspector", "0.1.0")

	// Tool 1: analyze_profile
	analyzeTool := mcp.NewTool("analyze_profile",
		mcp.WithDescription("Analyze a pprof profile (cpu, heap, etc.) to identify performance bottlenecks and hot paths. Requires a path or URL to the profile."),
		mcp.WithString("profile_path_or_url", mcp.Required(), mcp.Description("The local file path or remote HTTP URL to the pprof profile. Example: 'cpu.prof'")),
		mcp.WithString("sample_index", mcp.Description("Optional. The sample value to report (e.g. 'alloc_space', 'inuse_space' for heap profiles)")),
	)
	s.AddTool(analyzeTool, analyzeProfileHandler)

	// Tool 2: open_interactive_ui
	webUITool := mcp.NewTool("open_interactive_ui",
		mcp.WithDescription("Launches a background pprof embedded web server for a given profile and returns the local HTTP URL. Use this to provide Flamegraphs and visual graphs to the user."),
		mcp.WithString("profile_path_or_url", mcp.Required(), mcp.Description("The local file path or remote HTTP URL to the pprof profile.")),
	)
	s.AddTool(webUITool, openWebUIHandler)

	// Start standard standard input/output server
	log.Println("Starting Pprof MCP Server on stdio...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func analyzeProfileHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	profileSource, err := req.RequireString("profile_path_or_url")
	if err != nil || profileSource == "" {
		return mcp.NewToolResultError("profile_path_or_url is required"), nil
	}

	sampleIndex := req.GetString("sample_index", "")

	summary, err := profiler.Analyze(profileSource, sampleIndex)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to analyze profile: %v", err)), nil
	}

	return mcp.NewToolResultText(summary), nil
}

func openWebUIHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	profileSource, err := req.RequireString("profile_path_or_url")
	if err != nil || profileSource == "" {
		return mcp.NewToolResultError("profile_path_or_url is required"), nil
	}

	resultMsg, err := profiler.LaunchWebUI(profileSource)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to launch web UI: %v", err)), nil
	}

	return mcp.NewToolResultText(resultMsg), nil
}
