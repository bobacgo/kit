package types

import "time"

type Duration string

func (d Duration) TimeDuration() time.Duration {
	td, _ := time.ParseDuration(string(d))
	return td
}

func (d Duration) Check() error {
	_, err := time.ParseDuration(string(d))
	return err
}

func (d Duration) ToTimeDuration() time.Duration {
	return d.TimeDuration()
}
