package render

// Row represents a single row of text with an optional offset.
// The offset controls the horizontal position of the value within the row - see the Prefix() and Suffix() methods.
// To understand the idea of the offset, first let's imagine there's no offset at all and every row is rendered on some "virtual screen".
// Then the first character of every value would start at the same column with an index 0. The second character is at index 1, the third at index 2, and so on.
// Adding the offset to the row allows for shifting the row's value to the right or to the left by the value of the offset.
// Positive offset shifts the value to the right by the offset.
// Negative offset shifts to value to the left.
type Row struct {
	val    string
	offset int
}

func (rr Row) HasValue() bool {
	return rr.val != ""
}

func (rr Row) Length() int {
	return len(rr.val)
}

// StartIndex returns the relative position of the first character of the value.
// It is equal to the row's offset and it controls the behavior of Prefix() and Suffix() methods.
func (rr Row) StartIndex() int {
	return rr.offset
}

// EndIndex returns the position of the last character of the value.
// For a single-character string, the end index is the same as the start index.
// Note: For an empty string the end index is one less than the start index!
func (rr Row) EndIndex() int {
	return rr.offset + len(rr.val) - 1
}

func (rr Row) Value() string {
	return rr.val
}

// Prefix returns the part of the string before index 0 (it exist when the offset is negative).
func (rr Row) Prefix() string {
	if len(rr.val) == 0 {
		return ""
	}

	if rr.offset < 0 {
		if -rr.offset > rr.Length() {
			return rr.val + Spaces(-rr.offset-rr.Length())
		}
		return rr.val[:-rr.offset]
	}

	return ""
}

// Suffix returns the part of the string with non-negative indexes, i.e: it's a substring [startIndex..len(val)]
func (rr Row) Suffix() string {
	if rr.offset > 0 {
		return Spaces(rr.offset) + rr.val
	}
	if rr.offset == 0 {
		return rr.val
	}

	if rr.offset < 0 {
		if -rr.offset > rr.Length() {
			return ""
		}
		return rr.val[-rr.offset:]
	}

	return ""
}
