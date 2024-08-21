package render

import (
	"math"
	"strings"
)

// Rendering is a list of rows of text without gaps.
// To handle arbitrary indentation without adding spaces to the actual strings, every row has an offset.
// The offset allows to shift every string to the left or to the right by an arbitrary amount.
// If the offset is zero the strings are aligned to the left.
// Use the String() method to get the rendering as a single string.
type Rendering struct {
	Rows     []*Row //reverse order!
	minIndex int    // the smallest offset of all rows (optimization)
}

// Creates a new empty rendering.
func NewEmptyRendering() *Rendering {
	return &Rendering{
		Rows: []*Row{},
	}
}

// Creates a new rendering with a single row.
func NewPartialRendering(val string) *Rendering {
	return &Rendering{
		Rows: []*Row{
			{
				val: val,
			},
		},
	}
}

// ShiftTopBy offsets the top row (only that row!) by the given amount.
func (pr *Rendering) ShiftTopBy(thatMuch int) *Rendering {
	lastRow := pr.Rows[len(pr.Rows)-1]
	lastRow.offset += thatMuch
	if lastRow.offset < pr.minIndex {
		pr.minIndex = lastRow.offset
	}
	return pr
}

// AddOnTop adds a new row on top of the rendering.
func (pr *Rendering) AddOnTop(vals ...string) *Rendering {
	if len(vals) == 0 {
		return pr
	}
	var val string
	if len(vals) == 1 {
		val = vals[0]
	} else {
		val = strings.Join(vals, "")
	}
	pr.Rows = append(pr.Rows, &Row{
		val: val,
	})
	return pr
}

func (pr *Rendering) TopRow() Row {
	return pr.GetRow(0)
}

func (pr *Rendering) GetRow(n int) Row {
	if n < 0 || n >= len(pr.Rows) {
		panic("invalid row index")
	}
	idx := (len(pr.Rows) - 1) - n
	return *pr.Rows[idx]
}

func (pr *Rendering) String() string {
	var sb strings.Builder
	for i := len(pr.Rows) - 1; i >= 0; i-- {
		sb.WriteString(Spaces(-pr.minIndex + pr.Rows[i].offset))
		sb.WriteString(pr.Rows[i].val)
		sb.WriteString("\n")
	}
	return sb.String()
}

// NormalizeOffsets ensures that the top row starts at the zero offset i.e: the value has an offset of zero.
// This is done by shifting all the rows by the same amount (if necessary), so that the relative difference between the offsets of the rows remains the same.
// The effect is that the Prefix() of the top row is empty and the Suffix() of the top row is the entire value.
// NormalizeOffsets and it's opposite operation NormalizeOffsetsRev() are very useful when joining two renderings using the JoinRenderings() function.
func (pr *Rendering) NormalizeOffsets() *Rendering {
	var topRow = pr.TopRow()
	reindexBy := -topRow.StartIndex()

	if reindexBy == 0 {
		return pr
	}

	for i := 0; i < len(pr.Rows); i++ {
		pr.Rows[i].offset += reindexBy
		if pr.Rows[i].offset < pr.minIndex {
			pr.minIndex = pr.Rows[i].offset
		}
	}

	return pr
}

// NormalizeOffsetsRev ensures that the top row ends at the zero offset i.e: the value has an offset so that the last character is at the offset -1.
// This is, in a way, an operation that is an opposite of NormalizeOffset(): the Prefix() of the top row becomes the entire value and Suffix() of the top row becomes empty.
func (pr *Rendering) NormalizeOffsetsRev() *Rendering {
	var topRow = pr.TopRow()
	reindexBy := -topRow.EndIndex() - 1

	for i := 0; i < len(pr.Rows); i++ {
		pr.Rows[i].offset += reindexBy
		if pr.Rows[i].offset < pr.minIndex {
			pr.minIndex = pr.Rows[i].offset
		}
	}
	return pr
}

// Reverse reverses the vertical order of the rows in the rendering.
func (pr *Rendering) Reverse() {

	if len(pr.Rows) == 0 {
		return
	}

	var tmpRows []*Row

	for i := len(pr.Rows) - 1; i >= 0; i-- {
		tmpRows = append(tmpRows, pr.Rows[i])
	}

	pr.Rows = tmpRows
}

// AlignDistance returns the distance between zero-offsets of two renderings such that if this distance is used to print both renderings side by side, they will be adjacent to each other but not overlapping.
// See the JoinRenderings() function for more details.
// The distance is calculated between corresponding rows of both renderings, starting from the top, and the maximum of these distances is returned.
// For every single row the distance is calculated as follows:
// We start with the left and right row, each having it's own offset.
// Imagine we then align both left and right value as close to each other as possible without overlapping: there is no gap between the last character of the left value and the first character of the right value.
// Note we're not changing the individual offsets of left and right rows, we're just aligning the string values next to each other.
// Now regardless of the individual offsets of the left and right row, we can find the number of characters between the positions indexed as zero in both rows.
// This number plus one is the distance for a given row.
// Single row examples:

// 1) No shift example with distance 6:
// Actual renderings:
//
// "                           left rendering              right rendering
// "zero index position:            0                           0
// "                                a12345                      bcdef
//
// After "aligning" the renderings close to each other without overlapping:
// distance:                         <-5->
// zero index position:             0     0
// "                                a12345bcdef
//
// 2) Rows with shift example with distance 6:
// Actual renderings:
//
// "                           left rendering              right rendering
// "zero index position:            0                           0
// "                               ab123                      45cd
//
// After "shifting" the renderings close to each other without overlapping:
// distance:                            <-5->
// zero index position:                0     0
// "                                  ab12345cd
//
// 3) Normalized rows:
// Actual renderings:
//
// "                           left rendering              right rendering
// "zero index position:            0                           0
// "                          a12345                            6789b
//
// After "aligning" the renderings close to each other without overlapping:
// distance:                           -->0<--
// zero index position (overlaps):        0
// "                                a123456789b
func AlignDistance(left, right *Rendering) int {
	if left == nil || right == nil {
		panic("invalid rendering: nil")
	}

	var indexLeft = len(left.Rows) - 1
	var indexRight = len(right.Rows) - 1

	var maxDistance = math.MinInt
	var maybeMaxDistance = func(val int) {
		if val > maxDistance {
			maxDistance = val
		}
	}

	for indexLeft >= 0 && indexRight >= 0 {
		var rowLeft = left.Rows[indexLeft]
		var rowRight = right.Rows[indexRight]

		if rowLeft.HasValue() && rowRight.HasValue() {
			if rowLeft.EndIndex() >= 0 && rowRight.StartIndex() >= 0 {
				maybeMaxDistance(rowLeft.EndIndex() - rowRight.StartIndex() + 1)
			} else if rowLeft.EndIndex() < 0 && rowRight.StartIndex() >= 0 {
				maybeMaxDistance(rowLeft.EndIndex() - rowRight.StartIndex() + 1)
			} else if rowLeft.EndIndex() < 0 && rowRight.StartIndex() < 0 {
				maybeMaxDistance(rowLeft.EndIndex() - rowRight.StartIndex() + 1)
			} else {
				// rowLeft.EndIndex() >= 0 && rowRight.StartIndex() < 0
				maybeMaxDistance(rowLeft.EndIndex() - rowRight.StartIndex() + 1)
			}
		}

		indexLeft--
		indexRight--
	}

	return maxDistance
}

// JoinRenderings joins two renderings into a single new rendering.
// It joins the corresponding rows of both renderings side by side, starting from the top, with a given distance between them.
// You can calculate the distance required to have both renderings adjacent to each other without overlapping using the AlignDistance() function.
// If you use a distance greater than the one calculated by AlignDistance(), the renderings will be further apart - additional spaces are inserted between the corresponding values.
// If you use a distance smaller than the one calculated by AlignDistance(), the resulting rendering will have a different structure than the original renderings - so don't do it.
func JoinRenderings(left, right *Rendering, distance int) *Rendering {
	if left == nil || right == nil {
		panic("invalid rendering: nil")
	}

	var result = &Rendering{
		minIndex: left.minIndex,
	}

	var jLeft = 0
	var kRight = 0

	for jLeft < len(left.Rows) || kRight < len(right.Rows) {
		var leftVal string
		var rightVal string
		var joinedIndex int
		var numSpaces int

		if jLeft < len(left.Rows) && kRight < len(right.Rows) {
			left := left.GetRow(jLeft)
			right := right.GetRow(kRight)
			leftVal = left.Value()
			rightVal = right.Value()
			numSpaces = distance - (left.EndIndex() + 1) + right.StartIndex()
			joinedIndex = left.StartIndex()
		} else if jLeft < len(left.Rows) && kRight >= len(right.Rows) {
			left := left.GetRow(jLeft)
			leftVal = left.Value()
			joinedIndex = left.StartIndex()
			rightVal = ""
			numSpaces = 0
		} else if jLeft >= len(left.Rows) && kRight < len(right.Rows) {
			leftVal = ""
			right := right.GetRow(kRight)
			rightVal = right.Value()
			joinedIndex = distance + right.StartIndex()
			numSpaces = 0
		} else {
			panic("unexpected branch")
		}

		mergedRow := leftVal + Spaces(numSpaces) + rightVal
		result.AddOnTop(mergedRow).ShiftTopBy(joinedIndex)

		jLeft++
		kRight++
	}

	result.Reverse()
	return result
}
