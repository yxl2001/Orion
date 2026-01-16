package core

import "fmt"

type Pipeline struct {
	Name        string
	Steps       []ModuleKind
	TargetsFile string
}

type PipelineInput struct {
	TargetsFile string
}

func NewPipeline(name string) (Pipeline, error) {
	switch name {
	case "collect":
		return Pipeline{Name: name, Steps: []ModuleKind{ModuleCollector}}, nil
	case "probe":
		return Pipeline{Name: name, Steps: []ModuleKind{ModuleCollector, ModuleProbe}}, nil
	case "fp":
		return Pipeline{Name: name, Steps: []ModuleKind{ModuleCollector, ModuleProbe, ModuleFingerprint}}, nil
	case "scan":
		return Pipeline{Name: name, Steps: []ModuleKind{ModuleCollector, ModuleProbe, ModuleFingerprint, ModuleScanner}}, nil
	case "all":
		return Pipeline{Name: name, Steps: []ModuleKind{ModuleCollector, ModuleProbe, ModuleFingerprint, ModuleScanner}}, nil
	default:
		return Pipeline{}, fmt.Errorf("unknown pipeline %q", name)
	}
}
