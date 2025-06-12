package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileInfo struct {
	Name     string
	Size     int64
	IsDir    bool
	Path     string
	Children []*FileInfo
}

type SortType int

const (
	SortByName SortType = iota
	SortBySize
)

func main() {
	var (
		sortBy  = flag.String("sort", "name", "Sort method: name (by name) or size (by size)")
		reverse = flag.Bool("reverse", false, "Reverse sort order")
	)
	
	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [directory]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  directory\t\tTarget directory path (default: current directory)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s\t\t\tShow current directory\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s /path/to/dir\t\tShow specified directory\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s . -sort size\t\tSort by size\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s . -sort name -reverse\tReverse sort by name\n", os.Args[0])
	}
	
	flag.Parse()

	// Get target directory
	var targetDir string
	if len(flag.Args()) > 0 {
		targetDir = flag.Args()[0]
	} else {
		targetDir = "."
	}

	// Parse sort type
	var sortType SortType
	switch strings.ToLower(*sortBy) {
	case "size":
		sortType = SortBySize
	default:
		sortType = SortByName
	}

	// Build file tree
	root, err := buildFileTree(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Sort file tree
	sortFileTree(root, sortType, *reverse)

	// Display results
	fmt.Printf("%s\n", targetDir)
	if len(root.Children) > 0 {
		for i, child := range root.Children {
			isLast := i == len(root.Children)-1
			var connector string
			if isLast {
				connector = "└── "
			} else {
				connector = "├── "
			}
			sizeStr := formatSize(child.Size)
			if child.IsDir {
				fmt.Printf("%s%s/ (%s)\n", connector, child.Name, sizeStr)
			} else {
				fmt.Printf("%s%s (%s)\n", connector, child.Name, sizeStr)
			}
			
			// Print child nodes
			if len(child.Children) > 0 {
				var newPrefix string
				if isLast {
					newPrefix = "    "
				} else {
					newPrefix = "│   "
				}
				for j, grandChild := range child.Children {
					isGrandChildLast := j == len(child.Children)-1
					printFileTree(grandChild, newPrefix, isGrandChildLast)
				}
			}
		}
	}
}

func buildFileTree(dirPath string) (*FileInfo, error) {
	info, err := os.Stat(dirPath)
	if err != nil {
		return nil, err
	}

	root := &FileInfo{
		Name:  filepath.Base(dirPath),
		Path:  dirPath,
		IsDir: info.IsDir(),
		Size:  info.Size(),
	}

	if !info.IsDir() {
		return root, nil
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return root, err
	}

	for _, entry := range entries {
		childPath := filepath.Join(dirPath, entry.Name())
		childInfo, err := entry.Info()
		if err != nil {
			continue // Skip inaccessible files
		}

		child := &FileInfo{
			Name:  entry.Name(),
			Path:  childPath,
			IsDir: entry.IsDir(),
			Size:  childInfo.Size(),
		}

		if entry.IsDir() {
			// Recursively process subdirectories
			subTree, err := buildFileTree(childPath)
			if err == nil {
				child.Children = subTree.Children
				child.Size = calculateDirSize(child)
			}
		}

		root.Children = append(root.Children, child)
	}

	// Calculate total directory size
	root.Size = calculateDirSize(root)
	return root, nil
}

func calculateDirSize(dir *FileInfo) int64 {
	if !dir.IsDir {
		return dir.Size
	}

	var totalSize int64
	for _, child := range dir.Children {
		totalSize += child.Size
	}
	return totalSize
}

func sortFileTree(root *FileInfo, sortType SortType, reverse bool) {
	if root == nil || len(root.Children) == 0 {
		return
	}

	// Recursively sort child directories
	for _, child := range root.Children {
		if child.IsDir {
			sortFileTree(child, sortType, reverse)
		}
	}

	// Sort current level
	sort.Slice(root.Children, func(i, j int) bool {
		a, b := root.Children[i], root.Children[j]
		
		// Directories first
		if a.IsDir != b.IsDir {
			return a.IsDir
		}

		var result bool
		switch sortType {
		case SortBySize:
			result = a.Size > b.Size // Size descending
		default: // SortByName
			result = strings.ToLower(a.Name) < strings.ToLower(b.Name) // Name ascending
		}

		if reverse {
			return !result
		}
		return result
	})
}

func printFileTree(node *FileInfo, prefix string, isLast bool) {
	if node == nil {
		return
	}

	// Print current node
	var connector string
	if prefix == "" {
		connector = ""
	} else if isLast {
		connector = "└── "
	} else {
		connector = "├── "
	}

	sizeStr := formatSize(node.Size)
	if node.IsDir {
		fmt.Printf("%s%s%s/ (%s)\n", prefix, connector, node.Name, sizeStr)
	} else {
		fmt.Printf("%s%s%s (%s)\n", prefix, connector, node.Name, sizeStr)
	}

	// Print child nodes
	if len(node.Children) > 0 {
		var newPrefix string
		if prefix == "" {
			newPrefix = ""
		} else if isLast {
			newPrefix = prefix + "    "
		} else {
			newPrefix = prefix + "│   "
		}

		for i, child := range node.Children {
			isChildLast := i == len(node.Children)-1
			printFileTree(child, newPrefix, isChildLast)
		}
	}
}

func formatSize(size int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
