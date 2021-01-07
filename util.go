package main

import "time"

func epoch(t time.Time) int64 {
	return (t.Unix() - 1598306400) / 30
}
