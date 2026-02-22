# How I Stopped My AI from Guessing: Building a Polyglot Profiling MCP Server

If you’re anything like me, you’ve recently found yourself pair-programming heavily with AI assistants like Claude, GitHub Copilot, or Cursor. They are incredible at generating boilerplate, writing tests, and explaining complex concepts. 

But a few weeks ago, I hit a wall. 

I was working on a series of rigorous performance optimization tasks for a Go backend. I would describe the latency issue to the LLM, and it would immediately spit out five different suggestions: *"Maybe it's escaping to the heap here,"* or *"Perhaps this lock is under heavy contention,"* or *"Try pre-allocating this slice."*

It was just... blindly guessing. 

So the loop became agonizing:
1. Make the AI-suggested code change.
2. Manually run `go test -bench . -cpuprofile cpu.prof -memprofile mem.prof`.
3. Manually run `go tool pprof -top cpu.prof`.
4. Copy the terminal output.
5. Paste it back into the AI window.
6. The AI says, *"Ah, my mistake. The bottleneck is actually over *here*."*
7. Repeat.

I realized I wasn't just fixing code; I was acting as a highly inefficient, human meat-bridge between my command line and my AI.

That’s when I discovered **MCP (Model Context Protocol)**. And it changed my workflow completely.

---

## What is MCP?
The [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) is an open standard that allows Large Language Models (LLMs) to connect securely to local or remote development tools. 

Instead of copying and pasting data into a chat window, you write a standalone "Server" that exposes custom Tools to the AI. If the AI realizes it needs information (like, say, a database schema, the contents of a Jira ticket, or a performance profile), it seamlessly executes the tool on your machine and reads the real-time result.

## The Idea: The "Pprof Inspector"
I realized that to fix my optimization loop, my AI didn't need better prompts. It needed **eyes**. It needed to be able to run profiles and look at the flamegraphs itself.

So, I built precisely that. A local MCP server that transforms an AI agent into a data-driven Senior Performance Engineer.

---

## How I Built It (and Made it Polyglot)

Because MCP supports standard `stdio` communication, building a local server is incredibly lightweight. I chose Go and used the excellent `github.com/mark3labs/mcp-go` SDK.

### The Original Go (`pprof`) Implementation
My initial goal was simple: wrapper `go tool pprof`. I created a tool called `analyze_profile` that takes a `profile_path_or_url`. When the LLM triggers the tool, the Go server safely forks a subprocess:
```go
exec.Command("go", "tool", "pprof", "-top", source)
```
It reads the `stdout`, trims it so it doesn't overwhelm the LLM's context window, and feeds the textual summary right back into the chat. Boom. The LLM can now read CPU profiles.

But I wanted to go deeper. What if there was a memory leak? 
I updated the tool to accept an optional `sample_index` parameter. Now, the LLM could proactively pass `"alloc_space"` or `"inuse_objects"` to dive deep into heap metrics dynamically.

### Giving the AI Visual Powers (Flamegraphs)
Textual summaries are great for the LLM, but sometimes *I* still need to see the flamegraph. But LLMs can't easily parse visual SVGs or complex browser UIs.

So, I built a second tool: `open_interactive_ui`. 
Instead of trying to parse the graphs for the LLM, this tool runs `go tool pprof -http=localhost:0` as a background daemon, extracts the ephemeral port it spawned on, and returns the URL string to the LLM. 

Now, when the LLM finishes analyzing the text, it tells me: *"I've found the issue. To visually explore the Flamegraph yourself, click here: `http://localhost:54321`."*

### Scaling it to Python and Java (The Polyglot Refactor)
Once I had Go working flawlessly, I realized this concept wasn't limited to Go. 

I refactored the MCP server into a clean interface-based router. The server now checks file extensions and dynamically changes its behavior:
- **Python (`.py`)**: The server runs the script directly with Python's native `cProfile` module and organizes the result by cumulative execution time.
- **Java (`.jar`)**: It detects compiled Java archives and executes them with `-XX:StartFlightRecording`. It runs the profile securely and then guides the LLM (and me) to open the resulting `.jfr` file in Java Mission Control.

## The Result
My workflow has completely transformed. 

Now, when I hit a performance snag, I just tell my agent:
> *"Run my application with a high load, and then use the Pprof Inspector to analyze the CPU profile."*

The AI generates the profile, reads the functions taking the most time, analyzes the source code corresponding to those functions, and provides a **data-driven** optimization strategy. No guessing required.

If you are tired of acting as a clipboard manager for your AI, I highly recommend diving into MCP. Expanding the context window is great, but granting your agent the ability to physically interact with your tooling is true autonomy.

---
*You can find the full source code for the Polyglot Profiler MCP here:*
[https://github.com/Sarthak160/Profiler-MCP](https://github.com/Sarthak160/Profiler-MCP)
