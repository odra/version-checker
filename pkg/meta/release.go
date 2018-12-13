package meta

import (
	"time"
)

type ReleaseKind string

type Version struct {
	Ver  string
	Date time.Time
}

type ReleaseSpec interface {
	IsKind(kind ReleaseKind) bool
	LatestVersion() (*Version, error)
	HasNewVersion() (*Version, bool, error)
}
