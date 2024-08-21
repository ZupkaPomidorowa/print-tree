package main

import (
	"fmt"

	printer "github.com/ZupkaPomidorowa/print-tree"
)

func main() {
	fmt.Println(bigTree())
}

func bigTree() string {
	//Leaf nodes
	n1 := &printer.Node{
		Value: "1",
	}
	n5432 := &printer.Node{
		Value: "5432",
	}
	n5 := &printer.Node{
		Value: "5",
	}
	n2 := &printer.Node{
		Value: "2",
	}
	n345 := &printer.Node{
		Value: "345",
	}
	n6789 := &printer.Node{
		Value: "6789",
	}
	n9 := &printer.Node{
		Value: "9",
	}

	//Second-level nodes

	nPlus := &printer.Node{
		Value:      "+",
		LeftChild:  n2,
		RightChild: n345,
	}

	nBar := &printer.Node{
		Value:      "bar",
		LeftChild:  n6789,
		RightChild: n9,
	}

	// Third-level nodes

	nFoo := &printer.Node{
		Value:      "foo",
		LeftChild:  nPlus,
		RightChild: nBar,
	}

	leftRootChild := &printer.Node{
		Value:      "+",
		LeftChild:  n1,
		RightChild: nFoo,
	}
	rightRootChild := &printer.Node{
		Value:      "+",
		LeftChild:  n5432,
		RightChild: n5,
	}

	root := &printer.Node{
		Value:      "root",
		LeftChild:  leftRootChild,
		RightChild: rightRootChild,
	}

	return printer.PrintTree(root).String()
}
