package app

import (
	"context"
	"fmt"

	"orion/internal/config"
	"orion/internal/core"
	"orion/internal/modules/collector"
	"orion/internal/modules/fingerprint"
	"orion/internal/modules/probe"
	"orion/internal/modules/scanner"
	"orion/internal/storage"
)

type App struct {
	config config.Config
}

func New(cfg config.Config) *App {
	return &App{config: cfg}
}

func (a *App) Run(ctx context.Context, command string) error {
	bus := core.NewEventBus()
	store := storage.NewMemoryStore()

	engine := core.NewEngine(core.EngineConfig{
		Concurrency: a.config.Concurrency,
		RateLimit:   a.config.RateLimit,
		Timeout:     a.config.Timeout,
		Retry:       a.config.Retry,
	}, bus, store)

	engine.RegisterModule(collector.NewFileCollector(a.config.TargetsFile))
	engine.RegisterModule(probe.NewBasicProber())
	engine.RegisterModule(fingerprint.NewBasicFingerprint())
	engine.RegisterModule(scanner.NewNucleiAdapter())

	pipeline, err := core.NewPipeline(command)
	if err != nil {
		return err
	}

	result := engine.Run(ctx, pipeline)
	if result.Err != nil {
		return result.Err
	}

	fmt.Printf("completed pipeline %s with %d findings\n", pipeline.Name, len(result.Findings))
	return nil
}
