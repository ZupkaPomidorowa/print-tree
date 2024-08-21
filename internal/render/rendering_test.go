package render

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShiftTopBy(t *testing.T) {
	pr := NewPartialRendering("foo")
	pr.ShiftTopBy(-1)

	pr.AddOnTop("bar").ShiftTopBy(-2)
	pr.AddOnTop("baz").ShiftTopBy(-3)

	assert.Len(t, pr.Rows, 3)
	assert.Equal(t, -1, pr.Rows[0].offset)
	assert.Equal(t, -2, pr.Rows[1].offset)
	assert.Equal(t, -3, pr.Rows[2].offset)
}

func TestAddOnTopTwoRows(t *testing.T) {
	pr := NewPartialRendering("foo")
	pr.AddOnTop("bar").ShiftTopBy(-1)

	var expected = Nlnl(`
bar
 foo
`)
	assert.Equal(t, expected, pr.String())
}

func TestAddOnTopThreeRows(t *testing.T) {

	pr := NewPartialRendering("foo")
	pr.AddOnTop("bar").ShiftTopBy(-2)
	pr.AddOnTop("baz").ShiftTopBy(3)

	var expected = Nlnl(`
     baz
bar
  foo
`)
	assert.Equal(t, expected, pr.String())
}

func TestAlignDistanceNoShift(t *testing.T) {
	tests := []struct {
		leftValue  string
		rightValue string
		expected   int
	}{
		{"a12345", "bcdef", 6},
		{"123", "123", 3},
		{"123", "12", 3},
		{"123", "1", 3},
		{"12", "123", 2},
		{"1", "123", 1},
	}

	for _, tt := range tests {
		left := NewPartialRendering(tt.leftValue)
		right := NewPartialRendering(tt.rightValue)
		actual := AlignDistance(left, right)
		assert.Equal(t, tt.expected, actual)
	}
}

func TestAlignDistanceRightShifted(t *testing.T) {
	tests := []struct {
		leftValue   string
		rightPrefix string
		expected    int
	}{
		{"a123", "123", 7},
		{"a123", "12", 6},
		{"a123", "1", 5},
		{"a123", "", math.MinInt},
		{"12", "", math.MinInt},
		{"1", "", math.MinInt},
		{"", "", math.MinInt},
		{"", "1", math.MinInt},
		{"", "12", math.MinInt},
		{"", "123", math.MinInt},
		{"1", "123", 4},
		{"12", "123", 5},
		{"123", "123", 6},
	}

	for _, tt := range tests {
		left := NewPartialRendering(tt.leftValue)
		right := NewPartialRendering(tt.rightPrefix).ShiftTopBy(-len(tt.rightPrefix))

		actual := AlignDistance(left, right)
		assert.Equal(t, tt.expected, actual)
	}
}

func TestAlignDistanceBothShifted(t *testing.T) {
	tests := []struct {
		leftValue  string
		leftShift  int
		rightValue string
		rightShift int
		expected   int
	}{
		{"1", 0, "234", 0, 1},
		{"12", 0, "34", 0, 2},
		{"12", -1, "34", -1, 2},
		{"12", -2, "34", -2, 2},
		{"12", -2, "34", -1, 1},
		{"12", -2, "34", 0, 0},
		{"12", -3, "34", 0, -1},
		{"12", -4, "34", 0, -2},
		{"123", -3, "123", 0, 0},
		{"123", -2, "123", 0, 1},
		{"123", -1, "123", 0, 2},
		{"123", 0, "123", 0, 3},
		{"123", 0, "123", -1, 4},
		{"123", 0, "123", -2, 5},
		{"123", 0, "123", -3, 6},
		{"123", 0, "123", -4, 7},
		{"123", 0, "123", -5, 8},
		{"123", 1, "123", -5, 9},
		{"123", 2, "123", -5, 10},
		{"123", 3, "123", -5, 11},
		{"123", -3, "123", 1, -1},
		{"123", -3, "123", 2, -2},
		{"123", -3, "123", 3, -3},
		{"123", -4, "123", 3, -4},
		{"123", -5, "123", 3, -5},
	}

	for _, tt := range tests {
		left := NewPartialRendering(tt.leftValue).ShiftTopBy(tt.leftShift)
		right := NewPartialRendering(tt.rightValue).ShiftTopBy(tt.rightShift)

		actual := AlignDistance(left, right)
		assert.Equal(t, tt.expected, actual)
	}
}

func TestAlignDistance1(t *testing.T) {
	left := NewPartialRendering("1")
	right := NewPartialRendering("234")

	actual := AlignDistance(left, right) // 1234
	//                                    00     <- zero index positions, the distance is 1
	assert.Equal(t, 1, actual)
}

func TestAlignDistance2(t *testing.T) {
	left := NewPartialRendering("12")
	right := NewPartialRendering("34")

	actual := AlignDistance(left, right) // 1234
	//                                    0 0    <- zero index positions, the distance is 2
	assert.Equal(t, 2, actual) //
}

func TestAlignDistance3(t *testing.T) {
	left := NewPartialRendering("12").ShiftTopBy(-1)
	assert.Equal(t, "1", left.TopRow().Prefix())
	assert.Equal(t, "2", left.TopRow().Suffix())
	right := NewPartialRendering("34").ShiftTopBy(-1)
	assert.Equal(t, "3", right.TopRow().Prefix())
	assert.Equal(t, "4", right.TopRow().Suffix())

	actual := AlignDistance(left, right) // 1234
	//                                     0 0   <- zero index positions, the distance is 2
	assert.Equal(t, 2, actual)
}

func TestAlignDistance4(t *testing.T) {
	left := NewPartialRendering("12").NormalizeOffsetsRev()
	assert.Equal(t, "12", left.TopRow().Prefix())
	right := NewPartialRendering("34").NormalizeOffsets() // it's just for illustration purposes as the rendering is already normalized after creation with NewPartialRendering() function.
	assert.Equal(t, "34", right.TopRow().Suffix())

	actual := AlignDistance(left, right) // 1234
	//                                      0  <- when calculating the AlignDistance the strings are aligned without gaps and then the zero index positions of both strings occur at the same place.
	assert.Equal(t, 0, actual)
}

func TestAlignDistanceFromDoc(t *testing.T) {
	left := NewPartialRendering("ab123").ShiftTopBy(-1)
	right := NewPartialRendering("45cd").ShiftTopBy(-2)

	assert.Equal(t, "a", left.TopRow().Prefix())
	assert.Equal(t, "b123", left.TopRow().Suffix())

	assert.Equal(t, "45", right.TopRow().Prefix())
	assert.Equal(t, "cd", right.TopRow().Suffix())

	actual := AlignDistance(left, right)
	assert.Equal(t, 6, actual)
}

func TestAlignDistanceMultiRow(t *testing.T) {
	left := NewPartialRendering("12")
	left.AddOnTop("ab3").ShiftTopBy(-2)
	left.AddOnTop("c1234").ShiftTopBy(-1)

	right := NewPartialRendering("a")
	right.AddOnTop("12bc").ShiftTopBy(-2)
	right.AddOnTop("34d").ShiftTopBy(-2)
	right.AddOnTop("456d").ShiftTopBy(-3)

	actual := AlignDistance(left, right)
	assert.Equal(t, 7, actual)
}

func TestAlignDistanceMultiRow2(t *testing.T) {
	left := NewPartialRendering("12")
	left.AddOnTop("abc123").ShiftTopBy(-2)
	left.AddOnTop("c2345678").ShiftTopBy(-1)

	right := NewPartialRendering("a")
	right.AddOnTop("123456b3").ShiftTopBy(-6)
	right.AddOnTop("45678d").ShiftTopBy(-5)
	right.AddOnTop("9d").ShiftTopBy(-1)

	actual := AlignDistance(left, right)
	assert.Equal(t, 9, actual)
}

func TestNormalizeRev(t *testing.T) {
	pr := NewPartialRendering("foo")
	pr.AddOnTop("bar")
	pr.AddOnTop("baz")

	assert.Equal(t, "", pr.TopRow().Prefix())
	assert.Equal(t, "baz", pr.TopRow().Suffix())
	assert.Equal(t, "", pr.TopRow().Prefix())
	assert.Equal(t, "bar", pr.GetRow(1).Suffix())
	assert.Equal(t, "", pr.TopRow().Prefix())
	assert.Equal(t, "foo", pr.GetRow(2).Suffix())

	pr.NormalizeOffsetsRev()

	assert.Equal(t, "baz", pr.TopRow().Prefix())
	assert.Equal(t, "", pr.TopRow().Suffix())
	assert.Equal(t, "bar", pr.GetRow(1).Prefix())
	assert.Equal(t, "", pr.GetRow(1).Suffix())
	assert.Equal(t, "foo", pr.GetRow(2).Prefix())
	assert.Equal(t, "", pr.GetRow(2).Suffix())
}

func TestDoubleNormalizeRev(t *testing.T) {
	pr := NewPartialRendering("foo")
	pr.AddOnTop("bar")
	pr.AddOnTop("baz")

	assert.Equal(t, "", pr.TopRow().Prefix())
	assert.Equal(t, "baz", pr.TopRow().Suffix())
	assert.Equal(t, "", pr.TopRow().Prefix())
	assert.Equal(t, "bar", pr.GetRow(1).Suffix())
	assert.Equal(t, "", pr.TopRow().Prefix())
	assert.Equal(t, "foo", pr.GetRow(2).Suffix())

	pr.NormalizeOffsetsRev()

	assert.Equal(t, "baz", pr.TopRow().Prefix())
	assert.Equal(t, "", pr.TopRow().Suffix())
	assert.Equal(t, "bar", pr.GetRow(1).Prefix())
	assert.Equal(t, "", pr.GetRow(1).Suffix())
	assert.Equal(t, "foo", pr.GetRow(2).Prefix())
	assert.Equal(t, "", pr.GetRow(2).Suffix())

	pr.NormalizeOffsetsRev()

	assert.Equal(t, "baz", pr.TopRow().Prefix())
	assert.Equal(t, "", pr.TopRow().Suffix())
	assert.Equal(t, "bar", pr.GetRow(1).Prefix())
	assert.Equal(t, "", pr.GetRow(1).Suffix())
	assert.Equal(t, "foo", pr.GetRow(2).Prefix())
	assert.Equal(t, "", pr.GetRow(2).Suffix())
}

func TestNormalizeRevWithOffsets(t *testing.T) {
	pr := NewPartialRendering("foo")
	pr.AddOnTop("bar").ShiftTopBy(1)
	pr.AddOnTop("baz").ShiftTopBy(2)

	assert.Equal(t, "", pr.TopRow().Prefix())
	assert.Equal(t, "  baz", pr.TopRow().Suffix())
	assert.Equal(t, 2, pr.TopRow().StartIndex())
	assert.Equal(t, 4, pr.TopRow().EndIndex())

	assert.Equal(t, "", pr.GetRow(1).Prefix())
	assert.Equal(t, " bar", pr.GetRow(1).Suffix())
	assert.Equal(t, "", pr.GetRow(2).Prefix())
	assert.Equal(t, "foo", pr.GetRow(2).Suffix())

	pr.NormalizeOffsetsRev()

	assert.Equal(t, "", pr.TopRow().Suffix())
	assert.Equal(t, "baz", pr.TopRow().Prefix())
	assert.Equal(t, -3, pr.TopRow().StartIndex())
	assert.Equal(t, -1, pr.TopRow().EndIndex())

	assert.Equal(t, "bar ", pr.GetRow(1).Prefix())
	assert.Equal(t, "", pr.GetRow(1).Suffix())
	assert.Equal(t, "foo  ", pr.GetRow(2).Prefix())
	assert.Equal(t, "", pr.GetRow(2).Suffix())
}

func TestNormalize(t *testing.T) {
	pr := NewPartialRendering("foo")
	pr.AddOnTop("bar")
	pr.AddOnTop("baz")

	assert.Equal(t, "", pr.TopRow().Prefix())
	assert.Equal(t, "baz", pr.TopRow().Suffix())
	assert.Equal(t, "", pr.GetRow(1).Prefix())
	assert.Equal(t, "bar", pr.GetRow(1).Suffix())
	assert.Equal(t, "", pr.GetRow(2).Prefix())
	assert.Equal(t, "foo", pr.GetRow(2).Suffix())

	pr.NormalizeOffsetsRev()

	assert.Equal(t, "baz", pr.TopRow().Prefix())
	assert.Equal(t, "", pr.TopRow().Suffix())
	assert.Equal(t, "bar", pr.GetRow(1).Prefix())
	assert.Equal(t, "", pr.GetRow(1).Suffix())
	assert.Equal(t, "foo", pr.GetRow(2).Prefix())
	assert.Equal(t, "", pr.GetRow(2).Suffix())

	pr.NormalizeOffsets()

	assert.Equal(t, "", pr.TopRow().Prefix())
	assert.Equal(t, "baz", pr.TopRow().Suffix())
	assert.Equal(t, "", pr.GetRow(1).Prefix())
	assert.Equal(t, "bar", pr.GetRow(1).Suffix())
	assert.Equal(t, "", pr.GetRow(2).Prefix())
	assert.Equal(t, "foo", pr.GetRow(2).Suffix())
}

func TestDoubleNormalize(t *testing.T) {
	pr := NewPartialRendering("foo")
	pr.AddOnTop("bar")
	pr.AddOnTop("baz")

	assert.Equal(t, "", pr.TopRow().Prefix())
	assert.Equal(t, "baz", pr.TopRow().Suffix())
	assert.Equal(t, "", pr.GetRow(1).Prefix())
	assert.Equal(t, "bar", pr.GetRow(1).Suffix())
	assert.Equal(t, "", pr.GetRow(2).Prefix())
	assert.Equal(t, "foo", pr.GetRow(2).Suffix())

	pr.NormalizeOffsets()

	assert.Equal(t, "", pr.TopRow().Prefix())
	assert.Equal(t, "baz", pr.TopRow().Suffix())
	assert.Equal(t, "", pr.GetRow(1).Prefix())
	assert.Equal(t, "bar", pr.GetRow(1).Suffix())
	assert.Equal(t, "", pr.GetRow(2).Prefix())
	assert.Equal(t, "foo", pr.GetRow(2).Suffix())
}
