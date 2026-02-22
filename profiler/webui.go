package profiler

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"sync"
)

var (
	// Registry of running commands to clean up if/when server shuts down (optional but good practice)
	runningServers []*exec.Cmd
	serverMutex    sync.Mutex
)

// LaunchWebUI runs `go tool pprof -http=localhost:0 -no_browser <source>`
// and parses the ephemeral port returned in standard error/out.
func LaunchWebUI(source string) (string, error) {
	cmd := exec.Command("go", "tool", "pprof", "-http=localhost:0", "-no_browser", source)

	// pprof writes the serving url usually to stderr
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to get stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start pprof web ui: %v", err)
	}

	// Register so we *could* kill it later if MCP supported lifecycle hooks,
	// but generally they'll die when the parent dies or they stay as background helpers.
	serverMutex.Lock()
	runningServers = append(runningServers, cmd)
	serverMutex.Unlock()

	// Parse stderr to find the URL
	// E.g. "Serving web UI on http://localhost:54321"
	scanner := bufio.NewScanner(stderrPipe)
	urlRegex := regexp.MustCompile(`http://localhost:\d+`)

	urlChan := make(chan string, 1)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			match := urlRegex.FindString(line)
			if match != "" {
				urlChan <- match
				return // Found it, stop scanning
			}
		}
		close(urlChan)
	}()

	// Wait for the URL or process exit
	url, ok := <-urlChan
	if !ok {
		// It exited without printing a URL or scanner finished
		return "", fmt.Errorf("could not determine web UI url from pprof output")
	}

	return fmt.Sprintf("Interactive Web UI (including Flamegraphs) is now running at: %s", url), nil
}
