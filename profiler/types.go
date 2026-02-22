package profiler

// Analyzer abstracts the logic to execute a profiler and return a textual summary.
type Analyzer interface {
	// Analyze processes a target source file and returns its performance hotspots.
	// sampleIndex is purely optional and used by some implementations (like Go heap).
	Analyze(source string, sampleIndex string) (string, error)
}
