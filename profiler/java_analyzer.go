package profiler

import (
	"bytes"
	"fmt"
	"os/exec"
)

type JavaAnalyzer struct{}

func (j *JavaAnalyzer) Analyze(source string, sampleIndex string) (string, error) {
	// Execute JAR with Java Flight Recorder (JFR) active
	duration := "60s" // Defaulting to 60s for server safety if it's a daemon
	filename := "mcp-recording.jfr"
	jfrArgs := fmt.Sprintf("-XX:StartFlightRecording=duration=%s,filename=%s", duration, filename)

	cmd := exec.Command("java", jfrArgs, "-jar", source)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("java execution failed: %v\nStderr: %s", err, stderr.String())
	}

	// Java JFR outputs binary files meant for JMC or other parsers.
	// We give a custom message instead of the raw CLI output.
	result := fmt.Sprintf("Java Profiling Complete for %s.\n\nFlight Recording saved to: %s\n\nTo analyze Java profiling data visually, open this file using Java Mission Control (JDK Mission Control) or use 'jfr print %s'.\nStderr Output: %s", source, filename, filename, stderr.String())
	return result, nil
}
