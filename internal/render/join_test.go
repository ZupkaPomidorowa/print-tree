package render

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinSimple(t *testing.T) {
	left := NewEmptyRendering()
	left.AddOnTop("baz")
	left.AddOnTop("bar")
	left.AddOnTop("foo")

	var expected = Nlnl(`
foo
bar
baz
`)

	assert.Equal(t, expected, left.String())

	right := NewEmptyRendering()
	right.AddOnTop("789")
	right.AddOnTop("456")
	right.AddOnTop("123")

	expected = Nlnl(`
123
456
789
`)

	assert.Equal(t, expected, right.String())

	distance := AlignDistance(left, right)
	joined := JoinRenderings(left, right, distance)

	expected = Nlnl(`
foo123
bar456
baz789
`)
	assert.Equal(t, expected, joined.String())
	assert.Equal(t, "", joined.TopRow().Prefix())
	assert.Equal(t, joined.TopRow().Value(), joined.TopRow().Suffix())

	joined = JoinRenderings(left, right, distance+1)

	expected = Nlnl(`
foo 123
bar 456
baz 789
`)
	assert.Equal(t, expected, joined.String())
	assert.Equal(t, "", joined.TopRow().Prefix())
	assert.Equal(t, joined.TopRow().Value(), joined.TopRow().Suffix())

}

func TestJoinSymmetric(t *testing.T) {
	left := NewEmptyRendering()
	left.AddOnTop("22left22").ShiftTopBy(-2)
	left.AddOnTop("1left1").ShiftTopBy(-1)
	left.AddOnTop("left")

	var expected = Nlnl(`
  left
 1left1
22left22
`)

	assert.Equal(t, expected, left.String())

	right := NewPartialRendering("22right22").ShiftTopBy(-2)
	right.AddOnTop("1right1").ShiftTopBy(-1)
	right.AddOnTop("right")

	expected = Nlnl(`
  right
 1right1
22right22
`)

	assert.Equal(t, expected, right.String())

	distance := AlignDistance(left, right)
	joined := JoinRenderings(left, right, distance)

	expected = Nlnl(`
  left    right
 1left1  1right1
22left2222right22
`)
	assert.Equal(t, expected, joined.String())
	assert.Equal(t, "", joined.TopRow().Prefix())
	assert.Equal(t, joined.TopRow().Value(), joined.TopRow().Suffix())
}

func TestJoinSymmetricWithAdditionalDistance(t *testing.T) {
	left := NewEmptyRendering()
	left.AddOnTop("22left22").ShiftTopBy(-2)
	left.AddOnTop("1left1").ShiftTopBy(-1)
	left.AddOnTop("left")

	var expected = Nlnl(`
  left
 1left1
22left22
`)

	assert.Equal(t, expected, left.String())

	right := NewEmptyRendering()
	right.AddOnTop("22right22").ShiftTopBy(-2)
	right.AddOnTop("1right1").ShiftTopBy(-1)
	right.AddOnTop("right")

	expected = Nlnl(`
  right
 1right1
22right22
`)

	assert.Equal(t, expected, right.String())

	distance := AlignDistance(left, right)
	additionalDistance := 3
	joined := JoinRenderings(left, right, distance+additionalDistance)

	expected = Nlnl(`
  left       right
 1left1     1right1
22left22   22right22
`)
	assert.Equal(t, expected, joined.String())
	assert.Equal(t, "", joined.TopRow().Prefix())
	assert.Equal(t, joined.TopRow().Value(), joined.TopRow().Suffix())
}

func TestJoinLeftIsHigher(t *testing.T) {
	left := NewEmptyRendering()
	left.AddOnTop("666666left666666").ShiftTopBy(-6)
	left.AddOnTop("55555left55555").ShiftTopBy(-5)
	left.AddOnTop("4444left4444").ShiftTopBy(-4)
	left.AddOnTop("333left333").ShiftTopBy(-3)
	left.AddOnTop("22left22").ShiftTopBy(-2)
	left.AddOnTop("1left1").ShiftTopBy(-1)
	left.AddOnTop("left")

	var expected = Nlnl(`
      left
     1left1
    22left22
   333left333
  4444left4444
 55555left55555
666666left666666
`)

	assert.Equal(t, expected, left.String())

	right := NewEmptyRendering()
	right.AddOnTop("22right22").ShiftTopBy(-2)
	right.AddOnTop("1right1").ShiftTopBy(-1)
	right.AddOnTop("right")

	expected = Nlnl(`
  right
 1right1
22right22
`)

	assert.Equal(t, expected, right.String())

	distance := AlignDistance(left, right)
	additionalDistance := 3
	joined := JoinRenderings(left, right, distance+additionalDistance)

	expected = Nlnl(`
      left       right
     1left1     1right1
    22left22   22right22
   333left333
  4444left4444
 55555left55555
666666left666666
`)
	assert.Equal(t, expected, joined.String())
	assert.Equal(t, "", joined.TopRow().Prefix())
	assert.Equal(t, joined.TopRow().Value(), joined.TopRow().Suffix())
}

func TestJoinRightsHigher(t *testing.T) {
	left := NewEmptyRendering()
	left.AddOnTop("22left22").ShiftTopBy(-2)
	left.AddOnTop("1left1").ShiftTopBy(-1)
	left.AddOnTop("left")

	var expected = Nlnl(`
  left
 1left1
22left22
`)

	assert.Equal(t, expected, left.String())

	right := NewEmptyRendering()
	right.AddOnTop("666666right666666").ShiftTopBy(-6)
	right.AddOnTop("55555right55555").ShiftTopBy(-5)
	right.AddOnTop("4444right4444").ShiftTopBy(-4)
	right.AddOnTop("333right333").ShiftTopBy(-3)
	right.AddOnTop("22right22").ShiftTopBy(-2)
	right.AddOnTop("1right1").ShiftTopBy(-1)
	right.AddOnTop("right")

	expected = Nlnl(`
      right
     1right1
    22right22
   333right333
  4444right4444
 55555right55555
666666right666666
`)

	assert.Equal(t, expected, right.String())

	distance := AlignDistance(left, right)
	additionalDistance := 3
	joined := JoinRenderings(left, right, distance+additionalDistance)

	expected = Nlnl(`
  left       right
 1left1     1right1
22left22   22right22
          333right333
         4444right4444
        55555right55555
       666666right666666
`)
	assert.Equal(t, expected, joined.String())
	assert.Equal(t, "", joined.TopRow().Prefix())
	assert.Equal(t, joined.TopRow().Value(), joined.TopRow().Suffix())
}

func TestJoinSymmetricWithNormalize(t *testing.T) {
	left := NewEmptyRendering()
	left.AddOnTop("22left22").ShiftTopBy(-2)
	left.AddOnTop("1left1").ShiftTopBy(-1)
	left.AddOnTop("left").
		NormalizeOffsetsRev()

	var expected = Nlnl(`
  left
 1left1
22left22
`)

	assert.Equal(t, expected, left.String())

	right := NewEmptyRendering()
	right.AddOnTop("22right22").ShiftTopBy(-2)
	right.AddOnTop("1right1").ShiftTopBy(-1)
	right.AddOnTop("right").
		NormalizeOffsets()

	expected = Nlnl(`
  right
 1right1
22right22
`)

	assert.Equal(t, expected, right.String())
	distance := AlignDistance(left, right)
	joined := JoinRenderings(left, right, distance)

	expected = Nlnl(`
  left    right
 1left1  1right1
22left2222right22
`)
	fmt.Println(joined)
	assert.Equal(t, expected, joined.String())
	assert.Equal(t, "left", joined.TopRow().Prefix())
	assert.Equal(t, "    right", joined.TopRow().Suffix())
}
