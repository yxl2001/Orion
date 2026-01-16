package core

import (
	"context"
	"time"
)

type EngineConfig struct {
	Concurrency int
	RateLimit   int
	Timeout     time.Duration
	Retry       int
}

type Engine struct {
	config  EngineConfig
	bus     *EventBus
	store   Store
	modules map[ModuleKind][]Module
}

type RunResult struct {
	Findings []Finding
	Err      error
}

func NewEngine(cfg EngineConfig, bus *EventBus, store Store) *Engine {
	return &Engine{
		config:  cfg,
		bus:     bus,
		store:   store,
		modules: make(map[ModuleKind][]Module),
	}
}

func (e *Engine) RegisterModule(module Module) {
	e.modules[module.Kind()] = append(e.modules[module.Kind()], module)
}

func (e *Engine) Run(ctx context.Context, pipeline Pipeline) RunResult {
	input := &PipelineInput{
		TargetsFile: pipeline.TargetsFile,
	}

	var result RunResult

	for _, step := range pipeline.Steps {
		modules := e.modules[step]
		for _, module := range modules {
			moduleResult := module.Run(ctx, input, e.bus.Emit)
			if moduleResult.Err != nil {
				return RunResult{Err: moduleResult.Err}
			}

			if len(moduleResult.Observations) > 0 {
				e.store.SaveObservations(moduleResult.Observations)
			}
			if len(moduleResult.Fingerprints) > 0 {
				e.store.SaveFingerprints(moduleResult.Fingerprints)
			}
			if len(moduleResult.Findings) > 0 {
				e.store.SaveFindings(moduleResult.Findings)
				result.Findings = append(result.Findings, moduleResult.Findings...)
			}
		}
	}

	return result
}
