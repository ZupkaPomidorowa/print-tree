## Description
An exercise in ascii-art tree printing in golang.

The drawing style assumes/requires the following constraints:
- The tree is binary
- The values are strings
- If two values are printed in the same row, there must be at least three spaces between them
- The connector lines connect to the values from the "diagonal" directions (top-left, top-right, bottom-left, bottom-right)
- Every connector line spans three rows and no two characters of the connector line are in the same column
- The tree is "minimal" i.e. the nodes and subtrees are printed as close to each other as possible considering the above constraints

Besides the above constraints, except of node values, the tree itself is printed using just four characters: a space, slash, backslash and underscore.

**Note**: The printer was initially intended to be used for rendering trees for arithmetical expressions, but I decided to make it more generic and use arbitrary strings as values. After all "atan(x)" may be a valid part of expression :)

## Pre-requisites
- Go 1.22 or later

## Usage
```bash
go run cmd/main.go
```

## Output
```
               root
              /    \
          ___/      \___
         /              \
        +                +
       / \              / \
      /   \            /   \
     /     \          /     \
    1       foo   5432       5
           /   \
      ____/     \____
     /               \
    +                 bar
   / \               /   \
  /   \             /     \
 /     \           /       \
2       345    6789         9
```

