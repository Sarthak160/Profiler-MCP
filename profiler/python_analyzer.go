package profiler

import (
	"bytes"
	"fmt"
	"os/exec"
)

type PythonAnalyzer struct{}

func (p *PythonAnalyzer) Analyze(source string, sampleIndex string) (string, error) {
	// Profile a Python script on the fly
	cmd := exec.Command("python3", "-m", "cProfile", "-s", "cumtime", source)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("python cProfile execution failed: %v\nStderr: %s", err, stderr.String())
	}

	result := fmt.Sprintf("Python Profile Analysis for %s:\n\n%s", source, out.String())
	return result, nil
}
