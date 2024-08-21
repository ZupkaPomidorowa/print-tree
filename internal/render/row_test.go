package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartIndex(t *testing.T) {
	tests := []struct {
		val        string
		startIndex int
		expected   int
	}{

		{"", 0, 0},
		{"", 1, 1},
		{"", -1, -1},
		{"123", 0, 0},
		{"123", 1, 1},
		{"123", 2, 2},
		{"123", 3, 3},
		{"123", -1, -1},
		{"123", -2, -2},
		{"123", -3, -3},
	}

	for _, test := range tests {
		rr := Row{val: test.val, offset: test.startIndex}
		assert.Equal(t, test.expected, rr.StartIndex())
	}
}

func TestEndIndex(t *testing.T) {
	tests := []struct {
		val        string
		startIndex int
		expected   int
	}{
		{"", 0, -1},
		{"", 1, 0},
		{"", -1, -2},
		{"123", 0, 2},
		{"123", 1, 3},
		{"123", 2, 4},
		{"123", 3, 5},
		{"123", -1, 1},
		{"123", -2, 0},
		{"123", -3, -1},
	}

	for _, test := range tests {
		rr := Row{val: test.val, offset: test.startIndex}
		assert.Equal(t, test.expected, rr.EndIndex())
	}
}

func TestPrefix(t *testing.T) {
	tests := []struct {
		val        string
		startIndex int
		expected   string
	}{
		{"123", 2, ""},
		{"123", 1, ""},
		{"123", 0, ""},
		{"123", -1, "1"},
		{"123", -2, "12"},
		{"123", -3, "123"},
		{"123", -4, "123 "},
		{"123", -5, "123  "},
	}

	for _, test := range tests {
		rr := Row{val: test.val, offset: test.startIndex}
		assert.Equal(t, test.expected, rr.Prefix())
	}
}

func TestSuffix(t *testing.T) {
	tests := []struct {
		val        string
		startIndex int
		expected   string
	}{
		{"123", 0, "123"},
		{"123", 1, " 123"},
		{"123", 2, "  123"},
		{"123", 3, "   123"},
		{"123", 4, "    123"},
		{"123", -1, "23"},
		{"123", -2, "3"},
		{"123", -3, ""},
		{"123", -4, ""},
	}

	for _, test := range tests {
		rr := Row{val: test.val, offset: test.startIndex}
		assert.Equal(t, test.expected, rr.Suffix())
	}
}
