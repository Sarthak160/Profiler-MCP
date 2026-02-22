package profiler

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type GoAnalyzer struct{}

func (g *GoAnalyzer) Analyze(source string, sampleIndex string) (string, error) {
	// Run go tool pprof in batch mode
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

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("go pprof execution failed: %v\nStderr: %s", err, stderr.String())
	}

	// We only need the first ~30-50 lines to avoid overwhelming the LLM.
	outputStr := out.String()
	lines := strings.Split(outputStr, "\n")

	maxLines := 50
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, "... (truncated for brevity, remaining lines are less significant)")
	}

	result := fmt.Sprintf("Go Profile Analysis for %s:\n\n%s\n\nInterpretation hints:\n- flat: memory/time spent in the function itself\n- cum: memory/time spent in the function and its callees", source, strings.Join(lines, "\n"))
	return result, nil
}
