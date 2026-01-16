package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"orion/internal/app"
	"orion/internal/config"
)

func main() {
	var (
		configPath string
		command    string
	)

	flag.StringVar(&configPath, "config", "", "path to config file (json)")
	flag.StringVar(&command, "cmd", "all", "command to run: collect|probe|fp|scan|all")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config load error: %v\n", err)
		os.Exit(1)
	}

	application := app.New(cfg)
	ctx := context.Background()

	if err := application.Run(ctx, command); err != nil {
		fmt.Fprintf(os.Stderr, "run error: %v\n", err)
		os.Exit(1)
	}
}
