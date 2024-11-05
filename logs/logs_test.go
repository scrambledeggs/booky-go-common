// Log in JSON format for cloudwatch
// All level entries are the same except for Debug
// Level will be used for different parameters
package logs

import (
	"testing"

	"github.com/scrambledeggs/booky-go-common/assert"
)

func TestPrint(t *testing.T) {
	Debug("naknang", map[string]string{"naknang": "patatas"}, TO_SLACK)

	assert.ShouldBeNil(t, nil)
}
