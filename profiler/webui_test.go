package profiler_test

import (
	"strings"
	"testing"
	"time"

	"github.com/sarthak/pprof-mcp-server/profiler"
)

func TestLaunchWebUI(t *testing.T) {
	// The CPU profile generated in the sample dir
	res, err := profiler.LaunchWebUI("../sample/cpu.prof")
	if err != nil {
		t.Fatalf("Failed to launch web UI: %v", err)
	}

	if !strings.Contains(res, "http://localhost:") {
		t.Fatalf("Expected URL in result, got: %s", res)
	}
	t.Logf("Successfully launched UI at: %s", res)

	// Since it's a background process, we let it linger naturally
	// or wait a moment to ensure it doesn't crash immediately.
	time.Sleep(500 * time.Millisecond)
}
