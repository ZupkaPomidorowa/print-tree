An exercise in ascii-art tree printing in golang.

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

