package assert

import (
	"reflect"
	"testing"
)

func DeepEqual(t *testing.T, got any, want any, notes string) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s got <%s> wanted <%s>", notes, got, want)
	}
}
