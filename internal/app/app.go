package app

import "time"

type ExistingHash struct {
	EntryHash string
	ID        string
}

type ExistingFiles struct {
	CommittedAt time.Time
	ID          string
}
