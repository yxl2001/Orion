package core

type Target struct {
	Value  string
	Source string
	Labels map[string]string
}

type Asset struct {
	Host     string
	Port     int
	Protocol string
	Source   string
}

type Observation struct {
	Asset    Asset
	Type     string
	Payload  string
	Source   string
	Evidence string
}

type Fingerprint struct {
	Asset      Asset
	Product    string
	Version    string
	Confidence float64
	Evidence   []string
}

type Finding struct {
	Asset      Asset
	TemplateID string
	Severity   string
	MatchedAt  string
	Evidence   []string
}
