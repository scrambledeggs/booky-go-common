package assert

import (
	"testing"
)

func ShouldBeNil(t *testing.T, got any) {
	if got != nil {
		t.Errorf("got <%s> wanted <nil>", got)
	}
}
