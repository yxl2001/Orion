package collector

import (
	"bufio"
	"context"
	"os"
	"strings"

	"orion/internal/core"
)

type FileCollector struct {
	path string
}

func NewFileCollector(path string) *FileCollector {
	return &FileCollector{path: path}
}

func (c *FileCollector) Name() string {
	return "file-collector"
}

func (c *FileCollector) Kind() core.ModuleKind {
	return core.ModuleCollector
}

func (c *FileCollector) Run(ctx context.Context, input *core.PipelineInput, emit core.EmitFunc) core.ModuleResult {
	path := c.path
	if input.TargetsFile != "" {
		path = input.TargetsFile
	}

	file, err := os.Open(path)
	if err != nil {
		return core.ModuleResult{Err: err}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value := strings.TrimSpace(scanner.Text())
		if value == "" {
			continue
		}
		emit(core.Event{Type: "TargetDiscovered", Data: core.Target{Value: value, Source: c.Name()}})
	}

	if err := scanner.Err(); err != nil {
		return core.ModuleResult{Err: err}
	}

	return core.ModuleResult{}
}
