package model

import "time"

type Post struct {
	ID        int64
	Text      string
	Timestamp time.Time
}
