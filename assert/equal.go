package assert

import "testing"

func Equal(t *testing.T, got any, want any, notes string) {
	if got != want {
		t.Errorf("%s got <%s> wanted <%s>", notes, got, want)
	}
}
