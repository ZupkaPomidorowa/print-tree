package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	r := NewEmptyRendering()
	r.AddOnTop("22left22").ShiftTopBy(-2)
	r.AddOnTop("1left1").ShiftTopBy(-1)
	r.AddOnTop("left")

	var expected = Nlnl(`
  left
 1left1
22left22
`)

	assert.Equal(t, expected, r.String())
	r.Reverse()
	expected = Nlnl(`
22left22
 1left1
  left
`)
	assert.Equal(t, expected, r.String())
}
