package fingerprint

import (
	"context"

	"orion/internal/core"
)

type BasicFingerprint struct{}

func NewBasicFingerprint() *BasicFingerprint {
	return &BasicFingerprint{}
}

func (f *BasicFingerprint) Name() string {
	return "basic-fingerprint"
}

func (f *BasicFingerprint) Kind() core.ModuleKind {
	return core.ModuleFingerprint
}

func (f *BasicFingerprint) Run(ctx context.Context, input *core.PipelineInput, emit core.EmitFunc) core.ModuleResult {
	fingerprint := core.Fingerprint{
		Asset:      core.Asset{Host: "example", Port: 80, Protocol: "http"},
		Product:    "placeholder",
		Version:    "0.0.0",
		Confidence: 0.1,
		Evidence:   []string{"stub"},
	}
	return core.ModuleResult{Fingerprints: []core.Fingerprint{fingerprint}}
}
