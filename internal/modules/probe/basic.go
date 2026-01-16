package probe

import (
	"context"
	"fmt"

	"orion/internal/core"
)

type BasicProber struct{}

func NewBasicProber() *BasicProber {
	return &BasicProber{}
}

func (p *BasicProber) Name() string {
	return "basic-prober"
}

func (p *BasicProber) Kind() core.ModuleKind {
	return core.ModuleProbe
}

func (p *BasicProber) Run(ctx context.Context, input *core.PipelineInput, emit core.EmitFunc) core.ModuleResult {
	observation := core.Observation{
		Asset:   core.Asset{Host: "example", Port: 80, Protocol: "tcp"},
		Type:    "probe",
		Payload: fmt.Sprintf("placeholder probe output for %s", p.Name()),
		Source:  p.Name(),
	}
	return core.ModuleResult{Observations: []core.Observation{observation}}
}
