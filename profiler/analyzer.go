package profiler

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Analyze runs the `go tool pprof -top` command on the given profile source.
// The source can be a local file path or a remote URL.
// The sampleIndex can be used for things like `-alloc_space` or `-inuse_objects` for heap profiles.
func Analyze(source string, sampleIndex string) (string, error) {
	// Run go tool pprof in batch mode to get the top hotspots
	// Use -top to get the text report.
	args := []string{"tool", "pprof", "-top"}
	if sampleIndex != "" {
		args = append(args, "-sample_index="+sampleIndex)
	}
	args = append(args, source)

	cmd := exec.Command("go", args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Include stderr in the error message for better debugging
		return "", fmt.Errorf("pprof execution failed: %v\nStderr: %s", err, stderr.String())
	}

	// We only need the first ~30 lines to avoid overwhelming the LLM.
	outputStr := out.String()
	lines := strings.Split(outputStr, "\n")

	maxLines := 50
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, "... (truncated for brevity, remaining lines are less significant)")
	}

	result := fmt.Sprintf("Profile Analysis for %s:\n\n%s\n\nInterpretation hints:\n- flat: memory/time spent in the function itself\n- cum: memory/time spent in the function and its callees", source, strings.Join(lines, "\n"))
	return result, nil
}
