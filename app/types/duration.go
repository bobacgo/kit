package types

import "time"

// Duration
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
type Duration string

func (d Duration) TimeDuration() time.Duration {
	if d == "" {
		return 0
	}
	td, _ := time.ParseDuration(string(d))
	return td
}

func (d Duration) Check() error {
	_, err := time.ParseDuration(string(d))
	return err
}