package step

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"
	"github.com/stretchr/testify/assert"
)

func TestAlwaysPass_Run(t *testing.T) {
	step := NewAlwaysPass()

	results := step.Run(context.NewContext())

	assert.True(t, results.Pass)
}
