package context

import (
	"testing"
)

func TestRunWithConfigDir(t *testing.T) {
	RunWithConfigDir("./context", "www")
}