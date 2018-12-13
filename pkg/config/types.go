package config

import "github.com/integr8ly/version-checker/pkg/meta"

type Config struct {
	Channels []Channel
	Releases []Release
}

type Channel struct {
	Name string      `json:"name"`
	Kind string      `json:"kind"`
	Spec interface{} `json:"spec"`
}

type Release struct {
	Name    string           `json:"name"`
	Kind    string           `json:"kind"`
	Version string           `json:"version"`
	Spec    meta.ReleaseSpec `json:"spec"`
}

type Loader struct {
	data  []byte
	index int64
}
