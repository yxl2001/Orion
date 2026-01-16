package core

type Store interface {
	SaveObservations([]Observation)
	SaveFingerprints([]Fingerprint)
	SaveFindings([]Finding)
}
