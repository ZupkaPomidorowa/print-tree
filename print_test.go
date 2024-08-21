package printer

import (
	"testing"

	"github.com/ZupkaPomidorowa/print-tree/internal/render"
	"github.com/stretchr/testify/assert"
)

func TestPrintTreeOneCharRoot(t *testing.T) {
	leftChild := &Node{
		Value: "left",
	}
	rightChild := &Node{
		Value: "right",
	}

	root := &Node{
		Value:      "+",
		LeftChild:  leftChild,
		RightChild: rightChild,
	}

	actual := PrintTree(root)
	expected := render.Nlnl(`
       +
      / \
     /   \
    /     \
left       right
`)
	assert.Equal(t, expected, actual.String())
}

func TestPrintTreeTwoCharsRoot(t *testing.T) {
	leftChild := &Node{
		Value: "left",
	}
	rightChild := &Node{
		Value: "right",
	}

	root := &Node{
		Value:      "++",
		LeftChild:  leftChild,
		RightChild: rightChild,
	}

	actual := PrintTree(root)
	expected := render.Nlnl(`
       ++
      /  \
     /    \
    /      \
left        right
`)

	assert.Equal(t, expected, actual.String())
}

func TestPrintTreeTreeCharsRoot(t *testing.T) {
	leftChild := &Node{
		Value: "left",
	}
	rightChild := &Node{
		Value: "right",
	}

	root := &Node{
		Value:      "foo",
		LeftChild:  leftChild,
		RightChild: rightChild,
	}

	actual := PrintTree(root)
	expected := render.Nlnl(`
       foo
      /   \
     /     \
    /       \
left         right
`)

	assert.Equal(t, expected, actual.String())
}

func TestPrintTreeTwoLevels1(t *testing.T) {

	tree := buildTree("root")
	actual := PrintTree(tree)
	expected := render.Nlnl(`
                 root
                /    \
           ____/      \____
          /                \
       foo                  bar
      /   \                /   \
     /     \              /     \
    /       \            /       \
left         right   left         right
`)

	//         ^ three spaces up there
	assert.Equal(t, expected, actual.String())
}

// Note the similarity between this test case and the previous one.
// You can see that the root node value is shorter by one, so you could expect the rest of the tree is also "narrower" by one.
// But it's not what happens: Actually the tree gets wider by one!
// To understand why, focus on the horizontal connectors below the root node.
// In the previous test case, the connectors were 4 characters long.
// If we reduce the root node length by one, we could try to leave the connectors as they are - 4 characters long and "shift" the right tree one character to the left.
// But it's not possible, because then the values "right" and "left" in the middle of the bottom row would end up separated by only 2 spaces, which is not allowed - we require at least 3 spaces there!
// If we can't leave the connectors as they are (and of course we can't make them shorter - for the same reason), we have to make them longer.
// The next possible length for the connectors is 5 characters, and that's what happens here.
// All of this is because we insist on symmetry and both connectors must have the same length.
// Of course the algorithm doesn't literally follow the described logic: It achieves the desired result via simple arithmetic.
// Nevertheless the described logic was a driving factor behind the mentioned arithmetic :)
func TestPrintTreeTwoLevels2(t *testing.T) {

	tree := buildTree("baz")
	actual := PrintTree(tree)
	expected := render.Nlnl(`
                  baz
                 /   \
           _____/     \_____
          /                 \
       foo                   bar
      /   \                 /   \
     /     \               /     \
    /       \             /       \
left         right    left         right
`)
	//         ^ four spaces up there
	assert.Equal(t, expected, actual.String())
}

func buildTree(rootNodeVal string) *Node {
	leftLeaf := &Node{
		Value: "left",
	}
	rightLeaf := &Node{
		Value: "right",
	}

	leftChild := &Node{
		Value:      "foo",
		LeftChild:  leftLeaf,
		RightChild: rightLeaf,
	}
	rightChild := &Node{
		Value:      "bar",
		LeftChild:  leftLeaf,
		RightChild: rightLeaf,
	}

	return &Node{
		Value:      rootNodeVal,
		LeftChild:  leftChild,
		RightChild: rightChild,
	}
}
