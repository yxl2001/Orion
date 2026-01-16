package core

import "context"

type ModuleKind string

const (
	ModuleCollector   ModuleKind = "collector"
	ModuleProbe       ModuleKind = "probe"
	ModuleFingerprint ModuleKind = "fingerprint"
	ModuleScanner     ModuleKind = "scanner"
)

type Module interface {
	Name() string
	Kind() ModuleKind
	Run(ctx context.Context, input *PipelineInput, emit EmitFunc) ModuleResult
}

type ModuleResult struct {
	Observations []Observation
	Fingerprints []Fingerprint
	Findings     []Finding
	Err          error
}
