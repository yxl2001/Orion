# Orion

Orion is a Go-based security pipeline prototype for asset collection, probing, fingerprinting, and vulnerability scanning.

## Goals

- Provide a pluggable, composable pipeline for recon and scanning.
- Normalize assets, observations, fingerprints, and findings.
- Keep modules replaceable and core orchestration stable.

## Quick start

```bash
# Create a targets.txt file with one target per line.
printf "example.com\n" > targets.txt

# Run the default pipeline.
go run ./cmd/orion -cmd all
```

## Configuration

The prototype supports a JSON config file:

```json
{
  "concurrency": 20,
  "rate_limit": 200,
  "timeout": "10s",
  "retry": 1,
  "targets_file": "targets.txt",
  "output_path": "output.jsonl"
}
```

Run with:

```bash
go run ./cmd/orion -config config.json -cmd scan
```
