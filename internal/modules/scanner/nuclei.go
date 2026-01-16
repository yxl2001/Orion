package scanner

import (
	"context"

	"orion/internal/core"
)

type NucleiAdapter struct{}

func NewNucleiAdapter() *NucleiAdapter {
	return &NucleiAdapter{}
}

func (s *NucleiAdapter) Name() string {
	return "nuclei-adapter"
}

func (s *NucleiAdapter) Kind() core.ModuleKind {
	return core.ModuleScanner
}

func (s *NucleiAdapter) Run(ctx context.Context, input *core.PipelineInput, emit core.EmitFunc) core.ModuleResult {
	finding := core.Finding{
		Asset:      core.Asset{Host: "example", Port: 443, Protocol: "https"},
		TemplateID: "stub-template",
		Severity:   "info",
		MatchedAt:  "https://example",
		Evidence:   []string{"stub"},
	}
	return core.ModuleResult{Findings: []core.Finding{finding}}
}
