package profiler

import (
	"strings"
)

// Analyze determines the target language based on the file extension
// and delegates to the appropriate specialized Analyzer.
func Analyze(source string, sampleIndex string) (string, error) {
	var analyzer Analyzer

	if strings.HasSuffix(source, ".py") {
		analyzer = &PythonAnalyzer{}
	} else if strings.HasSuffix(source, ".jar") {
		analyzer = &JavaAnalyzer{}
	} else {
		// Default to Go pprof for .prof or generic binaries
		analyzer = &GoAnalyzer{}
	}

	return analyzer.Analyze(source, sampleIndex)
}
