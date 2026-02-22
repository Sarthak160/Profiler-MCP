package profiler_test

import (
	"fmt"
	"testing"

	"github.com/sarthak/pprof-mcp-server/profiler"
)

func TestAnalyze(t *testing.T) {
	// The CPU profile generated in the sample dir
	res, err := profiler.Analyze("../sample/cpu.prof", "")
	if err != nil {
		t.Fatalf("Failed to analyze: %v", err)
	}

	if len(res) == 0 {
		t.Fatal("Expected some output, got empty string")
	}

	fmt.Println("=== ANALYSIS OUTPUT ===")
	fmt.Println(res)
	fmt.Println("=======================")
}
