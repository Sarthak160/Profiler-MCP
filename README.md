# Profiler-MCP (Performance Inspector)

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server that gives your AI Coding Agent the ability to autonomously analyze Go application performance bottlenecks using `go tool pprof`.

Instead of blindly guessing why an application is slow or leaking memory, this server allows an LLM (like Claude or Copilot via an MCP client) to ingest actual profiling data, read the hot paths, and even spin up an interactive Flamegraph UI for you to inspect locally.

## Features

- 🕵️‍♂️ **CPU Profile Analysis**: Parses `cpu.prof` and returns a summarized, LLM-friendly textual analysis of where time is being spent.
- 🧠 **Memory/Heap Profiling**: Supports `-sample_index` filtering (`alloc_space`, `inuse_objects`, etc.) to track down memory leaks in `mem.prof`.
- 🔥 **Interactive Flamegraphs**: The AI can launch an ephemeral `pprof` web server in the background and hand you a clickable `http://localhost:<port>` link directly in the chat so you can explore the flamegraph visually.

## Pre-requisites

- [Go](https://go.dev/doc/install) (1.20+ recommended)
- To view visual graphs, you may optionally need [Graphviz](https://graphviz.org/download/) installed on your system.

## Setup / Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/Sarthak160/Profiler-MCP.git
   cd Profiler-MCP
   ```

2. Build the binary:
   ```bash
   go build -o pprof-mcp-server main.go
   ```

3. Note the **absolute path** to the compiled `pprof-mcp-server` binary.

## Integrating with AI Agents (MCP Clients)

Because this is an MCP server, it runs entirely locally on your machine over standard `stdio`. You just need to tell your AI client where the binary is located.

### 1. Claude Desktop
Add the following to your Claude Desktop configuration file (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "pprof-inspector": {
      "command": "/absolute/path/to/Profiler-MCP/pprof-mcp-server"
    }
  }
}
```

### 2. VS Code (Cline / Roo Code extensions)
Open the **MCP Servers** configuration inside the Cline or Roo Code extension settings (`VS Code -> Command Palette -> Cline: Open MCP Settings`) and add:

```json
{
  "mcpServers": {
    "pprof-inspector": {
      "command": "/absolute/path/to/Profiler-MCP/pprof-mcp-server"
    }
  }
}
```

## Available Tools

Once connected, your AI assistant will automatically have access to these tools:

1. **`analyze_profile`**: 
   - **Arguments**: `profile_path_or_url` (required), `sample_index` (optional)
   - **Description**: Analyzes a pprof profile (cpu, heap) to identify performance bottlenecks.
2. **`open_interactive_ui`**:
   - **Arguments**: `profile_path_or_url` (required)
   - **Description**: Launches a background pprof web server and returns the local HTTP URL to the LLM, allowing it to provide you a link for visual Flamegraph inspection.

## Example Prompts to try with your AI

* *"I generated a CPU profile at `./app/cpu.prof`. Can you use your pprof tools to tell me which function is consuming the most time?"*
* *"We have a memory leak. I just dumped `./app/mem.prof`. Parse it and look specifically at `-alloc_space` to tell me where the leak is originating."*
* *"Open an interactive flamegraph UI for `./app/cpu.prof` and give me the localhost link."*

## Testing Locally (Sample App)

If you want to test the server immediately, this repo contains a dummy app that intentionally spins the CPU and leaks memory.

```bash
cd sample
go run main.go

# This will generate both a cpu.prof and mem.prof locally!
```

Now you can point your configured AI agent at those newly generated `.prof` files.
