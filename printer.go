package printer

import (
	"github.com/ZupkaPomidorowa/print-tree/internal/render"
)

type Node struct {
	Value      string
	LeftChild  *Node
	RightChild *Node
}

func (n *Node) IsLeaf() bool {
	return n.LeftChild == nil && n.RightChild == nil
}

// Prints the given tree with root node at the top and children below it.
// Uses a slash/backslash/underscore for connector drawing and spaces for alignment.
// Returned Rendering is NOT normalized.
func PrintTree(curNode *Node) *render.Rendering {

	if curNode == nil || curNode.Value == "" {
		panic("nil or empty Node")
	}

	if curNode.IsLeaf() {
		return render.NewPartialRendering(curNode.Value)
	}

	// Print the left child. This will determine the position of the parent and the right child.
	leftChildRendering := PrintTree(curNode.LeftChild)

	// Rev-normalize the left child, because the top row value starts at the zero index, and we want the top value to END at the zero index.
	// This ensures that connector lines will be drawn correctly.
	leftChildRendering.NormalizeOffsetsRev()

	// Print the right child.
	rightChildRendering := PrintTree(curNode.RightChild).NormalizeOffsets()

	// Calculate the distance between the children. The distance corresponds to the number of characters between the zero index of the left child and the zero index of the right child,
	// when both child are placed as close as possible (touching but not overlapping).
	requiredChildDistance := render.AlignDistance(leftChildRendering, rightChildRendering) + 3 // 3 is to have some minimal space between the children values

	// minimal distance between the children depends on the drawing style and in our case it is related to the length of the current node (Fig. 1)
	minPossibleDistance := 6 + len(curNode.Value) + 1

	var distance int
	if minPossibleDistance >= requiredChildDistance {
		//just join the children using minimalDistance
		distance = minPossibleDistance - 1 //TODO: why -1?
	} else {
		// children require more space than minPossibleDistance. Join the children using 3 additiional spaces.
		distance = requiredChildDistance
	}

	// This is required to achieve symmetric rendering.
	if (distance-len(curNode.Value))%2 == 1 {
		distance += 1
	}

	result := render.JoinRenderings(leftChildRendering, rightChildRendering, distance)

	// Print the connectors. The connectors span three rows: lower connector row (closest to children), middle connector row and upper connector row (closest to parent).
	// lower row:
	lowerRowSpacesCnt := (distance - 2) // 2 is for the two slashes
	result.AddOnTop("/" + render.Spaces(lowerRowSpacesCnt) + "\\")

	// middle row:
	middleRowSpacesCnt := len(curNode.Value) + 2 // 2 is for the two slashes
	middleRowUnderscoreCount := (lowerRowSpacesCnt-middleRowSpacesCnt)/2 - 1
	result.AddOnTop(render.Underscores(middleRowUnderscoreCount) + "/" + render.Spaces(len(curNode.Value)+2) + "\\" + render.Underscores(middleRowUnderscoreCount)).ShiftTopBy(1)

	// upper row:
	shift := ((distance - len(curNode.Value)) / 2) - 1
	result.AddOnTop("/" + render.Spaces(len(curNode.Value)) + "\\").ShiftTopBy(shift)

	// parent value
	result.AddOnTop(curNode.Value).ShiftTopBy(shift + 1)

	return result
}

/*

------------------------------------------------------------
Fig. 1 - Minimal distance between children

      +
     / \
    /   \
   /     \
  1       2                  // <- in this case the children need exactly (len(parentValue) + 6 == 7) of spaces in between them.
   1234567                   //    This means the distance is 8: it takes 8 steps to go from the last character of the left child to the first character of the right child.


------------------------------------------------------------
Fig. 2

      foo
     /   \
    /     \
   /       \
  1         2                // <- in this case the children need exactly (len(parentValue) + 6 == 9) of spaces in between them.
   123456789                 //    This means the distance is 10: it takes 10 steps to go from the last character of the left child to the first character of the right child.


------------------------------------------------------------
Fig. 3

             +
            / \
         __/   \__
        /         \
       -           *         // <- in this case the children need more than len(parentValue) + 6 of spaces in between them. This is because their subtrees, when printed, require more space.
      / \         / \
     /   \       /   \
    /     \     /     \
   1       2   3       4
	 -->123<------------ required distance between any two nodes is 3 (a matter of drawing style)


------------------------------------------------------------
Fig. 4

                    *
                   / \
              ____/   \____
             /             \
            +               +
           / \             / \
          /   \           /   \
         /     \         /     \
        1       foo   432       5
               /   \
          ____/     \____                // <- this shows why we need three levels of connectors rows: If we just have two levels,
         /               \               //      then the horizontal "line" ("___") below the value "432" would be printed too close to this value.
        +                 bar            //      More space there results in a better look.
       / \               /   \
      /   \             /     \
     /     \           /       \
    2       345    6789         9

*/
