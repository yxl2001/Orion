package nuclei

import "context"

type Runner struct {
	BinaryPath  string
	Templates   string
	Concurrency int
	TimeoutSec  int
}

func (r *Runner) Run(ctx context.Context, target string) error {
	return nil
}
