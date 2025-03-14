package app

import "testing"

func TestXxx(t *testing.T) {
	m := make(map[string]string)
	m[""] = "123"
	t.Log(m[""])
}
