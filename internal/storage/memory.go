package storage

import "orion/internal/core"

type MemoryStore struct {
	Observations []core.Observation
	Fingerprints []core.Fingerprint
	Findings     []core.Finding
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) SaveObservations(items []core.Observation) {
	s.Observations = append(s.Observations, items...)
}

func (s *MemoryStore) SaveFingerprints(items []core.Fingerprint) {
	s.Fingerprints = append(s.Fingerprints, items...)
}

func (s *MemoryStore) SaveFindings(items []core.Finding) {
	s.Findings = append(s.Findings, items...)
}
