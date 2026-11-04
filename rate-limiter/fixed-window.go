package ratelimiter

import "time"

type FixedWindow struct {
	WindowStart time.Time
	WindowSize  time.Duration
	Count       int
	Limit       int
}

func (fw *FixedWindow) Allow() bool {
	now := time.Now()

	if now.Sub(fw.WindowStart) >= fw.WindowSize {
		fw.Count = 0
		fw.WindowStart = now.Truncate(fw.WindowSize)
	}

	if fw.Count < fw.Limit {
		fw.Count++
		return true
	}

	return false
}
