package mytree

import "fmt"

// used for building the FileTree
type FileTree struct {
	FileName         string
	IsDirectory bool
	Children    []*FileTree
}

func NewFileTree(path string, isDirectory bool) *FileTree {
	return &FileTree{path, isDirectory, nil}
}

// this is a bfs traversal
func TraverseTree(root *FileTree) {
	// A queue to hold the nodes to be visited
	queue := []*FileTree{root}

	for len(queue) > 0 {
		// Dequeue the first node
		node := queue[0]
		queue = queue[1:]

		fmt.Printf("Visiting: %s (Is directory: %t)\n", node.FileName, node.IsDirectory)

		// Enqueue all children of the current node
		for _, child := range node.Children {
			queue = append(queue, child)
		}
	}
}
